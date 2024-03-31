package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

const (
	DefaultBackgroundColor      = "#282828"
	RepoListColor               = "#b8bb26"
	PRListColor                 = "#fe8019"
	PRTLListColor               = "#d3869b"
	RevCmtListColor             = "#8ec07c"
	PrOpenColor                 = "#fabd2f"
	PrMergedColor               = "#b8bb26"
	PrClosedColor               = "#928374"
	AdditionsColor              = "#8ec07c"
	DeletionsColor              = "#fb4934"
	ReviewCommentedColor        = "#83a598"
	ReviewApprovedColor         = "#b8bb26"
	ReviewChangesRequestedColor = "#fabd2f"
	ReviewDismissedColor        = "#928374"
	DateColor                   = "#928374"
	NumReviewsColor             = "#665c54"
	NumCommentsColor            = "#83a598"
	DiffColor                   = "#83a598"
	RevCmtColor                 = "#d3869b"
)

var (
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color(DefaultBackgroundColor))

	modeStyle = baseStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#b8bb26"))

	helpVPTitleStyle = baseStyle.Copy().
				Bold(true).
				Background(lipgloss.Color("#8ec07c")).
				Align(lipgloss.Left)

	viewPortStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingRight(2).
			PaddingLeft(1).
			PaddingBottom(1)

	helpMsgStyle = baseStyle.Copy().
			Bold(true).
			Foreground(lipgloss.Color("#83a598"))

	dateStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color(DateColor))

	reviewCmtBodyStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color(RevCmtColor))

	filePathStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color(DateColor))

	numReviewsStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color(NumReviewsColor))

	numCommentsStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(lipgloss.Color(NumCommentsColor))

	diffStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color(DiffColor))

	linesChangedStyle = lipgloss.NewStyle().
				PaddingLeft(1)

	additionsStyle = linesChangedStyle.Copy().
			PaddingLeft(2).
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
			PaddingRight(1).
			Foreground(lipgloss.Color(color))

		return st
	}

	prStyle = func(state string) lipgloss.Style {
		st := lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Width(8).
			Bold(true).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color(DefaultBackgroundColor))

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
			PaddingRight(1).
			Bold(true).
			Align(lipgloss.Center)

		switch state {
		case ReviewCommented:
			st.Foreground(lipgloss.Color(ReviewCommentedColor))
		case ReviewApproved:
			st.Foreground(lipgloss.Color(ReviewApprovedColor))
		case ReviewChangesRequested:
			st.Foreground(lipgloss.Color(ReviewChangesRequestedColor))
		default:
			st.Foreground(lipgloss.Color(ReviewDismissedColor))
		}
		return st
	}
)
