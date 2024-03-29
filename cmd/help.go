package cmd

import "fmt"

var (
	configSampleFormat = `
pr-count: 20
sources:
- owner: owner
  repos:
    - name: repo-1
    - name: repo-2
    - name: repo-3
    - name: repo-4
`
	helpText = `prs lets you stay updated on the PRs you care about without leaving the terminal.

Keyboard shortcuts:
- ctrl+r: change repo
- ctrl+b: open PR in browser
- ctrl+b: open PR in browser
- ctrl+d: open diff using gh
- ctrl+v: view PR using gh

Usage: prs [flags]
`
)

func cfgErrSuggestion(msg string) string {
	return fmt.Sprintf(`%s

Make sure to structure the yml config file as follows:
%s
Use "prs -help" for more information`,
		msg,
		configSampleFormat,
	)
}
