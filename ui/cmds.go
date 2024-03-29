package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

func chooseRepo(repo string) tea.Cmd {
	return func() tea.Msg {
		return RepoChosenMsg{repo}
	}
}

func choosePR(prNumberStr string) tea.Cmd {
	return func() tea.Msg {
		prNumber, err := strconv.Atoi(prNumberStr)
		return PRChosenMsg{prNumber, err}
	}
}

func openPRInBrowser(url string) tea.Cmd {
	var openCmd string
	switch runtime.GOOS {
	case "darwin":
		openCmd = "open"
	default:
		openCmd = "xdg-open"
	}
	c := exec.Command(openCmd, url)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return PROpenedinBrowserMsg{url: url, err: err}
		}
		return tea.Msg(PROpenedinBrowserMsg{url: url})
	})
}

func showDiff(repoOwner, repoName string, prNumber int, pager *string) tea.Cmd {
	var pagerPrefix string
	if pager != nil {
		pagerPrefix = fmt.Sprintf("GH_PAGER='%s' ", *pager)

	}
	c := exec.Command("bash", "-c",
		fmt.Sprintf("%sgh --repo %s/%s pr diff %d",
			pagerPrefix,
			repoOwner,
			repoName,
			prNumber,
		))
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return PRDiffDoneMsg{err: err}
		}
		return tea.Msg(PRDiffDoneMsg{})
	})
}

func showPR(repoOwner, repoName string, prNumber int) tea.Cmd {
	c := exec.Command("bash", "-c",
		fmt.Sprintf("gh --repo %s/%s pr view --comments %d",
			repoOwner,
			repoName,
			prNumber,
		))
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return PRDiffDoneMsg{err: err}
		}
		return tea.Msg(PRDiffDoneMsg{})
	})
}

func hideHelp(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return HideHelpMsg{}
	})
}

func fetchPRS(ghClient *ghapi.GraphQLClient, repoOwner string, repoName string, prCount int) tea.Cmd {
	return func() tea.Msg {
		prs, err := GetPRs(ghClient, repoOwner, repoName, prCount)
		return PRsFetchedMsg{prs, err}
	}
}

func fetchPRTLItems(ghClient *ghapi.GraphQLClient, repoOwner string, repoName string, prNumber int, tlItemsCount int, commentsCount int) tea.Cmd {
	return func() tea.Msg {
		prTLItems, err := GetPRTL(ghClient, repoOwner, repoName, prNumber, tlItemsCount, commentsCount)
		return PRTLFetchedMsg{prNumber, prTLItems, err}
	}
}
