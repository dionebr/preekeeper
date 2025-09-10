package main

import (
    "fmt"
    "os"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type result struct {
    url    string
    status int
}

type model struct {
    results  []result
    quitting bool
}

type tickMsg result

var (
    okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))   // verde
    notFound   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))   // vermelho
    forbidden  = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))   // roxo
    neutral    = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))   // cinza
)

func main() {
    p := tea.NewProgram(initialModel())
    if err := p.Start(); err != nil {
        fmt.Println("Erro ao iniciar programa:", err)
        os.Exit(1)
    }
}

func initialModel() model {
    return model{
        results:  []result{},
        quitting: false,
    }
}

func (m model) Init() tea.Cmd {
    return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            m.quitting = true
            return m, tea.Quit
        }

    case tickMsg:
        m.results = append(m.results, result{
            url:    msg.url,
            status: msg.status,
        })
        return m, tickCmd()
    }
    return m, nil
}

func (m model) View() string {
    if m.quitting {
        return "Saindo...\n"
    }

    s := "🔍 Resultados do Scan (aperte 'q' para sair)\n\n"
    for _, r := range m.results {
        line := fmt.Sprintf("%s [%d]", r.url, r.status)
        s += "  • " + colorize(r.status, line) + "\n"
    }
    return s
}

func colorize(status int, text string) string {
    switch status {
    case 200:
        return okStyle.Render(text)
    case 403:
        return forbidden.Render(text)
    case 404:
        return notFound.Render(text)
    default:
        return neutral.Render(text)
    }
}

func tickCmd() tea.Cmd {
    return tea.Tick(time.Millisecond*600, func(t time.Time) tea.Msg {
        second := t.Second() % 3
        status := 200
        switch second {
        case 0:
            status = 200
        case 1:
            status = 403
        case 2:
            status = 404
        }
        return tickMsg{
            url:    fmt.Sprintf("http://site.com/%d", t.Second()),
            status: status,
        }
    })
}