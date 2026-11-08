package src

import "fyne.io/fyne/v2/widget"

// Software represents a single installable unit (which may contain multiple packages)
type Software struct {
	DisplayName string
	Manager     string   // "pacman" or "flatpak"
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
					DisplayName: "Zed Editor (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"zed"},
				},
				{
					DisplayName: "C / C++ Build Tools (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"base-devel", "gdb"},
				},
				{
					DisplayName: "CMake (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"cmake"},
				},
				{
					DisplayName: "OpenCV (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"opencv"},
				},
				{
					DisplayName: "Python & uv (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"python", "python-pip", "uv"},
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
					DisplayName: "Zig (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"zig"},
				},
				{
					DisplayName: "Kotlin (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"kotlin"},
				},
				{
					DisplayName: "Clojure (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"clojure"},
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
			Title: "Containers & Virtualization",
			Packages: []Software{
				{
					DisplayName: "Docker Engine & Compose (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"docker", "docker-compose"},
				},
				{
					DisplayName: "VirtualBox (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"virtualbox"},
				},
				{
					DisplayName: "GNOME Boxes (Flatpak)",
					Manager:     "flatpak",
					Packages:    []string{"org.gnome.Boxes"},
				},
			},
		},
		{
			Title: "AI & Machine Learning",
			Packages: []Software{
				{
					DisplayName: "Ollama (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"ollama"},
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
					Packages:    []string{"com.orcaslicer.OrcaSlicer"},
				},
			},
		},
		{
			Title: "Gaming",
			Packages: []Software{
				{
					DisplayName: "Steam (Flatpak)",
					Manager:     "flatpak",
					Packages:    []string{"com.valvesoftware.Steam"},
				},
				{
					DisplayName: "Wine (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"wine"},
				},
			},
		},
		{
			Title: "Downloads & P2P",
			Packages: []Software{
				{
					DisplayName: "qBittorrent (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"qbittorrent"},
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
					DisplayName: "OBS Studio (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"obs-studio"},
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
		{
			Title: "System & Snapshots",
			Packages: []Software{
				{
					DisplayName: "Snapper + auto-snapshots (Pacman)",
					Manager:     "pacman",
					Packages:    []string{"snapper", "snap-pac", "grub-btrfs"},
				},
			},
		},
	}
}