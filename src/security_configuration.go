package src

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// buildSecurityTab creates the UI for the security tab and returns a function to collect selected tasks
func BuildSecurityTab(installedPacman map[string]bool) (*container.Scroll, func() []InstallTask) {
	mainContentBox := container.NewVBox()

	securityTitle := widget.NewLabelWithStyle("Security", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	mainContentBox.Add(securityTitle)

	// UFW Elements
	ufwCheck := widget.NewCheck("UFW (Uncomplicated Firewall)", nil)
	ufwLabel := widget.NewLabel("Simple static setup. Blocks all incoming.\nRecommended for desktop PCs.")
	ufwCard := container.NewVBox(ufwCheck, ufwLabel)

	// Firewalld Elements
	fwldCheck := widget.NewCheck("Firewalld", nil)
	fwldLabel := widget.NewLabel("Dynamic zone manager.\nRecommended for Laptops & Public Wi-Fi.")

	// Firewalld extra options
	removeSSHCheck := widget.NewCheck("Remove SSH from 'public' zone", nil)
	blockPingCheck := widget.NewCheck("Block Ping (ICMP) in 'public' zone", nil)
	dropDefaultCheck := widget.NewCheck("Set default zone to 'drop' (Max Security)", nil)

	removeSSHCheck.Disable()
	blockPingCheck.Disable()
	dropDefaultCheck.Disable()

	// Mutual exclusion and dependency logic
	ufwCheck.OnChanged = func(checked bool) {
		if checked {
			fwldCheck.SetChecked(false)
		}
	}

	fwldCheck.OnChanged = func(checked bool) {
		if checked {
			ufwCheck.SetChecked(false)
			removeSSHCheck.Enable()
			blockPingCheck.Enable()
			dropDefaultCheck.Enable()
		} else {
			removeSSHCheck.SetChecked(false)
			removeSSHCheck.Disable()
			blockPingCheck.SetChecked(false)
			blockPingCheck.Disable()
			dropDefaultCheck.SetChecked(false)
			dropDefaultCheck.Disable()
		}
	}

	if installedPacman["ufw"] {
		ufwCheck.SetChecked(true)
		ufwCheck.Disable()
		fwldCheck.Disable()
	} else if installedPacman["firewalld"] {
		fwldCheck.SetChecked(true)
		fwldCheck.Disable()
		ufwCheck.Disable()
		removeSSHCheck.Enable()
		blockPingCheck.Enable()
		dropDefaultCheck.Enable()
	}

	fwldCard := container.NewVBox(
		fwldCheck,
		fwldLabel,
		widget.NewLabelWithStyle("Firewalld Extra Rules:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true}),
		removeSSHCheck,
		blockPingCheck,
		dropDefaultCheck,
	)

	securityGrid := container.NewGridWithColumns(2, container.NewPadded(ufwCard), container.NewPadded(fwldCard))
	mainContentBox.Add(securityGrid)

	// Task collection function
	getTasks := func() []InstallTask {
		var tasks []InstallTask

		// UFW Task
		if ufwCheck.Checked && !installedPacman["ufw"] {
			tasks = append(tasks, InstallTask{
				DisplayName: "UFW Firewall",
				Manager:     "bash",
				Script:      "pacman -S --noconfirm ufw && ufw default deny incoming && ufw default allow outgoing && ufw --force enable && systemctl enable ufw",
			})
		}

		// Firewalld Task
		if fwldCheck.Checked {
			script := ""

			// 1. Install and enable Firewalld if not present
			if !installedPacman["firewalld"] {
				script += `pacman -S --noconfirm firewalld && systemctl enable --now firewalld`
			}

			// 2. Firewalld specific security rules
			var fwCommands []string
			if removeSSHCheck.Checked {
				fwCommands = append(fwCommands, `firewall-cmd --permanent --zone=public --remove-service=ssh`)
			}
			if blockPingCheck.Checked {
				fwCommands = append(fwCommands, `firewall-cmd --permanent --zone=public --add-icmp-block=echo-request`)
			}
			if dropDefaultCheck.Checked {
				fwCommands = append(fwCommands, `firewall-cmd --set-default-zone=drop`)
			}

			// Append firewall rules to the script and reload if rules were added
			if len(fwCommands) > 0 {
				fwCommands = append(fwCommands, `firewall-cmd --reload`) // Apply permanent changes
				if script != "" {
					script += " && "
				}
				script += strings.Join(fwCommands, " && ")
			}

			// Add task if there's a script to run
			if script != "" {
				displayName := "Firewalld Configuration"
				tasks = append(tasks, InstallTask{
					DisplayName: displayName,
					Manager:     "bash",
					Script:      script,
				})
			}
		}

		return tasks
	}

	return container.NewVScroll(mainContentBox), getTasks
}
