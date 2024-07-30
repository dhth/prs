package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	viewPortWrapUpperLimit = 160
)

func (m model) View() string {
	var content string
	var footer string

	var statusBar string
	if m.message != "" {
		statusBar = RightPadTrim(m.message, m.terminalDetails.width)
	}

	switch m.activePane {
	case prListView:
		content = listStyle.Render(m.prsList.View())
	case prTLListView:
		content = listStyle.Render(m.prTLList.View())
	case repoListView:
		content = listStyle.Render(m.repoList.View())
	case prDetailsView:
		if !m.prRevCmtVPReady {
			content = "\n  Initializing..."
		} else {
			content = viewPortStyle.Render(fmt.Sprintf("  %s\n\n%s\n",
				prDetailsTitleStyle.Render(m.prDetailsTitle),
				m.prDetailsVP.View()))
		}
	case prRevCmtsView:
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

	var helpMsg string
	if m.showHelp {
		helpMsg = helpMsgStyle.Render("Press ? for help")
	}

	footerStr := fmt.Sprintf("%s%s",
		toolNameStyle.Render("prs"),
		helpMsg,
	)
	footer = footerStyle.Render(footerStr)

	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		statusBar,
		footer,
	)
}
