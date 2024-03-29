package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
)

const (
	PRStateOpen            = "OPEN"
	PRStateMerged          = "MERGED"
	PRStateClosed          = "CLOSED"
	TLItemPRCommit         = "PullRequestCommit"
	TLItemPRReview         = "PullRequestReview"
	TLItemMergedEvent      = "MergedEvent"
	ReviewPending          = "PENDING"
	ReviewCommented        = "COMMENTED"
	ReviewApproved         = "APPROVED"
	ReviewChangesRequested = "CHANGES_REQUESTED"
	ReviewDismissed        = "DISMISSED"
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

type delegateKeyMap struct {
	choose key.Binding
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
	CreatedAt string
	ClosedAt  string
	Author    struct {
		Login string
	}
	Url       string
	Additions int
	Deletions int
	Comments  struct {
		TotalCount int
	}
	Reviews struct {
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

type prTLItem struct {
	Type              string `graphql:"type: __typename"`
	PullRequestCommit struct {
		Commit struct {
			Oid             string
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
		CreatedAt string
		State     string
		Body      string
		Author    struct {
			Login string
		}
	} `graphql:"... on PullRequestReview"`
	MergedEvent struct {
		CreatedAt   string
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

	author := authorStyle(pr.Author.Login).Render(RightPadTrim(pr.Author.Login, 80))
	state := prStyle(pr.State).Render(pr.State)

	if pr.Additions > 0 {
		additions = additionsStyle.Render(fmt.Sprintf("+%d", pr.Additions))
	}
	if pr.Deletions > 0 {
		deletions = deletionsStyle.Render(fmt.Sprintf("-%d", pr.Deletions))
	}
	return fmt.Sprintf("%s%s%s%s", author, state, additions, deletions)
}

func (pr pr) FilterValue() string {
	return fmt.Sprintf("%d", pr.Number)
}

func (item prTLItem) Title() string {
	var title string
	switch item.Type {
	case TLItemPRCommit:
		if item.PullRequestCommit.Commit.Author.User != nil {
			author := authorStyle(item.PullRequestCommit.Commit.Author.User.Login).Render(Trim(item.PullRequestCommit.Commit.Author.User.Login, 50))
			title = fmt.Sprintf("%s pushed a commit", author)
		} else {
			title = fmt.Sprintf("%s pushed a commit", item.PullRequestCommit.Commit.Author.Name)
		}
	case TLItemPRReview:
		author := authorStyle(item.PullRequestReview.Author.Login).Render(Trim(item.PullRequestReview.Author.Login, 50))
		title = fmt.Sprintf("%sreviewed", author)
	case TLItemMergedEvent:
		author := authorStyle(item.MergedEvent.Actor.Login).Render(Trim(item.MergedEvent.Actor.Login, 50))
		title = fmt.Sprintf("%smerged the PR", author)
	}
	return title
}

func (item prTLItem) Description() string {
	var desc string
	switch item.Type {
	case TLItemPRCommit:
		desc = fmt.Sprintf("📧 %s", item.PullRequestCommit.Commit.MessageHeadline)
	case TLItemPRReview:
		reviewState := reviewStyle(item.PullRequestReview.State).Render(item.PullRequestReview.State)
		var comment string
		if item.PullRequestReview.Body != "" {
			comment = fmt.Sprintf(" with comment: %s", item.PullRequestReview.Body)
		}
		desc = fmt.Sprintf("🔎 %s%s", reviewState, comment)
	case TLItemMergedEvent:
		desc = fmt.Sprintf("🚀 message: %s", item.MergedEvent.MergeCommit.MessageHeadline)
	}
	return desc
}

func (item prTLItem) FilterValue() string {
	return item.Type
}
