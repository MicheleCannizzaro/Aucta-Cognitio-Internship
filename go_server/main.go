package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	parser "github.com/MicheleCannizzaro/Aucta-Cognitio-Internship/go_server/pkg"
	"github.com/MicheleCannizzaro/Aucta-Cognitio-Internship/go_server/structs"
	utility "github.com/MicheleCannizzaro/Aucta-Cognitio-Internship/go_server/tools"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
)

var (
	pgDumpOutput  = utility.ReadPgDumpJson("pg_dump.json")
	osdTreeOutput = utility.ReadOsdTreeJson("osd-tree.json")
	osdDumpOutput = utility.ReadOsdDumpJson("osd_dump.json")

	regexPattern = "(([0-9]{1,3}.){3}[0-9]{1,3})|(osd.[0-9]{1,2})|(sv[0-9]{1,2})|([a-zA-Z0-9-]*-site)|([a-zA-Z0-9-]*-region)|([a-zA-Z0-9-]*-zone)|([a-zA-Z0-9-]*-rack)|([a-zA-Z0-9]*-[a-zA-Z0-9]*)"

	//router declaration
	router *mux.Router

	faults           []string
	forecasting      Forecasting
	osdLifetimeInfos []structs.OsdLifetimeInfo

	//Prometheus Metrics
	datalossProbability = prometheus.NewGaugeVec(prometheus.GaugeOpts{ //Prometheus DataLossMetric
		Namespace: "dataloss_exporter",
		Name:      "pool_dataloss_probability",
		Help:      "Probability of data loss referred to the pool.",
	}, []string{"pool"})

	oneWeekDataLossForecasting = prometheus.NewGaugeVec(prometheus.GaugeOpts{ //Prometheus oneWeekMetric
		Namespace: "dataloss_exporter",
		Name:      "one_week_pool_dataloss_forecasting",
		Help:      "Forecasting of probability of data loss referred to the pool.",
	}, []string{"pool"})

	twoWeeksDataLossForecasting = prometheus.NewGaugeVec(prometheus.GaugeOpts{ //Prometheus twoWeeksMetric
		Namespace: "dataloss_exporter",
		Name:      "two_weeks_pool_dataloss_forecasting",
		Help:      "Forecasting of probability of data loss referred to the pool.",
	}, []string{"pool"})

	threeWeeksDataLossForecasting = prometheus.NewGaugeVec(prometheus.GaugeOpts{ //Prometheus threeWeeksMetric
		Namespace: "dataloss_exporter",
		Name:      "three_weeks_pool_dataloss_forecasting",
		Help:      "Forecasting of probability of data loss referred to the pool.",
	}, []string{"pool"})

	fourWeeksDataLossForecasting = prometheus.NewGaugeVec(prometheus.GaugeOpts{ //Prometheus fourWeeksMetric
		Namespace: "dataloss_exporter",
		Name:      "four_weeks_pool_dataloss_forecasting",
		Help:      "Forecasting of probability of data loss referred to the pool.",
	}, []string{"pool"})
)

type Forecasting struct {
	ForecastingTime  time.Time                 `json:"forecasting_time"`
	OsdLifetimeInfos []structs.OsdLifetimeInfo `json:"osds_lifetime_infos"`
}

type FaultsProActiveResponse struct {
	Faults                  []string           `json:"faults"`
	PoolDatalossProbability map[string]float64 `json:"pool_dataloss_probability"`
}

type OsdLifitimeForecastingResponse struct {
	OsdLifetimeForecasting  map[string]float64 `json:"osd_lifetime_forecasting"`
	WarningOsds             []string           `json:"warning_osds"`
	PoolDatalossProbability map[string]float64 `json:"pool_dataloss_probability_forecasting"`
}

type Args struct {
	Faults []string
}

type Args1 struct {
	ForecastingTime  string
	OsdLifetimeInfos []map[string]string //cambia
}

type RpcResponseStruct struct {
	FaultsProActiveResponse        FaultsProActiveResponse
	OsdLifitimeForecastingResponse OsdLifitimeForecastingResponse
}

type RpcServer RpcResponseStruct

// --- RPC methods ---
func (f *RpcServer) GiveFaults(args *Args, reply *RpcResponseStruct) error {
	fmt.Println("RPC Server: GiveFaults requested")

	//update pgDumpOutput, osdTreeOutput and osdDumpOutput
	clusterStatsGathering()

	for _, bucket := range args.Faults {
		// Check over the acquired inputs
		faultsRegex, err := regexp.Compile(regexPattern)
		if err != nil {
			return errors.New("error in regex compile")
		}

		if !faultsRegex.MatchString(bucket) {
			return errors.New("error in given faults %s")
		}

	}

	poolDataLossProbability, err := parser.GetPoolDataLossProbability(args.Faults, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		return errors.New("error in data loss probability calculation - check if bucket name exists")
	}

	response := FaultsProActiveResponse{
		Faults:                  args.Faults,
		PoolDatalossProbability: poolDataLossProbability,
	}

	res := RpcResponseStruct{
		FaultsProActiveResponse:        response,
		OsdLifitimeForecastingResponse: OsdLifitimeForecastingResponse{}, //zero value
	}

	fmt.Println("------------------------------------------")

	*reply = res
	return nil
}

func (f *RpcServer) GiveFaultsForecasting(args *Args1, reply *RpcResponseStruct) error {
	fmt.Println("RPC Server: GiveFaultsForecasting requested")

	//update pgDumpOutput, osdTreeOutput and osdDumpOutput
	clusterStatsGathering()

	osdLifetimeInfos := []structs.OsdLifetimeInfo{}

	//each arg is a map[string]string
	for _, arg := range args.OsdLifetimeInfos {

		currOsdLifetime, err := strconv.ParseFloat(arg["current_osd_lifetime"], 64)
		if err != nil {
			return errors.New("error in string to float64 conversion - currentOsdLifetime")
		}

		layout := "2006-01-02T15:04:05.000Z"
		initiationTime, err := time.Parse(layout, arg["initiation_date"])
		if err != nil {
			return errors.New("error in string to time.Time conversion - initiation_date")
		}

		osdLifetimeInfo := structs.OsdLifetimeInfo{
			OsdName:            arg["osd_name"],
			InitiationDate:     initiationTime,
			CurrentOsdLifetime: currOsdLifetime,
		}

		osdLifetimeInfos = append(osdLifetimeInfos, osdLifetimeInfo)
	}

	layout := "2006-01-02T15:04:05.000Z"
	forecastingTime, err := time.Parse(layout, args.ForecastingTime)
	if err != nil {
		return errors.New("error in string to time.Time conversion - initiation_date")
	}

	fmt.Printf("forecastingTime-> %v\n", forecastingTime)

	osdLifeForecastingMap, warningOsdSlice := parser.RiskFailureForecasting(osdLifetimeInfos, forecastingTime)
	fmt.Printf("osdLifeForecastingMap-> %v\nWarningOsdSlice-> %v\n", osdLifeForecastingMap, warningOsdSlice)

	//pool data loss probability calculation on WarningOsdSlice
	poolDataLossProbability, err := parser.GetPoolDataLossProbability(warningOsdSlice, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		return errors.New("error in data loss probability calculation - check if bucket name exists")
	}

	fmt.Printf("poolDataLossProbabilityMap-> %v\n", poolDataLossProbability)

	response := OsdLifitimeForecastingResponse{
		OsdLifetimeForecasting:  osdLifeForecastingMap,
		WarningOsds:             warningOsdSlice,
		PoolDatalossProbability: poolDataLossProbability,
	}

	res := RpcResponseStruct{
		FaultsProActiveResponse:        FaultsProActiveResponse{},
		OsdLifitimeForecastingResponse: response, //zero value
	}

	fmt.Println("------------------------------------------")

	*reply = res
	return nil
}

// --- REST API Endpoint Handlers --- pro-active (ceph administrator)->

// endpoint: /health
func getInfoHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REST Server: Get request arrived\n endpoint: /health")

	b, err := json.Marshal("System-Health: true ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //in case of conversion error
		log.Fatalln("There was an error in Marshalling request /Health")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b) //write the data to the connection

	fmt.Println("------------------------------------------")
}

// endpoint: /dataloss-prob/faults
func postFaultsProActive(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REST Server: Post request arrived\n endpoint: /dataloss-prob/faults")

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&faults)

	if err != nil {
		fmt.Println("Bad Request - There was an error in decoding the request body")

		response := "There was an error in decoding the request body"
		w.WriteHeader(http.StatusBadRequest)

		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			log.Fatalln("There was an error encoding the initialized struct")
		}
		fmt.Println("------------------------------------------")
		return
	}

	for _, bucket := range faults {
		// Check over the acquired inputs
		faultsRegex, err := regexp.Compile(regexPattern)
		if err != nil {
			log.Fatalln("There was an error in Regex compile")
			return
		}

		if !faultsRegex.MatchString(bucket) {
			fmt.Println("Bad Request - Error in requested data")

			response := "Error in request data no correct bucket form"
			w.WriteHeader(http.StatusBadRequest)

			err = json.NewEncoder(w).Encode(&response)
			if err != nil {
				log.Fatalln("There was an error encoding the initialized struct")
			}

			fmt.Println("------------------------------------------")
			return
		}

	}

	//pg_dump, osd_dump, osd-tree
	clusterStatsGathering()

	//pool data loss probability calculation
	poolDataLossProbability, err := parser.GetPoolDataLossProbability(faults, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		fmt.Println("Internal Server Error - Error calculating poolDataLossProbability  - check if bucket name exists")

		response := "Error in request data"
		w.WriteHeader(http.StatusInternalServerError)

		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			log.Fatalln("There was an error encoding the response")
		}

		fmt.Println("------------------------------------------")
		return
	}

	response := FaultsProActiveResponse{
		Faults:                  faults,
		PoolDatalossProbability: poolDataLossProbability,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		log.Fatalln("There was an error encoding the response")
	}

	fmt.Println("------------------------------------------")
}

// endpoint: /dataloss-prob/forecasting
func postForecastingProActive(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REST Server: Post request arrived\n endpoint: /dataloss-prob/forecasting")

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&forecasting)

	if err != nil {
		fmt.Println("Bad Request - There was an error in decoding the request body")

		response := "There was an error in decoding the request body"
		w.WriteHeader(http.StatusBadRequest)

		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			log.Fatalln("There was an error encoding the initialized struct")
		}
		fmt.Println("------------------------------------------")
		return
	}

	for _, info := range forecasting.OsdLifetimeInfos {
		// Check on OsdName field
		faultsRegex, err := regexp.Compile(regexPattern)
		if err != nil {
			log.Fatalln("There was an error in Regex compiling")
			return
		}

		if !faultsRegex.MatchString(info.OsdName) {
			fmt.Println("Bad Request - Error in requested data - osd_name")

			response := "Error in request data - osd_name"
			w.WriteHeader(http.StatusBadRequest)

			err = json.NewEncoder(w).Encode(&response)
			if err != nil {
				log.Fatalln("There was an error encoding the initialized struct")
			}

			fmt.Println("------------------------------------------")
			return
		}

	}

	//pg_dump, osd_dump, osd-tree
	clusterStatsGathering()

	forecastingTime := forecasting.ForecastingTime //time given by administrator through post
	fmt.Printf("ForecastingTime-> %v\n", forecastingTime)

	osdLifeForecastingMap, warningOsdSlice := parser.RiskFailureForecasting(forecasting.OsdLifetimeInfos, forecastingTime)
	fmt.Printf("OsdLifeForecastingMap-> %v\nWarningOsdSlice-> %v\n", osdLifeForecastingMap, warningOsdSlice)

	//pool data loss probability calculation on WarningOsdSlice
	poolDataLossProbability, err := parser.GetPoolDataLossProbability(warningOsdSlice, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		fmt.Println("Internal Server Error - Error calculating poolDataLossProbability  - check if bucket name exists")

		response := "Error in request data"
		w.WriteHeader(http.StatusInternalServerError)

		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			log.Fatalln("There was an error encoding the response")
		}

		fmt.Println("------------------------------------------")
		return
	}

	response := OsdLifitimeForecastingResponse{
		OsdLifetimeForecasting:  osdLifeForecastingMap,
		WarningOsds:             warningOsdSlice,
		PoolDatalossProbability: poolDataLossProbability,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		log.Fatalln("There was an error encoding the initialized struct")
	}

	fmt.Println("------------------------------------------")
}

// --- REST API Endpoint Handlers --- re-active  (python component) ->

// endpoint: /dataloss-prob/component/faults
func postFaultsReActive(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REST Server: Post request arrived\n endpoint: /dataloss-prob/component/faults")

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&faults)
	if err != nil {
		fmt.Println("Bad Request - There was an error in decoding the request body")

		response := "There was an error in decoding the request body"
		w.WriteHeader(http.StatusBadRequest)

		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			log.Fatalln("There was an error encoding the initialized struct")
		}
		fmt.Println("------------------------------------------")
		return
	}

	for _, bucket := range faults {
		// Check over the acquired inputs
		faultsRegex, err := regexp.Compile("(([0-9]{1,3}.){3}[0-9]{1,3})|(osd.[0-9]{1,2})|(sv[0-9]{1,2})|([a-zA-Z0-9-]*-site)|([a-zA-Z0-9-]*-region)|([a-zA-Z0-9-]*-zone)|([a-zA-Z0-9-]*-rack)|([a-zA-Z0-9]*-[a-zA-Z0-9]*)")
		if err != nil {
			log.Fatalln("There was an error in Regex compile")
			return
		}

		if !faultsRegex.MatchString(bucket) {
			fmt.Println("Bad Request - Error in requested data")

			response := "Error in request data no correct bucket form"
			w.WriteHeader(http.StatusBadRequest)

			err = json.NewEncoder(w).Encode(&response)
			if err != nil {
				log.Fatalln("There was an error encoding the initialized struct")
			}

			fmt.Println("------------------------------------------")
			return
		}

	}

	//pg_dump, osd_dump, osd-tree
	clusterStatsGathering()

	//pool data loss probability calculation
	poolDataLossProbability, err := parser.GetPoolDataLossProbability(faults, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		fmt.Println("Internal Server Error - Error calculating poolDataLossProbability - check if bucket name exists")

		response := "Error in request data"
		w.WriteHeader(http.StatusInternalServerError)

		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			log.Fatalln("There was an error encoding the response")
		}

		fmt.Println("------------------------------------------")
		return
	}

	//Prometheus metric handling/triggering
	pools := parser.GetPools(pgDumpOutput)
	for _, pool := range pools {
		datalossProbability.With(prometheus.Labels{"pool": "Pool: " + pool}).Set(poolDataLossProbability[pool])
	}

	w.WriteHeader(http.StatusOK)

	//err = json.NewEncoder(w).Encode(&faults)
	err = json.NewEncoder(w).Encode("PoolDataLossProbability metric sent to Prometheus endpoint")
	if err != nil {
		log.Fatalln("There was an error encoding the initialized struct")
	}

	fmt.Println("------------------------------------------")
}

// endpoint: /dataloss-prob/component/forecasting
func postForecastingReActive(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REST Server: Post request arrived\n endpoint: /dataloss-prob/component/forecasting")

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&osdLifetimeInfos)

	if err != nil {
		fmt.Println("Bad Request - There was an error in decoding the request body")

		response := "There was an error in decoding the request body"
		w.WriteHeader(http.StatusBadRequest)

		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			log.Fatalln("There was an error encoding the initialized struct")
		}
		fmt.Println("------------------------------------------")
		return
	}

	for _, info := range osdLifetimeInfos {
		// Check on OsdName field
		faultsRegex, err := regexp.Compile(regexPattern)
		if err != nil {
			log.Fatalln("There was an error in Regex compiling")
			return
		}

		if !faultsRegex.MatchString(info.OsdName) {
			fmt.Println("Bad Request - Error in requested data - osd_name")

			response := "Error in request data - osd_name"
			w.WriteHeader(http.StatusBadRequest)

			err = json.NewEncoder(w).Encode(&response)
			if err != nil {
				log.Fatalln("There was an error encoding the initialized struct")
			}

			fmt.Println("------------------------------------------")
			return
		}

	}

	//pg_dump, osd_dump, osd-tree
	clusterStatsGathering()

	progressiveWeekDays := []int{7, 14, 21, 28}

	for numWeek, numDays := range progressiveWeekDays {

		fmt.Printf("Forecasting week n.%d\n", numWeek+1)

		daysForwardTime := time.Now().AddDate(0, 0, numDays)

		_, warningOsdSlice := parser.RiskFailureForecasting(osdLifetimeInfos, daysForwardTime)

		//pool data loss probability calculation on WarningOsdSlice
		poolDataLossProbability, err := parser.GetPoolDataLossProbability(warningOsdSlice, pgDumpOutput, osdTreeOutput, osdDumpOutput)
		if err != nil {
			fmt.Println("Internal Server Error - Error calculating poolDataLossProbability - check if bucket name exists")

			response := "Error in request data"
			w.WriteHeader(http.StatusInternalServerError)

			err = json.NewEncoder(w).Encode(&response)
			if err != nil {
				log.Fatalln("There was an error encoding the response")
			}

			fmt.Println("------------------------------------------")
			return
		}

		pools := parser.GetPools(pgDumpOutput)

		//Prometheus metric handling/triggering
		switch numWeek {
		case 0:
			for _, pool := range pools {
				oneWeekDataLossForecasting.With(prometheus.Labels{"pool": "Pool: " + pool}).Set(poolDataLossProbability[pool])
			}
		case 1:
			for _, pool := range pools {
				twoWeeksDataLossForecasting.With(prometheus.Labels{"pool": "Pool: " + pool}).Set(poolDataLossProbability[pool])
			}
		case 2:
			for _, pool := range pools {
				threeWeeksDataLossForecasting.With(prometheus.Labels{"pool": "Pool: " + pool}).Set(poolDataLossProbability[pool])
			}
		case 3:
			for _, pool := range pools {
				fourWeeksDataLossForecasting.With(prometheus.Labels{"pool": "Pool: " + pool}).Set(poolDataLossProbability[pool])
			}
		}
	}

	w.WriteHeader(http.StatusOK)

	//err = json.NewEncoder(w).Encode(&osdLifetimeInfos)
	err = json.NewEncoder(w).Encode("Metrics '(one/two/three/four)WeeksDataLossForecasting' sent to Prometheus endpoint")
	if err != nil {
		log.Fatalln("There was an error encoding the initialized struct")
	}

	fmt.Println("------------------------------------------")
}

// Stats updater function
func clusterStatsGathering() error {

	f, err := os.Create("osd_dump.json")
	if err != nil {
		//fmt.Println(err)
		log.Fatal("Error in osd_dump.json creation", err)
	}
	defer f.Close()

	msg, err := exec.Command("sudo", "ceph", "osd", "dump", "--format=json").Output()
	if err != nil {
		//fmt.Println(err)
		log.Fatal("Error in executing \"osd dump --forma=json\" command", err)
	}

	f.Write(msg)

	f1, err := os.Create("pg_dump.json")
	if err != nil {
		//fmt.Println(err)
		log.Fatal("Error in pg_dump.json creation", err)
	}
	defer f1.Close()

	msg1, err := exec.Command("sudo", "ceph", "pg", "dump", "--format=json").Output()
	if err != nil {
		//fmt.Println(err)
		log.Fatal("Error in executing \"pg dump --format=json\" command", err)
	}

	f1.Write(msg1)

	f2, err := os.Create("osd-tree.json")
	if err != nil {
		//fmt.Println(err)
		log.Fatal("Error in osd_tree.json creation", err)
	}
	defer f2.Close()

	msg2, err := exec.Command("sudo", "ceph", "osd", "tree", "--format=json").Output()
	if err != nil {
		//fmt.Println(err)
		log.Fatal("Error in executing \"osd tree --format=json\" command", err)
	}

	f2.Write(msg2)

	return nil
}

func main() {
	//---------------------------------------------------------------------------------------------------------
	//                                             RPC SERVER CODE

	rpcServer := new(RpcServer)
	rpc.Register(rpcServer)

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", "localhost:1025")
	if err != nil {
		log.Fatal("Listener error: ", err)
	}

	go http.Serve(listener, nil)
	fmt.Println("Started RPC server on localhost:1205")

	//---------------------------------------------------------------------------------------------------------
	//											REST API SERVER CODE

	//router instantiation
	router = mux.NewRouter()

	//http Servers creation
	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	//Prometheus metrics registering
	prometheus.MustRegister(datalossProbability, oneWeekDataLossForecasting, twoWeeksDataLossForecasting,
		threeWeeksDataLossForecasting, fourWeeksDataLossForecasting)

	//creation of /metrics endpoint
	router.Handle("/metrics", promhttp.Handler())

	//creation of REST endpoints
	router.HandleFunc("/health", getInfoHealth).Methods("GET")

	router.HandleFunc("/dataloss-prob/component/faults", postFaultsReActive).Methods("POST")
	router.HandleFunc("/dataloss-prob/faults", postFaultsProActive).Methods("POST")

	router.HandleFunc("/dataloss-prob/forecasting", postForecastingProActive).Methods("POST")
	router.HandleFunc("/dataloss-prob/component/forecasting", postForecastingReActive).Methods("POST")

	//gorouting to properly handle server.ListenAndServe error
	go func() {
		fmt.Println("Started REST API server on localhost:8081")
		fmt.Println("Listening ...")
		fmt.Println("------------------------------------------")

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	fmt.Println("Listening ...")

	//Chanel creation of buffer_size = 1 to receive os.Signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
		//server.Close()  //to force the server to close if it doesn't close gracefully
	}
	log.Println("Graceful shutdown complete.")
}
