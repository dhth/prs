package ui

import "fmt"

var (
	helpText = fmt.Sprintf(`
  %s
%s
  %s

  %s
%s
  %s
%s
  %s
%s
  %s
%s
`,
		helpHeaderStyle.Render("prs Reference Manual"),
		helpSectionStyle.Render(`
  (scroll line by line with j/k/arrow keys or by half a page with <c-d>/<c-u>)

  prs has 5 views:
  - PR List View
  - PR Timeline List View
  - PR Review Comments View
  - Repo List View (only applicable when -mode=repos)
  - Help View (this one)
`),
		helpHeaderStyle.Render("Keyboard Shortcuts"),
		helpHeaderStyle.Render("General"),
		helpSectionStyle.Render(`
      <tab>                               Switch focus between PR List and PR Timeline Pane
      1                                   Switch focus to PR List View
      2                                   Switch focus to PR Timeline List View
      3                                   Switch focus to PR Review Comments View
      <ctrl+s>                            Switch focus to Repo List View
      ?                                   Switch focus to Help View
`),
		helpHeaderStyle.Render("PR List/Timeline List View"),
		helpSectionStyle.Render(`
      <ctrl+v>                            Show PR details
      <ctrl+d>                            Show PR diff
`),
		helpHeaderStyle.Render("PR List View"),
		helpSectionStyle.Render(`
      Indicators for current review
      decision:

      ±  implies                          CHANGES_REQUESTED                  
      🟡 implies                          REVIEW_REQUIRED                  
      ✅ implies                          APPROVED                  

      <ctrl+b>                            Open PR in the browser
      <ctrl+r>                            Reload PR list
      <enter>                             Switch focus to PR Timeline View for currently selected PR
      <enter>                             Show commit/revision range
`),
		helpHeaderStyle.Render("PR Timeline View"),
		helpSectionStyle.Render(`
      <ctrl+b>                            Open timeline item in browser
      <ctrl+r>                            Reload timeline list
      <enter>                             Switch focus to Review Comments View for currently selected item
`),
	)
)
