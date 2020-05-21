package cmd

import (
	"encoding/csv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var selectedColumns []int

var cutCmd = &cobra.Command{
	Use:   "cut",
	Short: "Select one or more columns from a CSV file",
	Run: func(cmd *cobra.Command, args []string) {
		csvReader := csv.NewReader(os.Stdin)
		csvWriter := csv.NewWriter(os.Stdout)

		for {
			inputRecord, err := csvReader.Read()
			if err == io.EOF { break }
			if err != nil {
				log.Fatal().Err(err).Msg("error parsing CSV")
			}

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

			err = csvWriter.Write(outputRecord)

			if err != nil {
				log.Fatal().
					Err(err).
					Strs("inputRecord", inputRecord).
					Strs("outputRecord", outputRecord).
					Msg("error writing CSV")
			}
		}

		csvWriter.Flush()
	},
}

func init() {
	cutCmd.Flags().IntSliceVarP(&selectedColumns, "columns", "c", nil, "a comma separated list of columns that should be extracted from the CSV")
	rootCmd.AddCommand(cutCmd)
}
