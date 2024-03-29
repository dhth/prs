package cmd

import (
	"os"
	"os/user"
	"strings"

	"github.com/dhth/prs/ui"
	"gopkg.in/yaml.v3"
)

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			os.Exit(1)
		}
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}
	return path
}

func readConfig(configFilePath string) (ui.Config, error) {
	localFile, err := os.ReadFile(configFilePath)
	if err != nil {
		os.Exit(1)
	}
	srcCfg := ui.SourceConfig{}
	err = yaml.Unmarshal(localFile, &srcCfg)
	if err != nil {
		return ui.Config{}, err
	}

	var repos []ui.Repo
	for _, source := range srcCfg.Sources {
		for _, repo := range source.Repos {
			repos = append(repos, ui.Repo{
				Owner: source.Owner,
				Name:  repo.Name,
			})
		}
	}
	var prCount = srcCfg.PRCount
	cfg := ui.Config{
		DiffPager: srcCfg.DiffPager,
		PRCount:   prCount,
		Repos:     repos,
	}
	return cfg, nil

}
