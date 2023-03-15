package main

import (
	"fmt"
	parser "go/pkg"
	utility "go/tools"
	"io"
	"log"
)

func main() {
	log.SetOutput(io.Discard) //disable log

	//UnMarshalling Json
	CephDumpOutput := utility.ReadJson("ceph_dump_beautified.json")
	//fmt.Println(CephDumpOutput.PgMap.PgStats[0].Pgid)

	//---getting different pg_ids---
	fmt.Println("------------------------------------------")
	pgSlice := parser.GetPgs(CephDumpOutput)
	fmt.Printf("#Different pgids  -> %d\n\n", len(pgSlice))
	fmt.Println("------------------------------------------")

	//---getting pools---
	poolSlice := parser.GetPools(CephDumpOutput)
	fmt.Printf("Pools -> %s\n", poolSlice)
	fmt.Printf("#Pools -> %d\n\n", len(poolSlice))
	fmt.Println("------------------------------------------")

	//---getting_osds---
	osds := parser.GetOsds(CephDumpOutput)
	fmt.Printf("OSDs ->%d\n", osds)
	fmt.Printf("#OSDs ->%d\n\n", len(osds))
	fmt.Println("------------------------------------------")

	//---getting_pgs_details---
	osdPgMap := parser.GetOsdPgMap(CephDumpOutput)
	osdPgNumberMap := parser.GetNumberOfAssociatedPgsPerOsdMap(osdPgMap)
	fmt.Printf("Number of mapped PGs per OSD ->%d\n\n", osdPgNumberMap)
	fmt.Printf("#PGs ->%d\n\n", parser.GetTotalNumberOfPgs(osdPgNumberMap))
	fmt.Println("------------------------------------------")

	//---getting_pgIdOsdMap---
	fmt.Println("-> pgIdOsdMap")
	pgIdOsdMap := parser.GetPgOsdMap(CephDumpOutput)
	fmt.Printf("pgid 31.0 is on these osds -> %d\n\n", pgIdOsdMap["31.0"])
	fmt.Println("------------------------------------------")

	//---osd_pgs_mapping---
	fmt.Println("-> osd_pgs_mapping")
	numberOfPg := parser.GetNumberOfAssociatedPgsPerOsdMap(osdPgMap)[3]
	fmt.Printf("OSD 3 contains these #%d PGs-> %s\n\n", numberOfPg, osdPgMap[3])
	fmt.Println("------------------------------------------")

	//------osd_pool_pgs_map------
	//	   {(osd,pool):[pgs]}
	fmt.Println("-> osd_pool_pgs_map   {(osd,pool):[pgs]}")
	osdPoolPgMap := parser.GetOsdPoolPgMap(CephDumpOutput)
	fmt.Printf("OSD 6 pool 31 contains these PGs -> %s\n\n", osdPoolPgMap[parser.OsdPoolTuple{Osd: 6, Pool: "31"}])

	//	   {(osd,pool):{pg:number_replicas},...}
	fmt.Println("-> osd_pool_pgs_map   {(osd,pool):{pg:number_replicas},...}")
	osdPollNumberPerPgsMap := parser.GetOsdPoolNumberPerPgsMap(CephDumpOutput)
	fmt.Printf("OSD 6 pool 31 contains these PGs -> %v\n\n", osdPollNumberPerPgsMap[parser.OsdPoolTuple{Osd: 6, Pool: "31"}])
	fmt.Println("------------------------------------------")

	//---osds_containing_pool---
	fmt.Println("->   Pool->OSDs")
	pool := "3"
	osdsContaingPool := parser.GetOsdsContainingPool(pool, CephDumpOutput)
	fmt.Printf("Pool: %s is spread between these OSDs-> %v\n\n", pool, osdsContaingPool)
	fmt.Println("------------------------------------------")

	//---osds_containing_pg---
	fmt.Println("->   PG->OSDs")
	givenPg := "19.1f"
	osdsContainingPg := parser.GetOsdsContainingPg(givenPg, CephDumpOutput)
	fmt.Printf("PG: %s is spread between these OSDs-> %d\n\n", givenPg, osdsContainingPg)
	fmt.Println("------------------------------------------")

	//---affected pool for osd crush---
	fmt.Println("-> Affected pools")
	faultOsd := 2
	affectedPools := parser.GetAffectedPools(faultOsd, CephDumpOutput)
	fmt.Printf("Affected pools for crush on OSD: %d  -> #%d ->  %v\n\n", faultOsd, len(affectedPools), affectedPools)
	fmt.Println("------------------------------------------")

	//---affected pg for osd crush---
	fmt.Println("-> Affected PGs")
	affectedPgs := parser.GetAffectedPgs(faultOsd, CephDumpOutput)
	fmt.Printf("Affected pgid for crush on OSD: %d -> %s\n\n", faultOsd, affectedPgs)
	fmt.Println("------------------------------------------")

	//---percentage of lost replicas for single pgid_item in fault_osd--- (OSD DEGRADATION)
	fmt.Printf("-> lost OSD: %d -> percentage%% of lost replicas for %s---(OSD DEGRADATION)\n", faultOsd, givenPg)
	percentage := parser.PercentageCalculationAffectedReplicasPg(faultOsd, givenPg, CephDumpOutput)
	fmt.Printf("\nPercentage of %s replicas lost -> %.2f%%\n\n", givenPg, percentage)
	fmt.Println("------------------------------------------")

	parser.WarningCheck(percentage, faultOsd, givenPg, CephDumpOutput)
	fmt.Printf("---------------------------------------------------------------\n\n")

	//---if these osds crush which are the affected Pgs and Pools?---
	faultOsdSlice := []int{2, 4, 3}
	fmt.Printf("-> If these osds:%d crush which Pgs and Pools are affected?\n", faultOsdSlice)
	totalAffectedPgs, totalAffectedPools := parser.GetTotalAffectedPgsAndPools(faultOsdSlice, CephDumpOutput)
	fmt.Printf("\n#%d Total affected pgs for osds:%d  -> (map is hidden, uncomment to view)", len(totalAffectedPgs), faultOsdSlice)
	//fmt.Printf("%s\n\n", totalAffectedPgs)
	fmt.Printf("\n#%d Total affected pools for osds:%d  -> %s\n\n", len(totalAffectedPools), faultOsdSlice, totalAffectedPools)

	//    {pg: number_affected_replicas}
	fmt.Printf("->   {pg: number_affected_replicas}\n")
	pgNumberOfAffectedReplicaMap := parser.GetPgNumberOfAffectedReplicaMap(faultOsdSlice, CephDumpOutput)
	fmt.Printf("%v\n\n", pgNumberOfAffectedReplicaMap)

	//inHealthPgs, goodPgs, warningPgs, lostPgds
	inHealthPgs, goodPgs, warningPgs, lostPgs := parser.GetPgsWithHighProbabilityOfLosingData(faultOsdSlice, CephDumpOutput)
	fmt.Printf("---------------------------------------------------------------\n")

	fmt.Printf("inHealthPgs: replicaLost=0%%\n\n%s\n", inHealthPgs)
	fmt.Printf("---------------------------------------------------------------\n")

	fmt.Printf("goodPgs (replicaLost<50%%):\n\n%s\n\n", goodPgs)
	goodPgsPoolSlice := parser.ExtractPoolsFromPgSlice(goodPgs)
	fmt.Printf("pools ->%s\n", goodPgsPoolSlice)
	fmt.Printf("previous pools are spread on these osds ->%d\n", parser.ExtractOsdsFromPoolSlice(goodPgsPoolSlice, CephDumpOutput))
	fmt.Printf("---------------------------------------------------------------\n")

	fmt.Printf("warningPgs (replicaLost>=50%% and replicaLost<100%%):\n\n%s\n\n", warningPgs)
	warningPgsPoolSlice := parser.ExtractPoolsFromPgSlice(warningPgs)
	fmt.Printf("pools ->%s\n", warningPgsPoolSlice)
	fmt.Printf("previous pools are spread on these osds ->%d\n", parser.ExtractOsdsFromPoolSlice(warningPgsPoolSlice, CephDumpOutput))
	fmt.Printf("---------------------------------------------------------------\n")

	fmt.Printf("lostPgs  (replicaLost=100%%):\n\n%s\n\n", lostPgs)
	lostPgsPoolSlice := parser.ExtractPoolsFromPgSlice(lostPgs)
	fmt.Printf("pools ->%s\n", lostPgsPoolSlice)
	fmt.Printf("previous pools are spread on these osds ->%d\n", parser.ExtractOsdsFromPoolSlice(lostPgsPoolSlice, CephDumpOutput))
	fmt.Printf("---------------------------------------------------------------\n")

	//checkTotal := len(inHealthPgs) + len(goodPgs) + len(warningPgs) + len(lostPgs)
	//fmt.Printf("\ncheck total: %d\n", checkTotal)

	//    {pool:[wPgs],...} -- Pool WarningPg Map
	fmt.Printf("->Pool WarningPg Map  {pool:[wPgs],...}   -- Fault Osds:%d\n\n", faultOsdSlice)
	poolWarningPgMap := parser.GetPoolWarningPgMap(faultOsdSlice, CephDumpOutput)
	fmt.Printf("%v\n", poolWarningPgMap)
	fmt.Printf("---------------------------------------------------------------\n")

	//    {pool:[lPgs],...}
	fmt.Printf("->Pool LostPg Map  {pool:[lPgs],...}   -- Fault Osds:%d\n\n", faultOsdSlice)
	poolLostPgMap := parser.GetPoolLostPgMap(faultOsdSlice, CephDumpOutput)
	fmt.Printf("%v\n", poolLostPgMap)
	fmt.Printf("---------------------------------------------------------------\n")

	//---if these pools crush which are the affected Pgs and osds?---(has it sense?)
	fmt.Printf("-> If these pools crush which are the affected Pgs and osds?\n\n")
	faultPools := []string{"1", "31"}
	affectedPgs, affectedOsds := parser.GetTotalAffectedPgsAndOsds(faultPools, CephDumpOutput)
	fmt.Printf("Affected OSDs: %d\n\n", affectedOsds)
	fmt.Printf("Affected Pgs: %d\n\n", len(affectedPgs))

}
