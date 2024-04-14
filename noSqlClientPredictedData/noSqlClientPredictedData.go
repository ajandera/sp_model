// Package noSqlClientPredictedData to handle influx communication
package noSqlClientPredictedData

import (
	"context"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"time"
)

// ClientData struct to store influx client
type ClientData struct {
	db influxdb2.Client
}

// NewConnect function to connect to influx
func NewConnect(url string, token string) ClientData {
	Client := influxdb2.NewClient(url, token)
	if Client == nil {
		panic("failed to connect influxdb")
	}
	return ClientData{Client}
}

// StoreData function to store data in influx bucket
// bucket is a store id
func (client *ClientData) StoreData(measurement string, dayIndex string, value int,
	setAverageOrderAmount float64, time time.Time, bucket string, org string) (bool, error) {
	buck, err := client.db.BucketsAPI().FindBucketByName(context.Background(), bucket)
	o, err := client.db.OrganizationsAPI().FindOrganizationByName(context.Background(), org)
	if err != nil || buck == nil {
		buck, err = client.db.BucketsAPI().CreateBucketWithNameWithID(context.Background(), *o.Id, bucket)
		if err != nil {
			return false, err
		}
	}

	writeAPI := client.db.WriteAPI(org, bucket)
	p := influxdb2.NewPointWithMeasurement(measurement).
		AddTag("daysToMeasurement", dayIndex).
		AddField("value", value).
		AddField("saoa", setAverageOrderAmount).
		SetTime(time)

	// write point immediately
	writeAPI.WritePoint(p)
	return true, nil
}

// Flush function to flush data for bucket
func (client *ClientData) Flush(bucket string, org string) (bool, error) {

	writeAPI := client.db.WriteAPI(org, bucket)

	// Force all unwritten data to be sent
	writeAPI.Flush()

	return true, nil
}

// GetData function to get predicted data by query
func (client *ClientData) GetData(query string, org string) (string, error) {
	// Get query client
	queryAPI := client.db.QueryAPI(org)

	// Query and get complete result as a string
	// Use default dialect
	result, err := queryAPI.QueryRaw(context.Background(), query, influxdb2.DefaultDialect())

	// Ensures background processes finishes
	client.db.Close()

	return result, err
}

// GetQuery function to get raw data from influx
func (client *ClientData) GetQuery(query string, org string) (*api.QueryTableResult, error) {

	// Get query client
	queryAPI := client.db.QueryAPI(org)

	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), query)

	// Ensures background processes finishes
	client.db.Close()

	return result, err
}
