//// +build integration

package goinflux

import (
	"os"
	"testing"
	"time"
)

func TestInflux(t *testing.T) {
	os.Setenv("INFLUX_HOST", "localhost")
	os.Setenv("INFLUX_PORT", "8086")
	os.Setenv("INFLUX_DB", "testing")

	tags := TagGroup{}
	tags["metric"] = "pulse"

	fields := FieldGroup{}
	fields["value"] = 1

	gi, err := NewGoInflux(os.Getenv("INFLUX_HOST"), os.Getenv("INFLUX_PORT"), 25, 25, 5)
	if err != nil {
		t.Fatalf("Error init: %v\n", err.Error())
		return
	}

	err = gi.AddPointError("test", tags, fields, time.Now().UnixNano())
	if err != nil {
		t.Fatalf("Error point: %v\n", err.Error())
		return
	}

	gi.Stahp()
}
