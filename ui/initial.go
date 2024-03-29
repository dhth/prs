package ui

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

func InitialModel(config Config) model {

	repoListItems := make([]list.Item, 0, len(config.Repos))
	for _, issue := range config.Repos {
		repoListItems = append(repoListItems, issue)
	}

	var repoListDelKeys = newRepoListDelKeyMap()
	repoListDel := newRepoListItemDel(repoListDelKeys)

	var prListDelKeys = newPRListDelKeyMap()
	prListDel := newPRListItemDel(prListDelKeys)

	var prTLListDelKeys = newPRTLListDelKeyMap()
	prTLListDel := newPRTLListItemDel(prTLListDelKeys)

	baseStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color("#282828"))

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
		Timeout:     5 * time.Second,
	}
	client, err := ghapi.NewGraphQLClient(opts)
	if err != nil {
		log.Fatalf("err getting client: %s", err.Error())
	}

	prTLCache := make(map[string][]prTLItem)

	m := model{
		config:        config,
		ghClient:      client,
		repoOwner:     config.Repos[0].Owner,
		repoName:      config.Repos[0].Name,
		prCount:       config.PRCount,
		repoList:      list.New(repoListItems, repoListDel, 0, 0),
		prsList:       list.New(nil, prListDel, 0, 0),
		prTLList:      list.New(nil, prTLListDel, 0, 0),
		prTLCache:     prTLCache,
		repoListStyle: repoListStyle,
		prListStyle:   prListStyle,
		prTLStyle:     prTLStyle,
		showHelp:      true,
	}

	m.repoList.Title = "Repos"
	m.repoList.SetStatusBarItemName("repo", "repos")
	m.repoList.DisableQuitKeybindings()
	m.repoList.SetShowHelp(false)
	m.repoList.Styles.Title.Background(lipgloss.Color("#fe8019"))
	m.repoList.Styles.Title.Foreground(lipgloss.Color("#282828"))
	m.repoList.Styles.Title.Bold(true)

	m.prsList.Title = "PRs (fetching...)"
	m.prsList.SetStatusBarItemName("PR", "PRS")
	m.prsList.DisableQuitKeybindings()
	m.prsList.SetShowHelp(false)
	m.prsList.Styles.Title.Background(lipgloss.Color("#fe8019"))
	m.prsList.Styles.Title.Foreground(lipgloss.Color("#282828"))
	m.prsList.Styles.Title.Bold(true)

	m.prTLList.Title = "PR Timeline"
	m.prTLList.SetStatusBarItemName("item", "items")
	m.prTLList.DisableQuitKeybindings()
	m.prTLList.SetShowHelp(false)
	m.prTLList.SetFilteringEnabled(false)
	m.prTLList.Styles.Title.Background(lipgloss.Color("#d3869b"))
	m.prTLList.Styles.Title.Foreground(lipgloss.Color("#282828"))
	m.prTLList.Styles.Title.Bold(true)

	return m
}
