package commands

import "github.com/spf13/cobra"

func GetMainCommand() *cobra.Command {
	command := &cobra.Command{
		Use: "monobuild for monorepo",
	}

	command.SilenceUsage = true
	command.AddCommand(GetVersionCommand(), GetBuild())

	return command
}
