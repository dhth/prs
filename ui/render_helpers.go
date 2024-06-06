package ui

import (
	"fmt"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

const (
	widthBudgetDefault    = 80
	responsiveWidthCutOff = 100
	wideScreenWidthFrac   = 0.9
)

func getPRTitle(pr *pr, terminalDetails terminalDetails) string {
	if pr == nil {
		return ""
	}
	var widthBudget int
	widthBudget = getFracInt(terminalDetails.width, wideScreenWidthFrac)

	var reviewDecision string

	if pr.ReviewDecision != nil {
		switch *pr.ReviewDecision {
		case "CHANGES_REQUESTED":
			reviewDecision = "± "
		case "APPROVED":
			reviewDecision = "✅ "
		case "REVIEW_REQUIRED":
			reviewDecision = "🟡 "
		}
	}
	return Trim(fmt.Sprintf("%s#%2d %s", reviewDecision, pr.Number, pr.PRTitle), widthBudget)
}

func getPRDesc(pr *pr, mode Mode, terminalDetails terminalDetails) string {
	if pr == nil {
		return ""
	}

	var widthBudget = widthBudgetDefault
	var widthFixed int
	var additions string
	var deletions string
	var reviews string
	var desc string

	switch mode {
	case RepoMode:
		widthFixed = 30
	case ReviewerMode:
		widthFixed = 22
	case AuthorMode:
		widthFixed = 20
	}

	if terminalDetails.width > responsiveWidthCutOff {
		widthBudget = getFracInt(terminalDetails.width, wideScreenWidthFrac) - widthFixed
	} else {
		widthBudget = terminalDetails.width - widthFixed
	}

	if widthBudget < 0 {
		widthBudget = widthBudgetDefault
	}

	if pr.Additions > 0 {
		additions = additionsStyle.Render(fmt.Sprintf("+%d", pr.Additions))
	}
	if pr.Deletions > 0 {
		deletions = deletionsStyle.Render(fmt.Sprintf("-%d", pr.Deletions))
	}

	if pr.Reviews.TotalCount > 0 {
		reviews = numReviewsStyle.Render(fmt.Sprintf("%dr", pr.Reviews.TotalCount))
	}

	switch mode {
	case RepoMode:
		updatedAt := dateStyle.Render(RightPadTrim("updated "+humanize.Time(pr.UpdatedAt), getFracInt(widthBudget, 0.3)))
		author := authorStyle(pr.Author.Login).Render(RightPadTrim(pr.Author.Login, getFracInt(widthBudget, 0.7)))
		state := prStyle(pr.State).Render(pr.State)

		desc = fmt.Sprintf("%s%s%s%s%s%s", author, updatedAt, state, additions, deletions, reviews)

	case ReviewerMode:
		updatedAt := dateStyle.Render(RightPadTrim("updated "+humanize.Time(pr.UpdatedAt), getFracInt(widthBudget, 0.3)))
		author := authorStyle(pr.Author.Login).Render(RightPadTrim(pr.Author.Login, getFracInt(widthBudget, 0.4)))
		repo := repoStyle.Render(RightPadTrim(pr.Repository.Name, getFracInt(widthBudget, 0.3)))

		desc = fmt.Sprintf("%s%s%s%s%s%s", author, repo, updatedAt, additions, deletions, reviews)

	case AuthorMode:
		updatedAt := dateStyle.Render(RightPadTrim("updated "+humanize.Time(pr.UpdatedAt), getFracInt(widthBudget, 0.3)))
		repo := repoStyle.Render(RightPadTrim(pr.Repository.Name, getFracInt(widthBudget, 0.7)))

		desc = fmt.Sprintf("%s%s%s%s%s", repo, updatedAt, additions, deletions, reviews)
	}

	return desc
}

func getPRTLItemTitle(item *prTLItem, terminalDetails terminalDetails) string {
	var title string
	var date string

	var widthBudget = widthBudgetDefault
	var widthFixed int

	switch item.Type {
	case tlItemPRCommit:
		widthFixed = 30
	case tlItemHeadRefForcePushed:
		widthFixed = 60
	case tlItemPRReadyForReview:
		widthFixed = 30
	case tlItemPRReviewRequested:
		widthFixed = 30
	case tlItemPRReview:
		widthFixed = 40
	case tlItemMergedEvent:
		widthFixed = 30
	}

	if terminalDetails.width > responsiveWidthCutOff {
		widthBudget = getFracInt(terminalDetails.width, wideScreenWidthFrac) - widthFixed
	} else {
		widthBudget = terminalDetails.width - widthFixed
	}

	if widthBudget < 0 {
		widthBudget = widthBudgetDefault
	}

	switch item.Type {
	case tlItemPRCommit:
		if item.PullRequestCommit.Commit.Author.User != nil {
			author := authorStyle(item.PullRequestCommit.Commit.Author.User.Login).Render(Trim(item.PullRequestCommit.Commit.Author.User.Login, widthBudget))
			date = dateStyle.Render(humanize.Time(item.PullRequestCommit.Commit.CommittedDate))
			title = fmt.Sprintf("%spushed a commit%s", author, date)
		} else {
			title = fmt.Sprintf("%s pushed a commit", Trim(item.PullRequestCommit.Commit.Author.Name, widthBudget))
		}
	case tlItemHeadRefForcePushed:
		actor := authorStyle(item.HeadRefForcePushed.Actor.Login).Render(Trim(item.HeadRefForcePushed.Actor.Login, widthBudget))
		beforeCommitHash := item.HeadRefForcePushed.BeforeCommit.Oid
		afterCommitHash := item.HeadRefForcePushed.AfterCommit.Oid
		if len(beforeCommitHash) >= commitHashLen {
			beforeCommitHash = beforeCommitHash[:commitHashLen]
		}
		if len(afterCommitHash) >= commitHashLen {
			afterCommitHash = afterCommitHash[:commitHashLen]
		}
		date = dateStyle.Render(humanize.Time(item.HeadRefForcePushed.CreatedAt))
		title = fmt.Sprintf("%sforce pushed head ref from %s to %s%s", actor, beforeCommitHash, afterCommitHash, date)
	case tlItemPRReadyForReview:
		actor := authorStyle(item.PullRequestReadyForReview.Actor.Login).Render(Trim(item.PullRequestReadyForReview.Actor.Login, widthBudget))
		title = fmt.Sprintf("%smarked PR as ready for review", actor)
	case tlItemPRReviewRequested:
		actor := authorStyle(item.PullRequestReviewRequested.Actor.Login).Render(Trim(item.PullRequestReviewRequested.Actor.Login, getFracInt(widthBudget, 0.5)))
		reviewer := authorStyle(item.PullRequestReviewRequested.RequestedReviewer.User.Login).Render(Trim(item.PullRequestReviewRequested.RequestedReviewer.User.Login, getFracInt(widthBudget, 0.5)))
		title = fmt.Sprintf("%srequested a review from %s", actor, reviewer)
	case tlItemPRReview:
		author := authorStyle(item.PullRequestReview.Author.Login).Render(Trim(item.PullRequestReview.Author.Login, widthBudget))
		date = dateStyle.Render(humanize.Time(item.PullRequestReview.CreatedAt))
		var comments string
		if item.PullRequestReview.Comments.TotalCount > 1 {
			comments = numCommentsStyle.Render(fmt.Sprintf("with %d comments", item.PullRequestReview.Comments.TotalCount))
		} else if item.PullRequestReview.Comments.TotalCount == 1 {
			comments = numCommentsStyle.Render("with 1 comment")
		}
		title = fmt.Sprintf("%sreviewed%s%s", author, comments, date)
	case tlItemMergedEvent:
		author := authorStyle(item.MergedEvent.Actor.Login).Render(Trim(item.MergedEvent.Actor.Login, widthBudget))
		date = dateStyle.Render(humanize.Time(item.MergedEvent.CreatedAt))
		title = fmt.Sprintf("%smerged the PR%s", author, date)
	}
	return title
}

func getPRTLItemDesc(item *prTLItem, terminalDetails terminalDetails) string {
	var widthBudget = widthBudgetDefault
	widthBudget = terminalDetails.width - 8

	var desc string
	switch item.Type {
	case tlItemPRCommit:
		desc = fmt.Sprintf("📧 %s", Trim(item.PullRequestCommit.Commit.MessageHeadline, widthBudget-10))
	case tlItemHeadRefForcePushed:
		desc = Trim(fmt.Sprintf("💪 %s", item.HeadRefForcePushed.AfterCommit.MessageHeadline), widthBudget)
	case tlItemPRReadyForReview:
		desc = fmt.Sprintf("🚦%s", dateStyle.Render(humanize.Time(item.PullRequestReadyForReview.CreatedAt)))
	case tlItemPRReviewRequested:
		desc = fmt.Sprintf("🙏%s", dateStyle.Render(humanize.Time(item.PullRequestReviewRequested.CreatedAt)))
	case tlItemPRReview:
		reviewState := reviewStyle(item.PullRequestReview.State).Render(item.PullRequestReview.State)
		var comment string
		if item.PullRequestReview.Body != "" {
			comment = Trim(fmt.Sprintf("with comment: %s", strings.Split(item.PullRequestReview.Body, "\r")[0]), widthBudget-14)
		}
		desc = fmt.Sprintf("🔎 %s%s", reviewState, comment)
	case tlItemMergedEvent:
		desc = Trim(fmt.Sprintf("🚀 message: %s", item.MergedEvent.MergeCommit.MessageHeadline), widthBudget)
	}
	return desc
}
