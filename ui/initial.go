package ui

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

func InitialModel(config Config, mode Mode) model {

	repoListItems := make([]list.Item, len(config.Repos))
	if mode == RepoMode {
		for i, issue := range config.Repos {
			repoListItems[i] = issue
		}
	}

	repoListDel := newRepoListItemDel()
	prListDel := newPRListItemDel()
	prTLListDel := newPRTLListItemDel()

	opts := ghapi.ClientOptions{
		EnableCache: true,
		CacheTTL:    time.Second * 30,
		Timeout:     5 * time.Second,
	}
	client, err := ghapi.NewGraphQLClient(opts)
	if err != nil {
		log.Fatalf("err getting client: %s", err.Error())
	}

	prTLCache := make(map[string][]*prTLItemResult)

	m := model{
		mode:            mode,
		config:          config,
		ghClient:        client,
		prCount:         config.PRCount,
		prsList:         list.New(nil, prListDel, 0, 0),
		prTLList:        list.New(nil, prTLListDel, 0, 0),
		prTLCache:       prTLCache,
		showHelp:        true,
		terminalDetails: terminalDetails{width: widthBudgetDefault},
	}

	switch m.mode {
	case RepoMode:
		m.repoList = list.New(repoListItems, repoListDel, 0, 0)
		m.repoList.Title = "Repos"
		m.repoList.SetStatusBarItemName("repo", "repos")
		m.repoList.DisableQuitKeybindings()
		m.repoList.SetShowHelp(false)
		m.repoList.SetFilteringEnabled(false)
		m.repoList.Styles.Title = m.repoList.Styles.Title.Background(lipgloss.Color(repoListColor)).
			Foreground(lipgloss.Color(defaultBackgroundColor)).
			Bold(true)
	case ReviewerMode, AuthorMode:
		m.activePane = prList
	}

	m.prsList.Title = "fetching..."
	m.prsList.SetStatusBarItemName("PR", "PRs")
	m.prsList.DisableQuitKeybindings()
	m.prsList.SetShowHelp(false)
	m.prsList.SetFilteringEnabled(false)
	m.prsList.Styles.Title = m.prsList.Styles.Title.Background(lipgloss.Color(prListColor)).
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Bold(true)

	m.prTLList.Title = "PR Timeline"
	m.prTLList.SetStatusBarItemName("item", "items")
	m.prTLList.DisableQuitKeybindings()
	m.prTLList.SetShowHelp(false)
	m.prTLList.SetFilteringEnabled(false)
	m.prTLList.Styles.Title = m.prTLList.Styles.Title.Background(lipgloss.Color(prTLListColor)).
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Bold(true)

	return m
}
