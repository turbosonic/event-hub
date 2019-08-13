package influxdb

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/turbosonic/event-hub/src/logging"
)

var (
	influxDBNAME string
	wg           sync.WaitGroup
	bp           client.BatchPoints
)

type influxdbLogger struct {
	client client.Client
}

func New() influxdbLogger {
	influxDBURL := os.Getenv("LOGGING_INFLUXDB_URL")
	if influxDBURL == "" {
		panic("No LOGGING_INFLUXDB_URL environment variable found")
	}

	influxDBNAME = os.Getenv("LOGGING_INFLUX_DB_NAME")
	if influxDBNAME == "" {
		panic("No LOGGING_INFLUX_DB_NAME environment variable found")
	}

	// create a client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: influxDBURL,
	})
	if err != nil {
		panic("Error creating InfluxDB Client")
	}

	// need to create the database here, unless it already exists
	q := client.NewQuery("CREATE DATABASE "+influxDBNAME, "", "")
	if response, err := c.Query(q); err != nil && response.Error() != nil {
		fmt.Println(response.Error())
	}

	bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxDBNAME,
		Precision: "s",
	})

	logger := influxdbLogger{
		client: c,
	}

	logger.startDataWriter()

	return logger
}

func (ifl influxdbLogger) LogEvent(l *logging.EventLog) {
	tags := map[string]string{
		"ID":        l.ID,
		"Name":      l.Name,
		"Source":    l.Source,
		"Valid":     strconv.FormatBool(l.Valid),
		"RequestID": l.RequestID,
	}
	fields := map[string]interface{}{
		"Events": 1,
	}
	pt, err := client.NewPoint("event", tags, fields, l.Timestamp)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	wg.Wait()
	bp.AddPoint(pt)
}

func (influxdb influxdbLogger) startDataWriter() {
	go func() {
		for {
			time.Sleep(time.Second * 5)
			wg.Add(1)
			pointsCount := len(bp.Points())
			if pointsCount > 0 {
				influxdb.client.Write(bp)
				bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
					Database:  influxDBNAME,
					Precision: "s",
				})

				fmt.Printf("%d points sent to InfluxDB\n", pointsCount)
			}
			wg.Done()
		}
	}()
}
