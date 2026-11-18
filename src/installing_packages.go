package src

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type InstallTask struct {
	DisplayName string
	Manager     string
	Packages    []string
	Script      string
}

// BuildPackagesTab creates the UI for the packages tab and returns:
// 1. Scroll container
// 2. getTasks function for batch install
// 3. refreshUI function to redraw when installation/uninstallation happens
func BuildPackagesTab(
	installedPacman, installedFlatpak map[string]bool,
	onUninstall func(task InstallTask, onSuccess func()),
	setStatus func(string),
) (*container.Scroll, func() []InstallTask, func()) {

	mainContentBox := container.NewVBox()

	catalogTitle := widget.NewLabelWithStyle("Installing packages/apps", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	mainContentBox.Add(catalogTitle)

	catalog := GetDefaultCatalog()

	var rebuild func()

	buildUI := func() {
		mainContentBox.Objects = nil
		mainContentBox.Add(catalogTitle)

		for i := range catalog {
			category := &catalog[i]

			title := widget.NewLabelWithStyle("  "+category.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
			mainContentBox.Add(title)

			checksContainer := container.NewGridWrap(fyne.NewSize(280, 40))

			for j := range category.Packages {
				pkg := &category.Packages[j]

				allPackagesInstalled := true
				for _, pkgName := range pkg.Packages {
					if pkg.Manager == "flatpak" {
						if !installedFlatpak[pkgName] {
							allPackagesInstalled = false
							break
						}
					} else {
						if !installedPacman[pkgName] {
							allPackagesInstalled = false
							break
						}
					}
				}

				pkg.IsInstalled = allPackagesInstalled

				if pkg.IsInstalled {
					// Installed: show label and delete/uninstall button
					label := widget.NewLabel(pkg.DisplayName)
					deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
						task := InstallTask{
							DisplayName: pkg.DisplayName,
							Manager:     pkg.Manager,
							Packages:    pkg.Packages,
						}
						setStatus("Uninstalling " + pkg.DisplayName + "...")
						onUninstall(task, func() {
							for _, pkgName := range pkg.Packages {
								if pkg.Manager == "flatpak" {
									delete(installedFlatpak, pkgName)
								} else {
									delete(installedPacman, pkgName)
								}
							}
							setStatus("Successfully uninstalled " + pkg.DisplayName + "!")
							fyne.Do(func() {
								rebuild()
							})
						})
					})
					deleteBtn.Importance = widget.LowImportance

					itemBox := container.NewBorder(nil, nil, nil, deleteBtn, container.NewPadded(label))
					checksContainer.Add(itemBox)
				} else {
					// Not installed: show checkbox
					check := widget.NewCheck(pkg.DisplayName, nil)
					// Ensure checkbox state starts clean or reflects current check state if needed
					pkg.Checkbox = check
					checksContainer.Add(check)
				}
			}

			mainContentBox.Add(container.NewPadded(checksContainer))
		}
		mainContentBox.Refresh()
	}

	rebuild = buildUI
	buildUI()

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

	return container.NewVScroll(mainContentBox), getTasks, buildUI
}