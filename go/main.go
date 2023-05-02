package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	parser "github.com/MicheleCannizzaro/Aucta-Cognitio-Internship/go/pkg"
	"github.com/MicheleCannizzaro/Aucta-Cognitio-Internship/go/structs"
	utility "github.com/MicheleCannizzaro/Aucta-Cognitio-Internship/go/tools"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
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

var (
	pgDumpOutput  = utility.ReadPgDumpJson("pg_dump.json")
	osdTreeOutput = utility.ReadOsdTreeJson("osd-tree.json")
	osdDumpOutput = utility.ReadOsdDumpJson("osd_dump.json")

	//router declaration
	router *mux.Router

	faults           []string
	forecasting      Forecasting
	osdLifetimeInfos []structs.OsdLifetimeInfo
	osdLifetime      map[string]map[string]interface{}

	//Prometheus Metrics
	datalossProbability = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "dataloss_exporter",
		Name:      "pool_dataloss_probability",
		Help:      "Probability of data loss referred to the pool.",
	}, []string{"pool"})

	oneWeekDataLossForecasting = prometheus.NewGaugeVec(prometheus.GaugeOpts{ //insert a structure with 4 metrics
		Namespace: "dataloss_exporter",
		Name:      "one_week_pool_dataloss_forecasting",
		Help:      "Forecasting of probability of data loss referred to the pool.",
	}, []string{"pool"})

	twoWeeksDataLossForecasting = prometheus.NewGaugeVec(prometheus.GaugeOpts{ //insert a structure with 4 metrics
		Namespace: "dataloss_exporter",
		Name:      "two_weeks_pool_dataloss_forecasting",
		Help:      "Forecasting of probability of data loss referred to the pool.",
	}, []string{"pool"})

	threeWeeksDataLossForecasting = prometheus.NewGaugeVec(prometheus.GaugeOpts{ //insert a structure with 4 metrics
		Namespace: "dataloss_exporter",
		Name:      "three_weeks_pool_dataloss_forecasting",
		Help:      "Forecasting of probability of data loss referred to the pool.",
	}, []string{"pool"})

	fourWeeksDataLossForecasting = prometheus.NewGaugeVec(prometheus.GaugeOpts{ //insert a structure with 4 metrics
		Namespace: "dataloss_exporter",
		Name:      "four_weeks_pool_dataloss_forecasting",
		Help:      "Forecasting of probability of data loss referred to the pool.",
	}, []string{"pool"})
)

func main() {
	//router instantiation
	router = mux.NewRouter()

	//http Servers creation
	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	//prometheus
	prometheus.MustRegister(datalossProbability, oneWeekDataLossForecasting, twoWeeksDataLossForecasting,
		threeWeeksDataLossForecasting, fourWeeksDataLossForecasting)

	//promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	//creation of /metrics endpoint
	router.Handle("/metrics", promhttp.Handler())

	//creation of REST endpoints
	router.HandleFunc("/health", getInfoHealth).Methods("GET")
	router.HandleFunc("/dataloss-prob", datalossProb).Methods("GET")
	router.HandleFunc("/dataloss-prob/component/faults", postFaultsReActive).Methods("POST")
	router.HandleFunc("/dataloss-prob/faults", postFaultsProActive).Methods("POST")
	router.HandleFunc("/dataloss-prob/forecasting", postForecastingProActive).Methods("POST")
	router.HandleFunc("/dataloss-prob/component/forecasting", postForecastingReActive).Methods("POST")

	//gorouting to properly handle server.ListenAndServe error
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

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

// pro-active -> ceph administrator
func getInfoHealth(w http.ResponseWriter, r *http.Request) {
	//convert data to json string
	b, err := json.Marshal("System-Health: true ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //in case of conversion error
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b) //write the data to the connection
}

func datalossProb(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("->%v", Faultsdb)
	//poolDataLossProbability, err := parser.GetPoolDataLossProbability(Faultsdb[len(Faultsdb)-1].Osds, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	//faults.Osds
	poolDataLossProbability, err := parser.GetPoolDataLossProbability(faults, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		return
	}

	//Prometheus metric handling
	pools := parser.GetPools(pgDumpOutput)
	for _, pool := range pools {
		datalossProbability.With(prometheus.Labels{"pool": "Pool: " + pool}).Set(poolDataLossProbability[pool])
	}

	//REST API handling
	b, err := json.Marshal(poolDataLossProbability)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //in case of conversion error
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b) //write the data to the connection
}

// re-active -> python component
func postFaultsReActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewDecoder(r.Body).Decode(&faults)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	// - must update dumps jsons
	fmt.Printf("faultOsds -> %v", faults)
	//pool data loss probability calculation						//faults.Osds
	poolDataLossProbability, err := parser.GetPoolDataLossProbability(faults, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		return
	}

	//Prometheus metric handling/triggering
	pools := parser.GetPools(pgDumpOutput)
	for _, pool := range pools {
		datalossProbability.With(prometheus.Labels{"pool": "Pool: " + pool}).Set(poolDataLossProbability[pool])
	}

	err = json.NewEncoder(w).Encode(&faults)
	if err != nil {
		log.Fatalln("There was an error encoding the initialized struct")
	}
}

// pro-active -> ceph administrator
func postFaultsProActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewDecoder(r.Body).Decode(&faults)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	// - must update dumps jsons

	//pool data loss probability calculation						//faults.Osds
	poolDataLossProbability, err := parser.GetPoolDataLossProbability(faults, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		return
	}

	response := FaultsProActiveResponse{
		Faults:                  faults,
		PoolDatalossProbability: poolDataLossProbability,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		log.Fatalln("There was an error encoding the initialized struct")
	}
}

// pro-active -> ceph administrator
func postForecastingProActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewDecoder(r.Body).Decode(&forecasting)

	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	// - must update dumps jsons
	forecastingTime := forecasting.ForecastingTime //time given by administrator through post
	//forecastingTime := time.Date(2023, time.Month(8), 7, 1, 10, 30, 0, time.UTC)
	fmt.Printf("forecastingTime-> %v\n forecasting.OsdLifetimeInfos-> %v\n", forecastingTime, forecasting.OsdLifetimeInfos)

	osdLifeForecastingMap, warningOsdSlice := parser.RiskFailureForecasting(forecasting.OsdLifetimeInfos, forecastingTime)
	fmt.Printf("osdLifeForecastingMap-> %v\n warningOsdSlice-> %v\n", osdLifeForecastingMap, warningOsdSlice)

	//pool data loss probability calculation on WarningOsdSlice
	poolDataLossProbability, err := parser.GetPoolDataLossProbability(warningOsdSlice, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		return
	}

	response := OsdLifitimeForecastingResponse{
		OsdLifetimeForecasting:  osdLifeForecastingMap,
		WarningOsds:             warningOsdSlice,
		PoolDatalossProbability: poolDataLossProbability,
	}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		log.Fatalln("There was an error encoding the initialized struct")
	}
}

// re-active -> python component
func postForecastingReActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewDecoder(r.Body).Decode(&osdLifetimeInfos)

	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	// - must update dumps jsons

	progressiveWeekDays := []int{7, 14, 21, 28}

	for numWeek, numDays := range progressiveWeekDays {

		daysForwardTime := time.Now().AddDate(0, 0, numDays)

		_, warningOsdSlice := parser.RiskFailureForecasting(osdLifetimeInfos, daysForwardTime)

		//pool data loss probability calculation on WarningOsdSlice
		poolDataLossProbability, err := parser.GetPoolDataLossProbability(warningOsdSlice, pgDumpOutput, osdTreeOutput, osdDumpOutput)
		if err != nil {
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

	err = json.NewEncoder(w).Encode(&osdLifetimeInfos)
	if err != nil {
		log.Fatalln("There was an error encoding the initialized struct")
	}
}
