package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"monobuild/pkg/config"
	"monobuild/pkg/deps"
	"os"
	"path"
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

			configPath := path.Join(pwd, entryPath, ".monobuild.yml")

			cfg, err := config.ParseConfigFromFile(configPath)
			if err != nil {
				return fmt.Errorf("cannot parse config file: %w", err)
			}

			for packName, pack := range cfg.Packages {
				fullPathPackage := path.Join(pwd, entryPath, packName, pack.Entry)

				depsFiles, err := deps.GetDepsAsFiles(fullPathPackage)
				if err != nil {
					return fmt.Errorf("cannot load deps: %w", err)
				}

				fmt.Println(depsFiles)
			}

			return nil
		},
	}

	command.Flags().StringVar(&entryPath, "entry", "", "use for specific config file (YAML format)")

	return command
}
