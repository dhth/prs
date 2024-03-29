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
	case prList:
		content = m.prListStyle.Render(m.prsList.View())
	case prTLList:
		content = m.prTLStyle.Render(m.prTLList.View())
	case repoList:
		content = m.repoListStyle.Render(m.repoList.View())
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#282828")).
		Background(lipgloss.Color("#7c6f64"))

	var helpMsg string
	if m.showHelp {
		helpMsg = " " + helpMsgStyle.Render("tab: switch focus; ctrl+r: change repo; ctrl+b: open in browser; ctrl+d: show diff; ctrl+v: view pr")
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
