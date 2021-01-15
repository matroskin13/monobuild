package commands

import (
	"bytes"
	"fmt"
	"github.com/matroskin13/monobuild/internal/errors"
	"github.com/matroskin13/monobuild/internal/slice"
	"github.com/matroskin13/monobuild/pkg/deps"
	"github.com/spf13/cobra"
	"html/template"
	"io/ioutil"
	"path"
)

func GetHasChangesCommand() *cobra.Command {
	var useTemplate string
	var outTemplateName string

	command := &cobra.Command{
		Use:   "changes",
		Short: "changes",
	}

	conf, preRunConf := setupConfigurationHook(command)

	command.PreRunE = preRunConf
	command.RunE = func(cmd *cobra.Command, args []string) error {
		buildMap := map[string][]string{}

		var outTemplate []byte
		var resultTemplate bytes.Buffer

		if useTemplate != "" {
			templateFile, err := ioutil.ReadFile(path.Join(conf.pwd, useTemplate))
			if err != nil {
				return errors.NewRichError("cannot read template file", err)
			}

			outTemplate = templateFile
		}

		for packName, pack := range conf.cfg.Packages {
			packName = packName
			fullPathPackage := pack.FullEntry

			if fullPathPackage == "" {
				fullPathPackage = path.Join(conf.applicationPath, pack.Entry)
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

				buildMap[pack.Entry] = report

				if useTemplate != "" {
					tpl, err := template.New("").Parse(string(outTemplate))
					if err != nil {
						return err
					}

					if err := tpl.Execute(&resultTemplate, map[string]interface{}{
						"serviceName": packName,
					}); err != nil {
						return err
					}

					resultTemplate.Write([]byte("\r\n\r\n"))
				}
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

		if useTemplate != "" {
			if err := ioutil.WriteFile(path.Join(conf.pwd, outTemplateName), resultTemplate.Bytes(), 0644); err != nil {
				return err
			}
		}

		return nil
	}

	command.Flags().StringVar(&useTemplate, "use-template", "", "use property for generate template")
	command.Flags().StringVar(&outTemplateName, "out-template", "", "specify output template name")

	return command
}
