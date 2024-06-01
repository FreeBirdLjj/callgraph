package subcmd

import (
	"context"
	"io"
	"os"

	"github.com/freebirdljj/immutable/comparator"
	"github.com/freebirdljj/immutable/slice"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/ssa"

	"github.com/freebirdljj/callgraph/gencallgraph"
	"github.com/freebirdljj/callgraph/statistics"
)

type (
	statisticsCommandFlags struct {
		output string
	}

	row struct {
		f        *ssa.Function
		refs     int
		refFuncs int
		refPkgs  int
	}
)

var (
	StyleNoColumns = table.Style{
		Name: "StyleNoColumns",
		Box: table.BoxStyle{
			BottomLeft:       "",
			BottomRight:      "",
			BottomSeparator:  "-",
			EmptySeparator:   text.RepeatAndTrim(" ", text.RuneWidthWithoutEscSequences("-")),
			Left:             "",
			LeftSeparator:    "",
			MiddleHorizontal: "-",
			MiddleSeparator:  "-",
			MiddleVertical:   " ",
			PaddingLeft:      "",
			PaddingRight:     "",
			PageSeparator:    "\n",
			Right:            "",
			RightSeparator:   "",
			TopLeft:          "",
			TopRight:         "",
			TopSeparator:     "-",
			UnfinishedRow:    " ~",
		},
		Color:   table.ColorOptionsDefault,
		Format:  table.FormatOptionsDefault,
		HTML:    table.DefaultHTMLOptions,
		Options: table.OptionsDefault,
		Title:   table.TitleOptionsDefault,
	}
)

func StatisticsCommand() *cobra.Command {

	const (
		cmdName        = "statistics"
		cmdDescription = "statistics calls of given packages"
	)

	statisticsCmdFlags := statisticsCommandFlags{}
	statisticsCmd := cobra.Command{
		Use:                   cmdName + " [flags] [packages]",
		DisableFlagsInUseLine: true,
		Short:                 cmdDescription,
		Long:                  cmdDescription,
		RunE:                  wrapStatisticsCmdRunE(&statisticsCmdFlags),
	}

	statisticsCmd.Flags().StringVarP(&statisticsCmdFlags.output, "output", "o", "-", "write to file")

	return &statisticsCmd
}

func wrapStatisticsCmdRunE(flags *statisticsCommandFlags) func(cmd *cobra.Command, args []string) error {
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

		return runStatistics(ctx, output, patterns)
	}
}

func runStatistics(ctx context.Context, output io.Writer, patterns []string) error {

	callg, err := gencallgraph.GenCallGraphForPackages(ctx, patterns)
	if err != nil {
		return err
	}

	funcLevelStatistics, err := statistics.StatisticsCallGraphAtFuncLevel(callg)
	if err != nil {
		return err
	}

	rows := make([]row, 0, len(funcLevelStatistics))
	for f, funcStatistics := range funcLevelStatistics {
		rows = append(rows, row{
			f:        f,
			refs:     funcStatistics.References,
			refFuncs: funcStatistics.ReferencedFuncs,
			refPkgs:  funcStatistics.ReferencedPkgs,
		})
	}

	drawFuncLevelStatisticsTable(rows, output)
	return nil
}

func drawFuncLevelStatisticsTable(rows []row, output io.Writer) {

	t := table.NewWriter()
	t.SetStyle(StyleNoColumns)
	t.SetOutputMirror(output)

	t.AppendHeader(table.Row{"func", "refs", "ref funcs", "ref pkgs"})

	for _, funcs := range slice.GroupBy(
		slice.FromGoSlice(rows).Filter(func(row row) bool {
			// only keep true source functions
			return row.f.Synthetic == ""
		}),
		comparator.CascadeComparator(
			comparator.OrderedComparator[string],
			func(row row) string {
				return row.f.Pkg.Pkg.Path()
			},
		),
	) {
		t.AppendRows(
			slice.Map(
				funcs.Sort(
					comparator.CascadeComparator(
						comparator.OrderedComparator[string],
						func(row row) string {
							return row.f.String()
						},
					),
				),
				func(row row) table.Row {
					return table.Row{row.f.String(), row.refs, row.refFuncs, row.refPkgs}
				},
			),
		)
		t.AppendSeparator()
	}

	t.Render()
}
