package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

const (
	prStateOpen                 = "OPEN"
	prStateMerged               = "MERGED"
	prStateClosed               = "CLOSED"
	prRevDecChangesReq          = "CHANGES_REQUESTED"
	prRevDecApproved            = "APPROVED"
	prRevDecRevReq              = "REVIEW_REQUIRED"
	tlItemPRCommit              = "PullRequestCommit"
	tlItemPRReadyForReview      = "ReadyForReviewEvent"
	tlItemPRReviewRequested     = "ReviewRequestedEvent"
	tlItemPRReview              = "PullRequestReview"
	tlItemMergedEvent           = "MergedEvent"
	tlItemHeadRefForcePushed    = "HeadRefForcePushedEvent"
	reviewPending               = "PENDING"
	reviewCommented             = "COMMENTED"
	reviewApproved              = "APPROVED"
	reviewChangesRequested      = "CHANGES_REQUESTED"
	reviewDismissed             = "DISMISSED"
	checkStatusStateCompleted   = "COMPLETED"
	checkRunType                = "CheckRun"
	statusContextType           = "StatusContext"
	checkConclusionStateSuccess = "SUCCESS"
	checkConclusionStateFailure = "FAILURE"
	checkConclusionStateError   = "ERROR"
	statusStateSuccess          = "SUCCESS"
	statusStateFailure          = "FAILURE"
	statusStateError            = "ERROR"
	requestedReviewerUser       = "User"
	prDetailsMetadataKeyPadding = 30
	checkNamePadding            = 40
	statusConclusionPadding     = 16
	reviewRequestsCount         = 20
	latestReviewsCount          = 30
	filesCount                  = 50
	labelsCount                 = 10
	assigneesCount              = 10
	issuesCount                 = 10
	participantsCount           = 30
	commentsCount               = 10
	commitsCount                = 30
	statusCheckContextsCount    = 50
	timeFormat                  = "2006/01/02 15:04"
	mergeableConflicting        = "CONFLICTING"
	noChecksHeader              = "## No Checks"
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
	PRCount int
	Repos   []Repo
	Query   *string
}

type prResult struct {
	pr          *pr
	title       string
	description string
	identifier  string
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
	LastEditedAt   *time.Time
	Author         struct {
		Login string
	}
	URL       string
	Additions int
	Deletions int
	Reviews   struct {
		TotalCount int
	}
}

type prDetails struct {
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
	LastEditedAt   *time.Time
	Author         struct {
		Login string
	}
	Additions      int
	Deletions      int
	ReviewRequests *struct {
		Nodes []struct {
			RequestedReviewer *struct {
				Type string `graphql:"type: __typename"`
				User struct {
					Login string
				} `graphql:"... on User "`
			}
		}
	} `graphql:"reviewRequests (first:$reviewRequestsCount)"`
	LatestReviews struct {
		Nodes []struct {
			Author struct {
				Login string
			}
			State string
		}
	} `graphql:"latestReviews (last: $latestReviewsCount)"`
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
	Assignees struct {
		Nodes []struct {
			Login string
		}
	} `graphql:"assignees (first: $assigneesCount)"`
	IssueReferences struct {
		Nodes []struct {
			Number int
			Title  string
			URL    string
		}
	} `graphql:"closingIssuesReferences (first: $issuesCount)"`
	Participants struct {
		Nodes []struct {
			Login string
		}
	} `graphql:"participants (first: $participantsCount)"`
	Comments struct {
		TotalCount int
		Nodes      []struct {
			Body      string
			UpdatedAt time.Time
			Author    struct {
				Login string
			}
		}
	} `graphql:"comments (first: $commentsCount)"`
	Commits struct {
		TotalCount int
		Nodes      []struct {
			Commit struct {
				AbbreviatedOid  string
				MessageHeadline string
				AuthoredDate    time.Time
				Author          struct {
					Name string
				}
			}
		}
	} `graphql:"commits (last: $commitsCount)"`
	MergedBy *struct {
		Login string
	}
	Milestone *struct {
		Title string
	}
	LastCommit struct {
		Nodes []struct {
			Commit struct {
				AbbreviatedOid    string
				StatusCheckRollup *struct {
					Contexts struct {
						Nodes []struct {
							Type     string `graphql:"type: __typename"`
							CheckRun struct {
								Status     string
								Conclusion *string
								Name       string
							} `graphql:"... on CheckRun"`
							StatusContext struct {
								State   string
								Context string
							} `graphql:"... on StatusContext"`
						}
					} `graphql:"contexts (first: $statusCheckContextsCount) "`
					State string
				}
			}
		}
	} `graphql:"lastCommit: commits(last: 1)"`
}

type PRDetailSection uint

const (
	PRMetadata PRDetailSection = iota
	PRDescription
	PRChecks
	PRReferences
	PRFilesChanged
	PRCommits
	PRComments
)

var PRDetailsSectionList = []PRDetailSection{
	PRMetadata,
	PRDescription,
	PRChecks,
	PRReferences,
	PRFilesChanged,
	PRCommits,
	PRComments,
}

type prReviewComment struct {
	CreatedAt time.Time
	Body      string
	Outdated  bool
	DiffHunk  string
	Path      string
	URL       string
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

type prDetailsQuery struct {
	RepositoryOwner struct {
		Repository struct {
			PullRequest prDetails `graphql:"pullRequest(number: $pullRequestNumber)"`
		} `graphql:"repository(name: $repositoryName)"`
	} `graphql:"repositoryOwner(login: $repositoryOwner)"`
}

type prTLItem struct {
	Type              string `graphql:"type: __typename"`
	PullRequestCommit struct {
		URL    string
		Commit struct {
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
			AbbreviatedOid string
		}
		AfterCommit struct {
			AbbreviatedOid  string
			URL             string
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
		URL       string
		CreatedAt time.Time
		State     string
		Body      string
		Comments  struct {
			TotalCount int
			Nodes      []prReviewComment
		} `graphql:"comments(first: 100)"`
		Author struct {
			Login string
		}
	} `graphql:"... on PullRequestReview"`
	MergedEvent struct {
		CreatedAt   time.Time
		URL         string
		MergeCommit struct {
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

func (pr prDetails) Metadata() string {
	var metadata []string

	metadata = append(metadata, fmt.Sprintf("- %s *%s*",
		RightPadTrim("State", prDetailsMetadataKeyPadding),
		pr.State,
	))

	metadata = append(metadata, fmt.Sprintf("- %s `@%s`",
		RightPadTrim("Author", prDetailsMetadataKeyPadding),
		pr.Author.Login,
	))

	if len(pr.Assignees.Nodes) > 0 {
		assignees := make([]string, len(pr.Assignees.Nodes))
		for i, l := range pr.Assignees.Nodes {
			assignees[i] = fmt.Sprintf("`@%s`", l.Login)
		}
		metadata = append(metadata, fmt.Sprintf("- %s %s",
			RightPadTrim("Assignees", prDetailsMetadataKeyPadding),
			strings.Join(assignees, ", "),
		))
	}

	if len(pr.Participants.Nodes) > 0 {
		participants := make([]string, len(pr.Participants.Nodes))
		for i, l := range pr.Participants.Nodes {
			participants[i] = fmt.Sprintf("`@%s`", l.Login)
		}
		metadata = append(metadata, fmt.Sprintf("- %s %s",
			RightPadTrim("Participants", prDetailsMetadataKeyPadding),
			strings.Join(participants, ", "),
		))
	}

	if pr.ReviewRequests != nil && len(pr.ReviewRequests.Nodes) > 0 {
		var requested []string
		for _, r := range pr.ReviewRequests.Nodes {
			if r.RequestedReviewer.Type != requestedReviewerUser {
				continue
			}
			requested = append(requested, fmt.Sprintf("`@%s`", r.RequestedReviewer.User.Login))
		}

		if len(requested) > 0 {
			metadata = append(metadata, fmt.Sprintf("- %s %s",
				RightPadTrim("Review requested from", prDetailsMetadataKeyPadding),
				strings.Join(requested, ", "),
			))
		}
	}

	metadata = append(metadata, fmt.Sprintf("- %s %s (%s)",
		RightPadTrim("Created at", prDetailsMetadataKeyPadding),
		pr.CreatedAt.Format(timeFormat),
		humanize.Time(pr.CreatedAt),
	))
	if pr.LastEditedAt != nil && *pr.LastEditedAt != pr.CreatedAt {
		metadata = append(metadata, fmt.Sprintf("- %s %s (%s)",
			RightPadTrim("Last edited at", prDetailsMetadataKeyPadding),
			pr.LastEditedAt.Format(timeFormat),
			humanize.Time(*pr.LastEditedAt),
		))
	}

	switch pr.State {
	case prStateClosed:
		if pr.ClosedAt != nil {
			metadata = append(metadata, fmt.Sprintf("- %s %s (%s)",
				RightPadTrim("Closed at", prDetailsMetadataKeyPadding),
				pr.ClosedAt.Format(timeFormat),
				humanize.Time(*pr.ClosedAt),
			))
		}
	case prStateMerged:
		metadata = append(metadata, fmt.Sprintf("- %s %s (%s) by `@%s`",
			RightPadTrim("Merged at", prDetailsMetadataKeyPadding),
			pr.MergedAt.Format(timeFormat),
			humanize.Time(*pr.MergedAt),
			pr.MergedBy.Login,
		))
	}

	if len(pr.Labels.Nodes) > 0 {
		labels := make([]string, len(pr.Labels.Nodes))
		for i, l := range pr.Labels.Nodes {
			labels[i] = fmt.Sprintf("*%s*", l.Name)
		}
		metadata = append(metadata, fmt.Sprintf("- %s %s",
			RightPadTrim("Labels", prDetailsMetadataKeyPadding),
			strings.Join(labels, " "),
		))
	}

	if pr.Commits.TotalCount > 0 {
		metadata = append(metadata, fmt.Sprintf("- %s %d",
			RightPadTrim("Commits", prDetailsMetadataKeyPadding),
			pr.Commits.TotalCount,
		))
	}

	if pr.Comments.TotalCount > 0 {
		metadata = append(metadata, fmt.Sprintf("- %s %d",
			RightPadTrim("Comments", prDetailsMetadataKeyPadding),
			pr.Comments.TotalCount,
		))
	}

	if pr.IsDraft {
		metadata = append(metadata, fmt.Sprintf("- %s `true`",
			RightPadTrim("Is draft",
				prDetailsMetadataKeyPadding),
		))
	}

	if pr.Mergeable == mergeableConflicting {
		metadata = append(metadata, fmt.Sprintf("- %s `true`", RightPadTrim("Has conflicts",
			prDetailsMetadataKeyPadding),
		))
	}

	if pr.Milestone != nil {
		metadata = append(metadata, fmt.Sprintf("- %s %s", RightPadTrim("Milestone",
			prDetailsMetadataKeyPadding),
			pr.Milestone.Title,
		))
	}

	if len(pr.LatestReviews.Nodes) > 0 {
		reviews := make([]string, len(pr.LatestReviews.Nodes))

		for i, r := range pr.LatestReviews.Nodes {
			var state string
			switch r.State {
			case reviewPending:
				state = "🟡"
			case reviewCommented:
				state = "💬"
			case reviewChangesRequested:
				state = "🔄"
			case reviewApproved:
				state = "✅"
			case reviewDismissed:
				state = "❌"
			}
			reviews[i] = fmt.Sprintf("`@%s` %s", r.Author.Login, state)
		}

		metadata = append(metadata, "\n---\n")

		metadata = append(metadata, fmt.Sprintf("- %s %s",
			RightPadTrim("Reviewed by", prDetailsMetadataKeyPadding),
			strings.Join(reviews, ", "),
		))
	}

	return fmt.Sprintf(`
## Metadata

%s`, strings.Join(metadata, "\n"))
}

func (pr prDetails) Description() string {
	return fmt.Sprintf(`
## Description

%s`, pr.Body)
}

func (pr prDetails) Checks() string {
	if len(pr.LastCommit.Nodes) == 0 {
		return noChecksHeader
	}
	if pr.LastCommit.Nodes[0].Commit.StatusCheckRollup == nil {
		return noChecksHeader
	}
	if len(pr.LastCommit.Nodes[0].Commit.StatusCheckRollup.Contexts.Nodes) == 0 {
		return noChecksHeader
	}

	var checks []string
	for _, n := range pr.LastCommit.Nodes[0].Commit.StatusCheckRollup.Contexts.Nodes {
		switch n.Type {
		case checkRunType:
			checkName := RightPadTrim(n.CheckRun.Name, checkNamePadding)
			if n.CheckRun.Conclusion != nil {
				var conclusionMarker string
				switch *n.CheckRun.Conclusion {
				case checkConclusionStateSuccess:
					conclusionMarker = " ✅"
				case checkConclusionStateFailure, checkConclusionStateError:
					conclusionMarker = " ❌"
				}
				checks = append(checks, fmt.Sprintf("- %s %s%s",
					checkName,
					RightPadTrim(fmt.Sprintf("`%s`", *n.CheckRun.Conclusion), statusConclusionPadding),
					conclusionMarker,
				))
			} else {
				checks = append(checks, fmt.Sprintf("- %s %s", checkName, n.CheckRun.Status))
			}
		case statusContextType:
			var stateMarker string
			switch n.StatusContext.State {
			case statusStateSuccess:
				stateMarker = " ✅"
			case statusStateFailure, statusStateError:
				stateMarker = " ❌"
			}
			checks = append(checks, fmt.Sprintf("- %s %s%s",
				RightPadTrim(n.StatusContext.Context, checkNamePadding),
				RightPadTrim(fmt.Sprintf("`%s`", n.StatusContext.State), statusConclusionPadding),
				stateMarker,
			))
		}
	}

	if len(checks) == 0 {
		return noChecksHeader
	}

	return fmt.Sprintf(`
## Checks

%s **%s**

%s`,
		RightPadTrim("> Status of latest commit", checkNamePadding+2),
		pr.LastCommit.Nodes[0].Commit.StatusCheckRollup.State,
		strings.Join(checks, "\n"))
}

func (pr prDetails) References() string {
	issues := make([]string, len(pr.IssueReferences.Nodes))
	for i, iss := range pr.IssueReferences.Nodes {
		issues[i] = fmt.Sprintf("- `#%d`: %s (%s)", iss.Number, iss.Title, iss.URL)
	}
	return fmt.Sprintf(`
## Referenced by

%s`, strings.Join(issues, "\n"))
}

func (pr prDetails) FilesChanged() string {
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
	return fmt.Sprintf(`
## Files changed

%s`, strings.Join(fc, "\n"))
}

func (pr prDetails) CommitsList() string {
	var commitsStr string

	commits := make([]string, len(pr.Commits.Nodes))
	for i, c := range pr.Commits.Nodes {
		hash := c.Commit.AbbreviatedOid

		commits[i] = fmt.Sprintf("- `%s`: %s **(%s)** `<%s>`",
			hash,
			c.Commit.MessageHeadline,
			humanize.Time(c.Commit.AuthoredDate),
			c.Commit.Author.Name,
		)
	}

	var commitsNumStr string
	if len(pr.Commits.Nodes) < pr.Commits.TotalCount {
		commitsNumStr = fmt.Sprintf(" (last %d out of %d)", len(pr.Commits.Nodes), pr.Commits.TotalCount)
	}

	commitsStr = fmt.Sprintf(`
## Commits%s

%s
`, commitsNumStr, strings.Join(commits, "\n"))

	return commitsStr
}

func (pr prDetails) CommentsList() string {
	comments := make([]string, len(pr.Comments.Nodes))
	for i, c := range pr.Comments.Nodes {
		comments[i] = fmt.Sprintf("`@%s` (%s):\n\n%s", c.Author.Login, humanize.Time(c.UpdatedAt), c.Body)
	}

	var commentsNumStr string
	if len(pr.Comments.Nodes) < pr.Comments.TotalCount {
		commentsNumStr = fmt.Sprintf(" (first %d out of %d)", len(pr.Comments.Nodes), pr.Comments.TotalCount)
	}

	return fmt.Sprintf(`
## Comments%s

%s
`, commentsNumStr, strings.Join(comments, "\n\n▬▬▬▬▬▬\n\n"))
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

func (prRes prResult) Title() string {
	return prRes.title
}

func (prRes prResult) Description() string {
	return prRes.description
}

func (prRes prResult) FilterValue() string {
	return fmt.Sprintf("%d", prRes.pr.Number)
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
