package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type doneMsg struct{}

type Project struct {
	Slug         string   `json:"slug"`
	Name         string   `json:"name"`
	Period       string   `json:"period"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
}

type Employer struct {
	Name     string    `json:"name"`
	Projects []Project `json:"projects"`
}

type CVData struct {
	Employers []Employer `json:"employers"`
}

type model struct {
	spinner             spinner.Model
	viewport            viewport.Model
	width               int
	height              int
	state               string
	cvData              CVData
	selectedEmployerIdx int
	selectedProjectIdx  int
	selectedProject     *Project
	ready               bool
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		spinner:             s,
		state:               "loading",
		cvData:              loadCV(),
		selectedEmployerIdx: 0,
		selectedProjectIdx:  0,
	}
}

func loadCV() CVData {
	file, _ := os.ReadFile("cv.json")
	var data CVData
	json.Unmarshal(file, &data)
	return data
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
			return doneMsg{}
		}),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerHeight := 6 
		footerHeight := 3
		if !m.ready {
			m.viewport = viewport.New(msg.Width-10, msg.Height-headerHeight-footerHeight)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 10
			m.viewport.Height = msg.Height - headerHeight - footerHeight
		}

	case doneMsg:
		m.state = "menu"

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "a":
			m.state = "about"
		case "p":
			m.state = "employer_select"
		case "s":
			m.state = "skills"
		case "l": 
			m.state = "cakeIsALie"
		case "up":
			if m.state == "employer_select" && m.selectedEmployerIdx > 0 {
				m.selectedEmployerIdx--
			} else if m.state == "projects" && m.selectedProjectIdx > 0 {
				m.selectedProjectIdx--
			} else {
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		case "down":
			if m.state == "employer_select" && m.selectedEmployerIdx < len(m.cvData.Employers)-1 {
				m.selectedEmployerIdx++
			} else if m.state == "projects" && m.selectedProjectIdx < len(m.cvData.Employers[m.selectedEmployerIdx].Projects)-1 {
				m.selectedProjectIdx++
			} else {
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		case "enter":
			if m.state == "employer_select" {
				m.state = "projects"
				m.selectedProjectIdx = 0
			} else if m.state == "projects" {
				m.selectedProject = &m.cvData.Employers[m.selectedEmployerIdx].Projects[m.selectedProjectIdx]
				m.state = "project_detail"
			}
		case "esc":
			if m.state == "project_detail" {
				m.state = "projects"
			} else if m.state == "projects" {
				m.state = "employer_select"
			} else {
				m.state = "menu"
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

var (
	headerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Padding(0, 1)
	menuStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true).Padding(0, 1)
	footerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	
	
	contentStyle = lipgloss.NewStyle().
		Padding(0, 2).
		Align(lipgloss.Left) 

	
	techStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Italic(true).
		Bold(true)
        
    
    cakeStyle = lipgloss.NewStyle().Align(lipgloss.Center)
)

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var headerText, menuText, bodyContent string

	
	
	const boxWidth = 70 

	cakeStyle := lipgloss.NewStyle().Width(boxWidth).Align(lipgloss.Center)
	techStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Italic(true)
	loadingStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#78e5a9")).
        Bold(true).
        Width(boxWidth). 
        Align(lipgloss.Center)
	
	switch m.state {
	case "loading":
		headerText = "Oscar Wendt"
		bodyContent = loadingStyle.Render(fmt.Sprintf("\n\n%s Loading Experience...", m.spinner.View()))
	case "menu":
		headerText = "Oscar Wendt"
		menuText = "(a) About   (p) Projects   (s) Skills   (q) Quit"
		bodyContent = cakeStyle.Render(cakeArt)
	case "cakeIsALie": 
		headerText = "The Cake is a lie"
		bodyContent = cakeStyle.Render(easterEgg)
	case "about":
		headerText = "About"
		menuText = "(esc) Back"
		bodyContent = "Oscar joined Layer 10 with excellent references from Ericsson...\n\nHe has shown remarkable potential for rapid growth."
	case "skills":
		headerText = "Skills"
		menuText = "(esc) Back"
		bodyContent = "JavaScript, React, TypeScript, Node.js, Next.js, React Native, Expo, Firebase, SQL, Docker, Tor..."
	case "employer_select":
		headerText = "Select Employer"
		menuText = "(↑/↓) Navigate   (enter) Select   (esc) Back"
		for i, emp := range m.cvData.Employers {
			cursor := "  "
			if i == m.selectedEmployerIdx {
				cursor = "> "
			}
			bodyContent += fmt.Sprintf("%s%s\n", cursor, emp.Name)
		}
	case "projects":
		headerText = m.cvData.Employers[m.selectedEmployerIdx].Name
		menuText = "(↑/↓) Navigate   (enter) Select   (esc) Back"
		for i, p := range m.cvData.Employers[m.selectedEmployerIdx].Projects {
			cursor := "  "
			if i == m.selectedProjectIdx {
				cursor = "> "
			}
			bodyContent += fmt.Sprintf("%s%s\n", cursor, p.Name)
		}
	case "project_detail":
		p := m.selectedProject
		headerText = p.Name
		menuText = "(↑/↓) Scroll   (esc) Back"
		
		
		techLabel := techStyle.Render("Technologies: ")
		bodyContent = fmt.Sprintf("%s\n\n%s\n\n%s%s", 
			p.Period, 
			p.Description, 
			techLabel, 
			strings.Join(p.Technologies, ", "))
	}
	
	centeredBody := contentStyle.Width(boxWidth).Render(bodyContent) 
	
m.viewport.SetContent(contentStyle.Width(boxWidth).Render(bodyContent))

	header := headerStyle.Render(headerText)
	menu := menuStyle.Render(menuText)
	footer := footerStyle.Render("Oscar Wendt CV v1.0.0")	

    uiBlock := lipgloss.JoinVertical(
        lipgloss.Center, 
        header,
        menu,
        "",
        centeredBody, 
        "",
        footer,
    )

  window := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Render(uiBlock)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, window)
}



func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}