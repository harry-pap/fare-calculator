package parser

import (
	"harry-pap/beat_assignment/calculator"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/afero"
)

func TestParseInputCSV(t *testing.T) {

	type args struct {
		channel chan []calculator.RidePart
		csvData string
		want    [][]calculator.RidePart
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"no rides",
			args{
				make(chan []calculator.RidePart, 500),
				"",
				[][]calculator.RidePart{},
			},
		},
		{
			"1 ride",
			args{
				make(chan []calculator.RidePart, 500),
				`1,37.966660,23.728308,1405594957
1,37.966627,23.728263,1405594966
1,37.966625,23.728263,1405594974
1,37.966613,23.728375,1405594984`,
				[][]calculator.RidePart{
					{
						{RideID: 1, Coordinate: calculator.Coordinate{Latitude: 37.966660, Longitude: 23.728308}, Timestamp: 1405594957},
						{RideID: 1, Coordinate: calculator.Coordinate{Latitude: 37.966627, Longitude: 23.728263}, Timestamp: 1405594966},
						{RideID: 1, Coordinate: calculator.Coordinate{Latitude: 37.966625, Longitude: 23.728263}, Timestamp: 1405594974},
						{RideID: 1, Coordinate: calculator.Coordinate{Latitude: 37.966613, Longitude: 23.728375}, Timestamp: 1405594984},
					},
				},
			},
		},

		{
			"4 rides",
			args{
				make(chan []calculator.RidePart, 500),
				`1,37.966660,23.728308,1405594957
1,37.966627,23.728263,1405594966
1,37.966625,23.728263,1405594974
1,37.966613,23.728375,1405594984
2,37.966627,23.728263,1405594966
2,37.966625,23.728263,1405594974
2,37.966613,23.728375,1405594984
3,37.966627,23.728263,1405594966
3,37.966625,23.728263,1405594974
4,37.966625,23.728263,1405594974`,
				[][]calculator.RidePart{
					{
						{RideID: 1, Coordinate: calculator.Coordinate{Latitude: 37.966660, Longitude: 23.728308}, Timestamp: 1405594957},
						{RideID: 1, Coordinate: calculator.Coordinate{Latitude: 37.966627, Longitude: 23.728263}, Timestamp: 1405594966},
						{RideID: 1, Coordinate: calculator.Coordinate{Latitude: 37.966625, Longitude: 23.728263}, Timestamp: 1405594974},
						{1, calculator.Coordinate{Latitude: 37.966613, Longitude: 23.728375}, 1405594984},
					}, {
						{RideID: 2, Coordinate: calculator.Coordinate{Latitude: 37.966627, Longitude: 23.728263}, Timestamp: 1405594966},
						{RideID: 2, Coordinate: calculator.Coordinate{Latitude: 37.966625, Longitude: 23.728263}, Timestamp: 1405594974},
						{RideID: 2, Coordinate: calculator.Coordinate{Latitude: 37.966613, Longitude: 23.728375}, Timestamp: 1405594984},
					}, {
						{RideID: 3, Coordinate: calculator.Coordinate{Latitude: 37.966627, Longitude: 23.728263}, Timestamp: 1405594966},
						{RideID: 3, Coordinate: calculator.Coordinate{Latitude: 37.966625, Longitude: 23.728263}, Timestamp: 1405594974},
					}, {
						{RideID: 4, Coordinate: calculator.Coordinate{Latitude: 37.966625, Longitude: 23.728263}, Timestamp: 1405594974},
					},
				},
			},
		},
	}

	appFS := afero.NewMemMapFs()
	fileName := "/beat/rates_output.csv"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			file, _ := appFS.Create(fileName)

			file.WriteString(tt.args.csvData)
			file.Close()
			file, _ = appFS.Open(fileName)

			ParseInputCSV(&file, tt.args.channel)

			got := make([][]calculator.RidePart, 0, 10)

			for i := 0; i < len(tt.args.want); i++ {
				got = append(got, <-tt.args.channel)
			}

			if !reflect.DeepEqual(got, tt.args.want) {
				t.Errorf("ParseInputCSV() pushed to channel %v, want %v", got, tt.args.want)
			}

			select {

			case <-time.After(150 * time.Millisecond):
				break

			case <-tt.args.channel:
				t.Errorf("ParseInputCSV pushed more messages to the work channel than expected")
			}

			file.Close()

			appFS.Remove(fileName)
		})
	}
}
