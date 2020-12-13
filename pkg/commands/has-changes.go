package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"monobuild/internal/slice"
	"monobuild/pkg/deps"
	"path"
)

func GetHasChangesCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "changes",
		Short: "changes",
	}

	conf, preRunConf := setupConfigurationHook(command)

	command.PreRunE = preRunConf
	command.RunE = func(cmd *cobra.Command, args []string) error {
		buildMap := map[string][]string{}

		for packName, pack := range conf.cfg.Packages {
			packName = packName
			fullPathPackage := pack.FullEntry

			if fullPathPackage == "" {
				fullPathPackage = path.Join(conf.applicationPath, packName, pack.Entry)
			}

			depsFiles, err := deps.GetDepsAsFiles(fullPathPackage)
			if err != nil {
				return fmt.Errorf("cannot load deps: %w", err)
			}

			isChangesDeps, changedDeps, err := deps.PackageChangeDeps(fullPathPackage, conf.revision)
			if err != nil {
				return fmt.Errorf("cannot get deps: %w", err)
			}

			diff := slice.Intersection(conf.diffFiles, depsFiles)

			if len(diff) > 0 || isChangesDeps {
				var report []string

				for _, diffFile := range diff {
					report = append(report, fmt.Sprintf("dependency file: %s", diffFile))
				}

				for _, dep := range changedDeps {
					report = append(report, fmt.Sprintf("require package: %s", dep.Mod))
				}

				buildMap[packName] = report
			}
		}

		if len(buildMap) > 0 {
			for packName, report := range buildMap {
				fmt.Printf("== Changes for %s ==\r\n", packName)
				for _, line := range report {
					fmt.Println(line)
				}
			}
		}

		return nil
	}

	return command
}
