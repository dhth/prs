package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newRepoListDelKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
		),
	}
}

func newRepoListItemDel(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(lipgloss.Color("#fe8019")).
		BorderLeftForeground(lipgloss.Color("#fe8019"))
	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle.
		Copy()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		switch msgType := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msgType,
				keys.choose):
				selected := m.SelectedItem()
				if selected != nil {
					return chooseRepo(selected.FilterValue())
				}
			}
		}
		return nil
	}

	return d
}
