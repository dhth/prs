package ui

type hideHelpMsg struct{}

type repoChosenMsg struct {
	repo string
}

type prChosenMsg struct {
	prNumber int
	err      error
}

type prsFetchedMsg struct {
	prs []pr
	err error
}

type reviewPRsFetchedMsg struct {
	prs []pr
	err error
}

type viewerLoginFetched struct {
	login string
	err   error
}

type prTLFetchedMsg struct {
	repoOwner string
	repoName  string
	prNumber  int
	prTLItems []prTLItem
	err       error
}

type urlOpenedinBrowserMsg struct {
	url string
	err error
}

type prDiffDoneMsg struct {
	err error
}

type prViewDoneMsg struct {
	err error
}
