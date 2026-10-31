package src

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type InstallTask struct {
	DisplayName string
	Manager     string
	Packages    []string
	Script      string
}

// BuildPackagesTab creates the UI for the packages tab and returns a function to collect selected tasks
func BuildPackagesTab(installedPacman, installedFlatpak map[string]bool) (*container.Scroll, func() []InstallTask) {
	mainContentBox := container.NewVBox()

	catalogTitle := widget.NewLabelWithStyle("Installing packages/apps", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	mainContentBox.Add(catalogTitle)

	catalog := GetDefaultCatalog()

	for i := range catalog {
		category := &catalog[i]

		title := widget.NewLabelWithStyle("  "+category.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
		mainContentBox.Add(title)

		checksContainer := container.NewGridWrap(fyne.NewSize(280, 40))

		for j := range category.Packages {
			pkg := &category.Packages[j]

			allPackagesInstalled := true
			for _, pkgName := range pkg.Packages {
				if pkg.Manager == "pacman" {
					if !installedPacman[pkgName] {
						allPackagesInstalled = false
						break
					}
				} else if pkg.Manager == "flatpak" {
					if !installedFlatpak[pkgName] {
						allPackagesInstalled = false
						break
					}
				}
			}

			pkg.IsInstalled = allPackagesInstalled

			check := widget.NewCheck(pkg.DisplayName, nil)
			if pkg.IsInstalled {
				check.SetChecked(true)
				check.Disable()
			}

			pkg.Checkbox = check
			checksContainer.Add(check)
		}

		mainContentBox.Add(container.NewPadded(checksContainer))
	}

	// This function will be called by main.go when the "Install" button is clicked
	getTasks := func() []InstallTask {
		var tasks []InstallTask
		for _, cat := range catalog {
			for _, pkg := range cat.Packages {
				if pkg.Checkbox != nil && pkg.Checkbox.Checked && !pkg.IsInstalled {
					tasks = append(tasks, InstallTask{
						DisplayName: pkg.DisplayName,
						Manager:     pkg.Manager,
						Packages:    pkg.Packages,
					})
				}
			}
		}
		return tasks
	}

	return container.NewVScroll(mainContentBox), getTasks
}
