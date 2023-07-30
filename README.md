This solution was originally implemented as an assignment for Beat back in early 2019,
a company based in Greece that had recently(at the time) been acquired by FreeNow, and which was shutdown by
FreeNow in 2021. Originally implemented in Go 1.11.1, has been updated to work with Go 1.20

# Problem Definition

Beat drivers perform thousands of rides per day. In order to ensure that our passengers always receive the
highest quality of service possible, we need to be able to identify mis-behaving drivers that overcharge
passengers and block them from our service.

Our servers keep detailed record of drivers' position throughout the progress of each ride. Due to the high
volume of daily rides we need to implement an automated way of estimating the fare for each ride so that
we may flag suspicious rides for review by our support teams.

Moreover, our drivers sometimes use low-cost devices that frequently report invalid/skewed GPS
coordinates to our servers. The fare estimation algorithm should be capable of detecting erroneous
coordinates and remove them before attempting to evaluate the ride fare.

In this exercise you are asked to create a **Fare Estimation script**. Its input will consist of a list of tuples of
the form **(id_ride, lat, lng, timestamp)** representing the position of the taxi-cab during a ride.
Two consecutive tuples belonging to the same ride form a segment S. For each segment, we define the
elapsed time **Δt** as the absolute difference of the segment endpoint timestamps and the distance covered
**Δs** as the Haversine distance of the segment endpoint coordinates.

Your first task is to filter the data on the input so that invalid entries are discarded before the fare
estimation process begins. Filtering should be performed as follows: consecutive tuples p1, p2 should be
used to calculate the segment’s speed U. If U > 100km/h, p2 should be removed from the set.
Once the data has been filtered you should proceed in estimating the fare by **aggregating the individual
estimates** for the ride segments using rules tabulated below:

| State               | Applicable when            | Fare amount                 |
|---------------------|----------------------------|-----------------------------|
| MOVING (U > 10km/h) | Time of day (05:00, 00:00] | 0.74 per km                 |
| MOVING (U > 10km/h) | Time of day (00:00, 05:00] | 1.30 per km                 |
| IDLE (U <= 10km/h)  | Always                     | 11.90 per hour of idle time |

At the start of each ride, a standard ‘flag’ amount of 1.30 is charged to the ride’s fare. The minimum ride
fare should be at least 3.47.

### Input Data
The sample data file contains one record per line (comma separated values). The input data is guaranteed
to contain continuous row blocks for each individual ride (i.e. the data will not be multiplexed). In addition,
the data is also pre-sorted for you in ascending timestamp order.

### Deliverables
1. A Golang 1.8+ program that processes input as provided by the sample data, filters out invalid points
   and produces an output comma separated value text file with the following format:
   id_ride, fare_estimate
Note that your solution should make good use of Golang’s concurrency facilities and result in a high-
performance script that is space and time efficient in its processing of the data.
2. A brief document with an overview of your design
3. A comprehensive unit and end-to-end test suit for your code.
We expect you to deliver:


### Notes
• In order to calculate the distance between two (lat,lng) pairs you can use the [Haversine distance
formula](https://en.wikipedia.org/wiki/Haversine_formula).

• Your code should be capable of ingesting and processing large datasets. Assume we will try to pipe
several GBs worth of samples to your script.

• Follow best-practices and be well documented throughout.
Good luck!

<br>
<br>
(end of problem definition)

## Run:
Under the current directory run:
`go run {{source_csv}} {{target_csv}}`

## DESIGN
The solution was implemented using the Fan-out/fan-in pattern. The main goroutine parses the input CSV,
and when all the parts of a ride are read, pushes them into a channel. Several worker goroutines read from this channel,
calculate the fare for the given ride, and push the result into a result channel. A single goroutine reads from
the result channel, and writes the ride estimates into a CSV file. Due to time limitations alternative concurrency 
implementations were not tried, most notably having distinct goroutines for reading and 
batching(with no parsing in place), and for parsing them into RideParts


## TESTS
There are unit tests in place for all the exported methods, and end-to-end tests in main_test.go
The end-to-end tests use relative paths for the input files, so if you're using GoLand you need to 
run the tests from this directory, as running the file/function will cause a FileNotFound panic


##### NOTE
Note that there are no checks in place for the validity of source_csv and target_csv arguments, and due to the fact it's
a small script, panics are used more than they normally would be.
