package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

func chooseRepo(repo string) tea.Cmd {
	return func() tea.Msg {
		return repoChosenMsg{repo}
	}
}

func openURLInBrowser(url string) tea.Cmd {
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
			return urlOpenedinBrowserMsg{url: url, err: err}
		}
		return tea.Msg(urlOpenedinBrowserMsg{url: url})
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
			return prDiffDoneMsg{err: err}
		}
		return tea.Msg(prDiffDoneMsg{})
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
			return prDiffDoneMsg{err: err}
		}
		return tea.Msg(prDiffDoneMsg{})
	})
}

func hideHelp(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return hideHelpMsg{}
	})
}

func fetchPRS(ghClient *ghapi.GraphQLClient, repoOwner string, repoName string, prCount int) tea.Cmd {
	return func() tea.Msg {
		prs, err := getPRs(ghClient, repoOwner, repoName, prCount)
		return prsFetchedMsg{prs, err}
	}
}

func fetchViewerLogin(ghClient *ghapi.GraphQLClient) tea.Cmd {
	return func() tea.Msg {
		login, err := getViewerLogin(ghClient)
		return viewerLoginFetched{login, err}
	}
}

func fetchPRsToReview(ghClient *ghapi.GraphQLClient, authorLogin string) tea.Cmd {
	return func() tea.Msg {
		prs, err := getPRsToReview(ghClient, authorLogin)
		return reviewPRsFetchedMsg{prs, err}
	}
}

func fetchAuthoredPRs(ghClient *ghapi.GraphQLClient, authorLogin string) tea.Cmd {
	return func() tea.Msg {
		prs, err := getAuthoredPRs(ghClient, authorLogin)
		return authoredPRsFetchedMsg{prs, err}
	}
}

func fetchPRTLItems(ghClient *ghapi.GraphQLClient, repoOwner string, repoName string, prNumber int, tlItemsCount int, setItems bool) tea.Cmd {
	return func() tea.Msg {
		prTLItems, err := getPRTL(ghClient, repoOwner, repoName, prNumber, tlItemsCount)
		return prTLFetchedMsg{repoOwner, repoName, prNumber, prTLItems, setItems, err}
	}
}
