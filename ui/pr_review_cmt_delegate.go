package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func newPRRevCmtListItemDel() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(lipgloss.Color(RevCmtListColor)).
		BorderLeftForeground(lipgloss.Color(RevCmtListColor))
	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle.
		Copy()

	return d
}
