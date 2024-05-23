package ui

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

func InitialModel(config Config, mode Mode) model {

	repoListItems := make([]list.Item, 0, len(config.Repos))
	if mode == RepoMode {
		for _, issue := range config.Repos {
			repoListItems = append(repoListItems, issue)
		}
	}

	repoListDel := newRepoListItemDel()
	prListDel := newPRListItemDel()
	prTLListDel := newPRTLListItemDel()

	baseStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color(defaultBackgroundColor))

	repoListStyle := baseStyle.Copy().
		PaddingTop(1).
		PaddingRight(2).
		PaddingLeft(1).
		PaddingBottom(1)

	prListStyle := repoListStyle.Copy()

	prTLStyle := repoListStyle.Copy().
		PaddingRight(1)

	opts := ghapi.ClientOptions{
		EnableCache: true,
		CacheTTL:    time.Second * 30,
		Timeout:     5 * time.Second,
	}
	client, err := ghapi.NewGraphQLClient(opts)
	if err != nil {
		log.Fatalf("err getting client: %s", err.Error())
	}

	prTLCache := make(map[string][]prTLItem)

	m := model{
		mode:          mode,
		config:        config,
		ghClient:      client,
		prCount:       config.PRCount,
		prsList:       list.New(nil, prListDel, 0, 0),
		prTLList:      list.New(nil, prTLListDel, 0, 0),
		prTLCache:     prTLCache,
		repoListStyle: repoListStyle,
		prListStyle:   prListStyle,
		prTLStyle:     prTLStyle,
		showHelp:      true,
	}

	switch m.mode {
	case RepoMode:
		m.repoList = list.New(repoListItems, repoListDel, 0, 0)
		m.repoList.Title = "Repos"
		m.repoList.SetStatusBarItemName("repo", "repos")
		m.repoList.DisableQuitKeybindings()
		m.repoList.SetShowHelp(false)
		m.repoList.SetFilteringEnabled(false)
		m.repoList.Styles.Title.Background(lipgloss.Color(repoListColor))
		m.repoList.Styles.Title.Foreground(lipgloss.Color(defaultBackgroundColor))
		m.repoList.Styles.Title.Bold(true)
	case ReviewerMode, AuthorMode:
		m.activePane = prList
	}

	m.prsList.Title = "fetching..."
	m.prsList.SetStatusBarItemName("PR", "PRs")
	m.prsList.DisableQuitKeybindings()
	m.prsList.SetShowHelp(false)
	m.prsList.SetFilteringEnabled(false)
	m.prsList.Styles.Title.Background(lipgloss.Color(prListColor))
	m.prsList.Styles.Title.Foreground(lipgloss.Color(defaultBackgroundColor))
	m.prsList.Styles.Title.Bold(true)

	m.prTLList.Title = "PR Timeline"
	m.prTLList.SetStatusBarItemName("item", "items")
	m.prTLList.DisableQuitKeybindings()
	m.prTLList.SetShowHelp(false)
	m.prTLList.SetFilteringEnabled(false)
	m.prTLList.Styles.Title.Background(lipgloss.Color(prTLListColor))
	m.prTLList.Styles.Title.Foreground(lipgloss.Color(defaultBackgroundColor))
	m.prTLList.Styles.Title.Bold(true)

	return m
}
