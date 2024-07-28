package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"flag"

	"github.com/dhth/prs/ui"
)

var (
	modeFlag  = flag.String("mode", "repos", "mode to run prs in; values: repos, query, reviewer, author")
	queryFlag = flag.String("query", "", "query to filter PRs by")
)

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func Execute() {
	currentUser, err := user.Current()
	var defaultConfigFilePath string
	if err == nil {
		defaultConfigFilePath = fmt.Sprintf("%s/.config/prs/prs.yml", currentUser.HomeDir)
	}
	configFilePath := flag.String("config-file", defaultConfigFilePath, "path of the config file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\nFlags:\n", helpText)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *configFilePath == "" {
		die("config-file cannot be empty")
	}

	configFilePathExp := expandTilde(*configFilePath)

	_, err = os.Stat(configFilePathExp)
	if os.IsNotExist(err) {
		die(fmt.Sprintf("Error: file doesn't exist at %q", configFilePathExp))
	}

	config, err := readConfig(configFilePathExp)
	if err != nil {
		die(fmt.Sprintf("Error reading config: %s", err.Error()))
	}

	if *queryFlag != "" {
		config.Query = queryFlag
	}

	if config.Query != nil {
		if strings.Contains(*config.Query, "type:issue") || strings.Contains(*config.Query, "type: issue") {
			die("type:issue cannot be used in the query")
		}

		if !strings.Contains(*config.Query, "type:pr") && !strings.Contains(*config.Query, "type: pr") {
			updatedQuery := fmt.Sprintf("type: pr %s", *config.Query)
			config.Query = &updatedQuery
		}
	}

	var mode ui.Mode
	switch *modeFlag {
	case "repos":
		mode = ui.RepoMode
	case "query":
		mode = ui.QueryMode
	case "reviewer":
		mode = ui.ReviewerMode
	case "author":
		mode = ui.AuthorMode
	default:
		die("unknown mode provided; possible values: repos, reviewer, author")
	}

	if mode == ui.RepoMode && len(config.Repos) == 0 {
		die("Error: no repos found in config file")
	}

	if mode == ui.QueryMode && config.Query == nil {
		die("Error: no query provided")
	}

	ui.RenderUI(config, mode)
}
