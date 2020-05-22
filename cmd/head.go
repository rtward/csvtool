package cmd

import (
	"encoding/csv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var headRecordCount int

var headCmd = &cobra.Command{
	Use:   "head",
	Short: "Select some number of records from the beginning of a CSV file",
	Run: func(cmd *cobra.Command, args []string) {
		csvReader := csv.NewReader(os.Stdin)
		csvWriter := csv.NewWriter(os.Stdout)

		hasHeader, err := cmd.Root().PersistentFlags().GetBool("header")
		if err != nil { log.Fatal().Err(err).Msg("error reading header command option") }

		numRecords := headRecordCount
		if hasHeader { numRecords += 1 }
		// Iterate through the file, and take recordCount + 1 lines, to account for the hasHeader
		for lineNum :=0; lineNum < numRecords; lineNum++  {
			inputRecord, err := csvReader.Read()
			if err == io.EOF { break }
			if err != nil {
				log.Fatal().Err(err).Msg("error parsing CSV")
			}

			err = csvWriter.Write(inputRecord)
			if err != nil {
				log.Fatal().
					Err(err).
					Strs("inputRecord", inputRecord).
					Msg("error writing CSV")
			}
		}

		csvWriter.Flush()
	},
}

func init() {
	headCmd.Flags().IntVarP(&headRecordCount, "number", "n", 10, "the amount of lines (excluding the header) that should be taken from the head of the input")
	rootCmd.AddCommand(headCmd)
}
