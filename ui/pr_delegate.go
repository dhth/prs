package ui

import (
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
)

func newPRListItemDel() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(lipgloss.Color(prListColor)).
		BorderLeftForeground(lipgloss.Color(prListColor))
	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	return d
}
