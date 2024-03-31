package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func newRepoListItemDel() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(lipgloss.Color(RepoListColor)).
		BorderLeftForeground(lipgloss.Color(RepoListColor))
	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle.
		Copy()

	return d
}
