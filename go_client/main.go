package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type FaultsProActiveResponse struct {
	Faults                  []string           `json:"faults"`
	PoolDatalossProbability map[string]float64 `json:"pool_dataloss_probability"`
}

type OsdLifitimeForecastingResponse struct {
	OsdLifetimeForecasting  map[string]float64 `json:"osd_lifetime_forecasting"`
	WarningOsds             []string           `json:"warning_osds"`
	PoolDatalossProbability map[string]float64 `json:"pool_dataloss_probability_forecasting"`
}

type RpcResponseStruct struct {
	FaultsProActiveResponse        FaultsProActiveResponse
	OsdLifitimeForecastingResponse OsdLifitimeForecastingResponse
}

type Args struct {
	Faults []string
}

type Args1 struct {
	ForecastingTime  string
	OsdLifetimeInfos []map[string]string
}

func main() {
	commandArgs := os.Args[1:]

	if commandArgs[0] == "faults" {
		// Address to this variable will be sent to the RPC server
		var reply RpcResponseStruct

		// givenFaults := []string{"osd.1", "osd.8"}
		givenFaults := strings.Split(commandArgs[1], ",")

		for _, bucket := range givenFaults {
			// Check over the acquired inputs
			faultsRegex, err := regexp.Compile("(([0-9]{1,3}.){3}[0-9]{1,3})|(osd.[0-9]{1,2})|(sv[0-9]{1,2})")

			if !faultsRegex.MatchString(bucket) {
				log.Fatal("Error in bucket acquisition", err)
			}

		}

		args := Args{
			Faults: givenFaults,
		}

		// DialHTTP connects to an HTTP RPC server at the specified network
		client, err := rpc.DialHTTP("tcp", "localhost:1025")
		if err != nil {
			log.Fatal("Client connection error: ", err)
		}

		// Invoke the remote function GiveFaults attached to TimeServer pointer
		err = client.Call("RpcServer.GiveFaults", args, &reply)

		if err != nil {
			log.Fatal("Client invocation error: ", err)
		}

		if !reflect.DeepEqual(reply.FaultsProActiveResponse, FaultsProActiveResponse{}) {

			// Print the reply from the server
			buffer := new(strings.Builder)
			err = json.NewEncoder(buffer).Encode(&reply.FaultsProActiveResponse)

			if err != nil {
				log.Fatalln("There was an error decoding the reply")
			}

			log.Printf("%v", buffer.String())
		}
	}

	if commandArgs[0] == "forecasting" {
		// Address to this variable will be sent to the RPC server
		var reply RpcResponseStruct

		givenForecastingTime := commandArgs[1]

		fmt.Printf("Forecasting Time Requested: %s\n", givenForecastingTime)

		// Check over the acquired forecasting time input
		givenFtRegex, err := regexp.Compile("([0-9]{4}-[0-9]{2}-[0-9]{2})T([0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{3}Z)")

		if !givenFtRegex.MatchString(givenForecastingTime) {
			log.Fatal("Error in Forecasting time layout, the correct layout is: \"2006-01-02T15:04:05.000Z\"\n", err)
		}

		//fake data filling
		osdLiftimeInfo := []map[string]string{}
		m1 := map[string]string{ //to fix - (read from a file)
			"osd_name":             "osd.1",
			"current_osd_lifetime": "30.0",
			"initiation_date":      "2018-10-14T02:53:00.000Z",
		}
		osdLiftimeInfo = append(osdLiftimeInfo, m1)

		m2 := map[string]string{
			"osd_name":             "osd.2",
			"current_osd_lifetime": "80.0",
			"initiation_date":      "2020-10-14T02:53:00.000Z",
		}
		osdLiftimeInfo = append(osdLiftimeInfo, m2)

		m3 := map[string]string{
			"osd_name":             "osd.3",
			"current_osd_lifetime": "69.0",
			"initiation_date":      "2023-02-14T02:53:00.000Z",
		}
		osdLiftimeInfo = append(osdLiftimeInfo, m3)

		args := Args1{
			ForecastingTime:  givenForecastingTime,
			OsdLifetimeInfos: osdLiftimeInfo,
		}

		// DialHTTP connects to an HTTP RPC server at the specified network
		client, err := rpc.DialHTTP("tcp", "localhost:1025")
		if err != nil {
			log.Fatal("Client connection error: ", err)
		}

		// Invoke the remote function GiveFaultsForecasting attached to TimeServer pointer
		err = client.Call("RpcServer.GiveFaultsForecasting", args, &reply)
		if err != nil {
			log.Fatal("Client invocation error: ", err)
		}

		if !reflect.DeepEqual(reply.OsdLifitimeForecastingResponse, OsdLifitimeForecastingResponse{}) {
			buffer := new(strings.Builder)
			err = json.NewEncoder(buffer).Encode(&reply.OsdLifitimeForecastingResponse)
			if err != nil {
				log.Fatalln("There was an error decoding the reply")
			}

			log.Printf("%v", buffer.String())
		}
	}
}