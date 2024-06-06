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
		statusBar = RightPadTrim(m.message, m.terminalDetails.width)
	}

	switch m.activePane {
	case prList:
		content = m.prListStyle.Render(m.prsList.View())
	case prTLList:
		content = m.prTLStyle.Render(m.prTLList.View())
	case repoList:
		content = m.repoListStyle.Render(m.repoList.View())
	case prRevCmts:
		var prRevCmtsVP string
		if !m.prRevCmtVPReady {
			prRevCmtsVP = "\n  Initializing..."
		} else {
			prRevCmtsVP = viewPortStyle.Render(fmt.Sprintf("  %s\n\n%s\n",
				helpVPTitleStyle.Render("Review Comments"),
				m.prRevCmtVP.View()))
		}
		content = prRevCmtsVP
	case helpView:
		var helpVP string
		if !m.helpVPReady {
			helpVP = "\n  Initializing..."
		} else {
			helpVP = viewPortStyle.Render(fmt.Sprintf("  %s\n\n%s\n",
				helpVPTitleStyle.Render("Help"),
				m.helpVP.View()))
		}
		content = helpVP
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Background(lipgloss.Color("#7c6f64"))

	var helpMsg string
	if m.showHelp {
		helpMsg = " " + helpMsgStyle.Render("Press ? for help")
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
