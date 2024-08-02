package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/dhth/prs/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envPrefix          = "PRS"
	author             = "@dhth"
	repoIssuesUrl      = "https://github.com/dhth/prs/issues"
	configFileName     = "prs/prs.yml"
	defaultSearchQuery = "is:pr author:@me sort:updated-desc state:open"
	defaultPRNum       = 20
	maxPRNum           = 50
)

var (
	errModeIncorrect          = errors.New("mode value is incorrect")
	errConfigFileDoesntExist  = errors.New("config file does not exist")
	errQuerySearchesForIssues = errors.New("searching for issues not supported")
	errNoReposProvided        = errors.New("no repos were provided")
)

func Execute(version string) {
	rootCmd, err := NewRootCommand()

	rootCmd.Version = version
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	_ = rootCmd.Execute()
}

func NewRootCommand() (*cobra.Command, error) {

	var (
		configFilePath string
		configPathFull string
		mode           ui.Mode
		modeInp        string
		repoStrs       []string
		repos          []ui.Repo
		searchQuery    string
		ghClient       *ghapi.GraphQLClient
		prNum          int
		diffPager      string
	)

	rootCmd := &cobra.Command{
		Use:          "prs",
		Short:        "prs lets you stay updated on pull requests from your terminal",
		Args:         cobra.MaximumNArgs(0),
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			configPathFull = expandTilde(configFilePath)

			if filepath.Ext(configPathFull) != ".yml" {
				return errConfigFileDoesntExist
			}
			_, err := os.Stat(configPathFull)

			fl := cmd.Flags()
			if fl != nil {
				cf := fl.Lookup("config-path")
				if cf != nil && cf.Changed && errors.Is(err, fs.ErrNotExist) {
					return errConfigFileDoesntExist
				}
			}

			var v *viper.Viper
			v, err = initializeConfig(cmd, configPathFull)
			if err != nil {
				return err
			}

			if prNum > maxPRNum {
				prNum = maxPRNum
			}

			switch modeInp {
			case "repos":
				mode = ui.RepoMode
			case "query":
				mode = ui.QueryMode
			case "reviewer":
				mode = ui.ReviewerMode
			case "author":
				mode = ui.AuthorMode
			default:
				return errModeIncorrect
			}

			if mode == ui.RepoMode {
				reposSl := v.GetStringSlice("repos")
				if len(reposSl) == 0 {
					return errNoReposProvided
				}

				for _, r := range reposSl {
					repoEls := strings.Split(r, "/")
					if len(repoEls) != 2 {
						return fmt.Errorf("Incorrect repo provided: %s", r)
					}

					repos = append(repos, ui.Repo{
						Owner: repoEls[0],
						Name:  repoEls[1],
					})
				}
			}

			opts := ghapi.ClientOptions{
				EnableCache: true,
				CacheTTL:    time.Second * 30,
				Timeout:     8 * time.Second,
			}

			ghClient, err = ghapi.NewGraphQLClient(opts)
			if err != nil {
				return fmt.Errorf(`Couldn't set up a Github client.
Is gh (https://github.com/cli/cli) installed and configured? (prs depends on gh for communicating with Github).

Error: %s`, err.Error())
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			config := ui.Config{
				DiffPager: &diffPager,
				PRCount:   prNum,
				Repos:     repos,
				Query:     &searchQuery,
			}
			return ui.RenderUI(ghClient, config, mode)
		},
	}
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s" .Version}}
`)

	ros := runtime.GOOS
	userCfgDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf(`Couldn't get your default config directory. This is a fatal error;
use --config-path to specify config file path manually.
Let %s know about this via %s.

Error: %s`, author, repoIssuesUrl, err)
	}

	var defaultConfigFilePath string
	switch ros {
	case "linux", "windows":
		defaultConfigFilePath = filepath.Join(userCfgDir, configFileName)
	default:
		// to use ~/.config instead of $HOME/Library/Application Support
		hd, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf(`Couldn't get your home directory. This is a fatal error;
use --config-path to specify config file path manually
Let %s know about this via %s.

Error: %s`, author, repoIssuesUrl, err)
		}
		defaultConfigFilePath = filepath.Join(hd, ".config", configFileName)
	}

	rootCmd.Flags().StringVarP(&configFilePath, "config-path", "c", defaultConfigFilePath, "location of prs's config file")
	rootCmd.Flags().StringVarP(&modeInp, "mode", "m", "query", "mode to run prs in; values: query, repos, reviewer, author")
	rootCmd.Flags().StringVarP(&searchQuery, "query", "q", defaultSearchQuery, "query to search PRs for")
	rootCmd.Flags().IntVarP(&prNum, "num", "n", defaultPRNum, "number of PRs to fetch")
	rootCmd.Flags().StringVar(&diffPager, "diff-pager", "", "pager to use for showing diffs")
	rootCmd.Flags().StringSliceVarP(&repoStrs, "repos", "r", nil, "repos to use for repo mode")

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd, nil
}

func initializeConfig(cmd *cobra.Command, configFile string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filepath.Base(configFile))
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Dir(configFile))

	var err error
	if err = v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return v, err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	err = bindFlags(cmd, v)
	if err != nil {
		return v, err
	}

	return v, nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	var err error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			fErr := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if fErr != nil {
				err = fErr
				return
			}
		}
	})
	return err
}
