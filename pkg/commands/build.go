package commands

import (
	"fmt"
	"github.com/matroskin13/monobuild/internal/env"
	"github.com/matroskin13/monobuild/internal/slice"
	"github.com/matroskin13/monobuild/pkg/build"
	"github.com/matroskin13/monobuild/pkg/deps"
	"github.com/spf13/cobra"
	"path"
)

func GetBuild() *cobra.Command {
	command := &cobra.Command{
		Use: "build your packages",
	}

	conf, preRunConf := setupConfigurationHook(command)

	command.PreRunE = preRunConf
	command.RunE = func(cmd *cobra.Command, args []string) error {
		dock := build.NewDocker()

		for packName, pack := range conf.cfg.Packages {
			fullPathPackage := path.Join(conf.applicationPath, pack.Entry)

			depsFiles, err := deps.GetDepsAsFiles(fullPathPackage)
			if err != nil {
				return fmt.Errorf("cannot load deps: %w", err)
			}

			needBuild := len(slice.Intersection(conf.diffFiles, depsFiles)) > 0
			image, err := env.ParseTemplateWithEnv(pack.Build.Docker.Image)
			if err != nil {
				return fmt.Errorf("invalid template: %w", err)
			}

			if needBuild {
				fmt.Printf("Package %q has changed, build has started...\r\n", packName)

				if pack.Build.Docker != nil {
					if err := dock.Build(cmd.Context(), conf.applicationPath, path.Join(pack.Entry), image); err != nil {
						return fmt.Errorf("cannot build image: %w", err)
					}

					fmt.Printf("Successfuly build docker image and push %q with image %q\r\n", packName, image)
				}
			} else {
				fmt.Printf("Package %q not changed \r\n", packName)
			}
		}

		return nil
	}

	return command
}
