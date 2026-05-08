package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
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
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Type        string    `json:"type"`
	Period      string    `json:"period"`
	Role        string    `json:"role,omitempty"`
	Description string    `json:"description,omitempty"`
	Projects    []Project `json:"projects,omitempty"`
}

type Contact struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type CVData struct {
	Contact   Contact    `json:"contact"`
	Skills    []string   `json:"skills"`
	Employers []Employer `json:"employers"`
}

type model struct {
	skillsFilter        textinput.Model
	spinner             spinner.Model
	viewport            viewport.Model
	width               int
	height              int
	state               string
	cvData              CVData
	selectedEmployerIdx int
	selectedProjectIdx  int
	selectedProject     *Project
	selectedEmployer    *Employer
	ready               bool
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = "Filter skills..."
	ti.CharLimit = 32
	ti.Width = 30

	return model{
		spinner:             s,
		state:               "loading",
		cvData:              loadCV(),
		selectedEmployerIdx: 0,
		selectedProjectIdx:  0,
		skillsFilter:        ti,
	}
}

func loadCV() CVData {
	file, _ := os.ReadFile("cv.json")
	var data CVData
	json.Unmarshal(file, &data)
	return data
}

func filterSkills(skills []string, query string) []string {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return skills
	}
	out := make([]string, 0, len(skills))
	for _, s := range skills {
		if strings.Contains(strings.ToLower(s), q) {
			out = append(out, s)
		}
	}
	return out
}

func (model model) Init() tea.Cmd {
	return tea.Batch(
		model.spinner.Tick,
		tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
			return doneMsg{}
		}),
	)
}

func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if model.state == "skills" {
		if kmsg, ok := msg.(tea.KeyMsg); ok {
			if kmsg.String() == "esc" {
				model.skillsFilter.Blur()
				model.skillsFilter.SetValue("")
				model.state = "menu"
				return model, nil
			}
		}
		var cmd tea.Cmd
		model.skillsFilter, cmd = model.skillsFilter.Update(msg)
		return model, cmd
	}

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		headerHeight := 6
		footerHeight := 3
		if !model.ready {
			model.viewport = viewport.New(msg.Width-10, msg.Height-headerHeight-footerHeight)
			model.ready = true
		} else {
			model.viewport.Width = msg.Width - 10
			model.viewport.Height = msg.Height - headerHeight - footerHeight
		}

	case doneMsg:
		model.state = "menu"

	case spinner.TickMsg:
		var cmd tea.Cmd
		model.spinner, cmd = model.spinner.Update(msg)
		return model, cmd

	case tea.KeyMsg:

		switch msg.String() {
		case "a":
			model.state = "about"

		case "p":
			model.state = "employer_select"
			model.selectedEmployerIdx = 0
			model.selectedProjectIdx = 0

		case "s":
			model.state = "skills"
			model.skillsFilter.SetValue("")
			return model, model.skillsFilter.Focus()
		case "q":
			return model, tea.Quit
		case "c":
			model.state = "contact"

		case "up", "k":
			if model.state == "employer_select" && model.selectedEmployerIdx > 0 {
				model.selectedEmployerIdx--
			} else if model.state == "projects" && model.selectedProjectIdx > 0 {
				model.selectedProjectIdx--
			} else {
				var cmd tea.Cmd
				model.viewport, cmd = model.viewport.Update(msg)
				return model, cmd
			}
		case "down", "j":
			if model.state == "employer_select" && model.selectedEmployerIdx < len(model.cvData.Employers)-1 {
				model.selectedEmployerIdx++
			} else if model.state == "projects" && model.selectedProjectIdx < len(model.cvData.Employers[model.selectedEmployerIdx].Projects)-1 {
				model.selectedProjectIdx++
			} else {
				var cmd tea.Cmd
				model.viewport, cmd = model.viewport.Update(msg)
				return model, cmd
			}
		case "enter":
			if model.state == "employer_select" {
				emp := &model.cvData.Employers[model.selectedEmployerIdx]
				if emp.Type == "tech" {
					model.state = "projects"
					model.selectedProjectIdx = 0
				} else {
					model.selectedEmployer = emp
					model.state = "employer_detail"
				}
			} else if model.state == "projects" {
				model.selectedProject = &model.cvData.Employers[model.selectedEmployerIdx].Projects[model.selectedProjectIdx]
				model.state = "project_detail"
			}
		case "esc":
			if model.state == "project_detail" {
				model.state = "projects"
			} else if model.state == "projects" {
				model.state = "employer_select"
			} else if model.state == "employer_detail" {
				model.state = "employer_select"
			} else {
				model.state = "menu"
			}
		}

	}

	return model, tea.Batch(cmds...)
}

var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true).Padding(0, 1)
	menuStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true).Padding(0, 1)
	footerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)

	contentStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Align(lipgloss.Left)

	techStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Italic(true).
			Bold(true)

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)

	cakeStyle = lipgloss.NewStyle().Align(lipgloss.Center)
)

func (model model) View() string {
	if !model.ready {
		return "Initializing..."
	}

	var headerText, menuText, bodyContent string

	const boxWidth = 70

	cakeStyle = lipgloss.NewStyle().Width(boxWidth).Align(lipgloss.Center)
	techStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Italic(true)

	loadingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#78e5a9")).
		Bold(true).
		Width(boxWidth).
		Align(lipgloss.Center)

	switch model.state {
	case "loading":
		headerText = "Oscar Wendt"
		bodyContent = loadingStyle.Render(fmt.Sprintf("\n\n%s Loading Experience...", model.spinner.View()))
	case "menu":
		headerText = "Oscar Wendt"
		menuText = "(a) About   (p) Projects   (s) Skills   (c) Contact   (q) Quit"
		bodyContent = cakeStyle.Render(cakeArt)
	case "cakeIsALie":
		headerText = "The Cake is a lie"
		bodyContent = cakeStyle.Render(easterEgg)
	case "about":
		headerText = "About"
		menuText = "(esc) Back"
		bodyContent = "Oscar joined Layer 10 with excellent references from Ericsson AB/Microwave, where he, over the course of nearly a year, was solely responsible for the modernization and redevelopment of a new, centralized system for test management and execution. \nHis time at Ericsson, which unfortunately came to an end due to downsizing in 2024, was preceded by studies in web development and security. During his studies, Oscar stood out as one of the few students who, despite limited prior experience, made significant progress particularly within the frontend domain. His genuine passion for programming, combined with experience from other industries such as service and sales, likely contributed to his steep learning curve. Oscar also has experience running his own business.  \n\nIn summary: Oscar is a web developer with a strong focus on modern frontend technologies. While his CV may not yet reflect many years in the field, he has shown remarkable potential for rapid growth toward a more senior role."
	case "contact":
		headerText = "Contact"
		menuText = "(esc) Back"
		c := model.cvData.Contact
		bodyContent = fmt.Sprintf(
			"\n%s\n%s\n\n%s %s\n%s %s",
			c.Name,
			lipgloss.NewStyle().Italic(true).Render(c.Role),
			labelStyle.Render("Phone:"), c.Phone,
			labelStyle.Render("Email:"), c.Email,
		)

	case "skills":
		headerText = "Skills"
		menuText = "(type to filter)   (esc) Back"
		filtered := filterSkills(model.cvData.Skills, model.skillsFilter.Value())
		var list string
		if len(filtered) == 0 {
			list = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true).Render("No skills match that filter.")
		} else {
			list = strings.Join(filtered, ", ")
		}
		count := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true).
			Render(fmt.Sprintf("%d / %d skills", len(filtered), len(model.cvData.Skills)))
		bodyContent = fmt.Sprintf("%s\n%s\n\n%s", model.skillsFilter.View(), count, list)
	case "employer_select":
		headerText = "Select Employer"
		menuText = "(↑/↓) Navigate   (enter) Select   (esc) Back"
		dividerEmitted := false
		for i, emp := range model.cvData.Employers {
			if emp.Type == "non-tech" && !dividerEmitted {
				bodyContent += dividerStyle.Render("  ── Previous (non-tech) employments ──") + "\n"
				dividerEmitted = true
			}
			cursor := "  "
			if i == model.selectedEmployerIdx {
				cursor = "> "
			}
			bodyContent += fmt.Sprintf("%s%s\n", cursor, emp.Name)
		}
	case "projects":
		headerText = model.cvData.Employers[model.selectedEmployerIdx].Name
		menuText = "(↑/↓) Navigate   (enter) Select   (esc) Back"
		for i, p := range model.cvData.Employers[model.selectedEmployerIdx].Projects {
			cursor := "  "
			if i == model.selectedProjectIdx {
				cursor = "> "
			}
			bodyContent += fmt.Sprintf("%s%s\n", cursor, p.Name)
		}
	case "project_detail":
		p := model.selectedProject
		headerText = p.Name
		menuText = "(↑/↓) Scroll   (esc) Back"

		techLabel := techStyle.Render("Technologies: ")
		bodyContent = fmt.Sprintf("%s\n\n%s\n\n%s%s",
			p.Period,
			p.Description,
			techLabel,
			strings.Join(p.Technologies, ", "))
	case "employer_detail":
		e := model.selectedEmployer
		headerText = e.Name
		menuText = "(esc) Back"
		location := e.Location
		if location == "" {
			location = "—"
		}
		role := e.Role
		if role == "" {
			role = "—"
		}
		bodyContent = fmt.Sprintf(
			"%s %s\n%s %s\n%s %s\n\n%s",
			labelStyle.Render("Period:  "), e.Period,
			labelStyle.Render("Location:"), location,
			labelStyle.Render("Role:    "), role,
			e.Description,
		)
	}

	centeredBody := contentStyle.Width(boxWidth).Render(bodyContent)

	model.viewport.SetContent(contentStyle.Width(boxWidth).Render(bodyContent))

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

	return lipgloss.Place(model.width, model.height, lipgloss.Center, lipgloss.Center, window)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
