package ui

import (
	"fmt"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
)

const (
	prStateOpen            = "OPEN"
	prStateMerged          = "MERGED"
	prStateClosed          = "CLOSED"
	tlItemPRCommit         = "PullRequestCommit"
	tlItemPRReview         = "PullRequestReview"
	tlItemMergedEvent      = "MergedEvent"
	reviewPending          = "PENDING"
	reviewCommented        = "COMMENTED"
	reviewApproved         = "APPROVED"
	reviewChangesRequested = "CHANGES_REQUESTED"
	reviewDismissed        = "DISMISSED"
)

type SourceConfig struct {
	DiffPager *string `yaml:"diff-pager"`
	PRCount   int     `yaml:"pr-count"`
	Sources   []struct {
		Owner string `yaml:"owner"`
		Repos []struct {
			Name string `yaml:"name"`
		} `yaml:"repos"`
	} `yaml:"sources"`
}

type Repo struct {
	Owner string
	Name  string
}

type Config struct {
	DiffPager *string
	PRCount   int
	Repos     []Repo
}

type pr struct {
	Number     int
	PRTitle    string `graphql:"prTitle: title"`
	Repository struct {
		Owner struct {
			Login string
		}
		Name string
	}
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  string
	Author    struct {
		Login string
	}
	Url       string
	Additions int
	Deletions int
	Reviews   struct {
		TotalCount int
	}
}

type prsQuery struct {
	RepositoryOwner struct {
		Repository struct {
			PullRequests struct {
				Nodes []pr
			} `graphql:"pullRequests(first: $pullRequestCount, states: [OPEN, MERGED, CLOSED], orderBy: {field: UPDATED_AT, direction: DESC})"`
		} `graphql:"repository(name: $repositoryName)"`
	} `graphql:"repositoryOwner(login: $repositoryOwner)"`
}

type prReviewComment struct {
	CreatedAt time.Time
	Body      string
	Outdated  bool
	DiffHunk  string
	Path      string
	Url       string
}

type prTLItem struct {
	Type              string `graphql:"type: __typename"`
	PullRequestCommit struct {
		Url    string
		Commit struct {
			Oid             string
			CommittedDate   time.Time
			MessageHeadline string
			Author          struct {
				Name string
				User *struct {
					Login string
				}
			}
		}
	} `graphql:"... on PullRequestCommit"`
	PullRequestReview struct {
		Url       string
		CreatedAt time.Time
		State     string
		Body      string
		Comments  struct {
			TotalCount int
			Nodes      []prReviewComment
		} `graphql:"comments(last: 100)"`
		Author struct {
			Login string
		}
	} `graphql:"... on PullRequestReview"`
	MergedEvent struct {
		CreatedAt   time.Time
		Url         string
		MergeCommit struct {
			Oid             string
			MessageHeadline string
		} `graphql:"mergeCommit: commit"`
		Actor struct {
			Login string
		}
	} `graphql:"... on MergedEvent"`
}

type prTLQuery struct {
	RepositoryOwner struct {
		Repository struct {
			PullRequest struct {
				TimelineItems struct {
					Nodes []prTLItem
				} `graphql:"timelineItems(last: $timelineItemsCount, itemTypes: [PULL_REQUEST_COMMIT, MERGED_EVENT, PULL_REQUEST_REVIEW])"`
			} `graphql:"pullRequest(number: $pullRequestNumber)"`
		} `graphql:"repository(name: $repositoryName)"`
	} `graphql:"repositoryOwner(login: $repositoryOwner)"`
}

func (repo Repo) Title() string {
	return repo.Name
}

func (repo Repo) Description() string {
	return repo.Owner
}

func (repo Repo) FilterValue() string {
	return fmt.Sprintf("%s:::%s", repo.Owner, repo.Name)
}

func (pr pr) Title() string {
	return fmt.Sprintf("#%2d %s", pr.Number, pr.PRTitle)
}

func (pr pr) Description() string {
	var additions string
	var deletions string
	var reviews string

	author := authorStyle(pr.Author.Login).Render(RightPadTrim(pr.Author.Login, 80))
	state := prStyle(pr.State).Render(pr.State)
	updatedAt := dateStyle.Render(RightPadTrim("updated "+humanize.Time(pr.UpdatedAt), 24))

	if pr.Additions > 0 {
		additions = additionsStyle.Render(fmt.Sprintf("+%d", pr.Additions))
	}
	if pr.Deletions > 0 {
		deletions = deletionsStyle.Render(fmt.Sprintf("-%d", pr.Deletions))
	}

	if pr.Reviews.TotalCount > 0 {
		reviews = numReviewsStyle.Render(fmt.Sprintf("%dr", pr.Reviews.TotalCount))
	}
	return fmt.Sprintf("%s%s%s%s%s%s", author, updatedAt, state, additions, deletions, reviews)
}

func (pr pr) FilterValue() string {
	return fmt.Sprintf("%d", pr.Number)
}

func (item prTLItem) Title() string {
	var title string
	var date string
	switch item.Type {
	case tlItemPRCommit:
		if item.PullRequestCommit.Commit.Author.User != nil {
			author := authorStyle(item.PullRequestCommit.Commit.Author.User.Login).Render(Trim(item.PullRequestCommit.Commit.Author.User.Login, 50))
			date = dateStyle.Render(humanize.Time(item.PullRequestCommit.Commit.CommittedDate))
			title = fmt.Sprintf("%s pushed a commit %s", author, date)
		} else {
			title = fmt.Sprintf("%s pushed a commit", item.PullRequestCommit.Commit.Author.Name)
		}
	case tlItemPRReview:
		author := authorStyle(item.PullRequestReview.Author.Login).Render(Trim(item.PullRequestReview.Author.Login, 50))
		date = dateStyle.Render(humanize.Time(item.PullRequestReview.CreatedAt))
		var comments string
		if item.PullRequestReview.Comments.TotalCount > 1 {
			comments = numCommentsStyle.Render(fmt.Sprintf("with %d comments", item.PullRequestReview.Comments.TotalCount))
		} else if item.PullRequestReview.Comments.TotalCount == 1 {
			comments = numCommentsStyle.Render("with 1 comment")
		}
		title = fmt.Sprintf("%sreviewed %s %s", author, comments, date)
	case tlItemMergedEvent:
		author := authorStyle(item.MergedEvent.Actor.Login).Render(Trim(item.MergedEvent.Actor.Login, 50))
		date = dateStyle.Render(humanize.Time(item.MergedEvent.CreatedAt))
		title = fmt.Sprintf("%smerged the PR %s", author, date)
	}
	return title
}

func (item prTLItem) Description() string {
	var desc string
	switch item.Type {
	case tlItemPRCommit:
		desc = fmt.Sprintf("📧 %s", item.PullRequestCommit.Commit.MessageHeadline)
	case tlItemPRReview:
		reviewState := reviewStyle(item.PullRequestReview.State).Render(item.PullRequestReview.State)
		var comment string
		if item.PullRequestReview.Body != "" {
			comment = fmt.Sprintf(" with comment: %s", Trim(strings.Split(item.PullRequestReview.Body, "\r")[0], 80))
		}
		desc = fmt.Sprintf("🔎 %s%s", reviewState, comment)
	case tlItemMergedEvent:
		desc = fmt.Sprintf("🚀 message: %s", item.MergedEvent.MergeCommit.MessageHeadline)
	}
	return desc
}

func (item prTLItem) FilterValue() string {
	return item.Type
}

func (cmt prReviewComment) render() string {
	var s string
	s += filePathStyle.Render("file: " + cmt.Path)
	s += "\n\n"
	if cmt.Outdated {
		s += outdatedStyle.Render("outdated")
		s += "\n\n"
	}
	s += reviewCmtBodyStyle.Render(cmt.Body)
	s += "\n\n"
	s += diffStyle.Render(cmt.DiffHunk)
	return s
}
