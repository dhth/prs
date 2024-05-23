package ui

type HideHelpMsg struct{}

type RepoChosenMsg struct {
	repo string
}

type PRChosenMsg struct {
	prNumber int
	err      error
}

type PRsFetchedMsg struct {
	prs []pr
	err error
}

type ReviewPRsFetchedMsg struct {
	prs []reviewPr
	err error
}

type ViewerLoginFetched struct {
	login string
	err   error
}

type PRTLFetchedMsg struct {
	repoOwner string
	repoName  string
	prNumber  int
	prTLItems []prTLItem
	err       error
}

type URLOpenedinBrowserMsg struct {
	url string
	err error
}

type PRDiffDoneMsg struct {
	err error
}

type PRViewDoneMsg struct {
	err error
}
