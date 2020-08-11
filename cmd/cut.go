package cmd

import (
	"encoding/csv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
)

func writeSelectedColumns(selectedColumns []int, csvWriter *csv.Writer, inputRecord []string) {
	outputRecord := make([]string, 0)

	for inputColumnIdx := range inputRecord {
		useColumn := false

		// If we don't list any columns, take all of them
		if len(selectedColumns) == 0 { useColumn = true }

		// If we do have some selected columns, then see if any of them match our current record index
		for selectedColumnIdx := range selectedColumns {
			selectedColumn := selectedColumns[selectedColumnIdx]
			// And if they do, then we'll take that column
			if selectedColumn == inputColumnIdx { useColumn = true }
		}

		if useColumn {
			outputRecord = append(outputRecord, inputRecord[inputColumnIdx])
		}
	}

	err := csvWriter.Write(outputRecord)

	if err != nil {
		log.Fatal().
			Err(err).
			Strs("inputRecord", inputRecord).
			Strs("outputRecord", outputRecord).
			Msg("error writing CSV header")
	}
}

var cutCmd = &cobra.Command{
	Use:   "cut",
	Short: "Select one or more columns from a CSV file",
	Run: func(cmd *cobra.Command, args []string) {
		selectedColumns, err := cmd.Flags().GetIntSlice("columns")
		if err != nil { log.Fatal().Err(err).Msg("error reading columns argument") }

		if hasHeader {
			log.Debug().Strs("header", header).Msg("writing header row")
			writeSelectedColumns(selectedColumns, csvWriter, header)
		}

		for {
			inputRecord, err := csvReader.Read()
			if err == io.EOF { break }
			if err != nil { log.Fatal().Err(err).Msg("error parsing CSV") }

			writeSelectedColumns(selectedColumns, csvWriter, inputRecord)
		}
	},
}

func init() {
	cutCmd.Flags().IntSliceP("columns", "c", nil, "a comma separated list of columns that should be extracted from the CSV")
	rootCmd.AddCommand(cutCmd)
}
