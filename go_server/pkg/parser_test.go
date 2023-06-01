package pkg

import (
	"fmt"
	"testing"

	utility "github.com/MicheleCannizzaro/Aucta-Cognitio-Internship/go_server/tools"
)

func TestGetPgOsdMap(t *testing.T) {
	//UnMarshalling Json
	pgDumpOutput := utility.ReadPgDumpJson("../pg_dump.json")

	actual := GetPgOsdMap(pgDumpOutput)
	expected := []int{4, 7, 5}

	if utility.SliceIntTestEquality(actual["31.0"], expected) {
		t.Logf("Test Succeeded: expected %d got %d\n\n", expected, actual["31.0"])

	} else {
		t.Errorf("Test Failed: expected %d got %d\n\n", expected, actual["31.0"])
	}

}

func TestIncrementalRiskCalculator(t *testing.T) {
	//UnMarshalling Json
	pgDumpOutput := utility.ReadPgDumpJson("../pg_dump.json")
	osdTreeOutput := utility.ReadOsdTreeJson("../osd-tree.json")
	osdDumpOutput := utility.ReadOsdDumpJson("../osd_dump.json")

	hosts := []string{"sv81", "sv82", "sv61", "newnamefor-sv53"}
	t.Logf("hosts -> %s\n\n", hosts)
	incrementalPgAffectedReplicaMap, err := IncrementalRiskCalculator(hosts, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		t.Errorf("Error in Incremental RiskCalculator")
	}
	inHealthPgs, goodPgs, warningPgs, lostPgs := GetPgsWithHighProbabilityOfLosingData(incrementalPgAffectedReplicaMap, pgDumpOutput, osdDumpOutput)
	//t.Logf("%v\n\n", incrementalPgAffectedReplicaMap)

	actualInHealthPgsLen := len(inHealthPgs)
	expectedZeroLen := 0

	actualGoodPgsLen := len(goodPgs)
	actualWarningPgsLen := len(warningPgs)
	actualLostPgsLen := len(lostPgs)

	if actualInHealthPgsLen == expectedZeroLen {
		t.Logf("(inHealthPgs) Test Succeeded: expected %d got %d\n\n", expectedZeroLen, actualInHealthPgsLen)
	} else {
		t.Errorf("(inHealthPgs) Test Failed: expected %d got %d\n\n", expectedZeroLen, actualInHealthPgsLen)
	}

	if actualGoodPgsLen == expectedZeroLen {
		t.Logf("(goodPgs) Test Succeeded: expected %d got %d\n\n", expectedZeroLen, actualInHealthPgsLen)
	} else {
		t.Errorf("(goodPgs) Test Failed expected %d got %d\n\n", expectedZeroLen, actualInHealthPgsLen)
	}

	if actualWarningPgsLen == 1 {
		t.Logf("(warningPgs) Test Succeeded: expected %d got %d\n\n", 1, actualWarningPgsLen)
	} else {
		t.Errorf("(warningPgs) Test Failed:  expected %d got %d\n\n", 1, actualWarningPgsLen)
	}

	if actualLostPgsLen == 736 {
		t.Logf("(lostPgs) Test Succeeded: expected %d got %d\n\n", 736, actualLostPgsLen)

	} else {
		t.Errorf("(lostPgs) Test Failed: expected %d got %d\n\n", 736, actualLostPgsLen)
	}

}

func TestPoolDataLossProbability(t *testing.T) {
	//UnMarshalling Json
	pgDumpOutput := utility.ReadPgDumpJson("../pg_dump.json")
	osdTreeOutput := utility.ReadOsdTreeJson("../osd-tree.json")
	osdDumpOutput := utility.ReadOsdDumpJson("../osd_dump.json")

	faults := []string{"default"} //rout bucket
	actualProbability, err := GetPoolDataLossProbability(faults, pgDumpOutput, osdTreeOutput, osdDumpOutput)
	if err != nil {
		t.Errorf("Error in GetPoolDataLossProbability")
	}

	//fmt.Printf("actualProbability %v", actualProbability)
	expectedProbability := 1.0

	for poolName := range actualProbability {

		if actualProbability[poolName] == expectedProbability {
			fmt.Println(poolName)
			t.Logf("%s Test Succeeded: expected %v got %v\n\n", poolName, expectedProbability, actualProbability[poolName])
		} else {
			fmt.Println(poolName)
			t.Errorf("%s Test Failed: expected %v got %v\n\n", poolName, expectedProbability, actualProbability[poolName])
		}

	}
}
