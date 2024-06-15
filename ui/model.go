package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

type Pane uint

const (
	repoList Pane = iota
	prList
	reviewPRList
	prTLList
	prRevCmts
	helpView
)

type Mode uint

const (
	RepoMode Mode = iota
	ReviewerMode
	AuthorMode
)

type model struct {
	mode            Mode
	config          Config
	ghClient        *ghapi.GraphQLClient
	repoOwner       string
	repoName        string
	prCount         int
	repoList        list.Model
	prsList         list.Model
	prTLList        list.Model
	prCache         []*prResult
	prRevCmtVP      viewport.Model
	prRevCmtVPReady bool
	prTLCache       map[string][]*prTLItemResult
	message         string
	helpVP          viewport.Model
	helpVPReady     bool
	activePane      Pane
	lastPane        Pane
	showHelp        bool
	repoChosen      bool
	userLogin       string
	terminalDetails terminalDetails
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, hideHelp(time.Minute*1))

	if m.mode == ReviewerMode || m.mode == AuthorMode {
		cmds = append(cmds, fetchViewerLogin(m.ghClient))
	}
	return tea.Batch(cmds...)
}
