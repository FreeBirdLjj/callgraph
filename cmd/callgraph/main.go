package main

import (
	"github.com/spf13/cobra"

	"github.com/freebirdljj/callgraph/cmd/callgraph/subcmd"
)

func rootCmd() *cobra.Command {

	const (
		cmdName        = "callgraph"
		cmdDescription = "A tool to analyze go program callgraph"
	)

	subcmds := []*cobra.Command{
		subcmd.DrawCommand(),
		subcmd.StatisticsCommand(),
	}

	rootCmd := cobra.Command{
		Use:   cmdName,
		Short: cmdDescription,
		Long:  cmdDescription,
	}
	rootCmd.AddCommand(subcmds...)
	return &rootCmd
}

func main() {
	cmd := rootCmd()
	cmd.Execute()
}
