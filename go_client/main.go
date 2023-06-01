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

	if len(os.Args[1:]) == 0 {

		fmt.Println("No parameter was provided: you can add \"faults\", \"forecasting\", \"current\"")
	} else {
		commandArgs := os.Args[1:]

		if commandArgs[0] == "faults" {
			// This variable will be sent to the RPC server
			var reply RpcResponseStruct

			if len(commandArgs[1:]) > 0 {

				givenFaults := strings.Split(commandArgs[1], ",")

				for _, bucket := range givenFaults {
					// Check over the acquired inputs
					//the last element of this pattern seem to render the controls useless, having a convection on the name of the chassis and the root solves this problem, by modifying the regex
					faultsRegex, err := regexp.Compile("(([0-9]{1,3}.){3}[0-9]{1,3})|(osd.[0-9]{1,2})|(sv[0-9]{1,2})|([a-zA-Z0-9-]*-site)|([a-zA-Z0-9-]*-region)|([a-zA-Z0-9-]*-zone)|([a-zA-Z0-9-]*-rack)|([a-zA-Z0-9]*-[a-zA-Z0-9]*)|.*")

					if !faultsRegex.MatchString(bucket) {
						log.Fatal("Error in bucket acquisition", err)
					}

				}

				args := Args{
					Faults: givenFaults,
				}

				// DialHTTP connects to an HTTP RPC server at the specified network
				client, err := rpc.DialHTTP("tcp", "localhost:1205")
				if err != nil {
					log.Fatal("Client connection error: ", err)
				}

				// Invoke the remote function GiveFaults
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
			} else {
				fmt.Println("No Buckets or Network Addresses were provided")
			}
		}

		if commandArgs[0] == "forecasting" {
			// This variable will be sent to the RPC server
			var reply RpcResponseStruct

			if len(commandArgs[1:]) > 0 {
				givenForecastingTime := commandArgs[1]

				fmt.Printf("Forecasting Time Requested: %s\n", givenForecastingTime)

				// Check over the acquired forecasting time input
				givenFtRegex, err := regexp.Compile("([0-9]{4}-[0-9]{2}-[0-9]{2})T([0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{3}Z)")

				if !givenFtRegex.MatchString(givenForecastingTime) {
					log.Fatal("Error in Forecasting time layout, the correct layout is (yyyy-MM-ddTHH:mm:ss.SSS'Z') e.g \"2006-01-02T15:04:05.000Z\"\n", err)
				}

				//---read osds infos fake data from json---

				//open file
				osds_infos_fake_f, err := os.Open("osds_infos_fake_data.json")
				if err != nil {
					fmt.Println(err)
				}

				defer osds_infos_fake_f.Close()

				//fake data filling
				osdLiftimeInfo := []map[string]string{}

				//read from json file and fill the map
				err = json.NewDecoder(osds_infos_fake_f).Decode(&osdLiftimeInfo)
				if err != nil {
					fmt.Println(err)
				}

				args := Args1{
					ForecastingTime:  givenForecastingTime,
					OsdLifetimeInfos: osdLiftimeInfo,
				}

				// DialHTTP connects to an HTTP RPC server at the specified network
				client, err := rpc.DialHTTP("tcp", "localhost:1205")
				if err != nil {
					log.Fatal("Client connection error: ", err)
				}

				// Invoke the remote function GiveFaultsForecasting
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
			} else {

				fmt.Println("No timestamp was provided (layout:yyyy-MM-ddTHH:mm:ss.SSS'Z') e.g \"2006-01-02T15:04:05.000Z\"")
			}
		}

		if commandArgs[0] == "current" {
			// This variable will be sent to the RPC server
			var reply RpcResponseStruct

			if len(commandArgs[1:]) <= 0 {
				args := Args{}

				// DialHTTP connects to an HTTP RPC server at the specified network
				client, err := rpc.DialHTTP("tcp", "localhost:1205")
				if err != nil {
					log.Fatal("Client connection error: ", err)
				}

				// Invoke the remote function GiveCurrentProbability
				err = client.Call("RpcServer.GiveCurrentProbability", args, &reply)

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

			} else {
				fmt.Println("No parameters required")
			}
		}
	}
}
