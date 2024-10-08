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
	projectHomePage    = "https://github.com/dhth/prs"
	issuesURL          = "https://github.com/dhth/prs/issues"
	configFileName     = "prs/prs.yml"
	defaultSearchQuery = "type:pr author:@me sort:updated-desc state:open"
	defaultPRNum       = 20
	maxPRNum           = 50
)

var (
	errCouldntGetHomeDir        = errors.New("couldn't get home directory")
	errCouldntGetConfigDir      = errors.New("couldn't get config directory")
	errModeIncorrect            = errors.New("incorrect mode provided")
	errConfigFileDoesntExist    = errors.New("config file does not exist")
	errNoReposProvided          = errors.New("no repos were provided")
	errIncorrectRepoProvided    = errors.New("incorrect repo provided")
	errCouldntSetupGithubClient = errors.New("couldn't set up a Github Client")
)

var reportIssueMsg = fmt.Sprintf("Let %s know about this error via %s.", author, issuesURL)

func Execute(version string) error {
	rootCmd, err := NewRootCommand(version)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		switch {
		case errors.Is(err, errCouldntGetHomeDir), errors.Is(err, errCouldntGetConfigDir):
			fmt.Printf(`
This is a fatal error; use --config-path to specify config file path manually.
%s
`, reportIssueMsg)
		}
		return err
	}

	err = rootCmd.Execute()

	switch {
	case errors.Is(err, errCouldntSetupGithubClient):
		fmt.Printf(`
If the error is due to misconfigured authentication, you can fix that by either of the following:
- Provide a valid Github token via $GH_TOKEN
- Have an authenticated instance of gh (https://github.com/cli/cli) available
`)
	}
	return err
}

func NewRootCommand(version string) (*cobra.Command, error) {
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
	)

	rootCmd := &cobra.Command{
		Use:   "prs",
		Short: "prs lets you stay updated on pull requests from your terminal",
		Long: fmt.Sprintf(`prs lets you stay updated on pull requests from your terminal.

Use it to query for specific pull requests based on a filter query (using
Github's search syntax), or have it let you pick a repository from a predefined
list.

Examples:
$ prs --query='type:pr repo:neovim/neovim state:open label:lua linked:issue'
$ prs -q 'type:pr author:@me state:open'
$ PRS_QUERY='type:pr user-review-requested:@me state:open' prs
$ prs # will read query from config file

$ prs --mode=repos --repos='dhth/prs,dhth/omm,dhth/hours'
$ PRS_REPOS='dhth/prs,dhth/omm,dhth/hours' prs --mode=repos
$ prs -m repos # will read repos from config file

Project home page: %s
`, projectHomePage),

		Args:         cobra.MaximumNArgs(0),
		SilenceUsage: true,
		Version:      version,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
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
			default:
				return errModeIncorrect
			}

			if mode == ui.RepoMode {
				var reposToUse []string
				// pretty ugly hack to get around the fact that
				// v.GetStringSlice("repos") always seems to prioritize the config file
				if len(repoStrs) > 0 && len(repoStrs[0]) > 0 && !strings.HasPrefix(repoStrs[0], "[") {
					reposToUse = repoStrs
				} else {
					reposToUse = v.GetStringSlice("repos")
				}

				if len(reposToUse) == 0 {
					return errNoReposProvided
				}

				for _, r := range reposToUse {
					repoEls := strings.Split(r, "/")
					// TODO: there can be more validations done here, maybe regex based
					if len(repoEls) != 2 {
						return fmt.Errorf("%w: %s", errIncorrectRepoProvided, r)
					}

					repos = append(repos, ui.Repo{
						Owner: strings.TrimSpace(repoEls[0]),
						Name:  strings.TrimSpace(repoEls[1]),
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
				return fmt.Errorf("%w: %s", errCouldntSetupGithubClient, err.Error())
			}
			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ui.Config{
				PRCount: prNum,
				Repos:   repos,
				Query:   &searchQuery,
			}
			return ui.RenderUI(ghClient, config, mode)
		},
	}
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s" .Version}}
`)

	ros := runtime.GOOS
	userCfgDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errCouldntGetConfigDir, err.Error())
	}

	var defaultConfigFilePath string
	switch ros {
	case "linux", "windows":
		defaultConfigFilePath = filepath.Join(userCfgDir, configFileName)
	default:
		// to use ~/.config instead of $HOME/Library/Application Support
		hd, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("%w: %s", errCouldntGetHomeDir, err.Error())
		}
		defaultConfigFilePath = filepath.Join(hd, ".config", configFileName)
	}

	rootCmd.Flags().StringVarP(&configFilePath, "config-path", "c", defaultConfigFilePath, "location of prs's config file")
	rootCmd.Flags().StringVarP(&modeInp, "mode", "m", "query", "mode to run prs in; values: query, repos")
	rootCmd.Flags().StringVarP(&searchQuery, "query", "q", defaultSearchQuery, "query to search PRs for")
	rootCmd.Flags().IntVarP(&prNum, "num", "n", defaultPRNum, "number of PRs to fetch")
	rootCmd.Flags().StringSliceVarP(&repoStrs, "repos", "r", nil, "comma separated list of repos to use for repo mode")

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd, nil
}

func initializeConfig(cmd *cobra.Command, configFile string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filepath.Base(configFile))
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Dir(configFile))

	err := v.ReadInConfig()
	if err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return v, err
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
