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
    init           bool
    pointBatchSize int
    heartBeat      time.Duration

    points   chan *client.Point
    sthapped chan int

    influx   client.Client
    bpConfig client.BatchPointsConfig
    bp       client.BatchPoints
}

// GoInflux interface to unexported type
type GoInflux interface {
    AddPoint(string, TagGroup, FieldGroup, int64) error
    Stahp() error
}

// NewGoInflux basically a constructor
func NewGoInflux(host string, port string, chanBuffer int, pointBatchSize int, heartBeat time.Duration) (GoInflux, error) {
    gi := goInflux{}
    var err error

    gi.init = true

    gi.bpConfig = client.BatchPointsConfig{
        Database:  os.Getenv("INFLUX_DB"),
        Precision: "us",
    }

    gi.heartBeat = heartBeat
    gi.pointBatchSize = pointBatchSize - 1 // so I can use > instead of >=
    gi.points = make(chan *client.Point, chanBuffer)
    gi.sthapped = make(chan int)

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

// NewGoInfluxDefaults constructor that uses env vars instead of parameters
func NewGoInfluxDefaults(chanBuffer int, pointBatchSize int, heartBeat time.Duration) (GoInflux, error) {
    return NewGoInflux(os.Getenv("INFLUX_HOST"), os.Getenv("INFLUX_PORT"), 1024, 1024, 1)
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

func (gi *goInflux) Stahp() error {
    if !gi.init {
        return errors.New("Influx connections not initialized")
    }

    // used to have a gi.sthap channel, but it got too complicated to catch in
    // the select statement, so we will shove a nil pointer into the points
    // channel to make sure all the points before this will get processed
    gi.points <- nil

    // give it a second!
    <-gi.sthapped

    return nil
}

func (gi *goInflux) managePoints() {
    delaySec := gi.heartBeat * time.Second
    delayChan := time.After(delaySec)

    for {
        select {
        case point := <-gi.points:
            if point == nil {
                // shutdown command received
                err := gi.writePoints()
                if err != nil {
                    fmt.Printf("Error in writing points: %v", err.Error())
                    return
                }
                break
            }

            gi.bp.AddPoint(point)
            if len(gi.bp.Points()) > gi.pointBatchSize {
                delayChan = time.After(delaySec)
                err := gi.writePoints()
                if err != nil {
                    fmt.Printf("Error in writing points: %v", err.Error())
                    return
                }
            }
            continue
        case <-delayChan:
            //fmt.Printf("(T)")
            delayChan = time.After(delaySec)
            err := gi.writePoints()
            if err != nil {
                fmt.Printf("Error in writing points: %v", err.Error())
                return
            }
            continue
        }
        break
    }

    //fmt.Println(len(gi.points)) // better be 0
    gi.sthapped <- 1
}

func (gi *goInflux) writePoints() error {
    var err error

    if len(gi.bp.Points()) > 0 {
        //fmt.Printf("%v... ", len(gi.bp.Points()))

        _, _, err = gi.influx.Ping(1 * time.Second)
        if err != nil {
            return err
        }

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
