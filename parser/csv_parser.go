package parser

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/afero"
	"harry-pap/beat_assignment/calculator"
	"io"
	"strconv"
)

// ParseInputCSV reads a given *afero.File CSV file, parses each line into a RidePart,
// batches the ride parts using the rideId, and pushes them to given channel
// Every 10,000 rides, a message is printed to standard output
func ParseInputCSV(file *afero.File, channel chan []calculator.RidePart) {
	var lastID int64
	var rides = make([]calculator.RidePart, 0, 512)
	var rideCounter int64

	reader := csv.NewReader(*file)
	initialized := false

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		entry := parseEntry(line)

		if initialized == false || entry.RideID == lastID {
			if initialized == false {
				lastID = entry.RideID
				initialized = true
			}
			rides = append(rides, entry)
		} else {
			rideCounter++

			if rideCounter%10000 == 0 {
				fmt.Printf("Processed %d rides\n", rideCounter)
			}

			channel <- rides

			lastID = entry.RideID

			rides = make([]calculator.RidePart, 0, 512)
			rides = append(rides, entry)
		}
	}

	if len(rides) > 0 {
		channel <- rides
	}
}

func parseEntry(line []string) calculator.RidePart {
	id, _ := strconv.ParseInt(line[0], 10, 32)
	lat, _ := strconv.ParseFloat(line[1], 64)
	long, _ := strconv.ParseFloat(line[2], 64)
	timestamp, _ := strconv.ParseInt(line[3], 10, 32)

	return calculator.RidePart{
		RideID:     id,
		Coordinate: calculator.Coordinate{Latitude: lat, Longitude: long},
		Timestamp:  int32(timestamp),
	}
}
