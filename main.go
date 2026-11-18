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
	var yayPkgs []string
	var names []string

	for _, t := range tasks {
		names = append(names, t.DisplayName)
		switch t.Manager {
		case "pacman":
			for _, p := range t.Packages {
				fmt.Fprintf(&rootScript, "pacman -S --noconfirm --needed %s\n", p)
			}
		case "yay":
			yayPkgs = append(yayPkgs, t.Packages...)
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
	if len(yayPkgs) > 0 {
		merged = append(merged, src.InstallTask{
			DisplayName: "AUR packages (" + strings.Join(yayPkgs, ", ") + ")",
			Manager:     "yay",
			Packages:    yayPkgs,
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
		args := append([]string{"pacman", "-S", "--noconfirm", "--needed"}, task.Packages...)
		cmd = exec.Command("pkexec", args...)
	case "yay":
		args := append([]string{"-S", "--noconfirm", "--needed"}, task.Packages...)
		cmd = exec.Command("yay", args...)
	case "flatpak":
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

func runUninstallCommand(task src.InstallTask) error {
	var cmd *exec.Cmd

	switch task.Manager {
	case "pacman", "yay":
		args := append([]string{"pacman", "-Rns", "--noconfirm"}, task.Packages...)
		cmd = exec.Command("pkexec", args...)
	case "flatpak":
		args := append([]string{"uninstall", "--user", "-y"}, task.Packages...)
		cmd = exec.Command("flatpak", args...)
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

	installedPacman := getInstalledPacmanPackages()
	installedFlatpak := getInstalledFlatpakPackages()

	setStatus := func(msg string) {
		fyne.Do(func() {
			statusLabel.SetText(msg)
		})
	}

	onUninstall := func(task src.InstallTask, onSuccess func()) {
		go func() {
			err := runUninstallCommand(task)
			if err != nil {
				setStatus("Failed to uninstall " + task.DisplayName + ": " + err.Error())
				return
			}
			onSuccess()
		}()
	}

	packagesTabUI, getPackagesTasks, refreshPackagesUI := src.BuildPackagesTab(
		installedPacman,
		installedFlatpak,
		onUninstall,
		setStatus,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Packages & Apps", packagesTabUI),
	)

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

			// Update installed maps using rawTasks so all individual packages (Pacman, Yay, Flatpak) are marked
			if !hasError {
				for _, t := range rawTasks {
					for _, pkgName := range t.Packages {
						if t.Manager == "flatpak" {
							installedFlatpak[pkgName] = true
						} else {
							installedPacman[pkgName] = true
						}
					}
				}
			}

			fyne.Do(func() {
				if !hasError {
					statusLabel.SetText("Installation complete!")
				} else {
					statusLabel.SetText("Installation failed. Please check the terminal output.")
				}

				progressBar.Stop()
				progressBar.Hide()
				downloadBtn.Enable()
				
				// Redraws the UI instantly, turning checkboxes into delete icons
				refreshPackagesUI()
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