package pkg

import (
	tools "go/tools"
	"testing"
)

func TestGetPgOsdMap(t *testing.T) {
	//UnMarshalling Json
	CephDumpOutput := tools.ReadJson("../ceph_dump_beautified.json")

	actual := GetPgOsdMap(CephDumpOutput)
	expected := []int{4, 7, 5}

	if tools.SliceIntTestEquality(actual["31.0"], expected) {
		t.Logf("Test Succeeded: expected %d got %d\n\n", expected, actual["31.0"])

	} else {
		t.Errorf("Test Failed: expected %d got %d\n\n", expected, actual["31.0"])
	}

}
