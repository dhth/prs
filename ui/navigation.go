package ui

import (
	"fmt"
	"strings"
)

const (
	maxCommentsForNavIndicator = 8
)

func (m *model) setPRDetailsContent(prDetails prDetails, section PRDetailSection) {
	content := fmt.Sprintf(`# %s (%s/%s/pull/%d)
`, prDetails.PRTitle, prDetails.Repository.Owner.Login, prDetails.Repository.Name, prDetails.Number,
	)

	switch section {
	case PRMetadata:
		content += prDetails.Metadata()
	case PRDescription:
		content += prDetails.Description()
	case PRChecks:
		content += prDetails.Checks()
	case PRReferences:
		content += prDetails.References()
	case PRFilesChanged:
		content += prDetails.FilesChanged()
	case PRCommits:
		content += prDetails.CommitsList()
	case PRComments:
		content += prDetails.CommentsList()
	}

	glErr := true
	if m.mdRenderer != nil {
		contentGl, err := m.mdRenderer.Render(content)
		if err == nil {
			m.prDetailsVP.SetContent(contentGl)
			glErr = false
		}
	}
	if glErr {
		m.prDetailsVP.SetContent(content)
	}

	sections := make([]string, len(PRDetailsSectionList))
	for i := 0; i < len(PRDetailsSectionList); i++ {
		sections[i] = "◯"
	}

	if prDetails.Body == "" {
		sections[PRDescription] = "◌"
	}
	// not foolproof, but should work in most cases
	// func (pr prDetails) Checks() will return with an appropriate message in
	// that case
	if !(len(prDetails.LastCommit.Nodes) > 0 && prDetails.LastCommit.Nodes[0].Commit.StatusCheckRollup != nil) {
		sections[PRChecks] = "◌"
	}
	if len(prDetails.IssueReferences.Nodes) == 0 {
		sections[PRReferences] = "◌"
	}
	if len(prDetails.Files.Nodes) == 0 {
		sections[PRFilesChanged] = "◌"
	}
	if len(prDetails.Commits.Nodes) == 0 {
		sections[PRCommits] = "◌"
	}
	if len(prDetails.Comments.Nodes) == 0 {
		sections[PRComments] = "◌"
	}

	sections[section] = "●"

	m.prDetailsTitle = fmt.Sprintf("PR Details%s", "  "+strings.Join(sections, " "))

	m.prDetailsVP.GotoTop()
}

func (m *model) GoToPRDetailSection(section uint) {
	if m.prDetailsCurrentSection == section {
		return
	}
	pr, ok := m.prsList.SelectedItem().(*prResult)
	if !ok {
		return
	}

	prDetails, ok := m.prDetailsCache[fmt.Sprintf("%s/%s:%d", pr.pr.Repository.Owner.Login, pr.pr.Repository.Name, pr.pr.Number)]
	if !ok {
		return
	}
	switch section {
	case 1:
		if prDetails.Body == "" {
			return
		}
	case 2:
		if !(len(prDetails.LastCommit.Nodes) > 0 && prDetails.LastCommit.Nodes[0].Commit.StatusCheckRollup != nil) {
			return
		}
	case 3:
		if len(prDetails.IssueReferences.Nodes) == 0 {
			return
		}
	case 4:
		if len(prDetails.Files.Nodes) == 0 {
			return
		}
	case 5:
		if len(prDetails.Commits.Nodes) == 0 {
			return
		}
	case 6:
		if len(prDetails.Comments.Nodes) == 0 {
			return
		}
	}

	m.setPRDetailsContent(prDetails, PRDetailsSectionList[section])
	m.prDetailsCurrentSection = section
}

func (m *model) setPRReviewCmt(tlItem *prTLItem, commentNum uint) {
	revCmts := tlItem.PullRequestReview.Comments.Nodes
	var sectionsStr string

	if len(revCmts) > maxCommentsForNavIndicator {
		sectionsStr = fmt.Sprintf("%d/%d", commentNum+1, len(revCmts))
	} else if len(revCmts) > 1 {
		sections := make([]string, len(revCmts))
		for i := 0; i < len(revCmts); i++ {
			sections[i] = "◯"
		}
		sections[commentNum] = "●"
		sectionsStr = "  " + strings.Join(sections, " ")
	}

	var outdated string
	if revCmts[commentNum].Outdated {
		outdated = " `(outdated)`"
	}

	content := fmt.Sprintf("# from @%s\n## %s%s\n%s\n```diff\n%s\n```", tlItem.PullRequestReview.Author.Login, revCmts[commentNum].Path, outdated, revCmts[commentNum].Body, revCmts[commentNum].DiffHunk)

	glErr := true
	if m.mdRenderer != nil {
		contentGl, err := m.mdRenderer.Render(content)
		if err == nil {
			m.prTLItemDetailVP.SetContent(contentGl)
			glErr = false
		}
	}
	if glErr {
		m.prDetailsVP.SetContent(content)
	}

	m.prTLItemDetailTitle = fmt.Sprintf("Review Comments%s", sectionsStr)
	m.prTLItemDetailVP.GotoTop()
}
