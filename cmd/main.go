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

	gi, err := goinflux.NewGoInfluxEnv()
	if err != nil {
		fmt.Printf("Error init: %v\n", err.Error())
		return
	}
	defer gi.Stahp()

	fmt.Printf("INIT\n")

	tags := goinflux.TagGroup{}
	tags["metric"] = "pulse"

	fields := goinflux.FieldGroup{}
	fields["value"] = 1

	// unleash a torrent of points
	//*
	for {
		err := gi.AddPoint("test", tags, fields, time.Now().UnixNano())
		if err != nil {
			fmt.Printf("Error point: %v\n", err.Error())
			return
		}
	}
	// */

	// trickle them in to make the timeout activate
	//*
	for {
		err := gi.AddPoint("test", tags, fields, time.Now().UnixNano())
		if err != nil {
			fmt.Printf("Error point: %v\n", err.Error())
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	// */
}
