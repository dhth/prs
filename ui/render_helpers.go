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

func getPRTitle(pr *pr) string {
	if pr == nil {
		return ""
	}

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
	return fmt.Sprintf("%s#%2d %s", reviewDecision, pr.Number, pr.PRTitle)
}

func getPRDesc(pr *pr, mode Mode, terminalDetails terminalDetails) string {
	if pr == nil {
		return ""
	}

	var widthBudget int
	var widthFixed int
	var additions string
	var deletions string
	var reviews string
	var desc string

	switch mode {
	case RepoMode:
		widthFixed = 30
	case QueryMode:
		widthFixed = 16
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
		author := getDynamicStyle(pr.Author.Login).Render(RightPadTrim(pr.Author.Login, getFracInt(widthBudget, 0.7)))
		state := prStyle(pr.State).Render(pr.State)

		desc = fmt.Sprintf("%s%s%s%s%s%s", author, updatedAt, state, additions, deletions, reviews)

	case QueryMode:
		repoStr := fmt.Sprintf("%s/%s", pr.Repository.Owner.Login, pr.Repository.Name)
		updatedAt := dateStyle.Render(RightPadTrim("updated "+humanize.Time(pr.UpdatedAt), getFracInt(widthBudget, 0.3)))
		author := getDynamicStyle(pr.Author.Login).Render(RightPadTrim(pr.Author.Login, getFracInt(widthBudget, 0.4)))
		state := prStyle(pr.State).Render(pr.State)
		repo := getDynamicStyle(repoStr).Render(RightPadTrim(repoStr, getFracInt(widthBudget, 0.3)))

		desc = fmt.Sprintf("%s%s%s%s%s%s%s", author, repo, updatedAt, state, additions, deletions, reviews)
	}

	return desc
}

func getPRTLItemTitle(item *prTLItem) string {
	var title string
	var date string

	switch item.Type {
	case tlItemPRCommit:
		if item.PullRequestCommit.Commit.Author.User != nil {
			author := getDynamicStyle(item.PullRequestCommit.Commit.Author.User.Login).Render(item.PullRequestCommit.Commit.Author.User.Login)
			date = dateStyle.Render(humanize.Time(item.PullRequestCommit.Commit.CommittedDate))
			title = fmt.Sprintf("%spushed a commit%s", author, date)
		} else {
			title = fmt.Sprintf("%s pushed a commit", item.PullRequestCommit.Commit.Author.Name)
		}

	case tlItemHeadRefForcePushed:
		actor := getDynamicStyle(item.HeadRefForcePushed.Actor.Login).Render(item.HeadRefForcePushed.Actor.Login)
		beforeCommitHash := item.HeadRefForcePushed.BeforeCommit.AbbreviatedOid
		afterCommitHash := item.HeadRefForcePushed.AfterCommit.AbbreviatedOid
		date = dateStyle.Render(humanize.Time(item.HeadRefForcePushed.CreatedAt))
		title = fmt.Sprintf("%sforce pushed head ref from %s to %s%s", actor, beforeCommitHash, afterCommitHash, date)

	case tlItemPRReadyForReview:
		actor := getDynamicStyle(item.PullRequestReadyForReview.Actor.Login).Render(item.PullRequestReadyForReview.Actor.Login)
		title = fmt.Sprintf("%smarked PR as ready for review", actor)

	case tlItemPRReviewRequested:
		actor := getDynamicStyle(item.PullRequestReviewRequested.Actor.Login).Render(item.PullRequestReviewRequested.Actor.Login)
		reviewer := getDynamicStyle(item.PullRequestReviewRequested.RequestedReviewer.User.Login).Render(item.PullRequestReviewRequested.RequestedReviewer.User.Login)
		title = fmt.Sprintf("%srequested a review from %s", actor, reviewer)

	case tlItemPRReview:
		author := getDynamicStyle(item.PullRequestReview.Author.Login).Render(item.PullRequestReview.Author.Login)
		date = dateStyle.Render(humanize.Time(item.PullRequestReview.CreatedAt))
		var comments string
		var more string
		if item.PullRequestReview.Comments.TotalCount > 0 {
			more = " ⏎"
		}
		if item.PullRequestReview.Comments.TotalCount > 1 {
			comments = numCommentsStyle.Render(fmt.Sprintf("with %d comments", item.PullRequestReview.Comments.TotalCount))
		} else if item.PullRequestReview.Comments.TotalCount == 1 {
			comments = numCommentsStyle.Render("with 1 comment")
		}
		title = fmt.Sprintf("%sreviewed%s%s%s", author, comments, date, more)

	case tlItemMergedEvent:
		author := getDynamicStyle(item.MergedEvent.Actor.Login).Render(item.MergedEvent.Actor.Login)
		date = dateStyle.Render(humanize.Time(item.MergedEvent.CreatedAt))
		title = fmt.Sprintf("%smerged the PR%s", author, date)
	}
	return title
}

func getPRTLItemDesc(item *prTLItem) string {
	var desc string
	switch item.Type {
	case tlItemPRCommit:
		desc = fmt.Sprintf("📧 %s", item.PullRequestCommit.Commit.MessageHeadline)
	case tlItemHeadRefForcePushed:
		desc = fmt.Sprintf("💪 %s", item.HeadRefForcePushed.AfterCommit.MessageHeadline)
	case tlItemPRReadyForReview:
		desc = fmt.Sprintf("🚦%s", dateStyle.Render(humanize.Time(item.PullRequestReadyForReview.CreatedAt)))
	case tlItemPRReviewRequested:
		desc = fmt.Sprintf("🙏%s", dateStyle.Render(humanize.Time(item.PullRequestReviewRequested.CreatedAt)))
	case tlItemPRReview:
		reviewState := reviewStyle(item.PullRequestReview.State).Render(item.PullRequestReview.State)
		var comment string
		if item.PullRequestReview.Body != "" {
			comment = fmt.Sprintf("with comment: %s", strings.Split(item.PullRequestReview.Body, "\r")[0])
		}
		desc = fmt.Sprintf("🔎 %s%s", reviewState, comment)
	case tlItemMergedEvent:
		desc = fmt.Sprintf("🚀 message: %s", item.MergedEvent.MergeCommit.MessageHeadline)
	}
	return desc
}
