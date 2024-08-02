package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

func InitialModel(ghClient *ghapi.GraphQLClient, config Config, mode Mode) model {

	prListDel := newPRListItemDel()
	prTLListDel := newPRTLListItemDel()

	prDetailsCache := make(map[string]prDetails)
	prTLCache := make(map[string][]*prTLItemResult)

	prDetailsCurSectionCache := make(map[string]uint)

	m := model{
		mode:                     mode,
		config:                   config,
		ghClient:                 ghClient,
		prsList:                  list.New(nil, prListDel, 0, 0),
		prTLList:                 list.New(nil, prTLListDel, 0, 0),
		prDetailsCache:           prDetailsCache,
		prTLCache:                prTLCache,
		showHelp:                 true,
		terminalDetails:          terminalDetails{width: widthBudgetDefault},
		prDetailsCurSectionCache: prDetailsCurSectionCache,
	}

	switch m.mode {
	case RepoMode:
		repoListItems := make([]list.Item, len(config.Repos))
		if mode == RepoMode {
			for i, issue := range config.Repos {
				repoListItems[i] = issue
			}
		}
		repoListDel := newRepoListItemDel()
		m.repoList = list.New(repoListItems, repoListDel, 0, 0)
		m.repoList.Title = "Repos"
		m.repoList.SetStatusBarItemName("repo", "repos")
		m.repoList.DisableQuitKeybindings()
		m.repoList.SetShowHelp(false)
		m.repoList.SetFilteringEnabled(false)
		m.repoList.Styles.Title = m.repoList.Styles.Title.Background(lipgloss.Color(repoListColor)).
			Foreground(lipgloss.Color(defaultBackgroundColor)).
			Bold(true)
		m.repoList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
		m.repoList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")
	case QueryMode, ReviewerMode, AuthorMode:
		m.activePane = prListView
	}

	m.prsList.Title = "fetching PRs..."
	m.prsList.SetStatusBarItemName("PR", "PRs")
	m.prsList.DisableQuitKeybindings()
	m.prsList.SetShowHelp(false)
	m.prsList.SetFilteringEnabled(false)
	m.prsList.Styles.Title = m.prsList.Styles.Title.Background(lipgloss.Color(fetchingColor)).
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Bold(true)
	m.prsList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	m.prsList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	m.prTLList.Title = "fetching timeline..."
	m.prTLList.SetStatusBarItemName("item", "items")
	m.prTLList.DisableQuitKeybindings()
	m.prTLList.SetShowHelp(false)
	m.prTLList.SetFilteringEnabled(false)
	m.prTLList.Styles.Title = m.prTLList.Styles.Title.Background(lipgloss.Color(prTLListColor)).
		Foreground(lipgloss.Color(defaultBackgroundColor)).
		Bold(true)
	m.prTLList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	m.prTLList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	return m
}
