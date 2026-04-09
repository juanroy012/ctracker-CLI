package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type screen int

const (
	screenProjects screen = iota
	screenCommits
)

type model struct {
	screen        screen
	projectTable  table.Model
	commitTable   table.Model
	projects      []project
	selectedIndex int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

		switch m.screen {
		case screenProjects:
			switch msg.String() {
			case "enter":
				m.selectedIndex = m.projectTable.Cursor()
				if m.selectedIndex >= 0 && m.selectedIndex < len(m.projects) {
					m.commitTable = buildCommitTable(m.projects[m.selectedIndex])
					m.screen = screenCommits
				}
				return m, nil
			}
		case screenCommits:
			switch msg.String() {
			case "esc", "backspace":
				m.screen = screenProjects
				return m, nil
			}
			m.commitTable, cmd = m.commitTable.Update(msg)
			return m, cmd
		}
	}

	if m.screen == screenProjects {
		m.projectTable, cmd = m.projectTable.Update(msg)
	} else {
		m.commitTable, cmd = m.commitTable.Update(msg)
	}
	return m, cmd
}

func (m model) View() tea.View {
	switch m.screen {
	case screenProjects:
		return tea.NewView(baseStyle.Render(
			m.projectTable.View() + "\n enter: open commits | q: quit\n",
		))
	case screenCommits:
		return tea.NewView(baseStyle.Render(
			m.commitTable.View() + "\n esc/backspace: back",
		))
	default:
		return tea.NewView("Unknown Screen")
	}
}

func buildCommitTable(p project) table.Model {
	columns := []table.Column{
		{Title: "No", Width: 5},
		{Title: "Commit", Width: 80},
	}

	var rows []table.Row
	for i, c := range p.commit {
		if c == "" {
			continue
		}
		row := table.Row{
			strconv.Itoa(len(p.commit) - i - 1),
			c,
		}
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(85),
	)

	applyTableStyle(&t)
	return t
}

func applyTableStyle(t *table.Model) {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
}

func main() {
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	input, err := getUserInput(scanner)
	if err != nil {
		log.Fatal(err)
	}
	dir, err := scanGitDir(input)
	if err != nil {
		log.Fatal(err)
	}
	var projects []project
	for _, d := range dir {
		ds, err := initializeProject(d)
		if err != nil {
			log.Fatal(err)
		}
		projects = append(projects, ds)
	}
	columns := []table.Column{
		{Title: "No", Width: 5},
		{Title: "Project", Width: 20},
		{Title: "Commits", Width: 10},
		{Title: "Path", Width: 80},
	}
	var rows []table.Row
	for i, f := range projects {
		row := table.Row{
			strconv.Itoa(i + 1),
			f.name,
			strconv.Itoa(len(f.commit) - 1),
			f.path,
		}
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
		table.WithWidth(115),
	)

	applyTableStyle(&t)

	m := model{screenProjects,
		t,
		t,
		projects,
		0,
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}
}
