package pkg

import (
	"errors"
	"fmt"
	s "go/structs"
	utility "go/tools"
	"log"
	"strconv"
	"strings"
	"time"
)

// ----------------------------STRUCTS----------------------------
type OsdPoolTuple struct {
	Osd, Pool interface{}
}

// ----------------INFORMATIONS GATHERING FUNCTIONS (pg dump)----------------
func GetPgs(pgDumpOutput s.PgDumpOutputStruct) []string {
	//defer utility.Duration(utility.Track("getPgIds"))

	pgStats := pgDumpOutput.PgMap.PgStats
	pgSlice := []string{}

	for _, s := range pgStats {
		pgSlice = append(pgSlice, s.Pgid)
	}

	return pgSlice
}

func GetPools(pgDumpOutput s.PgDumpOutputStruct) []string {
	//defer utility.Duration(utility.Track("getPools"))

	pgSlice := GetPgs(pgDumpOutput)
	poolSlice := []string{}

	for _, s := range pgSlice {
		poolId := strings.Split(s, ".")[0]

		poolSlice = append(poolSlice, poolId)
	}

	poolSlice = utility.RemoveDuplicateStr(poolSlice)

	return poolSlice
}

func GetOsds(pgDumpOutput s.PgDumpOutputStruct) []int {
	//defer utility.Duration(utility.Track("getOsds"))

	osdPgMap := GetOsdPgMap(pgDumpOutput)
	osds := []int{}

	for osd := range osdPgMap {
		osds = append(osds, osd)
	}

	return osds
}

func GetTotalNumberOfPgs(osdPgNumberMap map[int]int) int {
	//defer utility.Duration(utility.Track("getTotalNumberOfPgs"))

	totalPgs := 0

	for _, osdPgCount := range osdPgNumberMap {
		totalPgs += osdPgCount
	}

	return totalPgs
}

func GetOsdsContainingPool(pool string, pgDumpOutput s.PgDumpOutputStruct) []int {
	//defer utility.Duration(utility.Track("getOsdsContainingPool"))

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
	//defer utility.Duration(utility.Track("getOsdsContainingPg"))

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
	//defer utility.Duration(utility.Track("getAllPgInAllOsds"))

	osdPgMap := GetOsdPgMap(pgDumpOutput)

	allPgs := []string{}

	for _, pgSlice := range osdPgMap {
		allPgs = append(allPgs, pgSlice...)
	}

	return allPgs
}

// func GetPgTotalReplicas(givenPg string, pgDumpOutput s.PgDumpOutputStruct) float64 {
// 	//defer utility.Duration(utility.Track("getPgTotalReplicas"))
//
// 	allPg := GetAllPgInAllOsds(pgDumpOutput)
//
// 	//how many replicas of givenPg in allPg
// 	count := 0
// 	for _, pg := range allPg {
// 		if pg == givenPg {
// 			count += 1
// 		}
// 	}
//
// 	return float64(count)
// }

func ExtractPoolsFromPgSlice(pgSlice []string) []string {
	//defer utility.Duration(utility.Track("extractPoolsFromPgSlice"))

	poolSlice := []string{}

	for _, pg := range pgSlice {
		poolId := strings.Split(pg, ".")[0]
		poolSlice = append(poolSlice, poolId)
	}
	return utility.RemoveDuplicateStr(poolSlice)
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
	defer utility.Duration(utility.Track("getOsdPoolPgMap"))

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

func GetAffectedPools(faultOsd int, pgDumpOutput s.PgDumpOutputStruct) []string {
	//defer utility.Duration(utility.Track("getAffectedPools"))

	osdPgMap := GetOsdPgMap(pgDumpOutput)

	affectedPoolSlice := []string{}

	for _, pg := range osdPgMap[faultOsd] {
		poolId := strings.Split(pg, ".")[0]

		affectedPoolSlice = append(affectedPoolSlice, poolId)
	}

	return utility.RemoveDuplicateStr(affectedPoolSlice)
}

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

	//pgTotalReplicas := GetPgTotalReplicas(givenPg, pgDumpOutput)
	pgTotalReplicas := GetPgTotalReplicas(givenPg, osdDumpOutput)

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

	totalAffectedPools = utility.RemoveDuplicateStr(totalAffectedPools)

	return totalAffectedPgs, totalAffectedPools
}

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

func GetPgsWithHighProbabilityOfLosingData(pgNumberOfAffectedReplicaMap map[string]int, pgDumpOutput s.PgDumpOutputStruct, osdDumpOutput s.OsdDumpOutputStruct) (inHealthPgs []string, goodPgs []string, warningPgs []string, lostPgs []string) {
	//defer utility.Duration(utility.Track("getPgsWithHighProbabilityOfLosingData"))

	for pg, numLostReplicas := range pgNumberOfAffectedReplicaMap {
		//pgTotalReplicas := GetPgTotalReplicas(pg, pgDumpOutput)
		pgTotalReplicas := GetPgTotalReplicas(pg, osdDumpOutput)
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

	inHealthPgs = utility.Difference(utility.RemoveDuplicateStr(GetAllPgInAllOsds(pgDumpOutput)), allCompromisedPgs)

	return inHealthPgs, goodPgs, warningPgs, lostPgs
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

func GetPgTotalReplicas(givenPg string, osdDumpOutput s.OsdDumpOutputStruct) float64 {

	stringPoolId := strings.Split(givenPg, ".")[0]

	poolId, err := strconv.Atoi(stringPoolId)
	if err != nil {
		log.Fatal(err)
	}

	pool := GetPool(poolId, osdDumpOutput)

	numberOfPgReplica := pool.Size
	return float64(numberOfPgReplica)
}

// -----------------------OSD PG POOL (HOST) MAPPING FUNCTIONS--------------------------
func RiskCalculator(faultBucketName string, pgDumpOutput s.PgDumpOutputStruct, osdTreeOutput s.OsdTreeOutputStruct) map[string]int {

	faultBucket := GetNode(faultBucketName, osdTreeOutput)

	bucketDistributionMap := GetDistributionMap(faultBucket.Type, pgDumpOutput, osdTreeOutput)
	affectedOsds := bucketDistributionMap[faultBucketName].Osd.([]string)

	affectedOsdIds := GetOsdIdFromOsdNames(affectedOsds)

	pgNumberOfAffectedReplicaMap := GetPgNumberOfAffectedReplicaMap(affectedOsdIds, pgDumpOutput)
	return pgNumberOfAffectedReplicaMap
}

// important
func IncrementalRiskCalculator(faultBucketNames []string, pgDumpOutput s.PgDumpOutputStruct, osdTreeOutput s.OsdTreeOutputStruct) map[string]int {
	totalNumberOfAffectedReplicaMap := make(map[string]int)

	for _, faultBucketName := range faultBucketNames {
		NumberOfAffectedReplicaMap := RiskCalculator(faultBucketName, pgDumpOutput, osdTreeOutput)

		for pg, numReplicas := range NumberOfAffectedReplicaMap {
			if _, isPresent := totalNumberOfAffectedReplicaMap[pg]; !isPresent {
				totalNumberOfAffectedReplicaMap[pg] = numReplicas
			} else {
				totalNumberOfAffectedReplicaMap[pg] += numReplicas
			}
		}
	}

	return totalNumberOfAffectedReplicaMap
}

// -----------------------RISK FAILURE FORECASTING FUNCTION----------------------------
func GetOsdMeanDegradationRatePerWeek(currentTime time.Time, osdInitiationDate time.Time, currentOsdLifeTime float64) (meanDegradationRatePerWeek float64) {
	elapsedTime := currentTime.Sub(osdInitiationDate)
	elapsedWeeks := (elapsedTime.Hours() / 24) / 7

	meanDegradationRatePerWeek = currentOsdLifeTime / elapsedWeeks
	return
}

func GetOsdFaultTimePrediction(osdInitiationDate time.Time, currentOsdLifeTime float64) (time.Time, float64, error) {
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

func RiskFailureForecasting(osdMap map[string]map[string]interface{}, givenTime time.Time) (osdLifeForecastingMap map[string]float64, warningOsdSlice []string) {

	currentTime := time.Now().UTC()
	timeHaed := givenTime.Sub(currentTime)
	timeHaedInWeeks := (timeHaed.Hours() / 24) / 7

	osdLifeForecastingMap = make(map[string]float64)

	for osdName, value := range osdMap {
		meanDegradationRate := GetOsdMeanDegradationRatePerWeek(currentTime, value["initiationDate"].(time.Time), value["currentOsdLifeTime"].(float64))

		osdLifeForecasting := value["currentOsdLifeTime"].(float64) + meanDegradationRate*timeHaedInWeeks

		if osdLifeForecasting > 100.0 {
			osdLifeForecasting = 100.0
		}

		osdLifeForecastingMap[osdName] = utility.RoundFloat(osdLifeForecasting, 2)
	}

	for key, value := range osdLifeForecastingMap {
		if value >= 80 {
			warningOsdSlice = append(warningOsdSlice, key)
		}
	}

	return
}
