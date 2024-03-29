package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

type Pane uint

const (
	prList Pane = iota
	prTLList
	repoList
)

type model struct {
	config         Config
	ghClient       *ghapi.GraphQLClient
	repoOwner      string
	repoName       string
	prCount        int
	repoList       list.Model
	prsList        list.Model
	prTLList       list.Model
	prTLCache      map[string][]prTLItem
	message        string
	repoListStyle  lipgloss.Style
	prListStyle    lipgloss.Style
	prTLStyle      lipgloss.Style
	terminalHeight int
	terminalWidth  int
	activePane     Pane
	showHelp       bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		hideHelp(time.Minute*1),
		fetchPRS(m.ghClient, m.repoOwner, m.repoName, m.prCount),
	)
}
