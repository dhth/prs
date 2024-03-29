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

type PROpenedinBrowserMsg struct {
	url string
	err error
}
