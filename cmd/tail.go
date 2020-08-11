package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
)

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Select some number of records from the beginning of a CSV file",
	Run: func(cmd *cobra.Command, args []string) {
		recordCount, err := cmd.Flags().GetInt("number")
		if err != nil {
			log.Fatal().Err(err).Msg("error reading number argument")
		}

		if hasHeader {
			log.Debug().Strs("header", header).Msg("writing header row")
			err = csvWriter.Write(header)
			if err != nil {
				log.Fatal().Err(err).Msg("error writing CSV header")
			}
		}

		outputRecords := make([][]string, 0)

		// Iterate through the file, always keeping the last ten lines we've read
		for {
			inputRecord, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal().Err(err).Msg("error parsing CSV")
			}

			outputRecords = append(outputRecords, inputRecord)

			if len(outputRecords) > recordCount {
				outputRecords = outputRecords[len(outputRecords)-recordCount:]
			}
		}

		err = csvWriter.WriteAll(outputRecords)
		if err != nil {
			log.Fatal().Err(err).Msg("error writing output records")
		}
	},
}

func init() {
	tailCmd.Flags().IntP("number", "n", 10, "the amount of lines (excluding the header) that should be taken from the end of the input")
	rootCmd.AddCommand(tailCmd)
}
