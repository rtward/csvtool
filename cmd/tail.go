package cmd

import (
	"encoding/csv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var tailRecordCount int

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Select some number of records from the beginning of a CSV file",
	Run: func(cmd *cobra.Command, args []string) {
		csvReader := csv.NewReader(os.Stdin)
		csvWriter := csv.NewWriter(os.Stdout)

		headerRecord, err := csvReader.Read()
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("error reading CSV header")
		}
		err = csvWriter.Write(headerRecord)
		if err != nil {
			log.Fatal().
				Err(err).
				Strs("headerRecord", headerRecord).
				Msg("error writing CSV header")
		}

		outputRecords := make([][]string, 0)

		// Iterate through the file, always keeping the last ten lines we've read
		for  {
			inputRecord, err := csvReader.Read()
			if err == io.EOF { break }
			if err != nil {
				log.Fatal().Err(err).Msg("error parsing CSV")
			}

			outputRecords = append(outputRecords, inputRecord)

			if len(outputRecords) > tailRecordCount {
				outputRecords = outputRecords[len(outputRecords)-tailRecordCount:]
			}
		}

		err = csvWriter.WriteAll(outputRecords)
		if err != nil {
			log.Fatal().Err(err).Msg("error flushing CSV writer")
		}

		csvWriter.Flush()
	},
}

func init() {
	tailCmd.Flags().IntVarP(&tailRecordCount, "number", "n", 10, "the amount of lines (excluding the header) that should be taken from the end of the input")
	rootCmd.AddCommand(tailCmd)
}
