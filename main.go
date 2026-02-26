package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
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
	input      			textinput.Model
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

	ti := textinput.New()
    ti.Placeholder = "Enter your email..."
    ti.CharLimit = 64
    ti.Width = 30

	return model{
		spinner:             s,
		state:               "loading",
		cvData:              loadCV(),
		selectedEmployerIdx: 0,
		selectedProjectIdx:  0,
		input:               ti,
	}
}

func loadCV() CVData {
	file, _ := os.ReadFile("cv.json")
	var data CVData
	json.Unmarshal(file, &data)
	return data
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

	if model.state == "contact" {
		if kmsg, ok := msg.(tea.KeyMsg); ok {
			switch kmsg.String() {
			case "enter":
			userEmail := model.input.Value() 
            if userEmail != "" {
                
                
                
                go func(email string) {
                    botToken := os.Getenv("botToken")
                    chatID := os.Getenv("chatID")
                    message := fmt.Sprintf("🚀 New CV Contact!\nEmail: %s\nSent via OpenClaw CLI", email)
                    
                    
                     apiURL := fmt.Sprintf(                                                                 
       				"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",                   
       				botToken, chatID, url.QueryEscape(message))
                    
                
                    _, _ = http.Get(apiURL) 
                }(userEmail)
                

                model.input.SetValue("")
                model.state = "sending"
                return model, tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
                    return doneMsg{}
                })
            }
			case "esc":
				model.state = "menu"
				model.input.Blur()
				return model, nil
			default:
				// Only update input, do not handle global shortcuts
			}
		}
		var cmd tea.Cmd
		model.input, cmd = model.input.Update(msg)
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
        case "q":
            return model, tea.Quit
        case "c":
            model.state = "contact"
            return model, model.input.Focus()
    
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
				model.state = "projects"
				model.selectedProjectIdx = 0
			} else if model.state == "projects" {
				model.selectedProject = &model.cvData.Employers[model.selectedEmployerIdx].Projects[model.selectedProjectIdx]
				model.state = "project_detail"
			}
		case "esc":
			if model.state == "project_detail" {
				model.state = "projects"
			} else if model.state == "projects" {
				model.state = "employer_select"
			} else {
				model.state = "menu"
			}
		}


	}

	return model, tea.Batch(cmds...)
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
		bodyContent =  "Oscar joined Layer 10 with excellent references from Ericsson AB/Microwave, where he, over the course of nearly a year, was solely responsible for the modernization and redevelopment of a new, centralized system for test management and execution. \nHis time at Ericsson, which unfortunately came to an end due to downsizing in 2024, was preceded by studies in web development and security. During his studies, Oscar stood out as one of the few students who, despite limited prior experience, made significant progress particularly within the frontend domain. His genuine passion for programming, combined with experience from other industries such as service and sales, likely contributed to his steep learning curve. Oscar also has experience running his own business.  \n\nIn summary: Oscar is a web developer with a strong focus on modern frontend technologies. While his CV may not yet reflect many years in the field, he has shown remarkable potential for rapid growth toward a more senior role."
	case "contact":
    	headerText = "Contact"
    	menuText = "(enter) Send   (esc) Back"
    	bodyContent = "\nReach out to me:\n\n" + model.input.View()
	case "sending":
	    headerText = "OpenClaw Assistant"
    	menuText = "Processing..."
    	bodyContent = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#78e5a9")).
        Render("\n\nOpenClaw is now connecting us via email.\nOne moment please...")
		
	case "skills":
		headerText = "Skills"
		menuText = "(esc) Back"
		bodyContent = "JavaScript, React, TypeScript, Node.js, Next.js, React Native, Expo, Firebase, SQL, Docker, Tor..."
	case "employer_select":
		headerText = "Select Employer"
		menuText = "(↑/↓) Navigate   (enter) Select   (esc) Back"
		for i, emp := range model.cvData.Employers {
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
	_ = godotenv.Load()
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
