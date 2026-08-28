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

func runInstallCommand(task src.InstallTask) error {
	var cmd *exec.Cmd

	switch task.Manager {
	case "pacman":
		args := append([]string{"pacman", "-S", "--noconfirm"}, task.Packages...)
		cmd = exec.Command("pkexec", args...)
	case "flatpak":
		args := append([]string{"install", "--user", "-y", "flathub"}, task.Packages...)
		cmd = exec.Command("flatpak", args...)
	case "vscode":
		args := []string{}
		for _, ext := range task.Packages {
			args = append(args, "--install-extension", ext)
		}
		cmd = exec.Command("code", args...)
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

func getInstalledVSCodeExtensions() map[string]bool {
	installed := make(map[string]bool)
	cmd := exec.Command("code", "--list-extensions")
	output, err := cmd.Output()
	if err != nil {
		return installed
	}

	for _, ext := range strings.Split(string(output), "\n") {
		if ext = strings.TrimSpace(ext); ext != "" {
			installed[strings.ToLower(ext)] = true
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
	installedVSCode := getInstalledVSCodeExtensions()

	// 2. Build the tabs and get the task collector functions
	packagesTabUI, getPackagesTasks := src.BuildPackagesTab(installedPacman, installedFlatpak, installedVSCode)
	securityTabUI, getSecurityTasks := src.BuildSecurityTab(installedPacman)

	// Create a Tab widget with our two files' UIs
	tabs := container.NewAppTabs(
		container.NewTabItem("Packages & Apps", packagesTabUI),
		container.NewTabItem("Security", securityTabUI),
	)

	// ==========================================
	// BOTTOM CONTROL BAR
	// ==========================================
	statusLabel = widget.NewLabel("Ready")
	statusLabel.Alignment = fyne.TextAlignCenter

	progressBar = widget.NewProgressBarInfinite()
	progressBar.Hide()

	downloadBtn = widget.NewButton("Download & Install Selected", func() {
		var tasks []src.InstallTask

		// Combine tasks from both tabs
		tasks = append(tasks, getPackagesTasks()...)
		tasks = append(tasks, getSecurityTasks()...)

		if len(tasks) == 0 {
			statusLabel.SetText("No packages selected to install.")
			return
		}

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

			if !hasError {
				statusLabel.SetText("Installation complete! You can close the app.")
			} else {
				statusLabel.SetText("Installation failed. Please check the terminal output.")
			}

			progressBar.Stop()
			progressBar.Hide()
			downloadBtn.Enable()
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
