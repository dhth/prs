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

type PRTLFetchedMsg struct {
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
