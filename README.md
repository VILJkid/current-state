# current-state

A lightweight Terminal User Interface (TUI) application to display system information including user, memory, and disk statistics.

## Features

- **User Information**: Display the currently logged-in user
- **Memory Stats**: View available system memory
- **Disk Space**: Check disk space usage
- **Interactive TUI**: Navigate with keyboard shortcuts
- **Cross-platform**: Runs on Linux, macOS, and Windows

## Prerequisites

- **Go 1.25** or higher

## Installation

### Option 1: Using `go install`

```bash
go install github.com/VILJkid/current-state/cmd@latest
```

### Option 2: From Source

```bash
# Clone the repository
git clone https://github.com/VILJkid/current-state.git
cd current-state

# Build the binary
make build

# Run the application
make run
```

Or build and run in one command:

```bash
make all
```

## Usage

Simply run the application:

```bash
current-state
```

### Keyboard Shortcuts

| Key | Action |
| :---: | :--: |
| `m` | Select main menu |
| `a` | View available memory |
| `b` | View disk space |
| `c` | View current user |
| `q` | Quit the application |

## Dependencies

- [tview](https://github.com/rivo/tview) - Rich interactive TUI framework
- [tcell](https://github.com/gdamore/tcell) - Terminal handling
- [gopsutil](https://github.com/shirou/gopsutil) - System and process utilities

## Inspiration

This project was inspired by [Gerald Yerden](https://github.com/devhulk/)'s [YouTube video](https://youtu.be/9bDN2rrf-Pw). Thanks for creating that awesome tutorial.

## License

This project is open source and available under the MIT License.
