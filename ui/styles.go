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
	filePathColor               = "#d3869b"
	outdatedColor               = "#fabd2f"
	numReviewsColor             = "#665c54"
	numCommentsColor            = "#83a598"
	diffColor                   = "#83a598"
	revCmtColor                 = "#d3869b"
	revCmtDividerColor          = "#665c54"
	footerColor                 = "#7c6f64"
	helpMsgColor                = "#83a598"
	helpViewTitleColor          = "#83a598"
	helpHeaderColor             = "#83a598"
	helpSectionColor            = "#fabd2f"
	toolNameColor               = "#b8bb26"
)

var (
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color(defaultBackgroundColor))

	toolNameStyle = baseStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(toolNameColor))

	listStyle = baseStyle.
			PaddingTop(1).
			PaddingRight(2).
			PaddingLeft(1).
			PaddingBottom(1).
			Foreground(lipgloss.Color(defaultBackgroundColor))

	viewPortStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingRight(2).
			PaddingLeft(1).
			PaddingBottom(1)

	helpMsgStyle = baseStyle.
			Bold(true).
			Foreground(lipgloss.Color(helpMsgColor))

	dateStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color(dateColor))

	repoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(repoColor))

	reviewCmtBodyStyle = lipgloss.NewStyle()

	filePathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(filePathColor))

	outdatedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(outdatedColor))

	reviewCmtDividerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(revCmtDividerColor))

	numReviewsStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color(numReviewsColor))

	numCommentsStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(lipgloss.Color(numCommentsColor))

	diffStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(diffColor))

	linesChangedStyle = lipgloss.NewStyle().
				PaddingLeft(1)

	additionsStyle = linesChangedStyle.
			PaddingLeft(2).
			Foreground(lipgloss.Color(additionsColor))

	deletionsStyle = linesChangedStyle.
			Foreground(lipgloss.Color(deletionsColor))

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
			Foreground(lipgloss.Color(defaultBackgroundColor))

		switch state {
		case prStateOpen:
			st.Background(lipgloss.Color(prOpenColor))
		case prStateMerged:
			st.Background(lipgloss.Color(prMergedColor))
		default:
			st.Background(lipgloss.Color(prClosedColor))
		}
		return st
	}

	reviewStyle = func(state string) lipgloss.Style {
		st := lipgloss.NewStyle().
			PaddingRight(1).
			Bold(true).
			Align(lipgloss.Center)

		switch state {
		case reviewCommented:
			st.Foreground(lipgloss.Color(reviewCommentedColor))
		case reviewApproved:
			st.Foreground(lipgloss.Color(reviewApprovedColor))
		case reviewChangesRequested:
			st.Foreground(lipgloss.Color(reviewChangesRequestedColor))
		default:
			st.Foreground(lipgloss.Color(reviewDismissedColor))
		}
		return st
	}

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(defaultBackgroundColor)).
			Background(lipgloss.Color(footerColor))

	helpVPTitleStyle = baseStyle.
				Bold(true).
				Background(lipgloss.Color(helpViewTitleColor)).
				Align(lipgloss.Left)

	helpHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(helpHeaderColor))

	helpSectionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(helpSectionColor))
)
