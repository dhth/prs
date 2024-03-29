package cmd

import (
	"fmt"
	"os"
	"os/user"

	"flag"

	"github.com/dhth/prs/ui"
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
		die(cfgErrSuggestion(fmt.Sprintf("Error: file doesn't exist at %q", configFilePathExp)))
	}

	config, err := readConfig(configFilePathExp)
	if err != nil {
		die(cfgErrSuggestion(fmt.Sprintf("Error reading config: %s", err.Error())))
	}

	if len(config.Repos) == 0 {
		die(cfgErrSuggestion("Error: no repos found in config file"))
	}

	ui.RenderUI(config)
}
