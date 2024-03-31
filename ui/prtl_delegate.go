package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func newPRTLListItemDel() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(lipgloss.Color(PRTLListColor)).
		BorderLeftForeground(lipgloss.Color(PRTLListColor))
	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle.
		Copy()

	return d
}
