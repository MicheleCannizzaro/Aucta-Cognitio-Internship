package pkg

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	s "github.com/MicheleCannizzaro/Aucta-Cognitio-Internship/go_server/structs"
	utility "github.com/MicheleCannizzaro/Aucta-Cognitio-Internship/go_server/tools"
)

// ----------------------------STRUCTS----------------------------
type OsdPoolTuple struct {
	Osd, Pool interface{}
}

// ----------------INFORMATIONS GATHERING FUNCTIONS (pg dump)----------------
func GetPgs(pgDumpOutput s.PgDumpOutputStruct) []string {

	pgStats := pgDumpOutput.PgMap.PgStats
	pgSlice := []string{}

	for _, s := range pgStats {
		pgSlice = append(pgSlice, s.Pgid)
	}

	return pgSlice
}

func GetPools(pgDumpOutput s.PgDumpOutputStruct) []string {

	pgSlice := GetPgs(pgDumpOutput)
	poolSlice := []string{}

	for _, s := range pgSlice {
		poolId := strings.Split(s, ".")[0]

		poolSlice = append(poolSlice, poolId)
	}

	poolSlice = utility.RmvDuplStr(poolSlice)

	return poolSlice
}

func GetOsds(pgDumpOutput s.PgDumpOutputStruct) []int {

	osdPgMap := GetOsdPgMap(pgDumpOutput)
	osds := []int{}

	for osd := range osdPgMap {
		osds = append(osds, osd)
	}

	return osds
}

func GetTotalNumberOfPgs(osdPgNumberMap map[int]int) int {

	totalPgs := 0

	for _, osdPgCount := range osdPgNumberMap {
		totalPgs += osdPgCount
	}

	return totalPgs
}

func GetOsdsContainingPool(pool string, pgDumpOutput s.PgDumpOutputStruct) []int {

	osdPoolPgMap := GetOsdPoolPgMap(pgDumpOutput)

	poolOsdsSlice := []int{}

	for key := range osdPoolPgMap {
		if key.Pool == pool {
			poolOsdsSlice = append(poolOsdsSlice, key.Osd.(int))
		}

	}

	return utility.RemoveDuplicateInt(poolOsdsSlice)
}

func GetOsdsContainingPg(givenPg string, pgDumpOutput s.PgDumpOutputStruct) []int {

	osdPoolPgMap := GetOsdPoolPgMap(pgDumpOutput)

	pgOsdsSlice := []int{}

	for key, pgSlice := range osdPoolPgMap {
		for _, pg := range pgSlice {
			if pg == givenPg {
				pgOsdsSlice = append(pgOsdsSlice, key.Osd.(int))
			}
		}
	}

	return pgOsdsSlice
}

func GetAllPgInAllOsds(pgDumpOutput s.PgDumpOutputStruct) []string {

	osdPgMap := GetOsdPgMap(pgDumpOutput)

	allPgs := []string{}

	for _, pgSlice := range osdPgMap {
		allPgs = append(allPgs, pgSlice...)
	}

	return allPgs
}

func ExtractPoolsFromPgSlice(pgSlice []string) []string {
	//defer utility.Duration(utility.Track("extractPoolsFromPgSlice"))

	poolSlice := []string{}

	for _, pg := range pgSlice {
		poolId := strings.Split(pg, ".")[0]
		poolSlice = append(poolSlice, poolId)
	}
	return utility.RmvDuplStr(poolSlice)
}

func ExtractOsdsFromPoolSlice(poolSlice []string, pgDumpOutput s.PgDumpOutputStruct) []int {
	//defer utility.Duration(utility.Track("extractOsdsFromPoolSlice"))
	osdSlice := []int{}

	for _, pool := range poolSlice {
		osdsContainingPool := GetOsdsContainingPool(pool, pgDumpOutput)
		osdSlice = append(osdSlice, osdsContainingPool...)
	}

	return utility.RemoveDuplicateInt(osdSlice)
}

// -----------------------OSD PG POOL MAPPING FUNCTIONS--------------------------
func GetPgOsdMap(pgDumpOutput s.PgDumpOutputStruct) map[string][]int {
	//defer utility.Duration(utility.Track("getPgIdOsdMap"))

	pgStats := pgDumpOutput.PgMap.PgStats
	pgIdOsdMap := make(map[string][]int)

	for _, s := range pgStats {
		pgOsdSlice := s.Up
		pgIdOsdMap[s.Pgid] = pgOsdSlice
	}

	return pgIdOsdMap
}

func GetOsdPgMap(pgDumpOutput s.PgDumpOutputStruct) map[int][]string {
	//defer utility.Duration(utility.Track("getOsdPgMap"))

	pgIdOsdMap := GetPgOsdMap(pgDumpOutput)
	osdPgMap := make(map[int][]string)

	for key, osdSlice := range pgIdOsdMap {
		for _, osd := range osdSlice {
			if _, ok := osdPgMap[osd]; !ok {
				osdPgMap[osd] = []string{}
			}

			osdPgMap[osd] = append(osdPgMap[osd], key)
		}
	}

	return osdPgMap
}

func GetNumberOfAssociatedPgsPerOsdMap(osdPgMap map[int][]string) map[int]int {
	//defer utility.Duration(utility.Track("getNumberOfAssociatedPgsPerOsdMap"))

	osdPgNumberMap := make(map[int]int)

	for key, value := range osdPgMap {
		osdPgNumberMap[key] = len(value)
	}

	return osdPgNumberMap
}

func GetOsdPoolPgMap(pgDumpOutput s.PgDumpOutputStruct) map[OsdPoolTuple][]string {
	//defer utility.Duration(utility.Track("getOsdPoolPgMap"))

	osdPgMap := GetOsdPgMap(pgDumpOutput)
	osdPoolPgMap := make(map[OsdPoolTuple][]string)

	for osd := range osdPgMap {
		for _, pg := range osdPgMap[osd] {
			poolId := strings.Split(pg, ".")[0]
			osdPoolTuple := OsdPoolTuple{osd, poolId}

			if _, ok := osdPoolPgMap[osdPoolTuple]; !ok {
				osdPoolPgMap[osdPoolTuple] = []string{}
			}

			osdPoolPgMap[osdPoolTuple] = append(osdPoolPgMap[osdPoolTuple], pg)
		}
	}

	return osdPoolPgMap
}

func GetOsdPoolNumberPerPgsMap(pgDumpOutput s.PgDumpOutputStruct) map[OsdPoolTuple]map[string]int {
	//defer utility.Duration(utility.Track("getOsdPoolNumberPerPgsMap"))

	osdPoolPgMap := GetOsdPoolPgMap(pgDumpOutput)
	osdPoolNumberPerPgsMap := make(map[OsdPoolTuple]map[string]int)

	for tupleKey := range osdPoolPgMap {
		temp := make(map[string]int)

		for _, value := range osdPoolPgMap[tupleKey] {
			if _, ok := temp[value]; !ok {
				temp[value] = 1
			} else {
				temp[value] += 1
			}
			osdPoolNumberPerPgsMap[tupleKey] = temp
		}
	}

	return osdPoolNumberPerPgsMap
}

func GetPoolWarningPgMap(faultOsdSlice []int, pgDumpOutput s.PgDumpOutputStruct, osdDumpOutput s.OsdDumpOutputStruct) map[string][]string {
	//defer utility.Duration(utility.Track("getPoolWarningPgMap"))

	poolWarningPgMap := make(map[string][]string)

	pgNumberOfAffectedReplicaMap := GetPgNumberOfAffectedReplicaMap(faultOsdSlice, pgDumpOutput)
	_, _, warningPgs, _ := GetPgsWithHighProbabilityOfLosingData(pgNumberOfAffectedReplicaMap, pgDumpOutput, osdDumpOutput)

	for _, wPgs := range warningPgs {
		poolId := strings.Split(wPgs, ".")[0]

		if _, ok := poolWarningPgMap[poolId]; !ok {
			poolWarningPgMap[poolId] = []string{}
		}

		poolWarningPgMap[poolId] = append(poolWarningPgMap[poolId], wPgs)
	}

	return poolWarningPgMap
}

func GetPoolLostPgMap(faultOsdSlice []int, pgDumpOutput s.PgDumpOutputStruct, osdDumpOutput s.OsdDumpOutputStruct) map[string][]string {
	//defer utility.Duration(utility.Track("getPoolLostPgMap"))

	poolLostPgMap := make(map[string][]string)

	pgNumberOfAffectedReplicaMap := GetPgNumberOfAffectedReplicaMap(faultOsdSlice, pgDumpOutput)

	//_, _, _, lostPgs := GetPgsWithHighProbabilityOfLosingData(pgNumberOfAffectedReplicaMap, pgDumpOutput)
	_, _, _, lostPgs := GetPgsWithHighProbabilityOfLosingData(pgNumberOfAffectedReplicaMap, pgDumpOutput, osdDumpOutput)

	for _, wPgs := range lostPgs {
		poolId := strings.Split(wPgs, ".")[0]

		if _, ok := poolLostPgMap[poolId]; !ok {
			poolLostPgMap[poolId] = []string{}
		}

		poolLostPgMap[poolId] = append(poolLostPgMap[poolId], wPgs)
	}

	return poolLostPgMap
}

// ----------------OSDs FAULT CONSEQUENCES DEFINITION FUNCTIONS------------------

// important
func GetAffectedPools(faultOsd int, pgDumpOutput s.PgDumpOutputStruct) []string {
	//defer utility.Duration(utility.Track("getAffectedPools"))

	osdPgMap := GetOsdPgMap(pgDumpOutput)

	affectedPoolSlice := []string{}

	for _, pg := range osdPgMap[faultOsd] {
		poolId := strings.Split(pg, ".")[0]

		affectedPoolSlice = append(affectedPoolSlice, poolId)
	}

	return utility.RmvDuplStr(affectedPoolSlice)
}

// important
func GetAffectedPgs(faultOsd int, pgDumpOutput s.PgDumpOutputStruct) []string {
	//defer utility.Duration(utility.Track("getAffectedPgs"))

	osdPgMap := GetOsdPgMap(pgDumpOutput)

	affectedPgSlice := []string{}

	affectedPgSlice = append(affectedPgSlice, osdPgMap[faultOsd]...)

	return affectedPgSlice
}

func getNumAffectedReplicasPgItem(faultOsd int, givenPg string, pgDumpOutput s.PgDumpOutputStruct) float64 {
	//defer utility.Duration(utility.Track("getNumAffectedReplicasPgItem"))

	affectedPgs := GetAffectedPgs(faultOsd, pgDumpOutput)

	count := 0

	for _, pg := range affectedPgs {
		if pg == givenPg {
			count += 1
		}
	}

	return float64(count)
}

func PercentageCalculationAffectedReplicasPg(faultOsd int, givenPg string, pgDumpOutput s.PgDumpOutputStruct, osdDumpOutput s.OsdDumpOutputStruct) float64 {
	//defer utility.Duration(utility.Track("percentageCalculationAffectedReplicasPg"))

	numAffectedReplicaPgItem := getNumAffectedReplicasPgItem(faultOsd, givenPg, pgDumpOutput)

	pgTotalReplicas, _, _ := GetOsdNumberForPg(givenPg, osdDumpOutput)

	percentage := (numAffectedReplicaPgItem / pgTotalReplicas) * 100

	return utility.RoundFloat(percentage, 2)
}

func WarningCheck(percentage float64, faultOsd int, givenPg string, pgDumpOutput s.PgDumpOutputStruct, osdDumpOutput s.OsdDumpOutputStruct) {
	//defer utility.Duration(utility.Track("warningCheck"))
	//log.SetOutput(io.Discard) //disable log

	osds := GetOsds(pgDumpOutput)
	warningOsdSlice := make(map[int]float64)

	if percentage > 0.0 {
		for osd := range osds {
			//percentage = PercentageCalculationAffectedReplicasPg(osd, givenPg, pgDumpOutput)
			percentage = PercentageCalculationAffectedReplicasPg(osd, givenPg, pgDumpOutput, osdDumpOutput)
			if percentage > 0.0 && osd != faultOsd {
				warningOsdSlice[osd] = percentage
			}
		}
	}

	//log.SetOutput(os.Stdout) //renable log
	fmt.Printf("Warning: check these other OSDs and percentages for %s ->%v\n", givenPg, warningOsdSlice)
}

// important
func GetTotalAffectedPgsAndPools(faultOsdSlice []int, pgDumpOutput s.PgDumpOutputStruct) ([]string, []string) {
	//defer utility.Duration(utility.Track("getTotalAffectedPgsAndPools"))

	totalAffectedPgs := []string{}
	totalAffectedPools := []string{}

	for _, osd := range faultOsdSlice {
		affectedPgs := GetAffectedPgs(osd, pgDumpOutput)
		affectedPools := GetAffectedPools(osd, pgDumpOutput)

		totalAffectedPgs = append(totalAffectedPgs, affectedPgs...)
		totalAffectedPools = append(totalAffectedPools, affectedPools...)
	}

	totalAffectedPools = utility.RmvDuplStr(totalAffectedPools)

	return totalAffectedPgs, totalAffectedPools
}

// important
func GetPgNumberOfAffectedReplicaMap(faultOsdSlice []int, pgDumpOutput s.PgDumpOutputStruct) map[string]int {
	//defer utility.Duration(utility.Track("getPgNumberOfAffectedReplicaMap"))

	totalAffectedPgs, _ := GetTotalAffectedPgsAndPools(faultOsdSlice, pgDumpOutput)
	pgNumberOfAffectedReplicaMap := make(map[string]int)

	for _, pg := range totalAffectedPgs {
		if _, ok := pgNumberOfAffectedReplicaMap[pg]; !ok {
			pgNumberOfAffectedReplicaMap[pg] = 1
		} else {
			pgNumberOfAffectedReplicaMap[pg] += 1
		}
	}

	return pgNumberOfAffectedReplicaMap
}

// important
func GetPgsWithHighProbabilityOfLosingData(pgNumberOfAffectedReplicaMap map[string]int, pgDumpOutput s.PgDumpOutputStruct, osdDumpOutput s.OsdDumpOutputStruct) (inHealthPgs []string, goodPgs []string, warningPgs []string, lostPgs []string) {
	//defer utility.Duration(utility.Track("getPgsWithHighProbabilityOfLosingData"))

	for pg, numLostReplicas := range pgNumberOfAffectedReplicaMap {

		pgTotalReplicas, _, _ := GetOsdNumberForPg(pg, osdDumpOutput)
		percetageLostReplicas := (float64(numLostReplicas) / pgTotalReplicas) * 100

		if percetageLostReplicas > 0.0 && percetageLostReplicas < 50.0 {
			goodPgs = append(goodPgs, pg)
		} else if percetageLostReplicas >= 50.0 && percetageLostReplicas < 100.0 {
			warningPgs = append(warningPgs, pg)
		} else if percetageLostReplicas >= 100.0 {
			lostPgs = append(lostPgs, pg)
		}
	}

	compromisedPgs := [][]string{goodPgs, warningPgs, lostPgs}
	allCompromisedPgs := []string{}

	for _, ele := range compromisedPgs {
		allCompromisedPgs = append(allCompromisedPgs, ele...)
	}

	inHealthPgs = utility.Difference(utility.RmvDuplStr(GetAllPgInAllOsds(pgDumpOutput)), allCompromisedPgs)

	return inHealthPgs, goodPgs, warningPgs, lostPgs
}

// important
func GetLostAndWarningPgs(faultBucketOrRouter []string, pgDumpOutput s.PgDumpOutputStruct, osdTreeOutput s.OsdTreeOutputStruct, osdDumpOutput s.OsdDumpOutputStruct) ([]string, []string, []int, error) {
	warningPgsSlice := []string{}
	lostPgsSlice := []string{}

	incrementalPgAffectedReplicaMap, totalAffectedOsdIds, err := IncrementalRiskCalculator(faultBucketOrRouter, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		fmt.Println("GetLostAndWarningPgs: Error in IncrementalRiskCalculator")
		return nil, nil, nil, err
	}

	for pg, numPgAffectedReplicas := range incrementalPgAffectedReplicaMap {
		numPgTotalReplicas, k, m := GetOsdNumberForPg(pg, osdDumpOutput) //Equal to Pool Size, in other words is the number of OSDs in which is stored the PG

		numPgRemainingReplicas := int(numPgTotalReplicas) - numPgAffectedReplicas

		if k == 0 && m == 0 { //replicated pool
			if numPgRemainingReplicas <= 1 {
				if numPgRemainingReplicas == 0 {
					lostPgsSlice = append(lostPgsSlice, pg)
				}

				if numPgRemainingReplicas == 1 {
					warningPgsSlice = append(warningPgsSlice, pg)
				}
			}
		} else { //erasure coded pool
			if numPgRemainingReplicas <= k {
				if numPgRemainingReplicas < k {
					lostPgsSlice = append(lostPgsSlice, pg)
				}

				if numPgRemainingReplicas == k {
					warningPgsSlice = append(warningPgsSlice, pg)
				}
			}
		}
	}

	return lostPgsSlice, warningPgsSlice, totalAffectedOsdIds, nil
}

//GetPoolDataLossProbability
// important
// func GetPoolDataLossProbability(faultBucketOrRouter []string, pgDumpOutput s.PgDumpOutputStruct, osdTreeOutput s.OsdTreeOutputStruct, osdDumpOutput s.OsdDumpOutputStruct) (map[string]float64, error) {
//
// 	poolDataLossProbabilityMap := make(map[string]float64)
//
// 	//get the number of PGs inside each pool
// 	poolPgNumberMap, err := GetPoolPgNumberMap(pgDumpOutput, osdDumpOutput)
//
// 	if err != nil {
// 		error := errors.New("something went wrong in getting PoolPgNumberMap")
// 		return poolDataLossProbabilityMap, error
// 	}
//
// 	//extract Pools' values from GetLostAndWarningPgs
// 	poolAffectedPgsNumberMap := make(map[string]int)
// 	lostPgsSlice, warningPgsSlice, err := GetLostAndWarningPgs(faultBucketOrRouter, pgDumpOutput, osdTreeOutput, osdDumpOutput)
// 	if err != nil {
// 		fmt.Println("GetPoolDataLossProbability: Error in GetLostAndWarningPgs")
// 		return nil, err
// 	}
// 	lostPgsPoolSlice := ExtractPoolsFromPgSlice(lostPgsSlice)
//
// 	for _, pg := range warningPgsSlice {
// 		poolId := strings.Split(pg, ".")[0]
//
// 		if _, isPresent := poolAffectedPgsNumberMap[poolId]; !isPresent {
// 			poolAffectedPgsNumberMap[poolId] = 0
// 		}
//
// 		poolAffectedPgsNumberMap[poolId] += 1
// 	}
//
// 	pools := GetPools(pgDumpOutput)
//
// 	for _, pool := range pools {
//
// 		if utility.StringInSlice(pool, lostPgsPoolSlice) {
// 			poolAffectedPgsNumberMap[pool] = poolPgNumberMap[pool]
// 			//poolAffectedPgsNumberMap[pool] = -1 * poolPgNumberMap[pool]
// 			//poolAffectedPgsNumberMap[pool] = -1 * poolAffectedPgsNumberMap[pool]
// 		}
//
// 		poolDatalossProbability := float64(poolAffectedPgsNumberMap[pool]) / float64(poolPgNumberMap[pool])
//
// 		poolId, err := strconv.Atoi(pool)
// 		if err != nil {
// 			error := errors.New("something went wrong in converting poolId in string format to int")
// 			return nil, error
// 		}
// 		poolName := GetPool(poolId, osdDumpOutput).PoolName
//
// 		poolDataLossProbabilityMap[poolName] = utility.RoundFloat(poolDatalossProbability, 2)
// 	}
//
// 	return poolDataLossProbabilityMap, nil
// }

func GetPoolDataLossProbability(faultBucketOrRouter []string, pgDumpOutput s.PgDumpOutputStruct, osdTreeOutput s.OsdTreeOutputStruct, osdDumpOutput s.OsdDumpOutputStruct) (map[string]float64, error) {

	poolDataLossProbabilityMap := make(map[string]float64)
	UniqueOsdsInWarningPgSlicePerPool := make(map[string][]int)

	lostPgsSlice, warningPgsSlice, totalAffectedOsdIds, err := GetLostAndWarningPgs(faultBucketOrRouter, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		fmt.Println("GetPoolDataLossProbability: Error in GetLostAndWarningPgs")
		return nil, err
	}

	lostPgsPoolSlice := ExtractPoolsFromPgSlice(lostPgsSlice)
	fmt.Printf("lostPgsPoolSlice:= %v\n", lostPgsPoolSlice)

	//Prob(data is lost with next osd failure) = (num_unique_osds in warningPgSlice)/(num_good_osds)

	for _, pg := range warningPgsSlice {
		poolId := strings.Split(pg, ".")[0]

		actingSet := GetOsdsContainingPg(pg, pgDumpOutput)

		//fmt.Printf("Acting Set: %d\n", actingSet)

		for _, osdId := range actingSet {

			if !utility.IntInSlice(osdId, totalAffectedOsdIds) {
				if _, isPresent := UniqueOsdsInWarningPgSlicePerPool[poolId]; !isPresent {
					UniqueOsdsInWarningPgSlicePerPool[poolId] = []int{}
				}
				UniqueOsdsInWarningPgSlicePerPool[poolId] = append(UniqueOsdsInWarningPgSlicePerPool[poolId], osdId)
			}
		}
	}

	pools := GetPools(pgDumpOutput)
	for _, pool := range pools {
		UniqueOsdsInWarningPgSlicePerPool[pool] = utility.RemoveDuplicateInt(UniqueOsdsInWarningPgSlicePerPool[pool])
	}

	fmt.Printf("UniqueOsdsInWarningPgSlicePerPool:= %v\n", UniqueOsdsInWarningPgSlicePerPool)

	totalOsds := GetOsds(pgDumpOutput)
	remaingOSDs := utility.IntDifference(totalOsds, totalAffectedOsdIds)

	fmt.Printf("remaining OSDs:= %v\n", remaingOSDs)

	for _, pool := range pools {

		var poolDatalossProbability float64
		if len(remaingOSDs) == 0 {
			poolDatalossProbability = 1.0

		} else {

			if utility.StringInSlice(pool, lostPgsPoolSlice) {
				UniqueOsdsInWarningPgSlicePerPool[pool] = remaingOSDs
			}

			poolDatalossProbability = float64(len(UniqueOsdsInWarningPgSlicePerPool[pool])) / float64(len(remaingOSDs))
		}

		//convert poolIds in poolNames
		poolId, err := strconv.Atoi(pool)
		if err != nil {
			error := errors.New("something went wrong in converting poolId in string format to int")
			return nil, error
		}
		poolName := GetPool(poolId, osdDumpOutput).PoolName

		poolDataLossProbabilityMap[poolName] = utility.RoundFloat(poolDatalossProbability, 2)
	}

	return poolDataLossProbabilityMap, nil
}

func GetTotalAffectedPgsAndOsds(faultPoolSlice []string, pgDumpOutput s.PgDumpOutputStruct) ([]string, []int) {
	//defer utility.Duration(utility.Track("getTotalAffectedPgsAndOsds"))

	totalAffectedPgs := []string{}
	totalAffectedOsds := []int{}

	//osdPoolPgMap := getOsdPoolPgMap(pgDumpOutput)

	for _, pool := range faultPoolSlice {
		totalAffectedOsds = append(totalAffectedOsds, GetOsdsContainingPool(pool, pgDumpOutput)...)
	}

	totalAffectedOsds = utility.RemoveDuplicateInt(totalAffectedOsds)

	for _, osd := range totalAffectedOsds {
		totalAffectedPgs = append(totalAffectedPgs, GetAffectedPgs(osd, pgDumpOutput)...)
	}

	return totalAffectedPgs, totalAffectedOsds
}

// ----------------INFORMATIONS GATHERING FUNCTIONS (osd-tree)--------------------------

func GetBucketName(searchedDevice string, choosenBucketMap map[string][]string) (bucketName string) {
	for bucket, containedElements := range choosenBucketMap {
		for _, elementName := range containedElements {
			if elementName == searchedDevice {
				bucketName = bucket
			}
		}
	}
	return
}

func GetBucketQuantity(bucketType string, osdTreeOutput s.OsdTreeOutputStruct) (int, error) {
	osdTreeStats := osdTreeOutput.Nodes
	counter := 0

	bucketSlice := []string{"root", "region", "zone", "datacenter", "rack", "chassis", "host", "osd"}
	isBucketTypeValid := false

	for _, bucketName := range bucketSlice {
		if bucketName == bucketType {
			isBucketTypeValid = true
		}
	}

	if !isBucketTypeValid {
		error := errors.New("invalid Bucket Name")
		return counter, error
	}

	for _, node := range osdTreeStats {
		if node.Type == bucketType {
			counter += 1
		}
	}

	return counter, nil
}

// important
func GetDistributionMap(choosenMapKey string, pgDumpOutput s.PgDumpOutputStruct, osdTreeOutput s.OsdTreeOutputStruct) map[string]s.BucketDistribution {
	osdTreeStats := osdTreeOutput.Nodes

	bucketSlice := []string{"root", "region", "zone", "datacenter", "rack", "chassis", "host"}

	mappingSlice := []map[string][]string{}                  // example: map{bucketName1:[childrenName, childrenName, ...], bucketName2:[...], ...}
	distributionMap := make(map[string]s.BucketDistribution) // example: map{hostName1: bucketDistribution, hostName2: bucketDistribution, ...}

	for _, bucket := range bucketSlice {
		var bucketName string
		var bucketChildren []int
		tempMap := make(map[string][]string)

		for _, node := range osdTreeStats {
			if node.Type == bucket {
				bucketName = node.Name
				bucketChildren = append(bucketChildren, node.Children...)
			}

			for _, child := range bucketChildren {
				if node.ID == child {
					tempMap[bucketName] = append(tempMap[bucketName], node.Name)
				}
			}
		}

		mappingSlice = append(mappingSlice, tempMap)
	}

	rootMap := mappingSlice[0]
	regionMap := mappingSlice[1]
	zoneMap := mappingSlice[2]
	datacenterMap := mappingSlice[3]
	rackMap := mappingSlice[4]
	chassisMap := mappingSlice[5]
	hostMap := mappingSlice[6]

	switch choosenMapKey {
	case "root":
		for root := range rootMap {

			if _, isPresent := distributionMap[root]; !isPresent {
				regionSlice := rootMap[root]

				zoneSlice := []string{}
				datacenterSlice := []string{}
				rackSlice := []string{}
				chassisSlice := []string{}
				hostSlice := []string{}
				osdSlice := []string{}

				for _, region := range regionSlice {
					zoneSlice = append(zoneSlice, regionMap[region]...)
				}

				for _, zone := range zoneSlice {
					datacenterSlice = append(datacenterSlice, zoneMap[zone]...)
				}

				for _, datacenter := range datacenterSlice {
					rackSlice = append(rackSlice, datacenterMap[datacenter]...)
				}

				for _, rack := range rackSlice {
					chassisSlice = append(chassisSlice, rackMap[rack]...)
				}

				for _, chassis := range chassisSlice {
					hostSlice = append(hostSlice, chassisMap[chassis]...)
				}

				for _, host := range hostSlice {
					osdSlice = append(osdSlice, hostMap[host]...)

				}

				rootDistribution := s.BucketDistribution{Root: root, Region: regionSlice,
					Zone: zoneSlice, Datacenter: datacenterSlice, Rack: rackSlice,
					Chassis: chassisSlice, Host: hostSlice, Osd: osdSlice}

				distributionMap[root] = rootDistribution
			}
		}
	case "region":
		for region := range regionMap {

			regionName := region
			rootName := GetBucketName(regionName, rootMap)

			if _, isPresent := distributionMap[regionName]; !isPresent {
				zoneSlice := regionMap[regionName]

				datacenterSlice := []string{}
				rackSlice := []string{}
				chassisSlice := []string{}
				hostSlice := []string{}
				osdSlice := []string{}

				for _, zone := range zoneSlice {
					datacenterSlice = append(datacenterSlice, zoneMap[zone]...)
				}

				for _, datacenter := range datacenterSlice {
					rackSlice = append(rackSlice, datacenterMap[datacenter]...)
				}

				for _, rack := range rackSlice {
					chassisSlice = append(chassisSlice, rackMap[rack]...)
				}

				for _, chassis := range chassisSlice {
					hostSlice = append(hostSlice, chassisMap[chassis]...)
				}

				for _, host := range hostSlice {
					osdSlice = append(osdSlice, hostMap[host]...)

				}

				regionDistribution := s.BucketDistribution{Root: rootName, Region: regionName,
					Zone: zoneSlice, Datacenter: datacenterSlice, Rack: rackSlice,
					Chassis: chassisSlice, Host: hostSlice, Osd: osdSlice}

				distributionMap[regionName] = regionDistribution
			}
		}
	case "zone":
		for zone := range zoneMap {

			zoneName := zone
			regionName := GetBucketName(zoneName, regionMap)
			rootName := GetBucketName(regionName, rootMap)

			if _, isPresent := distributionMap[zoneName]; !isPresent {
				datacenterSlice := zoneMap[zoneName]

				rackSlice := []string{}
				chassisSlice := []string{}
				hostSlice := []string{}
				osdSlice := []string{}

				for _, datacenter := range datacenterSlice {
					rackSlice = append(rackSlice, datacenterMap[datacenter]...)
				}

				for _, rack := range rackSlice {
					chassisSlice = append(chassisSlice, rackMap[rack]...)
				}

				for _, chassis := range chassisSlice {
					hostSlice = append(hostSlice, chassisMap[chassis]...)
				}

				for _, host := range hostSlice {
					osdSlice = append(osdSlice, hostMap[host]...)

				}

				zoneDistribution := s.BucketDistribution{Root: rootName, Region: regionName,
					Zone: zoneName, Datacenter: datacenterSlice, Rack: rackSlice,
					Chassis: chassisSlice, Host: hostSlice, Osd: osdSlice}

				distributionMap[zoneName] = zoneDistribution
			}
		}
	case "datacenter":
		for datacenter := range datacenterMap {

			datacenterName := datacenter
			zoneName := GetBucketName(datacenterName, zoneMap)
			regionName := GetBucketName(zoneName, regionMap)
			rootName := GetBucketName(regionName, rootMap)

			if _, isPresent := distributionMap[datacenterName]; !isPresent {
				rackSlice := datacenterMap[datacenterName]

				chassisSlice := []string{}
				hostSlice := []string{}
				osdSlice := []string{}

				for _, rack := range rackSlice {
					chassisSlice = append(chassisSlice, rackMap[rack]...)
				}

				for _, chassis := range chassisSlice {
					hostSlice = append(hostSlice, chassisMap[chassis]...)
				}

				for _, host := range hostSlice {
					osdSlice = append(osdSlice, hostMap[host]...)

				}

				datacenterDistribution := s.BucketDistribution{Root: rootName, Region: regionName,
					Zone: zoneName, Datacenter: datacenterName, Rack: rackSlice,
					Chassis: chassisSlice, Host: hostSlice, Osd: osdSlice}

				distributionMap[datacenterName] = datacenterDistribution
			}
		}
	case "rack":
		for rack := range rackMap {

			rackName := rack
			datacenterName := GetBucketName(rackName, datacenterMap)
			zoneName := GetBucketName(datacenterName, zoneMap)
			regionName := GetBucketName(zoneName, regionMap)
			rootName := GetBucketName(regionName, rootMap)

			if _, isPresent := distributionMap[rackName]; !isPresent {
				chassisSlice := rackMap[rackName]

				hostSlice := []string{}
				osdSlice := []string{}

				for _, chassis := range chassisSlice {
					hostSlice = append(hostSlice, chassisMap[chassis]...)
				}

				for _, host := range hostSlice {
					osdSlice = append(osdSlice, hostMap[host]...)

				}

				rackDistribution := s.BucketDistribution{Root: rootName, Region: regionName,
					Zone: zoneName, Datacenter: datacenterName, Rack: rackName,
					Chassis: chassisSlice, Host: hostSlice, Osd: osdSlice}

				distributionMap[rackName] = rackDistribution
			}
		}
	case "chassis":
		for chassis := range chassisMap {

			chassisName := chassis
			rackName := GetBucketName(chassisName, rackMap)
			datacenterName := GetBucketName(rackName, datacenterMap)
			zoneName := GetBucketName(datacenterName, zoneMap)
			regionName := GetBucketName(zoneName, regionMap)
			rootName := GetBucketName(regionName, rootMap)

			if _, isPresent := distributionMap[chassisName]; !isPresent {
				hostSlice := chassisMap[chassisName]

				osdSlice := []string{}

				for _, host := range hostSlice {
					osdSlice = append(osdSlice, hostMap[host]...)
				}

				chassisDistribution := s.BucketDistribution{Root: rootName, Region: regionName,
					Zone: zoneName, Datacenter: datacenterName, Rack: rackName,
					Chassis: chassisName, Host: hostSlice, Osd: osdSlice}

				distributionMap[chassisName] = chassisDistribution
			}
		}
	case "host":
		for host := range hostMap {

			hostName := host
			chassisName := GetBucketName(hostName, chassisMap)
			rackName := GetBucketName(chassisName, rackMap)
			datacenterName := GetBucketName(rackName, datacenterMap)
			zoneName := GetBucketName(datacenterName, zoneMap)
			regionName := GetBucketName(zoneName, regionMap)
			rootName := GetBucketName(regionName, rootMap)

			if _, isPresent := distributionMap[hostName]; !isPresent {
				osdSlice := hostMap[hostName]
				hostDistribution := s.BucketDistribution{Root: rootName, Region: regionName,
					Zone: zoneName, Datacenter: datacenterName, Rack: rackName,
					Chassis: chassisName, Host: hostName, Osd: osdSlice}

				distributionMap[hostName] = hostDistribution
			}
		}
	case "osd":
		osds := GetOsds(pgDumpOutput)

		for _, osdId := range osds {
			osd := "osd." + strconv.Itoa(osdId)
			hostName := GetBucketName(osd, hostMap)
			chassisName := GetBucketName(hostName, chassisMap)
			rackName := GetBucketName(chassisName, rackMap)
			datacenterName := GetBucketName(rackName, datacenterMap)
			zoneName := GetBucketName(datacenterName, zoneMap)
			regionName := GetBucketName(zoneName, regionMap)
			rootName := GetBucketName(regionName, rootMap)

			if _, isPresent := distributionMap[osd]; !isPresent {
				osdDistribution := s.BucketDistribution{Root: rootName, Region: regionName,
					Zone: zoneName, Datacenter: datacenterName, Rack: rackName,
					Chassis: chassisName, Host: hostName, Osd: osd}

				distributionMap[osd] = osdDistribution
			}
		}

	}

	return distributionMap
}

func GetNode(nodeName string, osdTreeOutput s.OsdTreeOutputStruct) (nodeStruct s.NodeStruct) {
	nodeOsdTreeStats := osdTreeOutput.Nodes

	for _, node := range nodeOsdTreeStats {
		if node.Name == nodeName {
			nodeStruct = node
		}
	}
	return
}

// ----------------INFORMATIONS GATHERING FUNCTIONS (osd-dump)--------------------------
func GetPool(poolId int, osdDumpOuput s.OsdDumpOutputStruct) (poolStruct s.PoolStruct) {
	poolOsdDumpStats := osdDumpOuput.Pools

	for _, pool := range poolOsdDumpStats {
		if pool.Pool == poolId {
			poolStruct = pool
		}
	}

	return
}

func GetOsdNotIn(osdDumpOutput s.OsdDumpOutputStruct) []string {
	osdStats := osdDumpOutput.Osds
	osdNotIn := []string{}

	for _, osd := range osdStats {
		if osd.In == 0 {
			osdNotIn = append(osdNotIn, "osd."+strconv.Itoa(osd.Osd))
		}
	}
	return osdNotIn
}

// important
func GetErasureCodeProfileInfo(erasureCodeProfileName string, osdDumpOuput s.OsdDumpOutputStruct) (k int, m int, err error) {
	erasureCodeProfilesStruct := osdDumpOuput.ErasureCodeProfiles
	erasureCodeProfiles := reflect.Indirect(reflect.ValueOf(&erasureCodeProfilesStruct))
	numProfiles := erasureCodeProfiles.NumField()

	for i := 0; i < numProfiles; i++ {
		tag, ok := erasureCodeProfiles.Type().Field(i).Tag.Lookup("json")
		if !ok {
			error := errors.New("something went wrong in Tag.Lookup(\"json\")")
			return 0, 0, error
		}

		erasureCodeProfileTag := utility.ExtractTagNameFromString(tag)

		if erasureCodeProfileTag == erasureCodeProfileName {
			erasureCodeProfile := reflect.Indirect(erasureCodeProfiles).Field(i)
			numDataChunk := erasureCodeProfile.FieldByName("K").Interface().(string)
			numCodingChunk := erasureCodeProfile.FieldByName("M").Interface().(string)
			k, err = strconv.Atoi(numDataChunk)
			if err != nil {
				error := errors.New("something went wrong in converting numDataChunk to int")
				return 0, 0, error
			}
			m, err = strconv.Atoi(numCodingChunk)
			if err != nil {
				error := errors.New("something went wrong in converting numCodingChunk to int")
				return 0, 0, error
			}

		}
	}

	return k, m, nil
}

// important
func GetPoolPgNumberMap(pgDumpOutput s.PgDumpOutputStruct, osdDumpOutput s.OsdDumpOutputStruct) (map[string]int, error) {
	poolPgNumberMap := make(map[string]int)

	poolSlice := GetPools(pgDumpOutput)

	for _, poolId := range poolSlice {
		pId, err := strconv.Atoi(poolId)

		if err != nil {
			error := errors.New("something went wrong in converting pool_id")
			return poolPgNumberMap, error
		}

		pool := GetPool(pId, osdDumpOutput)
		pgp_num := pool.PgNumTarget

		poolPgNumberMap[poolId] = pgp_num
	}
	return poolPgNumberMap, nil
}

func GetOsdIdFromOsdNames(osdNames []string) (osdIds []int) {

	for _, name := range osdNames {
		osdId, err := strconv.Atoi(strings.Split(name, ".")[1])
		if err != nil {
			log.Fatal(err)
		}

		osdIds = append(osdIds, osdId)
	}

	return
}

// func GetPgTotalReplicas(givenPg string, osdDumpOutput s.OsdDumpOutputStruct) float64 {
//
// 	stringPoolId := strings.Split(givenPg, ".")[0]
//
// 	poolId, err := strconv.Atoi(stringPoolId)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	pool := GetPool(poolId, osdDumpOutput)
//
// 	numberOfPgReplica := pool.Size
// 	return float64(numberOfPgReplica)
// }

// important
func GetOsdNumberForPg(givenPg string, osdDumpOutput s.OsdDumpOutputStruct) (float64, int, int) {

	stringPoolId := strings.Split(givenPg, ".")[0]

	poolId, err := strconv.Atoi(stringPoolId)
	if err != nil {
		log.Fatal(err)
	}

	pool := GetPool(poolId, osdDumpOutput)

	var k int //number of data chunks
	var m int //number of coding chunks

	osdNumberForPg := pool.Size

	if pool.ErasureCodeProfile == "" {
		//replicated pool
		k = 0
		m = 0
	} else {
		//erasure coded pool
		k, m, err = GetErasureCodeProfileInfo(pool.ErasureCodeProfile, osdDumpOutput)

		if err != nil {
			return 0, 0, 0
		}
	}

	return float64(osdNumberForPg), k, m
}

// important
func GetPublicAddressOsdsMap(osdDumpOutput s.OsdDumpOutputStruct) map[string][]string {
	osds := osdDumpOutput.Osds
	publicAddressOsdsMap := make(map[string][]string)

	for _, osd := range osds {
		publicAddress := strings.Split(osd.PublicAddr, ":")[0]

		if _, isPresent := publicAddressOsdsMap[publicAddress]; !isPresent {
			publicAddressOsdsMap[publicAddress] = []string{}
		}

		publicAddressOsdsMap[publicAddress] = append(publicAddressOsdsMap[publicAddress], "osd."+strconv.Itoa(osd.Osd))
	}

	return publicAddressOsdsMap
}

// -----------------------OSD PG POOL (HOST) MAPPING FUNCTIONS--------------------------
// important
func RiskCalculator(faultBucketOrRouter string, pgDumpOutput s.PgDumpOutputStruct, osdTreeOutput s.OsdTreeOutputStruct, osdDumpOutput s.OsdDumpOutputStruct) ([]int, error) {
	var affectedOsds []string

	if faultBucketOrRouter == "" {

		affectedOsds = []string{}

	} else {
		addressRegex, _ := regexp.Compile("([0-9]{1,3}.){3}[0-9]{1,3}")

		if addressRegex.MatchString(faultBucketOrRouter) { //if name is router (address)

			publicAddressOsdsMap := GetPublicAddressOsdsMap(osdDumpOutput)
			affectedOsds = publicAddressOsdsMap[faultBucketOrRouter]

		} else { //if name is bucket

			faultBucket := GetNode(faultBucketOrRouter, osdTreeOutput)

			if faultBucket.Type == "osd" { //if bucket is osd

				affectedOsds = append(affectedOsds, faultBucketOrRouter)

			} else { //if bucket is one in the hierarchy from host to root

				if reflect.DeepEqual(faultBucket, s.NodeStruct{}) {

					fmt.Println("RiskCalculator: error no bucket with this name ")
					err := errors.New("error no bucket with this name")

					return nil, err

				} else {

					bucketDistributionMap := GetDistributionMap(faultBucket.Type, pgDumpOutput, osdTreeOutput)
					affectedOsds = bucketDistributionMap[faultBucketOrRouter].Osd.([]string)
				}
			}
		}
	}

	affectedOsdIds := GetOsdIdFromOsdNames(affectedOsds)

	return affectedOsdIds, nil
}

// important
func IncrementalRiskCalculator(faultBucketOrRouter []string, pgDumpOutput s.PgDumpOutputStruct, osdTreeOutput s.OsdTreeOutputStruct, osdDumpOutput s.OsdDumpOutputStruct) (map[string]int, []int, error) {
	totalNumberOfAffectedReplicaMap := make(map[string]int)
	totalAffectedOsdIds := []int{}

	for _, faultBucketOrRouterName := range faultBucketOrRouter {
		affectedOsdIds, err := RiskCalculator(faultBucketOrRouterName, pgDumpOutput, osdTreeOutput, osdDumpOutput)
		if err != nil {
			fmt.Println("IncrementalRiskCalculator: Error in RiskCalculator")
			return nil, nil, err
		}
		totalAffectedOsdIds = append(totalAffectedOsdIds, affectedOsdIds...)
	}

	//removing duplicate
	totalAffectedOsdIds = utility.RemoveDuplicateInt(totalAffectedOsdIds)
	fmt.Printf("Total Affected Osds -> %d\n", totalAffectedOsdIds)

	pgNumberOfAffectedReplicaMap := GetPgNumberOfAffectedReplicaMap(totalAffectedOsdIds, pgDumpOutput)

	for pg, numReplicas := range pgNumberOfAffectedReplicaMap {
		if _, isPresent := totalNumberOfAffectedReplicaMap[pg]; !isPresent {
			totalNumberOfAffectedReplicaMap[pg] = numReplicas
		} else {
			totalNumberOfAffectedReplicaMap[pg] += numReplicas
		}
	}

	return totalNumberOfAffectedReplicaMap, totalAffectedOsdIds, nil
}

// -----------------------RISK FAILURE FORECASTING FUNCTION----------------------------
func GetOsdMeanDegradationRatePerWeek(currentTime time.Time, osdInitiationDate time.Time, currentOsdLifeTime float64) (meanDegradationRatePerWeek float64) {
	elapsedTime := currentTime.Sub(osdInitiationDate)
	elapsedWeeks := (elapsedTime.Hours() / 24) / 7

	meanDegradationRatePerWeek = currentOsdLifeTime / elapsedWeeks
	return
}

func GetOsdFaultTimeForecasting(osdInitiationDate time.Time, currentOsdLifeTime float64) (time.Time, float64, error) {
	currentTime := time.Now().UTC()
	meanDegradationRatePerWeek := GetOsdMeanDegradationRatePerWeek(currentTime, osdInitiationDate, currentOsdLifeTime)
	faultPredictionWeeks := (100 - currentOsdLifeTime) / meanDegradationRatePerWeek

	var faultTimePrediction time.Time

	if faultPredictionWeeks >= 0 {

		faultPredictionDays := faultPredictionWeeks * 7
		faultTimePrediction = currentTime.AddDate(0, 0, int(faultPredictionDays))

	} else {
		error := errors.New("something went wrong in prediction - prediction in the past")
		return time.Time{}, 0.0, error
	}

	return faultTimePrediction, meanDegradationRatePerWeek, nil
}

// important
func RiskFailureForecasting(osdLifetimeInfoStructs []s.OsdLifetimeInfo, givenTime time.Time) (osdLifeForecastingMap map[string]float64, warningOsdSlice []string) {

	currentTime := time.Now().UTC()
	timeHaed := givenTime.Sub(currentTime)
	timeHaedInWeeks := (timeHaed.Hours() / 24) / 7

	osdLifeForecastingMap = make(map[string]float64)

	for _, osdLifetimeInfo := range osdLifetimeInfoStructs {
		meanDegradationRate := GetOsdMeanDegradationRatePerWeek(currentTime, osdLifetimeInfo.InitiationDate, osdLifetimeInfo.CurrentOsdLifetime)

		osdLifeForecasting := osdLifetimeInfo.CurrentOsdLifetime + meanDegradationRate*timeHaedInWeeks

		if osdLifeForecasting > 100.0 {
			osdLifeForecasting = 100.0
		}

		osdLifeForecastingMap[osdLifetimeInfo.OsdName] = utility.RoundFloat(osdLifeForecasting, 2)
	}

	for key, value := range osdLifeForecastingMap {
		if value >= 80 {
			warningOsdSlice = append(warningOsdSlice, key)
		}
	}

	return
}
