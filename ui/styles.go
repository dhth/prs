package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

const (
	PrOpenColor                 = "#fabd2f"
	PrMergedColor               = "#83c07c"
	PrClosedColor               = "#928374"
	AdditionsColor              = "#8ec07c"
	DeletionsColor              = "#fb4934"
	ReviewCommentedColor        = "#83a598"
	ReviewApprovedColor         = "#b8bb26"
	ReviewChangesRequestedColor = "#fabd2f"
	ReviewDismissedColor        = "#928374"
)

var (
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color("#282828"))

	modeStyle = baseStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#b8bb26"))

	helpMsgStyle = baseStyle.Copy().
			Bold(true).
			Foreground(lipgloss.Color("#83a598"))

	linesChangedStyle = lipgloss.NewStyle().
				PaddingLeft(1)
	additionsStyle = linesChangedStyle.Copy().
			Foreground(lipgloss.Color(AdditionsColor))

	deletionsStyle = linesChangedStyle.Copy().
			Foreground(lipgloss.Color(DeletionsColor))

	authorColors = []string{
		"#ccccff", // Lavender Blue
		"#ffa87d", // Light orange
		"#7385D8", // Light blue
		"#fabd2f", // Bright Yellow
		"#00abe5", // Deep Sky
		"#d3691e", // Chocolate
	}
	authorStyle = func(author string) lipgloss.Style {
		h := fnv.New32()
		h.Write([]byte(author))
		hash := h.Sum32()

		color := authorColors[int(hash)%len(authorColors)]

		st := lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color(color))

		return st
	}

	prStyle = func(state string) lipgloss.Style {
		st := lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Width(8).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#282828"))

		switch state {
		case PRStateOpen:
			st.Background(lipgloss.Color(PrOpenColor))
		case PRStateMerged:
			st.Background(lipgloss.Color(PrMergedColor))
		default:
			st.Background(lipgloss.Color(PrClosedColor))
		}
		return st
	}

	reviewStyle = func(state string) lipgloss.Style {
		st := lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Width(12).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#282828"))

		switch state {
		case ReviewCommented:
			st.Background(lipgloss.Color(ReviewCommentedColor))
		case ReviewApproved:
			st.Background(lipgloss.Color(ReviewApprovedColor))
		case ReviewChangesRequested:
			st.Background(lipgloss.Color(ReviewChangesRequestedColor))
		default:
			st.Background(lipgloss.Color(ReviewDismissedColor))
		}
		return st
	}
)
