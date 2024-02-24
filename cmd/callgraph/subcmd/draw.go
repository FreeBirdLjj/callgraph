package subcmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/tools/go/callgraph/static"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"

	"github.com/FreeBirdLjj/callgraph/draw/dot"
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

	pkgs, err := packages.Load(&packages.Config{
		Mode: -1,
	}, patterns...)
	if err != nil {
		return err
	}

	ssaProg, _ := ssautil.AllPackages(pkgs, ssa.SanityCheckFunctions)
	callg := static.CallGraph(ssaProg)

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
