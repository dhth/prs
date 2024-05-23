package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	ReviewMode
)

type model struct {
	mode            Mode
	config          Config
	ghClient        *ghapi.GraphQLClient
	repoOwner       string
	repoName        string
	activePRNumber  int
	prCount         int
	repoList        list.Model
	prsList         list.Model
	prTLList        list.Model
	prRevCmtVP      viewport.Model
	prRevCmtVPReady bool
	prTLCache       map[string][]prTLItem
	message         string
	repoListStyle   lipgloss.Style
	prListStyle     lipgloss.Style
	prTLStyle       lipgloss.Style
	helpVP          viewport.Model
	helpVPReady     bool
	terminalHeight  int
	terminalWidth   int
	activePane      Pane
	lastPane        Pane
	showHelp        bool
	repoChosen      bool
	userLogin       string
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, hideHelp(time.Minute*1))

	if m.mode == ReviewMode {
		cmds = append(cmds, fetchViewerLogin(m.ghClient))
	}
	return tea.Batch(cmds...)
}
