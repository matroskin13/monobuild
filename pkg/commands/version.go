package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func GetVersionCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "check version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("version: x.x.x")
		},
	}

	return command
}
