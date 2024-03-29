package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.message = ""

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.activePane == repoList {
				m.activePane = prList
			} else {
				return m, tea.Quit
			}
		case "tab", "shift+tab":
			if m.activePane == prList {
				m.activePane = prTLList
			} else {
				m.activePane = prList
			}
		case "ctrl+r":
			if m.activePane != repoList {
				m.activePane = repoList
			} else {
				m.activePane = prList
			}
		case "ctrl+b":
			if m.activePane == prList || m.activePane == prTLList {
				selected, ok := m.prsList.SelectedItem().(pr)
				if ok {
					cmds = append(cmds, openPRInBrowser(selected.Url))
				}
			}
		}
	case HideHelpMsg:
		m.showHelp = false
	case tea.WindowSizeMsg:
		_, h1 := m.prListStyle.GetFrameSize()
		_, h2 := m.prTLStyle.GetFrameSize()
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height

		m.repoList.SetHeight(msg.Height - h1 - 2)
		m.repoList.SetWidth(int(float64(msg.Width)*0.4) - 6)
		m.repoListStyle = m.repoListStyle.Width(int(float64(msg.Width)*0.4) - 2)

		m.prsList.SetHeight(msg.Height - h1 - 2)
		m.prsList.SetWidth(int(float64(msg.Width)*0.4) - 6)
		m.prListStyle = m.prListStyle.Width(int(float64(msg.Width)*0.4) - 2)

		m.prTLList.SetHeight(msg.Height - h2 - 2)
		m.prTLList.SetWidth(int(float64(msg.Width)*0.6) - 6)
		m.prTLStyle = m.prTLStyle.Width(int(float64(msg.Width)*0.6) - 2)
	case RepoChosenMsg:
		repoDetails := strings.Split(msg.repo, ":::")
		if len(repoDetails) != 2 {
			m.message = "Something went horribly wrong. Let @dhth know about this failure."
		} else {
			m.repoOwner = repoDetails[0]
			m.repoName = repoDetails[1]
			m.activePane = prList
			m.prsList.ResetSelected()
			m.prTLList.ResetSelected()
			cmds = append(cmds, fetchPRS(m.ghClient, m.repoOwner, m.repoName, m.prCount))
		}
	case PRChosenMsg:
		if msg.err != nil {
			m.message = "Something went wrong: " + msg.err.Error()
		} else {
			tlFromCache, ok := m.prTLCache[fmt.Sprintf("%s/%s:%d", m.repoOwner, m.repoName, msg.prNumber)]
			if !ok {
				cmds = append(cmds, fetchPRTLItems(m.ghClient, m.repoOwner, m.repoName, msg.prNumber, 100, 10))
			} else {
				tlItems := make([]list.Item, 0, len(tlFromCache))
				for _, issue := range tlFromCache {
					tlItems = append(tlItems, issue)
				}
				m.prTLList.SetItems(tlItems)
			}
		}
	case PRsFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
		} else {
			prs := make([]list.Item, 0, len(msg.prs))
			for _, issue := range msg.prs {
				prs = append(prs, issue)
			}
			m.prsList.SetItems(prs)
			m.prsList.Title = fmt.Sprintf("PRs (%s)", m.repoName)

			if len(msg.prs) > 0 {
				for i := 0; i < len(msg.prs); i++ {
					cmds = append(cmds, fetchPRTLItems(m.ghClient, m.repoOwner, m.repoName, msg.prs[i].Number, 100, 10))
				}
			}
		}
	case PRTLFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
		} else {
			tlItems := make([]list.Item, 0, len(msg.prTLItems))
			for _, issue := range msg.prTLItems {
				tlItems = append(tlItems, issue)
			}
			m.prTLList.SetItems(tlItems)
			m.prTLCache[fmt.Sprintf("%s/%s:%d", m.repoOwner, m.repoName, msg.prNumber)] = msg.prTLItems
		}
	case PROpenedinBrowserMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error opening url: %s", msg.err.Error())
		}
	}

	switch m.activePane {
	case prList:
		m.prsList, cmd = m.prsList.Update(msg)
		cmds = append(cmds, cmd)
	case prTLList:
		m.prTLList, cmd = m.prTLList.Update(msg)
		cmds = append(cmds, cmd)
	case repoList:
		m.repoList, cmd = m.repoList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
