package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"math"
	"math/rand"
	"time"
)

var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Select a random sampling of lines from a CSV file",
	Run: func(cmd *cobra.Command, args []string) {
		recordCount, err := cmd.Flags().GetInt("number")
		if err != nil { log.Fatal().Err(err).Msg("error reading number argument") }

		if hasHeader {
			log.Debug().Strs("header", header).Msg("writing header row")
			err = csvWriter.Write(header)
			if err != nil { log.Fatal().Err(err).Msg("error writing CSV header") }
		}

		// Implementation of the resivoir sampling algorithm from: https://en.wikipedia.org/wiki/Reservoir_sampling
		eof := false
		outputRecords := make([][]string, 0)

		// First fill up our output array with the requested number of items
		for len(outputRecords) < recordCount {
			inputRecord, err := csvReader.Read()
			if err == io.EOF {
				eof = true
				break
			}
			if err != nil {
				log.Fatal().Err(err).Msg("error parsing CSV")
			}

			outputRecords = append(outputRecords, inputRecord)
		}

		// If we've already run out of lines, there's nothing to do
		if !eof {
			// Generate an inital random seed
			rand.Seed(time.Now().UnixNano())
			w := math.Exp(math.Log(rand.Float64())/float64(recordCount))

			// Now go through the rest of the items, maybe swapping them out for one
			for {
				// Calculate a random number of items to skip
				skip := int(math.Floor(math.Log(rand.Float64())/math.Log(1-w)))
				for i := 0; i < skip; i++ {
					_, err := csvReader.Read()
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Fatal().Err(err).Msg("error parsing CSV")
					}
				}

				// Then take the next record and assign it to a random location
				inputRecord, err := csvReader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal().Err(err).Msg("error parsing CSV")
				}
				outputRecords[rand.Intn(len(outputRecords))] = inputRecord

				w = w * math.Exp(math.Log(rand.Float64())/float64(recordCount))
			}
		}

		err = csvWriter.WriteAll(outputRecords)
		if err != nil { log.Fatal().Err(err).Msg("error writing output records") }
	},
}

func init() {
	sampleCmd.Flags().IntP("number", "n", 10, "the amount of lines (excluding the header) that should be taken from the end of the input")
	rootCmd.AddCommand(sampleCmd)
}
