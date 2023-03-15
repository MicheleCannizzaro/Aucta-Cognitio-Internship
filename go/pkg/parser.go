package pkg

import (
	"fmt"
	s "go/structs"
	tools "go/tools"
	"strings"
)

// ----------------------------STRUCTS----------------------------
type OsdPoolTuple struct {
	Osd, Pool interface{}
}

// ----------------INFORMATIONS GATHERING FUNCTIONS----------------
func GetPgs(CephDumpOutput s.CephDumpOutputStruct) []string {
	defer tools.Duration(tools.Track("getPgIds"))

	pgStats := CephDumpOutput.PgMap.PgStats
	pgSlice := []string{}

	for _, s := range pgStats {
		pgSlice = append(pgSlice, s.Pgid)
	}

	return pgSlice
}

func GetPools(CephDumpOutput s.CephDumpOutputStruct) []string {
	defer tools.Duration(tools.Track("getPools"))

	pgSlice := GetPgs(CephDumpOutput)
	poolSlice := []string{}

	for _, s := range pgSlice {
		poolId := strings.Split(s, ".")[0]

		poolSlice = append(poolSlice, poolId)
	}

	poolSlice = tools.RemoveDuplicateStr(poolSlice)

	return poolSlice
}

func GetOsds(CephDumpOutput s.CephDumpOutputStruct) []int {
	defer tools.Duration(tools.Track("getOsds"))

	osdPgMap := GetOsdPgMap(CephDumpOutput)
	osds := []int{}

	for osd := range osdPgMap {
		osds = append(osds, osd)
	}

	return osds
}

func GetTotalNumberOfPgs(osdPgNumberMap map[int]int) int {
	defer tools.Duration(tools.Track("getTotalNumberOfPgs"))

	totalPgs := 0

	for _, osdPgCount := range osdPgNumberMap {
		totalPgs += osdPgCount
	}

	return totalPgs
}

func GetOsdsContainingPool(pool string, CephDumpOutput s.CephDumpOutputStruct) []int {
	defer tools.Duration(tools.Track("getOsdsContainingPool"))

	osdPoolPgMap := GetOsdPoolPgMap(CephDumpOutput)

	poolOsdsSlice := []int{}

	for key := range osdPoolPgMap {
		if key.Pool == pool {
			poolOsdsSlice = append(poolOsdsSlice, key.Osd.(int))
		}

	}

	return tools.RemoveDuplicateInt(poolOsdsSlice)
}

func GetOsdsContainingPg(givenPg string, CephDumpOutput s.CephDumpOutputStruct) []int {
	defer tools.Duration(tools.Track("getOsdsContainingPg"))

	osdPoolPgMap := GetOsdPoolPgMap(CephDumpOutput)

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

func GetAllPgInAllOsds(CephDumpOutput s.CephDumpOutputStruct) []string {
	defer tools.Duration(tools.Track("getAllPgInAllOsds"))

	osdPgMap := GetOsdPgMap(CephDumpOutput)

	allPgs := []string{}

	for _, pgSlice := range osdPgMap {
		allPgs = append(allPgs, pgSlice...)
	}

	return allPgs
}

func getPgTotalReplicas(givenPg string, CephDumpOutput s.CephDumpOutputStruct) float64 {
	defer tools.Duration(tools.Track("getPgTotalReplicas"))

	allPg := GetAllPgInAllOsds(CephDumpOutput)

	//how many replicas of givenPg in allPg
	count := 0
	for _, pg := range allPg {
		if pg == givenPg {
			count += 1
		}
	}

	return float64(count)
}

func ExtractPoolsFromPgSlice(pgSlice []string) []string {
	defer tools.Duration(tools.Track("extractPoolsFromPgSlice"))

	poolSlice := []string{}

	for _, pg := range pgSlice {
		poolId := strings.Split(pg, ".")[0]
		poolSlice = append(poolSlice, poolId)
	}
	return tools.RemoveDuplicateStr(poolSlice)
}

func ExtractOsdsFromPoolSlice(poolSlice []string, CephDumpOutput s.CephDumpOutputStruct) []int {
	defer tools.Duration(tools.Track("extractOsdsFromPoolSlice"))
	osdSlice := []int{}

	for _, pool := range poolSlice {
		osdsContainingPool := GetOsdsContainingPool(pool, CephDumpOutput)
		osdSlice = append(osdSlice, osdsContainingPool...)
	}

	return tools.RemoveDuplicateInt(osdSlice)
}

// -----------------------MAPPING FUNCTIONS--------------------------
func GetPgOsdMap(CephDumpOutput s.CephDumpOutputStruct) map[string][]int {
	defer tools.Duration(tools.Track("getPgIdOsdMap"))

	pgStats := CephDumpOutput.PgMap.PgStats
	pgIdOsdMap := make(map[string][]int)

	for _, s := range pgStats {
		pgOsdSlice := s.Up
		pgIdOsdMap[s.Pgid] = pgOsdSlice
	}

	return pgIdOsdMap
}

func GetOsdPgMap(CephDumpOutput s.CephDumpOutputStruct) map[int][]string {
	defer tools.Duration(tools.Track("getOsdPgMap"))

	pgIdOsdMap := GetPgOsdMap(CephDumpOutput)
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
	defer tools.Duration(tools.Track("getNumberOfAssociatedPgsPerOsdMap"))

	osdPgNumberMap := make(map[int]int)

	for key, value := range osdPgMap {
		osdPgNumberMap[key] = len(value)
	}

	return osdPgNumberMap
}

func GetOsdPoolPgMap(CephDumpOutput s.CephDumpOutputStruct) map[OsdPoolTuple][]string {
	defer tools.Duration(tools.Track("getOsdPoolPgMap"))

	osdPgMap := GetOsdPgMap(CephDumpOutput)
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

func GetOsdPoolNumberPerPgsMap(CephDumpOutput s.CephDumpOutputStruct) map[OsdPoolTuple]map[string]int {
	defer tools.Duration(tools.Track("getOsdPoolNumberPerPgsMap"))

	osdPoolPgMap := GetOsdPoolPgMap(CephDumpOutput)
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

func GetPoolWarningPgMap(faultOsdSlice []int, CephDumpOutput s.CephDumpOutputStruct) map[string][]string {
	defer tools.Duration(tools.Track("getPoolWarningPgMap"))

	poolWarningPgMap := make(map[string][]string)

	_, _, warningPgs, _ := GetPgsWithHighProbabilityOfLosingData(faultOsdSlice, CephDumpOutput)

	for _, wPgs := range warningPgs {
		poolId := strings.Split(wPgs, ".")[0]

		if _, ok := poolWarningPgMap[poolId]; !ok {
			poolWarningPgMap[poolId] = []string{}
		}

		poolWarningPgMap[poolId] = append(poolWarningPgMap[poolId], wPgs)
	}

	return poolWarningPgMap
}

func GetPoolLostPgMap(faultOsdSlice []int, CephDumpOutput s.CephDumpOutputStruct) map[string][]string {
	defer tools.Duration(tools.Track("getPoolLostPgMap"))

	poolLostPgMap := make(map[string][]string)

	_, _, _, lostPgs := GetPgsWithHighProbabilityOfLosingData(faultOsdSlice, CephDumpOutput)

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

func GetAffectedPools(faultOsd int, CephDumpOutput s.CephDumpOutputStruct) []string {
	defer tools.Duration(tools.Track("getAffectedPools"))

	osdPgMap := GetOsdPgMap(CephDumpOutput)

	affectedPoolSlice := []string{}

	for _, pg := range osdPgMap[faultOsd] {
		poolId := strings.Split(pg, ".")[0]

		affectedPoolSlice = append(affectedPoolSlice, poolId)
	}

	return tools.RemoveDuplicateStr(affectedPoolSlice)
}

func GetAffectedPgs(faultOsd int, CephDumpOutput s.CephDumpOutputStruct) []string {
	defer tools.Duration(tools.Track("getAffectedPgs"))

	osdPgMap := GetOsdPgMap(CephDumpOutput)

	affectedPgSlice := []string{}

	affectedPgSlice = append(affectedPgSlice, osdPgMap[faultOsd]...)

	return affectedPgSlice
}

func getNumAffectedReplicasPgItem(faultOsd int, givenPg string, CephDumpOutput s.CephDumpOutputStruct) float64 {
	defer tools.Duration(tools.Track("getNumAffectedReplicasPgItem"))

	affectedPgs := GetAffectedPgs(faultOsd, CephDumpOutput)

	count := 0

	for _, pg := range affectedPgs {
		if pg == givenPg {
			count += 1
		}
	}

	return float64(count)
}

func PercentageCalculationAffectedReplicasPg(faultOsd int, givenPg string, CephDumpOutput s.CephDumpOutputStruct) float64 {
	defer tools.Duration(tools.Track("percentageCalculationAffectedReplicasPg"))

	numAffectedReplicaPgItem := getNumAffectedReplicasPgItem(faultOsd, givenPg, CephDumpOutput)

	pgTotalReplicas := getPgTotalReplicas(givenPg, CephDumpOutput)

	percentage := (numAffectedReplicaPgItem / pgTotalReplicas) * 100

	return tools.RoundFloat(percentage, 2)
}

func WarningCheck(percentage float64, faultOsd int, givenPg string, CephDumpOutput s.CephDumpOutputStruct) {
	defer tools.Duration(tools.Track("warningCheck"))
	//log.SetOutput(io.Discard) //disable log

	osds := GetOsds(CephDumpOutput)
	warningOsdSlice := make(map[int]float64)

	if percentage > 0.0 {
		for osd := range osds {
			percentage = PercentageCalculationAffectedReplicasPg(osd, givenPg, CephDumpOutput)
			if percentage > 0.0 && osd != faultOsd {
				warningOsdSlice[osd] = percentage
			}
		}
	}

	//log.SetOutput(os.Stdout) //renable log
	fmt.Printf("Warning: check these other OSDs and percentages for %s ->%v\n", givenPg, warningOsdSlice)
}

func GetTotalAffectedPgsAndPools(faultOsdSlice []int, CephDumpOutput s.CephDumpOutputStruct) ([]string, []string) {
	defer tools.Duration(tools.Track("getTotalAffectedPgsAndPools"))

	totalAffectedPgs := []string{}
	totalAffectedPools := []string{}

	for _, osd := range faultOsdSlice {
		affectedPgs := GetAffectedPgs(osd, CephDumpOutput)
		affectedPools := GetAffectedPools(osd, CephDumpOutput)

		totalAffectedPgs = append(totalAffectedPgs, affectedPgs...)
		totalAffectedPools = append(totalAffectedPools, affectedPools...)
	}

	totalAffectedPools = tools.RemoveDuplicateStr(totalAffectedPools)

	return totalAffectedPgs, totalAffectedPools
}

func GetPgNumberOfAffectedReplicaMap(faultOsdSlice []int, CephDumpOutput s.CephDumpOutputStruct) map[string]int {
	defer tools.Duration(tools.Track("getPgNumberOfAffectedReplicaMap"))

	totalAffectedPgs, _ := GetTotalAffectedPgsAndPools(faultOsdSlice, CephDumpOutput)
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

func GetPgsWithHighProbabilityOfLosingData(faultOsdSlice []int, CephDumpOutput s.CephDumpOutputStruct) ([]string, []string, []string, []string) {
	defer tools.Duration(tools.Track("getPgsWithHighProbabilityOfLosingData"))

	pgNumberOfAffectedReplicaMap := GetPgNumberOfAffectedReplicaMap(faultOsdSlice, CephDumpOutput)

	goodPgs := []string{}
	warningPgs := []string{}
	lostPgs := []string{}

	for pg, numLostReplicas := range pgNumberOfAffectedReplicaMap {
		pgTotalReplicas := getPgTotalReplicas(pg, CephDumpOutput)
		percetageLostReplicas := (float64(numLostReplicas) / pgTotalReplicas) * 100

		//fmt.Println(percetageLostReplicas)
		if percetageLostReplicas > 0.0 && percetageLostReplicas < 50.0 {
			goodPgs = append(goodPgs, pg)
		} else if percetageLostReplicas >= 50.0 && percetageLostReplicas < 100.0 {
			warningPgs = append(warningPgs, pg)
		} else if percetageLostReplicas >= 50.0 && percetageLostReplicas == 100.0 {
			lostPgs = append(lostPgs, pg)
		}
	}

	compromisedPgs := [][]string{goodPgs, warningPgs, lostPgs}
	allCompromisedPgs := []string{}

	for _, ele := range compromisedPgs {
		allCompromisedPgs = append(allCompromisedPgs, ele...)
	}

	inHealthPgs := tools.Difference(tools.RemoveDuplicateStr(GetAllPgInAllOsds(CephDumpOutput)), allCompromisedPgs)

	return inHealthPgs, goodPgs, warningPgs, lostPgs
}

func GetTotalAffectedPgsAndOsds(faultPoolSlice []string, CephDumpOutput s.CephDumpOutputStruct) ([]string, []int) {
	defer tools.Duration(tools.Track("getTotalAffectedPgsAndOsds"))

	totalAffectedPgs := []string{}
	totalAffectedOsds := []int{}

	//osdPoolPgMap := getOsdPoolPgMap(CephDumpOutput)

	for _, pool := range faultPoolSlice {
		totalAffectedOsds = append(totalAffectedOsds, GetOsdsContainingPool(pool, CephDumpOutput)...)
	}

	totalAffectedOsds = tools.RemoveDuplicateInt(totalAffectedOsds)

	for _, osd := range totalAffectedOsds {
		totalAffectedPgs = append(totalAffectedPgs, GetAffectedPgs(osd, CephDumpOutput)...)
	}

	return totalAffectedPgs, totalAffectedOsds
}
