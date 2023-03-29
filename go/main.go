package main

import (
	"fmt"
	parser "go/pkg"
	utility "go/tools"
	"io"
	"log"
	"strconv"
	"time"
)

func main() {
	log.SetOutput(io.Discard) //disable log

	//UnMarshalling Json
	//pgDumpOutput := utility.ReadPgDumpJson("pg_dump.json")
	//osdTreeOutput := utility.ReadOsdTreeJson("osd-tree.json")
	//osdDumpOutput := utility.ReadOsdDumpJson("osd_dump.json")

	// fmt.Println("------------------------------------------")
	// fmt.Printf("  PG-DUMP JSON Information Gathering\n\n")
	// //---getting different pg_ids---
	// pgSlice := parser.GetPgs(pgDumpOutput)
	// fmt.Printf("#Different pgids  -> %d\n\n", len(pgSlice))
	// fmt.Println("------------------------------------------")

	// //---getting pools---
	// poolSlice := parser.GetPools(pgDumpOutput)
	// fmt.Printf("Pools -> %s\n", poolSlice)
	// fmt.Printf("#Pools -> %d\n\n", len(poolSlice))
	// fmt.Println("------------------------------------------")

	// //---getting_osds---
	// osds := parser.GetOsds(pgDumpOutput)
	// fmt.Printf("OSDs ->%d\n", osds)
	// fmt.Printf("#OSDs ->%d\n\n", len(osds))
	// fmt.Println("------------------------------------------")

	// //---getting_pgs_details---
	// osdPgMap := parser.GetOsdPgMap(pgDumpOutput)
	// osdPgNumberMap := parser.GetNumberOfAssociatedPgsPerOsdMap(osdPgMap)
	// fmt.Printf("Number of mapped PGs per OSD ->%d\n\n", osdPgNumberMap)
	// fmt.Printf("#PGs ->%d\n\n", parser.GetTotalNumberOfPgs(osdPgNumberMap))
	// fmt.Println("------------------------------------------")

	// //---getting_pgIdOsdMap---
	// fmt.Println("-> pgIdOsdMap")
	// pgIdOsdMap := parser.GetPgOsdMap(pgDumpOutput)
	// fmt.Printf("pgid 31.0 is on these osds -> %d\n\n", pgIdOsdMap["31.0"])
	// fmt.Println("------------------------------------------")

	// //---osd_pgs_mapping---
	// fmt.Println("-> osd_pgs_mapping")
	// numberOfPg := parser.GetNumberOfAssociatedPgsPerOsdMap(osdPgMap)[3]
	// fmt.Printf("OSD 3 contains these #%d PGs-> %s\n\n", numberOfPg, osdPgMap[3])
	// fmt.Println("------------------------------------------")

	// //------osd_pool_pgs_map------
	// //	   {(osd,pool):[pgs]}
	// fmt.Println("-> osd_pool_pgs_map   {(osd,pool):[pgs]}")
	// osdPoolPgMap := parser.GetOsdPoolPgMap(pgDumpOutput)
	// fmt.Printf("OSD 6 pool 31 contains these PGs -> %s\n\n", osdPoolPgMap[parser.OsdPoolTuple{Osd: 6, Pool: "31"}])

	// //	   {(osd,pool):{pg:number_replicas},...}
	// fmt.Println("-> osd_pool_pgs_map   {(osd,pool):{pg:number_replicas},...}")
	// osdPollNumberPerPgsMap := parser.GetOsdPoolNumberPerPgsMap(pgDumpOutput)
	// fmt.Printf("OSD 6 pool 31 contains these PGs -> %v\n\n", osdPollNumberPerPgsMap[parser.OsdPoolTuple{Osd: 6, Pool: "31"}])
	// fmt.Println("------------------------------------------")

	// //---osds_containing_pool---
	// fmt.Println("->   Pool->OSDs")
	// pool := "3"
	// osdsContaingPool := parser.GetOsdsContainingPool(pool, pgDumpOutput)
	// fmt.Printf("Pool: %s is spread between these OSDs-> %v\n\n", pool, osdsContaingPool)
	// fmt.Println("------------------------------------------")

	// //---osds_containing_pg---
	// fmt.Println("->   PG->OSDs")
	// givenPg := "19.1f"
	// osdsContainingPg := parser.GetOsdsContainingPg(givenPg, pgDumpOutput)
	// fmt.Printf("PG: %s is spread between these OSDs-> %d\n\n", givenPg, osdsContainingPg)
	// fmt.Println("------------------------------------------")

	// //---affected pool for osd crush---
	// fmt.Println("-> Affected pools")
	// faultOsd := 2
	// affectedPools := parser.GetAffectedPools(faultOsd, pgDumpOutput)
	// fmt.Printf("Affected pools for crush on OSD: %d  -> #%d ->  %v\n\n", faultOsd, len(affectedPools), affectedPools)
	// fmt.Println("------------------------------------------")

	// //---affected pg for osd crush---
	// fmt.Println("-> Affected PGs")
	// affectedPgs := parser.GetAffectedPgs(faultOsd, pgDumpOutput)
	// fmt.Printf("Affected pgid for crush on OSD: %d -> %s\n\n", faultOsd, affectedPgs)
	// fmt.Println("------------------------------------------")

	// //---percentage of lost replicas for single pgid_item in fault_osd--- (OSD DEGRADATION)
	// fmt.Printf("-> lost OSD: %d -> percentage%% of lost replicas for %s---(OSD DEGRADATION)\n", faultOsd, givenPg)
	// percentage := parser.PercentageCalculationAffectedReplicasPg(faultOsd, givenPg, pgDumpOutput, osdDumpOutput)
	// fmt.Printf("\nPercentage of %s replicas lost -> %.2f%%\n\n", givenPg, percentage)
	// fmt.Println("------------------------------------------")

	// parser.WarningCheck(percentage, faultOsd, givenPg, pgDumpOutput, osdDumpOutput)
	// fmt.Printf("---------------------------------------------------------------\n\n")

	// //---if these osds crush which are the affected Pgs and Pools?---
	// faultOsdSlice := []int{2, 4, 3}
	// fmt.Printf("-> If these osds:%d crush which Pgs and Pools are affected?\n", faultOsdSlice)
	// totalAffectedPgs, totalAffectedPools := parser.GetTotalAffectedPgsAndPools(faultOsdSlice, pgDumpOutput)
	// fmt.Printf("\n#%d Total affected pgs for osds:%d  -> (map is hidden, uncomment to view)", len(totalAffectedPgs), faultOsdSlice)
	// //fmt.Printf("%s\n\n", totalAffectedPgs)
	// fmt.Printf("\n#%d Total affected pools for osds:%d  -> %s\n\n", len(totalAffectedPools), faultOsdSlice, totalAffectedPools)

	// //    {pg: number_affected_replicas}
	// fmt.Printf("->   {pg: number_affected_replicas}\n")
	// pgNumberOfAffectedReplicaMap := parser.GetPgNumberOfAffectedReplicaMap(faultOsdSlice, pgDumpOutput)
	// fmt.Printf("%v\n\n", pgNumberOfAffectedReplicaMap)

	// //inHealthPgs, goodPgs, warningPgs, lostPgds
	// inHealthPgs, goodPgs, warningPgs, lostPgs := parser.GetPgsWithHighProbabilityOfLosingData(pgNumberOfAffectedReplicaMap, pgDumpOutput, osdDumpOutput)
	// fmt.Printf("---------------------------------------------------------------\n")

	// fmt.Printf("inHealthPgs (replicaLost=0%%):\n\n%s\n", inHealthPgs)
	// fmt.Printf("---------------------------------------------------------------\n")

	// fmt.Printf("goodPgs (0<replicaLost<50%%):\n\n%s\n\n", goodPgs)
	// goodPgsPoolSlice := parser.ExtractPoolsFromPgSlice(goodPgs)
	// fmt.Printf("pools ->%s\n", goodPgsPoolSlice)
	// fmt.Printf("previous pools are spread on these osds ->%d\n", parser.ExtractOsdsFromPoolSlice(goodPgsPoolSlice, pgDumpOutput))
	// fmt.Printf("---------------------------------------------------------------\n")

	// fmt.Printf("warningPgs (replicaLost>=50%% and replicaLost<100%%):\n\n%s\n\n", warningPgs)
	// warningPgsPoolSlice := parser.ExtractPoolsFromPgSlice(warningPgs)
	// fmt.Printf("pools ->%s\n", warningPgsPoolSlice)
	// fmt.Printf("previous pools are spread on these osds ->%d\n", parser.ExtractOsdsFromPoolSlice(warningPgsPoolSlice, pgDumpOutput))
	// fmt.Printf("---------------------------------------------------------------\n")

	// fmt.Printf("lostPgs  (replicaLost=100%%):\n\n%s\n\n", lostPgs)
	// lostPgsPoolSlice := parser.ExtractPoolsFromPgSlice(lostPgs)
	// fmt.Printf("pools ->%s\n", lostPgsPoolSlice)
	// fmt.Printf("previous pools are spread on these osds ->%d\n", parser.ExtractOsdsFromPoolSlice(lostPgsPoolSlice, pgDumpOutput))
	// fmt.Printf("---------------------------------------------------------------\n")

	// //checkTotal := len(inHealthPgs) + len(goodPgs) + len(warningPgs) + len(lostPgs)
	// //fmt.Printf("\ncheck total: %d\n", checkTotal)

	// //    {pool:[wPgs],...} -- Pool WarningPg Map
	// fmt.Printf("->Pool WarningPg Map  {pool:[wPgs],...}   -- Fault Osds:%d\n\n", faultOsdSlice)
	// poolWarningPgMap := parser.GetPoolWarningPgMap(faultOsdSlice, pgDumpOutput, osdDumpOutput)
	// fmt.Printf("%v\n", poolWarningPgMap)
	// fmt.Printf("---------------------------------------------------------------\n")

	// //    {pool:[lPgs],...}
	// fmt.Printf("->Pool LostPg Map  {pool:[lPgs],...}   -- Fault Osds:%d\n\n", faultOsdSlice)
	// poolLostPgMap := parser.GetPoolLostPgMap(faultOsdSlice, pgDumpOutput, osdDumpOutput)
	// fmt.Printf("%v\n", poolLostPgMap)
	// fmt.Printf("---------------------------------------------------------------\n")

	// //---if these pools crush which are the affected Pgs and osds?---(has it sense?)
	// fmt.Printf("-> If these pools crush which are the affected Pgs and osds?\n\n")
	// faultPools := []string{"1", "31"}
	// affectedPgs, affectedOsds := parser.GetTotalAffectedPgsAndOsds(faultPools, pgDumpOutput)
	// fmt.Printf("Affected OSDs: %d\n\n", affectedOsds)
	// fmt.Printf("Affected Pgs: %d\n\n", len(affectedPgs))
	// fmt.Printf("---------------------------------------------------------------\n")

	// //-----------------------osd-tree json-------------------------------------------
	// fmt.Printf("  OSD-TREE JSON Information Gathering\n\n")
	// fmt.Printf("-> Specialized Distribution Map information gathering\n\n")

	// hostName := "sv81" //map {hostName1: BucketDistribution , hostName2: BucketDistribution}
	// hostDistributionMap := parser.GetDistributionMap("host", pgDumpOutput, osdTreeOutput)
	// fmt.Printf("host %s contains these Osds-> %v\n\n", hostName, hostDistributionMap[hostName].Osd)
	// //---------------------------------

	// chassisName := "0a043220-0123456789"
	// chassisDistributionMap := parser.GetDistributionMap("chassis", pgDumpOutput, osdTreeOutput)
	// fmt.Printf("chassis %s contains these Hosts-> %v\n\n", chassisName, chassisDistributionMap[chassisName].Host)
	// fmt.Printf("chassis %s contains these Osds-> %v\n", chassisName, chassisDistributionMap[chassisName].Osd)
	// fmt.Printf("\n---------------------------------------------------------------\n")
	// //---------------------------------

	// fmt.Printf("-> Getting node information\n\n")
	// node := parser.GetNode("0a043220-0123456789", osdTreeOutput)
	// fmt.Printf("nodeName=%s, nodeId=%d, nodeType=%s", node.Name, node.ID, node.Type)
	// fmt.Printf("\n---------------------------------------------------------------\n")

	// //-----------------------osd-dump json-------------------------------------------
	// fmt.Printf("  OSD-DUMP JSON Information Gathering\n\n")
	// fmt.Printf("-> Getting pool information\n\n")
	// poolId := 3
	// poolStruct := parser.GetPool(poolId, osdDumpOutput)
	// fmt.Printf("pool=%d, poolName=%s, size=%d, min_size=%d", poolStruct.Pool, poolStruct.PoolName, poolStruct.Size, poolStruct.MinSize)
	// fmt.Printf("\n---------------------------------------------------------------\n")

	// //-----------------------Incremental Risk Calculator-------------------------------------------
	// fmt.Printf("-> If we switch off these Buckets how the risk of the pgs will change?\n\n")
	// fmt.Printf("		RISK CALCULATOR		\n\n")
	// host := "sv81"
	// pgNumberOfAffectedReplicaMap1 := parser.RiskCalculator(host, pgDumpOutput, osdTreeOutput)
	// fmt.Printf("(%s) NumberOfAffectedReplicaMap: %v\n\n", host, pgNumberOfAffectedReplicaMap1)

	// fmt.Printf("---------------------------------------------------------------\n")
	// fmt.Printf("		INCREMENTAL RISK CALCULATOR		\n\n")

	// //hosts := []string{"default"}  //buckets
	// hosts := []string{"sv81", "sv82", "sv61", "newnamefor-sv53"}
	// //hosts := []string{"newnamefor-sv53"}
	// fmt.Printf("hosts -> %s\n\n", hosts)
	// incrementalPgAffectedReplicaMap := parser.IncrementalRiskCalculator(hosts, pgDumpOutput, osdTreeOutput)
	// fmt.Printf("%v\n\n", incrementalPgAffectedReplicaMap)

	// fmt.Printf("[audit] len of incremental map %v\n\n", len(incrementalPgAffectedReplicaMap))
	// fmt.Printf("\n---------------------------------------------------------------\n")

	// inHealthPgs1, goodPgs1, warningPgs1, lostPgs1 := parser.GetPgsWithHighProbabilityOfLosingData(incrementalPgAffectedReplicaMap, pgDumpOutput, osdDumpOutput)
	// fmt.Printf("inHealthPgs (replicaLost=0%%):\n\n%s\n\n", inHealthPgs1)
	// fmt.Printf("goodPgs (0<replicaLost<50%%):\n\n%s\n\n", goodPgs1)
	// fmt.Printf("warningPgs (replicaLost>=50%% and replicaLost<100%%):\n\n%s\n\n", warningPgs1)
	// fmt.Printf("lostPgs  (replicaLost=100%%):\n\n%s\n", lostPgs1)
	// fmt.Printf("\n---------------------------------------------------------------\n")

	//-------------------------------------------------------------------------------------------

	osdInitiationDate := time.Date(2020, time.Month(3), 21, 1, 10, 30, 0, time.UTC)
	currentOsdLifeTime := utility.GetFloatRandomNumber(30, 90)

	if faultTimePrediction, meanDegradationRate, err := parser.GetOsdFaultTimePrediction(osdInitiationDate, currentOsdLifeTime); err == nil {
		fmt.Printf("the osd with this degradation rate:%f will reach the end of its optimal performance on date %v\n", meanDegradationRate, faultTimePrediction.Format("02-01-2006"))
	}

	//osdMap mock creation
	osdMap := make(map[string]map[string]interface{})
	for i := 1; i <= 8; i++ {
		osdMap["osd."+strconv.Itoa(i)] = map[string]interface{}{"initiationDate": time.Date(utility.GetIntRandomNumber(2018, 2023), time.Month(utility.GetIntRandomNumber(1, 12)), utility.GetIntRandomNumber(1, 30), utility.GetIntRandomNumber(0, 24), utility.GetIntRandomNumber(0, 59), 0, 0, time.UTC), "currentOsdLifeTime": utility.GetFloatRandomNumber(30, 80)}
	}

	forecastingTime := time.Date(2023, time.Month(4), 7, 1, 10, 30, 0, time.UTC)

	osdLifeForecastingMap, warningOsdSlice := parser.RiskFailureForecasting(osdMap, forecastingTime)
	fmt.Printf("\n%v\n", osdLifeForecastingMap)
	fmt.Printf("These osds are nearing the end of their optimal life (over 80%%) %s\n", warningOsdSlice)

	//Next step -> getting currentOsdLifeTime from disk lifetime
}
