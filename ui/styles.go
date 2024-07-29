package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

const (
	defaultBackgroundColor      = "#282828"
	repoListColor               = "#b8bb26"
	prListColor                 = "#fe8019"
	prTLListColor               = "#d3869b"
	revCmtListColor             = "#8ec07c"
	prOpenColor                 = "#fabd2f"
	prMergedColor               = "#b8bb26"
	prClosedColor               = "#928374"
	additionsColor              = "#8ec07c"
	deletionsColor              = "#fb4934"
	reviewCommentedColor        = "#83a598"
	reviewApprovedColor         = "#b8bb26"
	reviewChangesRequestedColor = "#fabd2f"
	reviewDismissedColor        = "#928374"
	dateColor                   = "#928374"
	repoColor                   = "#bdae93"
	numReviewsColor             = "#665c54"
	numCommentsColor            = "#83a598"
	footerColor                 = "#7c6f64"
	helpMsgColor                = "#83a598"
	helpViewTitleColor          = "#83a598"
	toolNameColor               = "#b8bb26"
	fetchingColor               = "#928374"
)

func getDynamicStyle(author string) lipgloss.Style {
	h := fnv.New32()
	h.Write([]byte(author))
	hash := h.Sum32()

	color := colors[int(hash)%len(colors)]

	st := lipgloss.NewStyle().
		PaddingRight(1).
		Foreground(lipgloss.Color(color))

	return st
}

var (
	baseStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(defaultBackgroundColor))

	toolNameStyle = baseStyle.
			Align(lipgloss.Center).
			PaddingLeft(1).
			PaddingRight(1).
			Bold(true).
			Background(lipgloss.Color(toolNameColor))

	listStyle = baseStyle.
			PaddingTop(1).
			PaddingBottom(1).
			Foreground(lipgloss.Color(defaultBackgroundColor))

	viewPortStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingRight(2).
			PaddingBottom(1)

	helpMsgStyle = baseStyle.
			PaddingLeft(2).
			Bold(true).
			Foreground(lipgloss.Color(helpMsgColor))

	dateStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color(dateColor))

	numReviewsStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color(numReviewsColor))

	numCommentsStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(lipgloss.Color(numCommentsColor))

	linesChangedStyle = lipgloss.NewStyle().
				PaddingLeft(1)

	additionsStyle = linesChangedStyle.
			PaddingLeft(2).
			Foreground(lipgloss.Color(additionsColor))

	deletionsStyle = linesChangedStyle.
			Foreground(lipgloss.Color(deletionsColor))

	prStyle = func(state string) lipgloss.Style {
		st := lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Width(8).
			Bold(true).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color(defaultBackgroundColor))

		var bgColor string
		switch state {
		case prStateOpen:
			bgColor = prOpenColor
		case prStateMerged:
			bgColor = prMergedColor
		default:
			bgColor = prClosedColor
		}
		return st.Background(lipgloss.Color(bgColor))
	}

	reviewStyle = func(state string) lipgloss.Style {
		st := lipgloss.NewStyle().
			PaddingRight(1).
			Bold(true).
			Align(lipgloss.Center)

		var bgColor string
		switch state {
		case reviewCommented:
			bgColor = reviewCommentedColor
		case reviewApproved:
			bgColor = reviewApprovedColor
		case reviewChangesRequested:
			bgColor = reviewChangesRequestedColor
		default:
			bgColor = reviewDismissedColor
		}
		return st.Foreground(lipgloss.Color(bgColor))
	}

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(defaultBackgroundColor)).
			Background(lipgloss.Color(footerColor))

	titleStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Bold(true).
			Foreground(lipgloss.Color(defaultBackgroundColor))

	helpVPTitleStyle = titleStyle.
				Background(lipgloss.Color(helpViewTitleColor))
)
