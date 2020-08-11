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

func getSelectedColumns(cmd *cobra.Command) []int {
	columnsByNum, err := cmd.Flags().GetIntSlice("columns")
	if err != nil { log.Fatal().Err(err).Msg("error reading columns argument") }

	if len(columnsByNum) > 0 {
		log.Debug().Ints("selectedColumns", columnsByNum).Msg("got columns from columns arg")
		return columnsByNum
	}

	columnsByName, err := cmd.Flags().GetStringSlice("column-names")
	if err != nil { log.Fatal().Err(err).Msg("error reading columnn-names argument") }

	if len(columnsByName) == 0 {
		log.Debug().Msg("no columns arg provided, defaulting to all columns")
		return nil
	}

	if !hasHeader { log.Fatal().Msg("cannot use column-names with CSV without a header") }

	selectedColumns := make([]int, len(columnsByName))
	for selectedColumnIdx := range columnsByName {
		columnName := columnsByName[selectedColumnIdx]

		found := false
		for headerColumnIdx := range header {
			headerColumn := header[headerColumnIdx]
			if headerColumn == columnName {
				selectedColumns = append(selectedColumns, headerColumnIdx)
				found = true
			}
		}

		if !found { log.Fatal().Str("column", columnName).Msg("unable to find column in the CSV header") }
	}

	log.Debug().Ints("selectedColumns", selectedColumns).Msg("got columns from column-names arg")
	return selectedColumns
}

var cutCmd = &cobra.Command{
	Use:   "cut",
	Short: "Select one or more columns from a CSV file",
	Run: func(cmd *cobra.Command, args []string) {
		selectedColumns := getSelectedColumns(cmd)

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
	cutCmd.Flags().StringSliceP("column-names", "n", nil, "a comma separated list of columns by name that should be extracted from the CSV")
	rootCmd.AddCommand(cutCmd)
}
