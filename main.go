package main

import (
	"fmt"
	"os/exec"
	"strings"

	"arch-quickstart-gui/src"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	statusLabel *widget.Label
	progressBar *widget.ProgressBarInfinite
	downloadBtn *widget.Button
)

func mergeTasks(tasks []src.InstallTask) []src.InstallTask {
    var rootScript strings.Builder
    var flatpakPkgs []string
    var names []string

    for _, t := range tasks {
        names = append(names, t.DisplayName)
        switch t.Manager {
        case "pacman":
            for _, p := range t.Packages {
                fmt.Fprintf(&rootScript, "pacman -S --noconfirm --needed %s\n", p)
            }
        case "bash":
            rootScript.WriteString(t.Script)
            rootScript.WriteString("\n")
        case "flatpak":
            flatpakPkgs = append(flatpakPkgs, t.Packages...)
        }
    }

    var merged []src.InstallTask
    if rootScript.Len() > 0 {
        merged = append(merged, src.InstallTask{
            DisplayName: strings.Join(names, ", "),
            Manager:     "bash",
            Script:      rootScript.String(),
        })
    }
    if len(flatpakPkgs) > 0 {
        merged = append(merged, src.InstallTask{
            DisplayName: "Flatpak apps",
            Manager:     "flatpak",
            Packages:    flatpakPkgs,
        })
    }
    return merged
}

func runInstallCommand(task src.InstallTask) error {
	var cmd *exec.Cmd

	switch task.Manager {
	case "pacman":
		args := append([]string{"pacman", "-S", "--noconfirm"}, task.Packages...)
		cmd = exec.Command("pkexec", args...)
	case "flatpak":
		// Ensure the flathub remote exists for user to avoid "No remote refs found for ‘flathub’" on fresh installs
		remoteCmd := exec.Command("flatpak", "remote-add", "--if-not-exists", "--user", "flathub", "https://dl.flathub.org/repo/flathub.flatpakrepo")
		_ = remoteCmd.Run()

		args := append([]string{"install", "--user", "-y", "flathub"}, task.Packages...)
		cmd = exec.Command("flatpak", args...)
	case "bash":
		cmd = exec.Command("pkexec", "bash", "-c", task.Script)
	default:
		return fmt.Errorf("unknown package manager")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func getInstalledPacmanPackages() map[string]bool {
	installed := make(map[string]bool)
	cmd := exec.Command("pacman", "-Qq")
	output, err := cmd.Output()
	if err != nil {
		return installed
	}
	for _, pkg := range strings.Split(string(output), "\n") {
		if pkg = strings.TrimSpace(pkg); pkg != "" {
			installed[pkg] = true
		}
	}
	return installed
}

func getInstalledFlatpakPackages() map[string]bool {
	installed := make(map[string]bool)
	cmdUser := exec.Command("flatpak", "list", "--user", "--app", "--columns=application")
	outUser, _ := cmdUser.Output()
	cmdSystem := exec.Command("flatpak", "list", "--system", "--app", "--columns=application")
	outSystem, _ := cmdSystem.Output()

	allOutput := string(outUser) + "\n" + string(outSystem)
	for _, pkg := range strings.Split(allOutput, "\n") {
		if pkg = strings.TrimSpace(pkg); pkg != "" {
			installed[pkg] = true
		}
	}
	return installed
}

func main() {
	a := app.NewWithID("com.T0rquemada.archquickstart")
	w := a.NewWindow("Arch Quickstart")
	w.Resize(fyne.NewSize(850, 750))

	// 1. Scan the system once at startup
	installedPacman := getInstalledPacmanPackages()
	installedFlatpak := getInstalledFlatpakPackages()

	// 2. Build the packages tab
	packagesTabUI, getPackagesTasks := src.BuildPackagesTab(installedPacman, installedFlatpak)

	// Single tab container
	tabs := container.NewAppTabs(
		container.NewTabItem("Packages & Apps", packagesTabUI),
	)

	// ==========================================
	// BOTTOM CONTROL BAR
	// ==========================================
	statusLabel = widget.NewLabel("Ready")
	statusLabel.Alignment = fyne.TextAlignCenter

	progressBar = widget.NewProgressBarInfinite()
	progressBar.Hide()

	downloadBtn = widget.NewButton("Download & Install Selected", func() {
		rawTasks := getPackagesTasks()

		if len(rawTasks) == 0 {
			statusLabel.SetText("No packages selected to install.")
			return
		}

		// Combine everything so the password is requested only once
		tasks := mergeTasks(rawTasks)

		downloadBtn.Disable()
		progressBar.Show()
		progressBar.Start()

		statusLabel.SetText("Installing packages... Please wait, this may take a while. Check terminal for details.")

		go func() {
			hasError := false

			for _, task := range tasks {
				fmt.Printf("Installing: %s...\n", task.DisplayName)
				err := runInstallCommand(task)

				if err != nil {
					fmt.Printf("Failed to install %s: %v\n", task.DisplayName, err)
					hasError = true
					break
				}
				fmt.Printf("Successfully installed %s!\n", task.DisplayName)
			}

			// Safely dispatch UI updates back to Fyne's main thread
			fyne.Do(func() {
				if !hasError {
					statusLabel.SetText("Installation complete! You can close the app.")
				} else {
					statusLabel.SetText("Installation failed. Please check the terminal output.")
				}

				progressBar.Stop()
				progressBar.Hide()
				downloadBtn.Enable()
			})
		}()
	})

	statusArea := container.NewVBox(
		statusLabel,
		progressBar,
		downloadBtn,
	)

	buttonRow := container.New(layout.NewPaddedLayout(), statusArea)

	mainLayout := container.NewBorder(
		nil,
		buttonRow,
		nil,
		nil,
		tabs,
	)

	w.SetContent(mainLayout)
	w.ShowAndRun()
}