package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	var content string
	var footer string

	var statusBar string
	if m.message != "" {
		statusBar = RightPadTrim(m.message, m.terminalWidth)
	}

	switch m.activePane {
	case prList, prTLList:
		content = lipgloss.JoinHorizontal(lipgloss.Top,
			m.prListStyle.Render(m.prsList.View()),
			m.prTLStyle.Render(m.prTLList.View()),
		)
	case repoList:
		content = m.repoListStyle.Render(m.repoList.View())
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#282828")).
		Background(lipgloss.Color("#7c6f64"))

	var helpMsg string
	if m.showHelp {
		helpMsg = " " + helpMsgStyle.Render("tab: switch focus; ctrl+r: change repo; q: quit")
	}

	footerStr := fmt.Sprintf("%s%s",
		modeStyle.Render("prs"),
		helpMsg,
	)
	footer = footerStyle.Render(footerStr)

	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		statusBar,
		footer,
	)
}
