package main

import (
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

// End-to-end tests
// Note that a relative path is used for the test input('../testdata/',), and the result is written in /tmp/output.csv
// Unfortunately, this is OS dependent
func Test_main(t *testing.T) {
	type args struct {
		inputFile string
		want      []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"single artificial CSV",
			args{
				"testdata/sample.csv",
				[]string{"", "1,27.62"},
			},
		},
		{
			"provided CSV",
			args{
				"testdata/paths.csv",
				[]string{"", //last newline
					"1,11.34", "2,13.10", "3,33.84", "4,3.47", "5,22.78", "6,9.41",
					"7,30.01", "8,9.21", "9,6.35"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := os.Args
			outputCsv, err := os.CreateTemp("", "output.csv")

			if err != nil {
				t.Errorf("Failed to create temp file, %s", err)
				return
			}

			args[1] = tt.args.inputFile
			args[2] = outputCsv.Name()

			finished := make(chan string)
			go func() {
				main()
				finished <- "done"
			}()

			select {

			case <-time.After(1 * time.Second):
				t.Errorf("Main failed to complete after 1 second")
			case <-finished:
				break
			}

			bytes, _ := os.ReadFile(outputCsv.Name())

			got := strings.Split(string(bytes), "\n")

			sort.Strings(got)

			os.Remove(outputCsv.Name())
			if !reflect.DeepEqual(got, tt.args.want) {
				t.Errorf("End to end got = %v, want %v", got, tt.args.want)
			}
		})
	}
}
