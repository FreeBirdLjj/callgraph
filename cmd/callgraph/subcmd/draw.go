package subcmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/freebirdljj/callgraph/draw/dot"
	"github.com/freebirdljj/callgraph/gencallgraph"
)

type (
	drawCommandFlags struct {
		name   string
		output string
	}
)

func DrawCommand() *cobra.Command {

	const (
		cmdName        = "draw"
		cmdDescription = "draw a callgraph of given packages"
	)

	drawCmdFlags := drawCommandFlags{}
	drawCmd := cobra.Command{
		Use:                   cmdName + " [flags] [packages]",
		DisableFlagsInUseLine: true,
		Short:                 cmdDescription,
		Long:                  cmdDescription,
		RunE:                  wrapDrawCmdRunE(&drawCmdFlags),
	}

	drawCmd.Flags().StringVar(&drawCmdFlags.name, "name", "G", "graph name")
	drawCmd.Flags().StringVarP(&drawCmdFlags.output, "output", "o", "-", "write to file")

	return &drawCmd
}

func wrapDrawCmdRunE(flags *drawCommandFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		ctx := cmd.Context()

		output := cmd.OutOrStdout()
		if flags.output != "-" {
			outputFile, err := os.Create(flags.output)
			if err != nil {
				return err
			}
			defer outputFile.Close()
			output = outputFile
		}

		patterns := args
		if len(patterns) == 0 {
			patterns = []string{"."}
		}

		return draw(ctx, flags.name, output, patterns)
	}
}

func draw(ctx context.Context, name string, output io.Writer, patterns []string) error {

	callg, err := gencallgraph.GenCallGraphForPackages(ctx, patterns)
	if err != nil {
		return err
	}

	graph, err := dot.DrawCallGraphAsDotDigraph(callg, name)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(output, graph)
	if err != nil {
		return err
	}

	return nil
}
