package goinflux

import (
	"errors"
	"fmt"
	"os"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

//TagGroup are indexable
type TagGroup map[string]string

// FieldGroup are selectable
type FieldGroup map[string]interface{}

type goInflux struct {
	init bool

	points chan *client.Point

	influx   client.Client
	bpConfig client.BatchPointsConfig
	bp       client.BatchPoints
}

// GoInflux interface to unexported type
type GoInflux interface {
	AddPoint(string, TagGroup, FieldGroup, int64) error
}

// NewGoInflux basically a constructor
func NewGoInflux(host string, port string) (GoInflux, error) {
	gi := goInflux{}
	var err error

	gi.init = true

	gi.bpConfig = client.BatchPointsConfig{
		Database:  os.Getenv("INFLUX_DB"),
		Precision: "us",
	}

	gi.points = make(chan *client.Point, 10000)
	gi.bp, err = client.NewBatchPoints(gi.bpConfig)
	if err != nil {
		return &gi, err
	}

	gi.influx, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://" + os.Getenv("INFLUX_HOST") + ":" + os.Getenv("INFLUX_PORT"),
	})
	if err != nil {
		return &gi, err
	}

	go gi.managePoints()

	return &gi, nil
}

// NewGoInfluxEnv constructor that uses env vars instead of parameters
func NewGoInfluxEnv() (GoInflux, error) {
	return NewGoInflux(os.Getenv("INFLUX_HOST"), os.Getenv("INFLUX_PORT"))
}

func (gi *goInflux) AddPoint(measurement string, tags TagGroup, fields FieldGroup, ts int64) error {
	if !gi.init {
		return errors.New("Influx connections not initialized")
	}

	point, err := client.NewPoint(measurement, tags, fields, time.Unix(0, ts))
	if err != nil {
		return err
	}

	gi.points <- point

	return nil
}

func (gi *goInflux) managePoints() {
	delaySec := 10 * time.Second
	delayChan := time.After(delaySec)

	for {
		select {
		case point := <-gi.points:
			gi.bp.AddPoint(point)
			if len(gi.bp.Points()) > 10000 {
				delayChan = time.After(delaySec)
				err := gi.writePoints()
				if err != nil {
					fmt.Printf("Error in writing points: %v", err.Error())
					return
				}
			}
		case <-delayChan:
			fmt.Printf("Time's up!  Writing %v points.\n", len(gi.bp.Points()))
			delayChan = time.After(delaySec)
			err := gi.writePoints()
			if err != nil {
				fmt.Printf("Error in writing points: %v", err.Error())
				return
			}
		}
	}
}

func (gi *goInflux) writePoints() error {
	var err error

	if len(gi.bp.Points()) > 0 {
		//fmt.Printf("Writing %v points.\n", len(gi.bp.Points()))

		err = gi.influx.Write(gi.bp)
		if err != nil {
			return err
		}

		gi.bp, err = client.NewBatchPoints(gi.bpConfig)
		if err != nil {
			return err
		}
	}

	return nil
}
