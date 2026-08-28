# Arch Quickstart GUI

## About
### Functions
1. Installing some apps
2. Install and add some rules for security (firewalld)

## Usage
Run the following command in your terminal to download and launch the GUI:
```bash
curl -LO https://github.com/T0rquemada/arch-quickstart-gui/releases/latest/download/arch-quickstart-gui && chmod +x arch-quickstart-gui && ./arch-quickstart-gui
```

## Prerequisites for development
- Golang

## Development commands
Run the program:
```sh
go run .
```
We use "." instead "main.go" because use multiple files.

## ToDo
- [ ] Flutter
- [ ] Update packages button
- [ ] Visually divide categories
- [X] Let two columns or more in row (Flexible design)
- [X] Launch by command in new system
- [X] Progress bar & status
- [ ] Debloatware for GNOME DE (maps, weather, contacts)
