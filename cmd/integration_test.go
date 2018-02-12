// +build integration

package main

import (
	"os"
	"testing"
)

func TestWritePoints(t *testing.T) {
	os.Setenv("INFLUX_HOST", "localhost")
	os.Setenv("INFLUX_PORT", "8086")
	os.Setenv("INFLUX_DB", "testing")

	writePoints(10, 32, 32, 0, "pulse")

	writePoints(10, 32, 32, 5, "pulse")
}
