package ui

type hideHelpMsg struct{}

type repoChosenMsg struct {
	repo string
}

type prsFetchedMsg struct {
	prs []pr
	err error
}

type prMetadataFetchedMsg struct {
	repoOwner string
	repoName  string
	prNumber  int
	metadata  prMetadata
	err       error
}

type reviewPRsFetchedMsg prsFetchedMsg

type authoredPRsFetchedMsg prsFetchedMsg

type viewerLoginFetched struct {
	login string
	err   error
}

type prTLFetchedMsg struct {
	repoOwner string
	repoName  string
	prNumber  int
	prTLItems []prTLItem
	setItems  bool
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
