package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

var verboseCount int

var rootCmd = &cobra.Command{
	Use:   "csvtool",
	Short: "CSV Tool - CSV Swiss Army Knife",
	Long: fmt.Sprintf(
		`A small command line tool for performing operations on CSV files.

Examples:

Get the first column from the CSV
cat myfile.csv | csvtool cut -c1
`,
	),
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("error executing root command")
	}
}

func init() {
	cobra.OnInitialize(func() {
		logWriter := zerolog.ConsoleWriter{Out: os.Stderr}
		logWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("csvtool-%s", i))
		}

		log.Logger = log.Output(logWriter)

		if verboseCount == 0 { zerolog.SetGlobalLevel(zerolog.WarnLevel) }
		if verboseCount == 1 { zerolog.SetGlobalLevel(zerolog.InfoLevel) }
		if verboseCount >= 2 { zerolog.SetGlobalLevel(zerolog.DebugLevel) }
	})

	rootCmd.PersistentFlags().CountVarP(&verboseCount, "verbose", "v", "Print additional logging info, may be used multiple times")
}

