package ui

import (
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
)

func newRepoListItemDel() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(lipgloss.Color(repoListColor)).
		BorderLeftForeground(lipgloss.Color(repoListColor))
	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	return d
}
