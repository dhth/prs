package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	viewPortWrapUpperLimit = 160
	vpNotReadyMsg          = "Initializing..."
)

func (m Model) View() tea.View {
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
		if !m.prTLItemDetailVPReady {
			content = vpNotReadyMsg
		} else {
			content = viewPortStyle.Render(fmt.Sprintf("  %s\n\n%s\n",
				prDetailsTitleStyle.Render(m.prDetailsTitle),
				m.prDetailsVP.View()))
		}
	case prTLItemDetailView:
		var prRevCmtsVP string
		if !m.prTLItemDetailVPReady {
			prRevCmtsVP = vpNotReadyMsg
		} else {
			prRevCmtsVP = viewPortStyle.Render(fmt.Sprintf("  %s\n\n%s\n",
				helpVPTitleStyle.Render(m.prTLItemDetailTitle),
				m.prTLItemDetailVP.View()))
		}
		content = prRevCmtsVP
	case helpView:
		var helpVP string
		if !m.helpVPReady {
			helpVP = vpNotReadyMsg
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

	rendered := lipgloss.JoinVertical(lipgloss.Left,
		content,
		statusBar,
		footer,
	)

	v := tea.NewView(rendered)
	v.AltScreen = true

	return v
}
