package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var hasHeader bool
var header []string
var csvReader *csv.Reader
var csvWriter *csv.Writer

var rootCmd = &cobra.Command{
	Use:   "csvtool",
	Short: "CSV Tool - CSV Swiss Army Knife",
	Long: "A small command line tool for performing operations on CSV files.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verboseCount, err := cmd.Root().PersistentFlags().GetCount("verbose")
		if err != nil { log.Fatal().Err(err).Msg("error reading verbose count flag") }

		logWriter := zerolog.ConsoleWriter{Out: os.Stderr}
		logWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("csvtool-%s", i))
		}

		log.Logger = log.Output(logWriter)

		if verboseCount == 0 { zerolog.SetGlobalLevel(zerolog.WarnLevel) }
		if verboseCount == 1 { zerolog.SetGlobalLevel(zerolog.InfoLevel) }
		if verboseCount >= 2 { zerolog.SetGlobalLevel(zerolog.DebugLevel) }

		log.Debug().Int("verboseCount", verboseCount).Msg("zerolog configured")

		log.Debug().Msg("creating CSV reader")
		csvReader = csv.NewReader(os.Stdin)
		log.Debug().Msg("created CSV reader")

		log.Debug().Msg("creating CSV writer")
		csvWriter = csv.NewWriter(os.Stdout)
		log.Debug().Msg("created CSV writer")

		noHeader, err := cmd.Root().PersistentFlags().GetBool("no-header")
		if err != nil { log.Fatal().Err(err).Msg("error reading header command flag") }

		hasHeader = !noHeader

		if hasHeader {
			log.Debug().Msg("reading header row")
			header, err = csvReader.Read()
			if err != nil { log.Fatal().Err(err).Msg("error reading header row") }
			log.Info().Strs("header", header).Msg("read header row")
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Debug().Msg("flushing CSV output")
		csvWriter.Flush()
		err := csvWriter.Error()
		if err != nil { log.Fatal().Err(err).Msg("error flushing CSV output") }
		log.Debug().Msg("flushed CSV output")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil { log.Fatal().Err(err).Msg("error executing root command") }
}

func init() {
	rootCmd.PersistentFlags().CountP("verbose", "v", "Print additional logging info, may be used multiple times")
	rootCmd.PersistentFlags().BoolP("no-header", "H", false, "Treat the first line as input instead of a header")
}

