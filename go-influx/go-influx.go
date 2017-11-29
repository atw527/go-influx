package goinflux

import (
	"fmt"
	"os"
	"time"

	"github.com/atw527/intelligence/shared"
	client "github.com/influxdata/influxdb/client/v2"
)

//TagGroup are indexable
type TagGroup map[string]string

// FieldGroup are selectable
type FieldGroup map[string]interface{}

type goInflux struct {
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
func NewGoInflux() (GoInflux, error) {
	gi := goInflux{}
	var err error

	gi.bpConfig = client.BatchPointsConfig{
		Database:  os.Getenv("INFLUX_DB"),
		Precision: "us",
	}

	gi.points = make(chan *client.Point, 1000)
	gi.bp, err = client.NewBatchPoints(gi.bpConfig)
	if err != nil {
		return &gi, err
	}

	gi.influx, err = shared.GetInfluxEnv()
	if err != nil {
		return &gi, err
	}

	go gi.managePoints()

	return &gi, nil
}

func (gi *goInflux) AddPoint(measurement string, tags TagGroup, fields FieldGroup, ts int64) error {
	point, err := client.NewPoint(measurement, tags, fields, time.Unix(0, ts))
	if err != nil {
		return err
	}

	gi.points <- point

	return nil
}

func (gi *goInflux) managePoints() {
	for {
		delay := time.After(10 * time.Second)

		select {
		case point := <-gi.points:
			gi.bp.AddPoint(point)
			if len(gi.bp.Points()) > 1000 {
				err := gi.writePoints()
				if err != nil {
					fmt.Printf("Error in writing points: %v", err.Error())
					return
				}
			}
		case <-delay:
			fmt.Println("Time's up!")
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
