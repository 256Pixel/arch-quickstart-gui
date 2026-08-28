package src

import "fyne.io/fyne/v2/widget"

// Software represents a single installable unit (which may contain multiple packages)
type Software struct {
	DisplayName string
	Manager     string   // "pacman", "flatpak", or "vscode"
	Packages    []string // List of package IDs to install
	IsInstalled bool
	Checkbox    *widget.Check
}

// Category groups software together
type Category struct {
	Title    string
	Packages []Software
}

// GetDefaultCatalog returns the predefined software list
func GetDefaultCatalog() []Category {
	return []Category{
		{
			Title: "Browsers",
			Packages: []Software{
				{
					DisplayName: "Brave (Flatpak)",
					Manager:     "flatpak",
					Packages:    []string{"com.brave.Browser"},
				},
				{
					DisplayName: "Firefox (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"firefox"},
				},
				{
					DisplayName: "LibreWolf (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"librewolf"},
				},
			},
		},
		{
			Title: "Development & Languages",
			Packages: []Software{
				{
					DisplayName: "Git (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"git"},
				},
				{
					DisplayName: "VS Code OSS Build (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"code"},
				},
				{
					DisplayName: "Zed Editor (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"zed"},
				},
				{
					DisplayName: "C / C++ Build Tools (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"base-devel", "cmake", "gdb"},
				},
				{
					DisplayName: "Python (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"python", "python-pip"},
				},
				{
					DisplayName: "Go (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"go"},
				},
				{
					DisplayName: "Rust (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"rust"},
				},
				{
					DisplayName: "Dart (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"dart"},
				},
				{
					DisplayName: "JavaScript / Node.js (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"nodejs", "npm"},
				},
				{
					DisplayName: "PHP (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"php", "composer"},
				},
			},
		},
		{
			Title: "VS Code Extensions",
			Packages: []Software{
				{
					DisplayName: "C/C++ Tools (VS Code)",
					Manager:     "vscode",
					Packages:    []string{"ms-vscode.cpptools"},
				},
				{
					DisplayName: "Python & Pylance (VS Code)",
					Manager:     "vscode",
					Packages:    []string{"ms-python.python", "ms-python.vscode-pylance"},
				},
				{
					DisplayName: "Go (VS Code)",
					Manager:     "vscode",
					Packages:    []string{"golang.go"},
				},
				{
					DisplayName: "PlatformIO IDE (VS Code)",
					Manager:     "vscode",
					Packages:    []string{"platformio.platformio-ide"},
				},
			},
		},
		{
			Title: "Containers & Virtualization",
			Packages: []Software{
				{
					DisplayName: "Docker Engine & Compose (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"docker", "docker-compose"},
				},
			},
		},
		{
			Title: "Hardware & Maker",
			Packages: []Software{
				{
					DisplayName: "Bambu Studio (Flatpak)",
					Manager:     "flatpak",
					Packages:    []string{"com.bambulab.BambuStudio"},
				},
				{
					DisplayName: "OrcaSlicer (Flatpak)",
					Manager:     "flatpak",
					Packages:    []string{"com.softfever3d.OrcaSlicer"},
				},
			},
		},
		{
			Title: "Security & Privacy",
			Packages: []Software{
				{
					DisplayName: "1Password (Flatpak)",
					Manager:     "flatpak",
					Packages:    []string{"com.onepassword.OnePassword"},
				},
				{
					DisplayName: "KeePassXC (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"keepassxc"},
				},
				{
					DisplayName: "Bitwarden (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"bitwarden"},
				},
				{
					DisplayName: "VeraCrypt (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"veracrypt"},
				},
			},
		},
		{
			Title: "Office & Productivity",
			Packages: []Software{
				{
					DisplayName: "LibreOffice Writer (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"libreoffice-fresh"},
				},
				{
					DisplayName: "LibreOffice Calc (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"libreoffice-fresh"},
				},
				{
					DisplayName: "Gnumeric (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"gnumeric"},
				},
			},
		},
		{
			Title: "Personal",
			Packages: []Software{
				{
					DisplayName: "VLC Media Player (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"vlc"},
				},
				{
					DisplayName: "Signal Desktop (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"signal-desktop"},
				},
				{
					DisplayName: "Discord (Flatpak)",
					Manager:     "flatpak",
					Packages:    []string{"com.discordapp.Discord"},
				},
				{
					DisplayName: "Obsidian (Flatpak)",
					Manager:     "flatpak",
					Packages:    []string{"md.obsidian.Obsidian"},
				},
			},
		},
	}
}
