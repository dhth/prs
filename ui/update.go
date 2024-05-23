package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

const useHighPerformanceRenderer = false

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.message = ""

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.activePane == repoList {
				if !m.repoChosen {
					return m, tea.Quit
				}
				m.repoList.ResetSelected()
				m.activePane = m.lastPane
			} else if m.activePane == helpView {
				m.activePane = m.lastPane
			} else if m.activePane == prRevCmts {
				m.prRevCmtVP.GotoTop()
				m.activePane = prTLList
			} else if m.activePane == prTLList {
				m.prTLList.ResetSelected()
				m.activePane = prList
			} else {
				return m, tea.Quit
			}
		case "ctrl+r":
			switch m.activePane {
			case prList:
				m.prsList.ResetSelected()
				m.prTLList.ResetSelected()
				switch m.mode {
				case RepoMode:
					cmds = append(cmds, fetchPRS(m.ghClient, m.repoOwner, m.repoName, m.prCount))
				case ReviewerMode:
					cmds = append(cmds, fetchPRsToReview(m.ghClient, m.userLogin))
				case AuthorMode:
					cmds = append(cmds, fetchAuthoredPRs(m.ghClient, m.userLogin))
				}
			case prTLList:
				var repoOwner, repoName string
				var prNumber int

				switch m.mode {
				case RepoMode:
					repoOwner = m.repoOwner
					repoName = m.repoName
					prNumber = m.activePRNumber
				case ReviewerMode:
					pr, reviewPrOk := m.prsList.SelectedItem().(prResult)
					if reviewPrOk {
						repoOwner = pr.details.Repository.Owner.Login
						repoName = pr.details.Repository.Name
						prNumber = pr.details.Number
					}
				}
				if repoOwner != "" && repoName != "" {
					cmds = append(cmds, fetchPRTLItems(m.ghClient, repoOwner, repoName, prNumber, 100))
					m.prTLList.ResetSelected()
				}
			}
		case "[":
			if m.activePane == prTLList {
				m.prsList.CursorUp()
				selected := m.prsList.SelectedItem()
				if selected != nil {
					cmds = append(cmds, choosePR(selected.FilterValue()))
				}
			}
		case "]":
			if m.activePane == prTLList {
				m.prsList.CursorDown()
				selected := m.prsList.SelectedItem()
				if selected != nil {
					cmds = append(cmds, choosePR(selected.FilterValue()))
				}
			}
		case "1":
			if m.activePane != prList {
				m.activePane = prList
			}
		case "d":
			if m.activePane != prList {
				selected, ok := m.prsList.SelectedItem().(prResult)
				if ok {
					m.message = fmt.Sprintf("context: %v", selected.context)
				}
			}
		case "enter":
			switch m.activePane {
			case prList:
				setTlCmd, ok := m.setTL()
				if !ok {
					m.message = "Could't get repo/pr details. Inform @dhth on github."
				} else {
					if setTlCmd != nil {
						cmds = append(cmds, setTlCmd)
					}
					m.activePane = prTLList
				}
			case prTLList:
				item, ok := m.prTLList.SelectedItem().(prTLItem)
				if ok {
					if item.Type == tlItemPRReview {
						revCmts := item.PullRequestReview.Comments.Nodes
						if len(revCmts) > 0 {
							var prReviewCmts string
							for _, cmt := range revCmts {
								prReviewCmts += cmt.render()
								prReviewCmts += "\n\n"
								prReviewCmts += reviewCmtDividerStyle.Render(strings.Repeat("-", int(float64(m.terminalWidth)*0.6)))
								prReviewCmts += "\n\n"
							}
							prReviewCmts += "\n\n"
							m.prRevCmtVP.SetContent(prReviewCmts)
							m.activePane = prRevCmts
						}
					}
				}
			case repoList:
				selected := m.repoList.SelectedItem()
				if selected != nil {
					cmds = append(cmds, chooseRepo(selected.FilterValue()))
				}
			}
		case "2":
			if m.activePane != prTLList {
				setTlCmd, ok := m.setTL()
				if !ok {
					m.message = "Could't get repo/pr details. Inform @dhth on github."
				} else {
					if setTlCmd != nil {
						cmds = append(cmds, setTlCmd)
					}
					m.activePane = prTLList
				}
			}
		case "3":
			if m.activePane == prTLList {
				item, ok := m.prTLList.SelectedItem().(prTLItem)
				if ok {
					if item.Type == tlItemPRReview {
						revCmts := item.PullRequestReview.Comments.Nodes
						if len(revCmts) > 0 {
							var prReviewCmts string
							for _, cmt := range revCmts {
								prReviewCmts += cmt.render()
								prReviewCmts += "\n\n"
								prReviewCmts += reviewCmtDividerStyle.Render(strings.Repeat("*", int(float64(m.terminalWidth)*0.6)))
								prReviewCmts += "\n\n"
							}
							prReviewCmts += "\n\n"
							m.prRevCmtVP.SetContent(prReviewCmts)
							m.activePane = prRevCmts
						}
					}
				}
			}
		case "tab", "shift+tab":
			if m.activePane == prList {
				setTlCmd, ok := m.setTL()
				if !ok {
					m.message = "Could't get repo/pr details. Inform @dhth on github."
				} else {
					if setTlCmd != nil {
						cmds = append(cmds, setTlCmd)
					}
					m.activePane = prTLList
				}
			} else {
				m.activePane = prList
			}
		case "ctrl+s":
			if m.mode == RepoMode {
				if m.activePane != repoList {
					m.lastPane = m.activePane
					m.activePane = repoList
				} else {
					m.activePane = m.lastPane
				}
			}
		case "ctrl+b":
			switch m.activePane {
			case prList:
				var url string
				switch m.mode {
				case RepoMode:
					pr, ok := m.prsList.SelectedItem().(prResult)
					if ok {
						url = pr.details.Url
					}
				case ReviewerMode:
					pr, ok := m.prsList.SelectedItem().(prResult)
					if ok {
						url = pr.details.Url
					}
				}
				cmds = append(cmds, openURLInBrowser(url))
			case prTLList, prRevCmts:
				item, ok := m.prTLList.SelectedItem().(prTLItem)
				if ok {
					switch item.Type {
					case tlItemPRCommit:
						cmds = append(cmds, openURLInBrowser(item.PullRequestCommit.Url))
					case tlItemHeadRefForcePushed:
						cmds = append(cmds, openURLInBrowser(item.HeadRefForcePushed.AfterCommit.Url))
					case tlItemPRReview:
						cmds = append(cmds, openURLInBrowser(item.PullRequestReview.Url))
					case tlItemMergedEvent:
						cmds = append(cmds, openURLInBrowser(item.MergedEvent.Url))
					}
				}
			}
		case "ctrl+d":
			if m.activePane == prList || m.activePane == prTLList {
				switch m.mode {
				case RepoMode:
					pr, ok := m.prsList.SelectedItem().(prResult)
					if ok {
						cmds = append(cmds, showDiff(m.repoOwner,
							m.repoName,
							pr.details.Number,
							m.config.DiffPager))
					}
				case ReviewerMode:
					pr, ok := m.prsList.SelectedItem().(prResult)
					if ok {
						cmds = append(cmds, showDiff(pr.details.Repository.Owner.Login,
							pr.details.Repository.Name,
							pr.details.Number,
							m.config.DiffPager))
					}
				}
			}
		case "ctrl+v":
			if m.activePane == prList || m.activePane == prTLList {
				switch m.mode {
				case RepoMode:
					pr, ok := m.prsList.SelectedItem().(prResult)
					if ok {
						cmds = append(cmds, showPR(m.repoOwner, m.repoName, pr.details.Number))
					}
				case ReviewerMode:
					pr, ok := m.prsList.SelectedItem().(prResult)
					if ok {
						cmds = append(cmds, showPR(pr.details.Repository.Owner.Login,
							pr.details.Repository.Name,
							pr.details.Number))
					}
				}
			}
		case "g":
			if m.activePane == prRevCmts {
				m.prRevCmtVP.GotoTop()
			}
		case "G":
			if m.activePane == prRevCmts {
				m.prRevCmtVP.GotoBottom()
			}
		case "?":
			m.lastPane = m.activePane
			m.activePane = helpView
		}
	case hideHelpMsg:
		m.showHelp = false
	case tea.WindowSizeMsg:
		_, h1 := m.prListStyle.GetFrameSize()
		_, h2 := m.prTLStyle.GetFrameSize()
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height

		if m.mode == RepoMode {
			m.repoList.SetHeight(msg.Height - h1 - 2)
			m.repoList.SetWidth(msg.Width)
			m.repoListStyle = m.repoListStyle.Width(msg.Width)
		}

		m.prsList.SetHeight(msg.Height - h1 - 2)
		m.prsList.SetWidth(msg.Width)
		m.prListStyle = m.prListStyle.Width(msg.Width)

		m.prTLList.SetHeight(msg.Height - h2 - 2)
		m.prTLList.SetWidth(msg.Width)
		m.prTLStyle = m.prTLStyle.Width(msg.Width)

		if !m.prRevCmtVPReady {
			m.prRevCmtVP = viewport.New(int(float64(msg.Width)*0.9), msg.Height-7)
			m.prRevCmtVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.prRevCmtVPReady = true
		} else {
			m.prRevCmtVP.Width = int(float64(msg.Width) * 0.9)
			m.prRevCmtVP.Height = msg.Height - 7
		}

		if !m.helpVPReady {
			m.helpVP = viewport.New(msg.Width, msg.Height-7)
			m.helpVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.helpVP.SetContent(helpText)
			m.helpVPReady = true
		} else {
			m.helpVP.Width = msg.Width
			m.helpVP.Height = msg.Height - 7
		}
	case repoChosenMsg:
		repoDetails := strings.Split(msg.repo, ":::")
		if len(repoDetails) != 2 {
			m.message = "Something went horribly wrong. Let @dhth know about this failure."
		} else {
			m.repoChosen = true
			m.repoOwner = repoDetails[0]
			m.repoName = repoDetails[1]
			m.activePane = prList
			m.prsList.ResetSelected()
			m.prTLList.ResetSelected()
			cmds = append(cmds, fetchPRS(m.ghClient, m.repoOwner, m.repoName, m.prCount))
		}
	case viewerLoginFetched:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error fetching gh username: %s", msg.err)
		} else {
			m.userLogin = msg.login
			switch m.mode {
			case ReviewerMode:
				cmds = append(cmds, fetchPRsToReview(m.ghClient, m.userLogin))
			case AuthorMode:
				cmds = append(cmds, fetchAuthoredPRs(m.ghClient, m.userLogin))
			}
		}
	case prChosenMsg:
		if msg.err != nil {
			m.message = "Something went wrong: " + msg.err.Error()
		} else {
			m.activePRNumber = msg.prNumber
		}
	case prsFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
		} else {
			prs := make([]list.Item, 0, len(msg.prs))
			for _, pr := range msg.prs {
				prs = append(prs, prResult{context: prContextRepo, details: pr})
			}
			m.prsList.SetItems(prs)
			m.prsList.Title = fmt.Sprintf("PRs (%s)", m.repoName)

			if len(msg.prs) > 0 {
				for _, pr := range msg.prs {
					cmds = append(cmds, fetchPRTLItems(m.ghClient,
						pr.Repository.Owner.Login,
						pr.Repository.Name,
						pr.Number,
						100,
					))
				}
				firstPRNumber := strconv.Itoa(msg.prs[0].Number)
				cmds = append(cmds, choosePR(firstPRNumber))
			}
		}
	case reviewPRsFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
		} else {
			prs := make([]list.Item, 0, len(msg.prs))
			for _, pr := range msg.prs {
				prs = append(prs, prResult{context: prContextReviewer, details: pr})
			}
			m.prsList.SetItems(prs)
			m.prsList.Title = "PRs requesting your review"

			if len(msg.prs) > 0 {
				for _, pr := range msg.prs {
					cmds = append(cmds, fetchPRTLItems(m.ghClient, pr.Repository.Owner.Login, pr.Repository.Name, pr.Number, 100))
				}
				firstPRNumber := strconv.Itoa(msg.prs[0].Number)
				cmds = append(cmds, choosePR(firstPRNumber))
			}
		}
	case authoredPRsFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
		} else {
			prs := make([]list.Item, 0, len(msg.prs))
			for _, pr := range msg.prs {
				prs = append(prs, prResult{context: prContextAuthor, details: pr})
			}
			m.prsList.SetItems(prs)
			m.prsList.Title = "Open PRs authored by you"

			if len(msg.prs) > 0 {
				for _, pr := range msg.prs {
					cmds = append(cmds, fetchPRTLItems(m.ghClient, pr.Repository.Owner.Login, pr.Repository.Name, pr.Number, 100))
				}
				firstPRNumber := strconv.Itoa(msg.prs[0].Number)
				cmds = append(cmds, choosePR(firstPRNumber))
			}
		}
	case prTLFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
		} else {
			tlItems := make([]list.Item, 0, len(msg.prTLItems))
			for _, issue := range msg.prTLItems {
				tlItems = append(tlItems, issue)
			}
			m.prTLList.SetItems(tlItems)
			m.prTLList.Title = fmt.Sprintf("PR #%d Timeline", msg.prNumber)
			m.prTLCache[fmt.Sprintf("%s/%s:%d", msg.repoOwner, msg.repoName, msg.prNumber)] = msg.prTLItems
		}
	case urlOpenedinBrowserMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error opening url: %s", msg.err.Error())
		}
	case prDiffDoneMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error opening diff: %s", msg.err.Error())
		}
	case prViewDoneMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error showing PR: %s", msg.err.Error())
		}
	}

	switch m.activePane {
	case prList:
		m.prsList, cmd = m.prsList.Update(msg)
		cmds = append(cmds, cmd)
	case prTLList:
		m.prTLList, cmd = m.prTLList.Update(msg)
		cmds = append(cmds, cmd)
	case prRevCmts:
		m.prRevCmtVP, cmd = m.prRevCmtVP.Update(msg)
		cmds = append(cmds, cmd)
	case repoList:
		m.repoList, cmd = m.repoList.Update(msg)
		cmds = append(cmds, cmd)
	case helpView:
		m.helpVP, cmd = m.helpVP.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) setTL() (tea.Cmd, bool) {
	var cmd tea.Cmd
	var repoOwner, repoName string
	var prNumber int

	switch m.mode {
	case RepoMode:
		repoOwner = m.repoOwner
		repoName = m.repoName
		prNumber = m.activePRNumber
	case ReviewerMode, AuthorMode:
		pr, reviewPrOk := m.prsList.SelectedItem().(prResult)
		if reviewPrOk {
			repoOwner = pr.details.Repository.Owner.Login
			repoName = pr.details.Repository.Name
			prNumber = pr.details.Number
		}
	}

	if repoOwner == "" || repoName == "" {
		return nil, false
	}

	tlFromCache, ok := m.prTLCache[fmt.Sprintf("%s/%s:%d", repoOwner, repoName, prNumber)]
	if !ok {
		cmd = fetchPRTLItems(m.ghClient, repoOwner, repoName, prNumber, 100)
		return cmd, true
	}

	tlItems := make([]list.Item, 0, len(tlFromCache))
	for _, issue := range tlFromCache {
		tlItems = append(tlItems, issue)
	}
	m.prTLList.SetItems(tlItems)
	m.prTLList.Title = fmt.Sprintf("PR #%d Timeline", prNumber)

	return nil, true
}
