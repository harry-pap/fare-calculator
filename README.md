This project was developed using Go 1.11.1


The only dependency is `github.com/spf13/afero`, for (file) testing purposes

## Usage:
Under the current directory run:
`go run {{source_csv}} {{target_csv}}`


## DESIGN
The main goroutine parses the input CSV, and when all the parts of a ride are read, pushes them into 
a channel. Several worker goroutines read from this channel, calculate the fare for the given ride, and push the 
result into a result channel. A single goroutine reads from the result channel, and writes the ride estimates
into a CSV file. Due to time limitations alternative concurrency implementations were not tried, 
most notably having distinct goroutines for reading and batching(with no parsing in place),
 and for parsing them into RideParts


## TESTS
There are unit tests in place for all the exported methods, and end to end tests in main_test.go
The end to end tests use relative paths for the input files, so if you're using GoLand you need to 
run the tests from this directory, as running the file/function will cause a FileNotFound panic


##### NOTE
Note that there are no checks in place for the validity of source_csv and target_csv arguments, and due to the fact it's
a small script, panics are used more than they normally would be.
