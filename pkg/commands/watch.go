package commands

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"monobuild/internal/errors"
	"monobuild/pkg/config"
	"monobuild/pkg/deps"
	"monobuild/pkg/service"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"
)

func GetWatch() *cobra.Command {
	var entryPath string

	command := &cobra.Command{
		Use: "watch",
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

			wg := sync.WaitGroup{}

			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			for packName, pack := range cfg.Packages {
				packName = packName

				fullPathPackage := path.Join(applicationPath, packName, pack.Entry)

				depsFiles, err := deps.GetDepsAsFiles(fullPathPackage)
				if err != nil {
					return fmt.Errorf("cannot load deps: %w", err)
				}

				wg.Add(1)

				watcher, err := fsnotify.NewWatcher()
				if err != nil {
					return err
				}

				go func() {
					// listen files
					defer watcher.Close()
					defer wg.Done()

					if err := func() error {
						for _, dep := range depsFiles {
							if err := watcher.Add(dep); err != nil {
								return err
							}
						}

						w := service.NewWriter(packName)
						packageService, err := service.NewService(fullPathPackage, w)
						if err != nil {
							return err
						}

						for {
							select {
							case event, ok := <-watcher.Events:
								if !ok {
									return err
								}

								if event.Op&fsnotify.Write == fsnotify.Write {
									if err := packageService.Reload(); err != nil {
										return fmt.Errorf("cannot reload service: %w", err)
									}
								}
							case <-done:
								return packageService.Stop()
							}
						}
					}(); err != nil {
						fmt.Println(err)
						return
					}
				}()
			}

			wg.Wait()

			return nil
		},
	}

	command.Flags().StringVar(&entryPath, "entry", "", "use for specific config file (YAML format)")

	return command
}
