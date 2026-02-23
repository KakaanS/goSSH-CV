package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type doneMsg struct{}

type CVData struct {
    Projects []struct {
        Name string `json:"name"`
    } `json:"projects"`
    Employers []struct {
        Name string `json:"name"`
    } `json:"employers"`
}

type model struct {
	spinner  spinner.Model
	width    int
	height   int
	state    string // "loading", "menu", "about", "projects", "employers"
	cvData   CVData
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		spinner: s,
		state:   "loading",
		cvData:  loadCV(),
	}
	}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
			return doneMsg{}
		}),
	)
}

func loadCV() CVData {
    file, _ := os.ReadFile("cv.json")
    var data CVData
    json.Unmarshal(file, &data)
    return data
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case doneMsg:
		m.state = "menu"
		return m, nil

	case tea.KeyMsg:
		switch msg.String(){
		
		case "q":
			return m, tea.Quit

		case "a":
			m.state="about"

		case "p": 
			m.state="projects"

		case "e":
			m.state="employers"

		case "esc": 
			m.state="menu"
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	}

	return m, nil
}

func (m model) viewLoading() string {
	return fmt.Sprintf("%s Oscar Wendt", m.spinner.View())
}

func (m model) viewMenu() string {
	return `
	Welcome to Oscar Wendt's interactive resume!
	
	Press: 
	(a) About
	(p) Projects
	(e) Employers
	(q) Quit
	`
}

func (m model) viewAbout() string {
	return `
	I'm a software developer with a passion for creating innovative solutions. 
	If you don't find me working or exploring new technologies, you can find me in the woods, sailing or if the mood is right, enjoying some counter-strike with friends. 

	Press (esc) to go back`
}

func (m model) viewProjects() string {
    s := "PROJECTS\n\n"
    for _, p := range m.cvData.Projects {
        s += fmt.Sprintf("- %s\n", p.Name)
    }
    s += "\nPress (esc) to go back"
    return s
}

func (m model) viewEmployers() string {
    s := "EMPLOYERS\n\n"
    for _, e := range m.cvData.Employers {
        s += fmt.Sprintf("- %s\n", e.Name)
    }
    s += "\nPress (esc) to go back"
    return s
}
var (
	// Style for the main container to match the terminal-shop vibe
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			Padding(0, 1)

	contentStyle = lipgloss.NewStyle().
			Padding(1, 2)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)
)

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	var content string
	var header string = " OSCAR WENDT "

	// Switch content based on state
	switch m.state {
	case "loading":
		content = fmt.Sprintf("%s Loading Experience...", m.spinner.View())
	case "menu":
		content = "Welcome to the interactive resume.\n\n(a) about  (p) projects  (e) employers  (q) quit"
	case "about":
		header = " ABOUT "
		content = "Software developer with a focus on modern frontend technologies [cite: 9].\nEnjoys sailing, the woods, and Counter-Strike[cite: 119]."
	case "projects":
		header = " PROJECTS "
		for _, p := range m.cvData.Projects {
			content += fmt.Sprintf("• %s\n", p.Name)
		}
	case "employers":
		header = " EMPLOYERS "
		for _, e := range m.cvData.Employers {
			content += fmt.Sprintf("• %s\n", e.Name)
		}
	}

	// Build the view components
	styledHeader := headerStyle.Render(header)
	styledContent := contentStyle.Render(content)
	
	var footer string
	if m.state != "menu" && m.state != "loading" {
		footer = footerStyle.Render("press (esc) to return")
	} else if m.state == "menu" {
		footer = footerStyle.Render("v1.0.0")
	}

	// Combine components into a vertical stack
	fullView := lipgloss.JoinVertical(
		lipgloss.Center,
		styledHeader,
		styledContent,
		footer,
	)

	// This is the magic part: it centers the entire "page" 
	// and ensures no overflow or scrolling ghosting.
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		fullView,
	)
}

func main() {
	
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}