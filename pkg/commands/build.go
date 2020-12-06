package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"monobuild/internal/env"
	"monobuild/internal/errors"
	"monobuild/internal/slice"
	"monobuild/pkg/build"
	"monobuild/pkg/config"
	"monobuild/pkg/deps"
	"monobuild/pkg/git"
	"os"
	"path"
	"strings"
)

func GetBuild() *cobra.Command {
	var entryPath string

	command := &cobra.Command{
		Use: "build your packages",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			}

			cfg, err := config.ParseDefaultConfig(applicationPath)
			if err != nil {
				return errors.NewRichError("Incorrect format of configuration file", err)
			}

			repoPath, err := git.ResolveGitPath(applicationPath)
			if err != nil {
				return err
			}

			dock := build.NewDocker()

			diff, err := git.GetDiffFiles(repoPath)
			if err != nil {
				return err
			}

			for packName, pack := range cfg.Packages {
				fullPathPackage := path.Join(applicationPath, packName, pack.Entry)

				depsFiles, err := deps.GetDepsAsFiles(fullPathPackage)
				if err != nil {
					return fmt.Errorf("cannot load deps: %w", err)
				}

				needBuild := len(slice.Intersection(diff, depsFiles)) > 0
				image, err := env.ParseTemplateWithEnv(pack.Build.Docker.Image)
				if err != nil {
					return fmt.Errorf("invalid template: %w", err)
				}

				if needBuild {
					fmt.Printf("Package %q has changed, build has started...\r\n", packName)

					if pack.Build.Docker != nil {
						if err := dock.Build(cmd.Context(), applicationPath, path.Join(packName, pack.Entry), image); err != nil {
							return fmt.Errorf("cannot build image: %w", err)
						}

						fmt.Printf("Successfuly build docker image and push %q with image %q\r\n", packName, image)
					}
				} else {
					fmt.Printf("Package %q not changed \r\n", packName)
				}
			}

			return nil
		},
	}

	command.Flags().StringVar(&entryPath, "entry", "", "use for specific config file (YAML format)")

	return command
}
