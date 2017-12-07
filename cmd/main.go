package main

import (
	"fmt"
	"os"
	"time"

	goinflux "github.com/atw527/goinflux/goinflux"
)

func main() {
	os.Setenv("INFLUX_HOST", "localhost")
	os.Setenv("INFLUX_PORT", "8086")
	os.Setenv("INFLUX_DB", "testing")

	fmt.Printf("INIT\n\n")

	//
	// Pulsing/trickling points will test the cleanup and timer functions
	//

	// these would tend to leave points behind on the shutdown
	writePoints(10, 32, 32, 0, "pulse")
	writePoints(10, 16, 32, 0, "pulse")
	writePoints(100, 16, 32, 0, "pulse")

	// timer should kick in for this one
	writePoints(20, 32, 32, 1000, "pulse")

	//
	// Flooding points will test throughput
	//

	// scaling up pointBatchSize (moderate improvement)
	writePoints(1000000, 16, 1024, 0, "flood")
	writePoints(1000000, 16, 4096, 0, "flood")
	writePoints(1000000, 16, 16384, 0, "flood")
	writePoints(1000000, 16, 65536, 0, "flood")

	// scaling up chanBuffer (flat improvement)
	writePoints(1000000, 1024, 1024, 0, "flood")
	writePoints(1000000, 4096, 1024, 0, "flood")
	writePoints(1000000, 16384, 1024, 0, "flood")
	writePoints(1000000, 65536, 1024, 0, "flood")

	// scaling up both in parallel (best improvement)
	writePoints(1000000, 1024, 1024, 0, "flood")
	writePoints(1000000, 4096, 4096, 0, "flood")
	writePoints(1000000, 16384, 16384, 0, "flood")
	writePoints(1000000, 65536, 65536, 0, "flood")
}

func writePoints(points int, chanBuffer int, pointBatchSize int, interval time.Duration, metric string) {
	tags := goinflux.TagGroup{}
	tags["metric"] = metric

	fields := goinflux.FieldGroup{}
	fields["value"] = 1

	fmt.Printf("Writing %v (%vM) points, chanBuffer=%v, pointBatchSize=%v\n", points, points/1000000, chanBuffer, pointBatchSize)
	gi, err := goinflux.NewGoInflux(os.Getenv("INFLUX_HOST"), os.Getenv("INFLUX_PORT"), chanBuffer, pointBatchSize, 5)
	if err != nil {
		fmt.Printf("Error init: %v\n", err.Error())
		return
	}

	// unleash a torrent of points
	for i := 0; i < points; i++ {
		err := gi.AddPoint("test", tags, fields, time.Now().UnixNano())
		if err != nil {
			fmt.Printf("Error point: %v\n", err.Error())
			return
		}

		if interval > 0 {
			time.Sleep(interval * time.Millisecond)
		}
	}

	gi.Stahp()

	time.Sleep(5 * time.Second) // this sleep period puts space in the graph
}
