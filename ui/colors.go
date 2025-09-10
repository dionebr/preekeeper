package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Status colors
	StatusOK        = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))   // verde
	StatusRedirect  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))   // amarelo  
	StatusClient    = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))   // roxo
	StatusServer    = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))   // vermelho
	StatusNeutral   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))   // cinza
	
	// UI elements
	HeaderStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	BorderStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	InfoStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	SuccessStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	ErrorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	ProgressStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	
	// Banner
	BannerStyle     = lipgloss.NewStyle().
						Foreground(lipgloss.Color("5")).
						Bold(true).
						Border(lipgloss.RoundedBorder()).
						BorderForeground(lipgloss.Color("6")).
						Padding(0, 1)
)

func GetStatusColor(status int) lipgloss.Style {
	switch {
	case status >= 200 && status < 300:
		return StatusOK
	case status >= 300 && status < 400:
		return StatusRedirect
	case status >= 400 && status < 500:
		return StatusClient
	case status >= 500:
		return StatusServer
	default:
		return StatusNeutral
	}
}