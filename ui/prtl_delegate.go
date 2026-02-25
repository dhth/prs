package ui

import (
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
)

func newPRTLListItemDel() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(lipgloss.Color(prTLListColor)).
		BorderLeftForeground(lipgloss.Color(prTLListColor))
	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	return d
}
