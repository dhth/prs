# prs Reference Manual

(scroll line by line with j/k/arrow keys or by half a page with c-d/c-u)

## Views

prs has 6 views:

- PR List View
- PR Details View
- PR Timeline List View
- PR Review Comments View
- Repo List View (only applicable when --mode=repos)
- Help View (this one)

## Keyboard Shortcuts

### General

```text
  tab/shift+tab                     Switch focus between PR List and PR Timeline Pane
  1                                 Switch focus to PR List View
  2                                 Switch focus to PR Timeline List View
  3                                 Switch focus to PR Review Comments View
  ctrl+s                            Switch focus to Repo List View
  ?                                 Switch focus to Help View
```

### PR List/Timeline List View


```text
  ctrl+v                            Show PR details
  ctrl+d                            Show PR diff
```

### PR List View

```text
  Indicators for current review decision:

  ±  implies                        CHANGES_REQUESTED
  🟡 implies                        REVIEW_REQUIRED
  ✅ implies                        APPROVED

  d                                 Open PR Details View
  ctrl+b                            Open PR in browser
  ctrl+r                            Reload PR list
  enter                             Switch focus to PR Timeline View for currently selected PR
  enter                             Show commit/revision range
```

### PR Timeline View

```text
  ctrl+b                            Open timeline item in browser
  ctrl+r                            Reload timeline list and PR details
  enter                             Switch focus to Review Comments View for currently selected item
```

### PR Details View

```text
  h/N/←                             Go to previous section
  l/n/→                             Go to next section
  1/2/3...                          Go to specific section
  J/]                               Go to next PR
  K/[                               Go to previous PR
  d                                 Go back to PR List View
  ctrl+b                            Open PR in browser
```
