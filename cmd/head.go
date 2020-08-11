package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
)

var headCmd = &cobra.Command{
	Use:   "head",
	Short: "Select some number of records from the beginning of a CSV file",
	Run: func(cmd *cobra.Command, args []string) {
		recordCount, err := cmd.Flags().GetInt("number")
		if err != nil { log.Fatal().Err(err).Msg("error reading number argument") }

		if hasHeader {
			log.Debug().Strs("header", header).Msg("writing header row")
			err = csvWriter.Write(header)
			if err != nil { log.Fatal().Err(err).Msg("error writing CSV header") }
		}

		// Iterate through the file, and take recordCount + 1 lines, to account for the hasHeader
		for lineNum :=0; lineNum < recordCount; lineNum++  {
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

	},
}

func init() {
	headCmd.Flags().IntP("number", "n", 10, "the amount of lines (excluding the header) that should be taken from the head of the input")
	rootCmd.AddCommand(headCmd)
}
