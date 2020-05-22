package cmd

import (
	"encoding/csv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"math"
	"math/rand"
	"os"
	"time"
)

var sampleRecordCount int

var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Select a random sampling of lines from a CSV file",
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

		// Implementation of the resivoir sampling algorithm from: https://en.wikipedia.org/wiki/Reservoir_sampling
		eof := false
		outputRecords := make([][]string, 0)

		// First fill up our output array with the requested number of items
		for len(outputRecords) < sampleRecordCount {
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
			w := math.Exp(math.Log(rand.Float64())/float64(sampleRecordCount))

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

				w = w * math.Exp(math.Log(rand.Float64())/float64(sampleRecordCount))
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
	sampleCmd.Flags().IntVarP(&sampleRecordCount, "number", "n", 10, "the amount of lines (excluding the header) that should be taken from the end of the input")
	rootCmd.AddCommand(sampleCmd)
}
