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
	var args []string

	switch runtime.GOOS {
	case "windows":
		openCmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		openCmd = "open"
		args = []string{url}
	default:
		openCmd = "xdg-open"
		args = []string{url}
	}
	c := exec.Command(openCmd, args...)
	err := c.Run()
	return func() tea.Msg {
		return urlOpenedinBrowserMsg{url: url, err: err}
	}
}

func showDiff(repoOwner, repoName string, prNumber int) tea.Cmd {
	cmd := []string{
		"gh",
		"--repo",
		fmt.Sprintf("%s/%s", repoOwner, repoName),
		"pr",
		"diff",
		fmt.Sprintf("%d", prNumber),
	}
	c := exec.Command(cmd[0], cmd[1:]...)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return prDiffDoneMsg{err: err}
		}
		return tea.Msg(prDiffDoneMsg{})
	})
}

func showPR(repoOwner, repoName string, prNumber int) tea.Cmd {
	cmd := []string{
		"gh",
		"--repo",
		fmt.Sprintf("%s/%s", repoOwner, repoName),
		"pr",
		"view",
		"--comments",
		fmt.Sprintf("%d", prNumber),
	}
	c := exec.Command(cmd[0], cmd[1:]...)

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

func fetchPRSFromQuery(ghClient *ghapi.GraphQLClient, queryStr string, prCount int) tea.Cmd {
	return func() tea.Msg {
		prs, err := getPRDataFromQuery(ghClient, queryStr, prCount)
		return prsFetchedMsg{prs, err}
	}
}

func fetchPRSForRepo(ghClient *ghapi.GraphQLClient, repoOwner string, repoName string, prCount int) tea.Cmd {
	return func() tea.Msg {
		queryStr := fmt.Sprintf("type:pr repo:%s/%s sort:updated-desc", repoOwner, repoName)
		prs, err := getPRDataFromQuery(ghClient, queryStr, prCount)
		return prsFetchedMsg{prs, err}
	}
}

func fetchPRMetadata(ghClient *ghapi.GraphQLClient, repoOwner, repoName string, prNumber int) tea.Cmd {
	return func() tea.Msg {
		metadata, err := getPRMetadata(ghClient, repoOwner, repoName, prNumber)
		return prMetadataFetchedMsg{repoOwner, repoName, prNumber, metadata, err}
	}
}

func fetchPRTLItems(ghClient *ghapi.GraphQLClient, repoOwner string, repoName string, prNumber int, tlItemsCount int, setItems bool) tea.Cmd {
	return func() tea.Msg {
		prTLItems, err := getPRTLData(ghClient, repoOwner, repoName, prNumber, tlItemsCount)
		return prTLFetchedMsg{repoOwner, repoName, prNumber, prTLItems, setItems, err}
	}
}
