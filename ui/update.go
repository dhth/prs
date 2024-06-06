package ui

import (
	"fmt"
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
				pr, ok := m.prsList.SelectedItem().(*prResult)
				if ok {
					repoOwner := pr.pr.Repository.Owner.Login
					repoName := pr.pr.Repository.Name
					prNumber := pr.pr.Number
					cmds = append(cmds, fetchPRTLItems(m.ghClient, repoOwner, repoName, prNumber, 100, false))
					m.prTLList.ResetSelected()
				}
			}
		case "1":
			if m.activePane != prList {
				m.activePane = prList
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
				}
			case prTLList:
				item, ok := m.prTLList.SelectedItem().(*prTLItemResult)
				if ok {
					if item.item.Type == tlItemPRReview {
						revCmts := item.item.PullRequestReview.Comments.Nodes
						if len(revCmts) > 0 {
							var prReviewCmts string
							for _, cmt := range revCmts {
								prReviewCmts += cmt.render()
								prReviewCmts += "\n\n"
								prReviewCmts += reviewCmtDividerStyle.Render(strings.Repeat("-", int(float64(m.terminalDetails.width)*0.6)))
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
				}
			}
		case "3":
			if m.activePane == prTLList {
				item, ok := m.prTLList.SelectedItem().(*prTLItemResult)
				if ok {
					if item.item.Type == tlItemPRReview {
						revCmts := item.item.PullRequestReview.Comments.Nodes
						if len(revCmts) > 0 {
							var prReviewCmts string
							for _, cmt := range revCmts {
								prReviewCmts += cmt.render()
								prReviewCmts += "\n\n"
								prReviewCmts += reviewCmtDividerStyle.Render(strings.Repeat("*", int(float64(m.terminalDetails.width)*0.6)))
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
					pr, ok := m.prsList.SelectedItem().(*prResult)
					if ok {
						url = pr.pr.Url
					}
				case ReviewerMode:
					pr, ok := m.prsList.SelectedItem().(*prResult)
					if ok {
						url = pr.pr.Url
					}
				}
				cmds = append(cmds, openURLInBrowser(url))
			case prTLList, prRevCmts:
				item, ok := m.prTLList.SelectedItem().(*prTLItemResult)
				if ok {
					switch item.item.Type {
					case tlItemPRCommit:
						cmds = append(cmds, openURLInBrowser(item.item.PullRequestCommit.Url))
					case tlItemHeadRefForcePushed:
						cmds = append(cmds, openURLInBrowser(item.item.HeadRefForcePushed.AfterCommit.Url))
					case tlItemPRReview:
						cmds = append(cmds, openURLInBrowser(item.item.PullRequestReview.Url))
					case tlItemMergedEvent:
						cmds = append(cmds, openURLInBrowser(item.item.MergedEvent.Url))
					}
				}
			}
		case "ctrl+d":
			if m.activePane == prList || m.activePane == prTLList {
				switch m.mode {
				case RepoMode:
					pr, ok := m.prsList.SelectedItem().(*prResult)
					if ok {
						cmds = append(cmds, showDiff(m.repoOwner,
							m.repoName,
							pr.pr.Number,
							m.config.DiffPager))
					}
				case ReviewerMode:
					pr, ok := m.prsList.SelectedItem().(*prResult)
					if ok {
						cmds = append(cmds, showDiff(pr.pr.Repository.Owner.Login,
							pr.pr.Repository.Name,
							pr.pr.Number,
							m.config.DiffPager))
					}
				}
			}
		case "ctrl+v":
			if m.activePane == prList || m.activePane == prTLList {
				switch m.mode {
				case RepoMode:
					pr, ok := m.prsList.SelectedItem().(*prResult)
					if ok {
						cmds = append(cmds, showPR(m.repoOwner, m.repoName, pr.pr.Number))
					}
				case ReviewerMode:
					pr, ok := m.prsList.SelectedItem().(*prResult)
					if ok {
						cmds = append(cmds, showPR(pr.pr.Repository.Owner.Login,
							pr.pr.Repository.Name,
							pr.pr.Number))
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
		m.terminalDetails.width = msg.Width

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
			m.prRevCmtVP = viewport.New(msg.Width-2, msg.Height-7)
			m.prRevCmtVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.prRevCmtVPReady = true
		} else {
			m.prRevCmtVP.Width = msg.Width - 2
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

		prs := make([]list.Item, len(m.prCache))
		for i := 0; i < len(m.prCache); i++ {
			m.prCache[i].title = getPRTitle(m.prCache[i].pr, m.terminalDetails)
			m.prCache[i].description = getPRDesc(m.prCache[i].pr, m.mode, m.terminalDetails)
			prs[i] = m.prCache[i]
		}
		m.prsList.SetItems(prs)

		if m.activePane == prTLList {
			m.setTL()
		}

	case repoChosenMsg:
		repoDetails := strings.Split(msg.repo, ":::")
		if len(repoDetails) != 2 {
			m.message = "Something went horribly wrong. Let @dhth know about this failure."
		} else {
			m.repoChosen = true
			m.prsList.Title = "fetching..."
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
	case prsFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
			m.prsList.Title = "error"
		} else {
			prs := make([]list.Item, len(msg.prs))
			prResults := make([]*prResult, len(msg.prs))

			for i, pr := range msg.prs {
				prResults[i] = &prResult{
					pr:          &pr,
					title:       getPRTitle(&pr, m.terminalDetails),
					description: getPRDesc(&pr, RepoMode, m.terminalDetails),
				}
				prs[i] = prResults[i]
			}

			m.prCache = prResults
			m.prsList.SetItems(prs)
			m.prsList.Title = fmt.Sprintf("PRs (%s)", m.repoName)

			if len(msg.prs) > 0 {
				for _, pr := range msg.prs {
					cmds = append(cmds, fetchPRTLItems(m.ghClient,
						pr.Repository.Owner.Login,
						pr.Repository.Name,
						pr.Number,
						100,
						false,
					))
				}
			}
		}
	case reviewPRsFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
		} else {
			prs := make([]list.Item, len(msg.prs))

			prResults := make([]*prResult, len(msg.prs))

			for i, pr := range msg.prs {
				prResults[i] = &prResult{
					pr:          &pr,
					title:       getPRTitle(&pr, m.terminalDetails),
					description: getPRDesc(&pr, m.mode, m.terminalDetails),
				}
				prs[i] = prResults[i]
			}

			m.prCache = prResults
			m.prsList.SetItems(prs)
			m.prsList.Title = "PRs requesting your review"

			if len(msg.prs) > 0 {
				for _, pr := range msg.prs {
					cmds = append(cmds, fetchPRTLItems(m.ghClient, pr.Repository.Owner.Login, pr.Repository.Name, pr.Number, 100, false))
				}
			}
		}
	case authoredPRsFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
		} else {
			prs := make([]list.Item, len(msg.prs))

			prResults := make([]*prResult, len(msg.prs))

			for i, pr := range msg.prs {
				prResults[i] = &prResult{
					pr:          &pr,
					title:       getPRTitle(&pr, m.terminalDetails),
					description: getPRDesc(&pr, m.mode, m.terminalDetails),
				}
				prs[i] = prResults[i]
			}

			m.prCache = prResults
			m.prsList.SetItems(prs)
			m.prsList.Title = "Open PRs authored by you"

			if len(msg.prs) > 0 {
				for _, pr := range msg.prs {
					cmds = append(cmds, fetchPRTLItems(m.ghClient, pr.Repository.Owner.Login, pr.Repository.Name, pr.Number, 100, false))
				}
			}
		}
	case prTLFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
		} else {

			tlItemsResult := make([]*prTLItemResult, len(msg.prTLItems))

			for i, item := range msg.prTLItems {
				tlItemsResult[i] = &prTLItemResult{
					item:        &item,
					title:       getPRTLItemTitle(&item, m.terminalDetails),
					description: getPRTLItemDesc(&item, m.terminalDetails),
				}
			}
			m.prTLCache[fmt.Sprintf("%s/%s:%d", msg.repoOwner, msg.repoName, msg.prNumber)] = tlItemsResult

			if msg.setItems {
				prTLItems := make([]list.Item, len(msg.prTLItems))
				for i, result := range tlItemsResult {
					prTLItems[i] = result
				}
				m.prTLList.SetItems(prTLItems)
				m.prTLList.Title = fmt.Sprintf("PR #%d Timeline", msg.prNumber)
				m.activePane = prTLList
			}
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

	pr, prOk := m.prsList.SelectedItem().(*prResult)
	if !prOk {
		return nil, false
	}

	repoOwner = pr.pr.Repository.Owner.Login
	repoName = pr.pr.Repository.Name
	prNumber = pr.pr.Number

	tlFromCache, ok := m.prTLCache[fmt.Sprintf("%s/%s:%d", repoOwner, repoName, prNumber)]
	if !ok {
		cmd = fetchPRTLItems(m.ghClient, repoOwner, repoName, prNumber, 100, true)
		return cmd, true
	}

	tlItems := make([]list.Item, len(tlFromCache))

	// this list always get rerendered as it seems to be preferrable over recomputing the string rep of every item in
	// every list in m.prTLCache when the terminal window is resized
	for i, result := range tlFromCache {
		title := getPRTLItemTitle(result.item, m.terminalDetails)
		description := getPRTLItemDesc(result.item, m.terminalDetails)

		result.title = title
		result.description = description

		tlItems[i] = result
	}

	m.prTLList.SetItems(tlItems)
	m.prTLList.Title = fmt.Sprintf("PR #%d Timeline", prNumber)
	m.activePane = prTLList

	return nil, true
}
