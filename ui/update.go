package ui

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/prs/internal/utils"
)

const (
	useHighPerformanceRenderer = false
	viewPortMoveLineCount      = 5
)

var (
	//go:embed assets/help.md
	helpStr string

	ErrPRDetailsNotCached = errors.New("PR details were not saved")
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.message = ""

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "Q":
			return m, tea.Quit
		case "ctrl+c", "q", "esc":
			switch m.activePane {
			case repoListView:
				if !m.repoChosen {
					return m, tea.Quit
				}
				m.activePane = m.lastPane
			case helpView:
				m.activePane = m.lastPane
			case prRevCmtsView:
				m.prRevCmtVP.GotoTop()
				m.activePane = prTLListView
			case prTLListView:
				m.prTLList.ResetSelected()
				m.activePane = prListView
			case prDetailsView:
				if m.lastPane == m.activePane {
					m.activePane = m.secondLastActivePane
					m.lastPane = prDetailsView
					break
				}
				m.activePane = m.lastPane
				m.lastPane = prDetailsView
			case prListView:
				if m.mode == RepoMode {
					m.activePane = repoListView
					m.repoChosen = false
					break
				}
				return m, tea.Quit
			default:
				return m, tea.Quit
			}

		case "ctrl+r":
			switch m.activePane {
			case prListView:

				switch m.mode {
				case RepoMode:
					cmds = append(cmds, fetchPRSForRepo(m.ghClient, m.repoOwner, m.repoName, m.config.PRCount))
				case QueryMode:
					cmds = append(cmds, fetchPRSFromQuery(m.ghClient, *m.config.Query, m.config.PRCount))
				case ReviewerMode:
					cmds = append(cmds, fetchPRsToReview(m.ghClient, m.userLogin, m.config.PRCount))
				case AuthorMode:
					cmds = append(cmds, fetchAuthoredPRs(m.ghClient, m.userLogin, m.config.PRCount))
				}
				m.prsList.Title = "fetching PRs..."
				m.prsList.Styles.Title = m.prsList.Styles.Title.Background(lipgloss.Color(fetchingColor))

			case prTLListView:
				pr, ok := m.prsList.SelectedItem().(*prResult)
				if !ok {
					break
				}

				repoOwner := pr.pr.Repository.Owner.Login
				repoName := pr.pr.Repository.Name
				prNumber := pr.pr.Number
				cmds = append(cmds, fetchPRTLItems(m.ghClient, repoOwner, repoName, prNumber, 100, true))
				m.prTLList.Title = "fetching timeline..."
				m.prTLList.Styles.Title = m.prTLList.Styles.Title.Background(lipgloss.Color(fetchingColor))
			}
		case "1":
			if m.activePane != prTLListView && m.activePane != prRevCmtsView && m.activePane != prDetailsView {
				break
			}

			switch m.activePane {
			case prDetailsView:
				m.GoToPRDetailSection(0)
			default:
				m.activePane = prListView
			}

		case "enter":
			switch m.activePane {
			case prListView:
				setTlCmd, ok := m.setTL()
				if !ok {
					m.message = "Could't get repo/pr details. Inform @dhth on github."
				} else {
					if setTlCmd != nil {
						cmds = append(cmds, setTlCmd)
					}
				}
			case prTLListView:
				item, ok := m.prTLList.SelectedItem().(*prTLItemResult)
				if !ok {
					break
				}

				if item.item.Type != tlItemPRReview {
					break
				}

				if len(item.item.PullRequestReview.Comments.Nodes) == 0 {
					break
				}

				m.setPRReviewCmt(item.item, 0)
				m.prRevCurCmtNum = 0
				m.activePane = prRevCmtsView

			case repoListView:
				selected := m.repoList.SelectedItem()
				if selected != nil {
					cmds = append(cmds, chooseRepo(selected.FilterValue()))
				}
			}
		case "2":
			if m.activePane != prListView && m.activePane != prRevCmtsView && m.activePane != prDetailsView {
				break
			}

			switch m.activePane {
			case prDetailsView:
				m.GoToPRDetailSection(1)
			default:
				setTlCmd, ok := m.setTL()
				if !ok {
					m.message = "Could't get repo/pr details. Inform @dhth on github."
					break
				}

				if setTlCmd != nil {
					cmds = append(cmds, setTlCmd)
				}
			}

		case "3":
			if m.activePane != prTLListView && m.activePane != prDetailsView {
				break
			}

			switch m.activePane {
			case prDetailsView:
				m.GoToPRDetailSection(2)
			default:
				tlItem, ok := m.prTLList.SelectedItem().(*prTLItemResult)
				if !ok {
					break
				}

				if tlItem.item.Type != tlItemPRReview {
					break
				}

				if len(tlItem.item.PullRequestReview.Comments.Nodes) == 0 {
					break
				}

				m.setPRReviewCmt(tlItem.item, 0)
				m.activePane = prRevCmtsView
			}

		case "4":
			if m.activePane != prDetailsView {
				break
			}

			m.GoToPRDetailSection(3)

		case "5":
			if m.activePane != prDetailsView {
				break
			}

			m.GoToPRDetailSection(4)

		case "6":
			if m.activePane != prDetailsView {
				break
			}

			m.GoToPRDetailSection(5)

		case "j", "down":
			if m.activePane != prRevCmtsView && m.activePane != helpView && m.activePane != prDetailsView {
				break
			}

			switch m.activePane {
			case prRevCmtsView:
				if m.prRevCmtVP.AtBottom() {
					break
				}
				m.prRevCmtVP.LineDown(viewPortMoveLineCount)
			case prDetailsView:
				if m.prDetailsVP.AtBottom() {
					break
				}
				m.prDetailsVP.LineDown(viewPortMoveLineCount)
			case helpView:
				if m.helpVP.AtBottom() {
					break
				}
				m.helpVP.LineDown(viewPortMoveLineCount)
			}

		case "k", "up":
			if m.activePane != prRevCmtsView && m.activePane != helpView && m.activePane != prDetailsView {
				break
			}

			switch m.activePane {
			case prRevCmtsView:
				if m.prRevCmtVP.AtTop() {
					break
				}
				m.prRevCmtVP.LineUp(viewPortMoveLineCount)
			case prDetailsView:
				if m.prDetailsVP.AtTop() {
					break
				}
				m.prDetailsVP.LineUp(viewPortMoveLineCount)
			case helpView:
				if m.helpVP.AtTop() {
					break
				}
				m.helpVP.LineUp(viewPortMoveLineCount)
			}

		case "tab", "shift+tab":
			if m.activePane == helpView || m.activePane == prDetailsView {
				break
			}

			if m.activePane == prListView {
				setTlCmd, ok := m.setTL()
				if !ok {
					m.message = "Could't get repo/pr details. Inform @dhth on github."
				} else {
					if setTlCmd != nil {
						cmds = append(cmds, setTlCmd)
					}
				}
			} else {
				m.activePane = prListView
			}
		case "ctrl+s":
			if m.mode == RepoMode {
				if m.activePane != repoListView {
					m.lastPane = m.activePane
					m.activePane = repoListView
				} else {
					m.activePane = m.lastPane
				}
			}

		case "ctrl+b":
			switch m.activePane {
			case prListView, prDetailsView:
				pr, ok := m.prsList.SelectedItem().(*prResult)
				if !ok {
					break
				}

				cmds = append(cmds, openURLInBrowser(pr.pr.Url))
			case prTLListView, prRevCmtsView:
				item, ok := m.prTLList.SelectedItem().(*prTLItemResult)
				if !ok {
					break
				}

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

		case "ctrl+d":
			if m.activePane != prListView && m.activePane != prTLListView && m.activePane != prRevCmtsView {
				break
			}

			pr, ok := m.prsList.SelectedItem().(*prResult)
			if !ok {
				break
			}

			cmds = append(cmds, showDiff(pr.pr.Repository.Owner.Login,
				pr.pr.Repository.Name,
				pr.pr.Number,
				m.config.DiffPager))

		case "ctrl+v":
			if m.activePane == helpView {
				break
			}
			pr, ok := m.prsList.SelectedItem().(*prResult)
			if !ok {
				break
			}

			cmds = append(cmds, showPR(pr.pr.Repository.Owner.Login,
				pr.pr.Repository.Name,
				pr.pr.Number))

		case "g":
			switch m.activePane {
			case prRevCmtsView:
				m.prRevCmtVP.GotoTop()
			case prDetailsView:
				m.prDetailsVP.GotoTop()
			case helpView:
				m.helpVP.GotoTop()
			}
		case "G":
			switch m.activePane {
			case prRevCmtsView:
				m.prRevCmtVP.GotoBottom()
			case prDetailsView:
				m.prDetailsVP.GotoBottom()
			case helpView:
				m.helpVP.GotoBottom()
			}

		case "K", "[":
			if m.activePane != prDetailsView {
				break
			}

			m.prsList.CursorUp()
			pr, ok := m.prsList.SelectedItem().(*prResult)
			if !ok {
				break
			}

			prDetails, ok := m.prDetailsCache[fmt.Sprintf("%s/%s:%d", pr.pr.Repository.Owner.Login, pr.pr.Repository.Name, pr.pr.Number)]
			if !ok {
				break
			}

			m.setPRDetailsContent(prDetails, 0)
			m.prDetailsCurrentSection = 0

		case "J", "]":
			if m.activePane != prDetailsView {
				break
			}

			m.prsList.CursorDown()
			pr, ok := m.prsList.SelectedItem().(*prResult)
			if !ok {
				break
			}

			prDetails, ok := m.prDetailsCache[fmt.Sprintf("%s/%s:%d", pr.pr.Repository.Owner.Login, pr.pr.Repository.Name, pr.pr.Number)]
			if !ok {
				break
			}

			m.setPRDetailsContent(prDetails, 0)
			m.prDetailsCurrentSection = 0

		case "d":
			if m.activePane != prListView && m.activePane != prDetailsView && m.activePane != prTLListView && m.activePane != prRevCmtsView {
				break
			}

			if m.activePane == prDetailsView {
				m.activePane = m.lastPane
				break
			}

			pr, ok := m.prsList.SelectedItem().(*prResult)
			if !ok {
				break
			}

			prDetails, ok := m.prDetailsCache[fmt.Sprintf("%s/%s:%d", pr.pr.Repository.Owner.Login, pr.pr.Repository.Name, pr.pr.Number)]
			if !ok {
				break
			}

			m.setPRDetailsContent(prDetails, 0)
			m.prDetailsCurrentSection = 0

			m.prDetailsVP.GotoTop()
			m.lastPane = m.activePane
			m.activePane = prDetailsView

		case "l", "n", "right":
			if m.activePane != prDetailsView && m.activePane != prRevCmtsView {
				break
			}

			switch m.activePane {
			case prDetailsView:
				pr, ok := m.prsList.SelectedItem().(*prResult)
				if !ok {
					break
				}

				prDetails, ok := m.prDetailsCache[fmt.Sprintf("%s/%s:%d", pr.pr.Repository.Owner.Login, pr.pr.Repository.Name, pr.pr.Number)]
				if !ok {
					break
				}

				nextSectionFound := false
				var nextSection uint
				if m.prDetailsCurrentSection == uint(len(PRDetailsSectionList)-1) {
					nextSection = 0
				} else {
					nextSection = m.prDetailsCurrentSection + 1
				}

				for {
					switch nextSection {
					case 0:
						nextSectionFound = true
					case 1:
						if prDetails.Body != "" {
							nextSectionFound = true
						}
					case 2:
						if len(prDetails.IssueReferences.Nodes) > 0 {
							nextSectionFound = true
						}
					case 3:
						if len(prDetails.Files.Nodes) > 0 {
							nextSectionFound = true
						}
					case 4:
						if len(prDetails.Commits.Nodes) > 0 {
							nextSectionFound = true
						}
					case 5:
						if len(prDetails.Comments.Nodes) > 0 {
							nextSectionFound = true
						} else {
							nextSection = 0
							nextSectionFound = true
						}
					}

					if nextSectionFound {
						break
					}

					nextSection += 1
				}

				if !nextSectionFound {
					break
				}

				if nextSection > uint(len(PRDetailsSectionList)-1) {
					m.message = "Something went wrong"
					break
				}

				m.setPRDetailsContent(prDetails, PRDetailsSectionList[nextSection])
				m.prDetailsCurrentSection = nextSection

			case prRevCmtsView:
				tlItem, ok := m.prTLList.SelectedItem().(*prTLItemResult)
				if !ok {
					break
				}

				if tlItem.item.Type != tlItemPRReview {
					break
				}

				if len(tlItem.item.PullRequestReview.Comments.Nodes) <= 1 {
					break
				}

				nextCommentIndex := m.prRevCurCmtNum + 1
				if nextCommentIndex > uint(len(tlItem.item.PullRequestReview.Comments.Nodes))-1 {
					nextCommentIndex = 0
				}

				m.setPRReviewCmt(tlItem.item, nextCommentIndex)
				m.prRevCurCmtNum = nextCommentIndex
			}

		case "h", "N", "left":
			if m.activePane != prDetailsView && m.activePane != prRevCmtsView {
				break
			}

			switch m.activePane {
			case prDetailsView:
				pr, ok := m.prsList.SelectedItem().(*prResult)
				if !ok {
					break
				}

				prDetails, ok := m.prDetailsCache[fmt.Sprintf("%s/%s:%d", pr.pr.Repository.Owner.Login, pr.pr.Repository.Name, pr.pr.Number)]
				if !ok {
					break
				}

				prevSectionFound := false
				var prevSection uint
				if m.prDetailsCurrentSection == 0 {
					prevSection = uint(len(PRDetailsSectionList) - 1)
				} else {
					prevSection = m.prDetailsCurrentSection - 1
				}

				for {
					switch prevSection {
					case 0:
						prevSectionFound = true
					case 1:
						if prDetails.Body != "" {
							prevSectionFound = true
						}
					case 2:
						if len(prDetails.IssueReferences.Nodes) > 0 {
							prevSectionFound = true
						}
					case 3:
						if len(prDetails.Files.Nodes) > 0 {
							prevSectionFound = true
						}
					case 4:
						if len(prDetails.Commits.Nodes) > 0 {
							prevSectionFound = true
						}
					case 5:
						if len(prDetails.Comments.Nodes) > 0 {
							prevSectionFound = true
						}
					}

					if prevSectionFound {
						break
					}

					prevSection -= 1
				}

				m.setPRDetailsContent(prDetails, PRDetailsSectionList[prevSection])
				m.prDetailsCurrentSection = prevSection

			case prRevCmtsView:
				tlItem, ok := m.prTLList.SelectedItem().(*prTLItemResult)
				if !ok {
					break
				}

				if tlItem.item.Type != tlItemPRReview {
					break
				}

				if len(tlItem.item.PullRequestReview.Comments.Nodes) <= 1 {
					break
				}

				var prevCommentIndex uint
				if m.prRevCurCmtNum == 0 {
					prevCommentIndex = uint(len(tlItem.item.PullRequestReview.Comments.Nodes) - 1)
				} else {
					prevCommentIndex = m.prRevCurCmtNum - 1
				}

				m.setPRReviewCmt(tlItem.item, prevCommentIndex)
				m.prRevCurCmtNum = prevCommentIndex
			}

		case "?":
			if m.activePane == helpView {
				m.activePane = m.lastPane
				break
			}
			if m.activePane == prDetailsView {
				m.secondLastActivePane = m.lastPane
			}

			m.lastPane = m.activePane
			m.activePane = helpView
		}
	case hideHelpMsg:
		m.showHelp = false
	case tea.WindowSizeMsg:
		w, h := listStyle.GetFrameSize()
		m.terminalDetails.width = msg.Width

		if m.mode == RepoMode {
			m.repoList.SetHeight(msg.Height - h - 2)
			m.repoList.SetWidth(msg.Width - w)
		}

		m.prsList.SetHeight(msg.Height - h - 2)
		m.prsList.SetWidth(msg.Width - w)

		m.prTLList.SetHeight(msg.Height - h - 2)
		m.prTLList.SetWidth(msg.Width - w)

		if !m.prRevCmtVPReady {
			m.prRevCmtVP = viewport.New(msg.Width-2, msg.Height-7)
			m.prRevCmtVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.prRevCmtVPReady = true
			m.prRevCmtVP.KeyMap.HalfPageDown.SetKeys("ctrl+d")
			m.prRevCmtVP.KeyMap.Up.SetEnabled(false)
			m.prRevCmtVP.KeyMap.Down.SetEnabled(false)
		} else {
			m.prRevCmtVP.Width = msg.Width - 2
			m.prRevCmtVP.Height = msg.Height - 7
		}

		if !m.prDetailsVPReady {
			m.prDetailsVP = viewport.New(msg.Width-2, msg.Height-7)
			m.prDetailsVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.prDetailsVPReady = true
			m.prDetailsVP.KeyMap.HalfPageDown.SetKeys("ctrl+d")
			m.prDetailsVP.KeyMap.Up.SetEnabled(false)
			m.prDetailsVP.KeyMap.Down.SetEnabled(false)
		} else {
			m.prDetailsVP.Width = msg.Width - 2
			m.prDetailsVP.Height = msg.Height - 7
		}

		vpWrap := (msg.Width - 4)
		if vpWrap > viewPortWrapUpperLimit {
			vpWrap = viewPortWrapUpperLimit
		}

		m.mdRenderer, _ = utils.GetMarkDownRenderer(vpWrap)

		helpToRender := helpStr
		switch m.mdRenderer {
		case nil:
			break
		default:
			helpStrGl, err := m.mdRenderer.Render(helpStr)
			if err != nil {
				break
			}
			helpToRender = helpStrGl
		}

		if !m.helpVPReady {
			m.helpVP = viewport.New(msg.Width, msg.Height-7)
			m.helpVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.helpVP.SetContent(helpToRender)
			m.helpVP.KeyMap.Up.SetEnabled(false)
			m.helpVP.KeyMap.Down.SetEnabled(false)
			m.helpVPReady = true
		} else {
			m.helpVP.Width = msg.Width
			m.helpVP.Height = msg.Height - 7
		}

		prs := make([]list.Item, len(m.prCache))
		for i := 0; i < len(m.prCache); i++ {
			m.prCache[i].title = getPRTitle(m.prCache[i].pr)
			m.prCache[i].description = getPRDesc(m.prCache[i].pr, m.mode, m.terminalDetails)
			prs[i] = m.prCache[i]
		}
		m.prsList.SetItems(prs)

		if m.activePane == prTLListView {
			m.setTL()
		}

	case repoChosenMsg:
		repoDetails := strings.Split(msg.repo, ":::")
		if len(repoDetails) != 2 {
			m.message = "Something went horribly wrong. Let @dhth know about this failure."
		} else {
			m.repoChosen = true
			m.prsList.Title = "fetching PRs..."
			m.prsList.Styles.Title = m.prsList.Styles.Title.Background(lipgloss.Color(fetchingColor))
			m.repoOwner = repoDetails[0]
			m.repoName = repoDetails[1]
			m.activePane = prListView
			m.prsList.ResetSelected()
			m.prTLList.ResetSelected()
			cmds = append(cmds, fetchPRSForRepo(m.ghClient, m.repoOwner, m.repoName, m.config.PRCount))
		}
	case viewerLoginFetched:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error fetching gh username: %s", msg.err)
		} else {
			m.userLogin = msg.login
			switch m.mode {
			case ReviewerMode:
				cmds = append(cmds, fetchPRsToReview(m.ghClient, m.userLogin, m.config.PRCount))
			case AuthorMode:
				cmds = append(cmds, fetchAuthoredPRs(m.ghClient, m.userLogin, m.config.PRCount))
			}
		}
	case prsFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
			m.prsList.Title = "error"
			break
		}

		prs := make([]list.Item, len(msg.prs))
		prResults := make([]*prResult, len(msg.prs))

		for i, pr := range msg.prs {
			prResults[i] = &prResult{
				pr:          &pr,
				title:       getPRTitle(&pr),
				description: getPRDesc(&pr, m.mode, m.terminalDetails),
			}
			prs[i] = prResults[i]
		}

		m.prCache = prResults
		m.prsList.SetItems(prs)

		switch m.mode {
		case RepoMode:
			m.prsList.Title = fmt.Sprintf("PRs (%s)", m.repoName)
		case QueryMode:
			m.prsList.Title = "Results"
		}

		m.prsList.ResetSelected()
		m.prsList.Styles.Title = m.prsList.Styles.Title.Background(lipgloss.Color(prListColor))

		for _, pr := range msg.prs {
			cmds = append(cmds, fetchPRTLItems(m.ghClient,
				pr.Repository.Owner.Login,
				pr.Repository.Name,
				pr.Number,
				100,
				false,
			))
			cmds = append(cmds, fetchPRMetadata(m.ghClient,
				pr.Repository.Owner.Login,
				pr.Repository.Name,
				pr.Number,
			))
		}

	case reviewPRsFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
			break
		}

		prs := make([]list.Item, len(msg.prs))

		prResults := make([]*prResult, len(msg.prs))

		for i, pr := range msg.prs {
			prResults[i] = &prResult{
				pr:          &pr,
				title:       getPRTitle(&pr),
				description: getPRDesc(&pr, m.mode, m.terminalDetails),
			}
			prs[i] = prResults[i]
		}

		m.prCache = prResults
		m.prsList.SetItems(prs)
		m.prsList.ResetSelected()
		m.prsList.Title = "Open PRs requesting your review"
		m.prsList.Styles.Title = m.prsList.Styles.Title.Background(lipgloss.Color(prListColor))

		if len(msg.prs) > 0 {
			for _, pr := range msg.prs {
				cmds = append(cmds, fetchPRTLItems(m.ghClient, pr.Repository.Owner.Login, pr.Repository.Name, pr.Number, 100, false))
				cmds = append(cmds, fetchPRMetadata(m.ghClient,
					pr.Repository.Owner.Login,
					pr.Repository.Name,
					pr.Number,
				))
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
					title:       getPRTitle(&pr),
					description: getPRDesc(&pr, m.mode, m.terminalDetails),
				}
				prs[i] = prResults[i]
			}

			m.prCache = prResults
			m.prsList.SetItems(prs)
			m.prsList.Title = "Open PRs authored by you"
			m.prsList.ResetSelected()
			m.prsList.Styles.Title = m.prsList.Styles.Title.Background(lipgloss.Color(prListColor))

			if len(msg.prs) > 0 {
				for _, pr := range msg.prs {
					cmds = append(cmds, fetchPRTLItems(m.ghClient, pr.Repository.Owner.Login, pr.Repository.Name, pr.Number, 100, false))
					cmds = append(cmds, fetchPRMetadata(m.ghClient,
						pr.Repository.Owner.Login,
						pr.Repository.Name,
						pr.Number,
					))
				}
			}
		}

	case prMetadataFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
			break
		}

		m.prDetailsCache[fmt.Sprintf("%s/%s:%d", msg.repoOwner, msg.repoName, msg.prNumber)] = msg.metadata

	case prTLFetchedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
			break
		}

		tlItemsResult := make([]*prTLItemResult, len(msg.prTLItems))

		for i, item := range msg.prTLItems {
			tlItemsResult[i] = &prTLItemResult{
				item:        &item,
				title:       getPRTLItemTitle(&item),
				description: getPRTLItemDesc(&item),
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
			m.prTLList.Styles.Title = m.prTLList.Styles.Title.Background(lipgloss.Color(prTLListColor))
			m.activePane = prTLListView
		}

		m.prTLList.ResetSelected()

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
	case prListView:
		m.prsList, cmd = m.prsList.Update(msg)
		cmds = append(cmds, cmd)
	case prTLListView:
		m.prTLList, cmd = m.prTLList.Update(msg)
		cmds = append(cmds, cmd)
	case prDetailsView:
		m.prDetailsVP, cmd = m.prDetailsVP.Update(msg)
		cmds = append(cmds, cmd)
	case prRevCmtsView:
		m.prRevCmtVP, cmd = m.prRevCmtVP.Update(msg)
		cmds = append(cmds, cmd)
	case repoListView:
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

	prItem, prOk := m.prsList.SelectedItem().(*prResult)
	if !prOk {
		return nil, false
	}

	repoOwner = prItem.pr.Repository.Owner.Login
	repoName = prItem.pr.Repository.Name
	prNumber = prItem.pr.Number

	tlFromCache, ok := m.prTLCache[fmt.Sprintf("%s/%s:%d", repoOwner, repoName, prNumber)]
	if !ok {
		cmd = fetchPRTLItems(m.ghClient, repoOwner, repoName, prNumber, 100, true)
		return cmd, true
	}

	tlItems := make([]list.Item, len(tlFromCache))

	// this list always get rerendered as it seems to be preferrable over recomputing the string rep of every item in
	// every list in m.prTLCache when the terminal window is resized
	for i, result := range tlFromCache {
		title := getPRTLItemTitle(result.item)
		description := getPRTLItemDesc(result.item)

		result.title = title
		result.description = description

		tlItems[i] = result
	}

	m.prTLList.SetItems(tlItems)
	m.prTLList.Title = fmt.Sprintf("PR #%d Timeline", prNumber)
	m.activePane = prTLListView

	return nil, true
}
