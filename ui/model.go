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
	prRevCmtsView
	helpView
)

type Mode uint

const (
	QueryMode Mode = iota
	RepoMode
	ReviewerMode
	AuthorMode
)

type model struct {
	mode                    Mode
	config                  Config
	ghClient                *ghapi.GraphQLClient
	repoOwner               string
	repoName                string
	repoList                list.Model
	prsList                 list.Model
	prTLList                list.Model
	prCache                 []*prResult
	prRevCmtVP              viewport.Model
	prRevCmtVPReady         bool
	prDetailsTitle          string
	prDetailsVP             viewport.Model
	prDetailsVPReady        bool
	prDetailsCache          map[string]prDetails
	prTLCache               map[string][]*prTLItemResult
	message                 string
	helpVP                  viewport.Model
	helpVPReady             bool
	activePane              Pane
	lastPane                Pane
	secondLastActivePane    Pane
	showHelp                bool
	repoChosen              bool
	userLogin               string
	terminalDetails         terminalDetails
	mdRenderer              *glamour.TermRenderer
	prDetailsCurrentSection uint
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, hideHelp(time.Minute*1))

	if m.mode == QueryMode {
		cmds = append(cmds, fetchPRSFromQuery(m.ghClient, *m.config.Query, m.config.PRCount))
	}

	if m.mode == ReviewerMode || m.mode == AuthorMode {
		cmds = append(cmds, fetchViewerLogin(m.ghClient))
	}
	return tea.Batch(cmds...)
}
