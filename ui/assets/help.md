# prs Reference Manual

(scroll line by line with j/k/arrow keys or by half a page with c-d/c-u)

## Views

prs has 6 views:

- PR List View
- PR Details View
- PR Timeline List View
- PR Timeline Item Detail View
- Repo List View (only applicable when --mode=repos)
- Help View (this one)

## Keyboard Shortcuts

### General

```text
  q/esc/ctrl+c                      go back
  Q                                 quit from anywhere
  ?                                 Open Help View
  d                                 Open PR Details View
  ctrl+v                            Show PR details using gh
```

### PR List View

```text
  Indicators for current review decision:

  ±  implies                        CHANGES_REQUESTED
  🟡 implies                        REVIEW_REQUIRED
  ✅ implies                        APPROVED

  ⏎/tab/shift+tab/2                 Switch focus to PR Timeline View
  ctrl+s                            Switch focus to Repo List View (when --mode=repos)
  ctrl+d                            Show PR diff
  ctrl+r                            Reload PR list
  ctrl+b                            Open PR in browser
```

### PR Details View

```text
  h/N/←                             Go to previous section
  l/n/→                             Go to next section
  1/2/3...                          Go to specific section
  J/]                               Go to next PR
  K/[                               Go to previous PR
  d                                 Go back to last view
  ctrl+b                            Open PR in browser
```

### Timeline List View


```text
  tab/shift+tab/1                   Switch focus to PR List View
  ⏎/3                               Show details for PR timeline item (when applicable)
  ctrl+d                            Show PR diff
  ctrl+b                            Open timeline item in browser
  ctrl+r                            Reload PR timeline
```

### Timeline Item Detail View


```text
  1                                 Switch focus to PR List View
  2                                 Switch focus to PR Timeline List View
  ctrl+d                            Show PR diff
  ctrl+b                            Open timeline item in browser
  h/N/←                             Go to previous section
  l/n/→                             Go to next section
```
