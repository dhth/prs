package ui

var (
	HelpText = `
  prs Reference Manual (scroll with j/k or arrow keys)

  prs has 5 views:
  - PR List View
  - PR Timeline List View
  - PR Review Comments View
  - Repo List View
  - Help View (this one)

  Keyboard Shortcuts:

  General
      <tab>       Switch focus between PR List and PR Timeline Pane
      1           Switch focus to PR List View
      2           Switch focus to PR Timeline List View
      <ctrl+r>    Switch focus to Repo List View
      ?           Switch focus to Help View

  PR List/Timeline List View
      <ctrl+v>    Show PR details
      <ctrl+d>    Show PR diff

  PR List View
      <ctrl+b>    Open PR in the browser
      <enter>     Switch focus to PR Timeline View for currently selected PR

  PR Timeline View
      <ctrl+b>    Open timeline item in browser
      <enter>     Switch focus to Review Comments View for currently selected item
`
)
