package ui

import (
	"fmt"
	"time"
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

	commitHashLen = 7
)

type terminalDetails struct {
	width int
}

type SourceConfig struct {
	DiffPager *string `yaml:"diff-pager"`
	PRCount   int     `yaml:"pr-count"`
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
	ReviewDecision *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ClosedAt       string
	Author         struct {
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
