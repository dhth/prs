package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

const (
	prStateOpen              = "OPEN"
	prStateMerged            = "MERGED"
	prStateClosed            = "CLOSED"
	prRevDecChangesReq       = "CHANGES_REQUESTED"
	prRevDecApproved         = "APPROVED"
	prRevDecRevReq           = "REVIEW_REQUIRED"
	tlItemPRCommit           = "PullRequestCommit"
	tlItemPRReadyForReview   = "ReadyForReviewEvent"
	tlItemPRReviewRequested  = "ReviewRequestedEvent"
	tlItemPRReview           = "PullRequestReview"
	tlItemMergedEvent        = "MergedEvent"
	tlItemHeadRefForcePushed = "HeadRefForcePushedEvent"
	reviewPending            = "PENDING"
	reviewCommented          = "COMMENTED"
	reviewApproved           = "APPROVED"
	reviewChangesRequested   = "CHANGES_REQUESTED"
	reviewDismissed          = "DISMISSED"

	mergeableConflicting = "CONFLICTING"

	commitHashLen = 7

	timeFormat = "2006/01/02 15:04"
)

type terminalDetails struct {
	width int
}

type SourceConfig struct {
	DiffPager *string `yaml:"diff-pager"`
	PRCount   *int    `yaml:"pr-count"`
	Sources   *[]struct {
		Owner string `yaml:"owner"`
		Repos []struct {
			Name string `yaml:"name"`
		} `yaml:"repos"`
	} `yaml:"sources"`
	Query *string `yaml:"query"`
}

type Repo struct {
	Owner string
	Name  string
}

type Config struct {
	DiffPager *string
	PRCount   int
	Repos     []Repo
	Query     *string
}

type prResult struct {
	pr          *pr
	title       string
	description string
}

type prTLItemResult struct {
	item        *prTLItem
	title       string
	description string
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
	State          string
	Mergeable      string
	IsDraft        bool
	ReviewDecision *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ClosedAt       *time.Time
	MergedAt       *time.Time
	Author         struct {
		Login string
	}
	Url       string
	Additions int
	Deletions int
	Reviews   struct {
		TotalCount int
	}
	Body  string
	Files struct {
		Nodes []struct {
			Path      string
			Additions int
			Deletions int
		}
	} `graphql:"files (first: $filesCount)"`
	Labels struct {
		Nodes []struct {
			Name string
		}
	} `graphql:"labels (first: $labelsCount)"`
	MergedBy *struct {
		Login string
	}
}

type prReviewComment struct {
	CreatedAt time.Time
	Body      string
	Outdated  bool
	DiffHunk  string
	Path      string
	Url       string
}

type userLoginQuery struct {
	Viewer struct {
		Login string
	}
}

type prSearchQuery struct {
	Search struct {
		Edges []struct {
			Node struct {
				Type string `graphql:"type: __typename"`
				pr   `graphql:"... on PullRequest"`
			}
		}
	} `graphql:"search(query: $query, type: ISSUE, first: $count)"`
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
	HeadRefForcePushed struct {
		CreatedAt time.Time
		Actor     struct {
			Login string
		}
		BeforeCommit struct {
			Oid string
		}
		AfterCommit struct {
			Oid             string
			Url             string
			MessageHeadline string
		}
	} `graphql:"... on HeadRefForcePushedEvent"`
	PullRequestReadyForReview struct {
		CreatedAt time.Time
		Actor     struct {
			Login string
		}
	} `graphql:"... on ReadyForReviewEvent"`
	PullRequestReviewRequested struct {
		CreatedAt time.Time
		Actor     struct {
			Login string
		}
		RequestedReviewer struct {
			User struct {
				Login string
			} `graphql:"... on User"`
		}
	} `graphql:"... on ReviewRequestedEvent"`
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
				} `graphql:"timelineItems(last: $timelineItemsCount, itemTypes: [PULL_REQUEST_COMMIT, READY_FOR_REVIEW_EVENT, REVIEW_REQUESTED_EVENT, MERGED_EVENT, PULL_REQUEST_REVIEW, HEAD_REF_FORCE_PUSHED_EVENT])"`
			} `graphql:"pullRequest(number: $pullRequestNumber)"`
		} `graphql:"repository(name: $repositoryName)"`
	} `graphql:"repositoryOwner(login: $repositoryOwner)"`
}

func (pr pr) DetailsMd() string {
	var metadata []string
	metadata = append(metadata, fmt.Sprintf("- Author      :  `@%s`", pr.Author.Login))
	metadata = append(metadata, fmt.Sprintf("- Created at  :  %s (%s)", pr.CreatedAt.Format(timeFormat), humanize.Time(pr.CreatedAt)))

	switch pr.State {
	case prStateClosed:
		if pr.ClosedAt != nil {
			metadata = append(metadata, fmt.Sprintf("- Closed at   :  %s (%s)", pr.ClosedAt.Format(timeFormat), humanize.Time(*pr.ClosedAt)))
		}
	case prStateMerged:
		metadata = append(metadata, fmt.Sprintf("- Merged at   :  %s (%s) by `@%s`", pr.MergedAt.Format(timeFormat), humanize.Time(*pr.MergedAt), pr.MergedBy.Login))
	}

	if len(pr.Labels.Nodes) > 0 {
		labels := make([]string, len(pr.Labels.Nodes))
		for i, l := range pr.Labels.Nodes {
			labels[i] = fmt.Sprintf("*%s*", l.Name)
		}
		metadata = append(metadata, fmt.Sprintf("- Labels      :  %s", strings.Join(labels, " ")))
	}

	if pr.IsDraft {
		metadata = append(metadata, "- Draft")
	}

	if pr.Mergeable == mergeableConflicting {
		metadata = append(metadata, "- Has conflicts")
	}

	var fcStr string
	if len(pr.Files.Nodes) > 0 {
		fc := make([]string, len(pr.Files.Nodes))
		for i, f := range pr.Files.Nodes {
			var additions string
			var deletions string

			if f.Additions > 0 {
				additions = fmt.Sprintf(" `+%d`", f.Additions)
			}

			if f.Deletions > 0 {
				deletions = fmt.Sprintf(" `-%d`", f.Deletions)
			}

			fc[i] = fmt.Sprintf("- %s%s%s", f.Path, additions, deletions)
		}
		fcStr = fmt.Sprintf(`
## Files changed

%s
`, strings.Join(fc, "\n"))
	}

	details := fmt.Sprintf(`# %d: %s

%s

---

## Description

%s

---
%s`, pr.Number, pr.PRTitle, strings.Join(metadata, "\n"), pr.Body, fcStr)

	return details
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

func (pr prResult) Title() string {
	return pr.title
}

func (pr prResult) Description() string {
	return pr.description
}

func (pr prResult) FilterValue() string {
	return fmt.Sprintf("%d", pr.pr.Number)
}

func (ir prTLItemResult) Title() string {
	return ir.title
}
func (ir prTLItemResult) Description() string {
	return ir.description
}

func (ir prTLItemResult) FilterValue() string {
	return ir.title
}
