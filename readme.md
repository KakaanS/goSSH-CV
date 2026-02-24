# Terminal CV Viewer

This is a terminal-based CV viewer written in Go, using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework. It displays your CV interactively in the terminal, loading data from a simple `cv.json` file you can easily edit.

## Features

- Interactive terminal UI for browsing employers and projects
- Easy to customize with your own CV data (edit `cv.json`)
- Minimal dependencies, fast and portable

## Quick Start

```sh
go run main.go || go run .
```

### 1. Build

```sh
go build -o ssh-cv main.go
```

This will create an executable named `ssh-cv` in the project directory.

### 2. Run

```sh
./ssh-cv
```

### 3. Customize Your CV

Edit the `cv.json` file to add your own employers, projects, and skills. The structure is simple and easy to follow. Here is a minimal example:

```json
{
  "employers": [
    {
      "name": "Your Company",
      "projects": [
        {
          "slug": "project-1",
          "name": "Project Name",
          "period": "2024-01 - 2024-06",
          "description": "Short project description.",
          "technologies": ["Go", "React"]
        }
      ]
    }
  ]
}
```

Keep it simple! Add more employers or projects as needed, but you don't need to fill in every possible detail.

## Dependencies

This project uses the following Go libraries:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)

All dependencies are managed via Go modules (`go.mod`).
