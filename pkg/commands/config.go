package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"monobuild/internal/errors"
	"monobuild/pkg/config"
	"monobuild/pkg/git"
	"os"
	"path"
	"strings"
)

type configuration struct {
	applicationPath string
	cfg             *config.Config
	repoPath        string
	diffFiles       []string
	revision        string
	pwd             string
}

func setupConfigurationHook(cmd *cobra.Command) (*configuration, func(cmd *cobra.Command, args []string) error) {
	var entryPath string
	var modulePath string
	var confCommand configuration
	var revision string

	cmd.Flags().StringVar(&entryPath, "entry", "", "use for specific config file (YAML format)")
	cmd.Flags().StringVar(&modulePath, "module", "", "specify the module outside of config")
	cmd.Flags().StringVar(&revision, "revision", "HEAD", "specify git revision")

	return &confCommand, func(cmd *cobra.Command, args []string) error {
		var cfg *config.Config

		pwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot get pwd path: %w", err)
		}

		applicationPath := pwd
		if entryPath != "" {
			applicationPath = path.Join(pwd, entryPath)
			if strings.Index(entryPath, "/") == 0 {
				applicationPath = entryPath
			}

			_cfg, err := config.ParseDefaultConfig(applicationPath)
			if err != nil {
				return errors.NewRichError("Incorrect format of configuration file", err)
			}

			cfg = _cfg
		} else if modulePath != "" {
			fullModulePath := path.Join(pwd, modulePath)
			applicationPath = fullModulePath

			cfg = &config.Config{
				Packages: []config.Package{
					{FullEntry: fullModulePath},
				},
			}
		}

		repoPath, err := git.ResolveGitPath(applicationPath)
		if err != nil {
			return err
		}

		diff, err := git.GetDiffFiles(repoPath, revision)
		if err != nil {
			return err
		}

		confCommand = configuration{
			applicationPath: applicationPath,
			cfg:             cfg,
			repoPath:        repoPath,
			diffFiles:       diff,
			revision:        revision,
			pwd:             pwd,
		}

		return nil
	}
}
