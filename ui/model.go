package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

type Pane uint

const (
	repoListView Pane = iota
	prListView
	prDetailsView
	reviewPRListView
	prTLListView
	prTLItemDetailView
	helpView
)

type Mode uint

const (
	QueryMode Mode = iota
	RepoMode
)

type Model struct {
	mode                     Mode
	config                   Config
	ghClient                 *ghapi.GraphQLClient
	repoOwner                string
	repoName                 string
	repoList                 list.Model
	prsList                  list.Model
	prTLList                 list.Model
	prCache                  []*prResult
	prTLItemDetailVP         viewport.Model
	prTLItemDetailVPReady    bool
	prDetailsTitle           string
	prTLItemDetailTitle      string
	prDetailsVP              viewport.Model
	prDetailsVPReady         bool
	prDetailsCache           map[string]prDetails
	prTLCache                map[string][]*prTLItemResult
	message                  string
	helpVP                   viewport.Model
	helpVPReady              bool
	activePane               Pane
	lastPane                 Pane
	secondLastActivePane     Pane
	showHelp                 bool
	repoChosen               bool
	terminalDetails          terminalDetails
	mdRenderer               *glamour.TermRenderer
	prDetailsCurrentSection  uint
	prDetailsCurSectionCache map[string]uint
	prRevCurCmtNum           uint
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, hideHelp(time.Minute*1))

	if m.mode == QueryMode {
		cmds = append(cmds, fetchPRSFromQuery(m.ghClient, *m.config.Query, m.config.PRCount))
	}

	return tea.Batch(cmds...)
}
