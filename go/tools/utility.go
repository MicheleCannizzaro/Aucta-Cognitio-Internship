package tools

import (
	"encoding/json"
	"fmt"
	s "go/structs"
	"io"
	"log"
	"math"
	"os"
	"time"
)

func ReadJson(jsonFileName string) s.CephDumpOutputStruct {
	jsonFile, err := os.Open(jsonFileName)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	var pgStruc s.CephDumpOutputStruct

	json.Unmarshal(byteValue, &pgStruc)

	return pgStruc
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	slice := []string{}

	for _, item := range strSlice {
		if _, ok := allKeys[item]; !ok {
			allKeys[item] = true
			slice = append(slice, item)
		}
	}
	return slice
}

func RemoveDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, ok := allKeys[item]; !ok {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func Track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func Duration(msg string, start time.Time) {
	log.Printf("%v: %v\n\n", msg, time.Since(start))
}

// func durationTracking(f func(a ...any)) (interface{}, time.Duration) {
// 	start := time.Now()
// 	result := f()
// 	duration := time.Since(start)
// 	return result, duration
// }

func Difference(a, b []string) []string {
	diff := []string{}
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}

func SliceIntTestEquality(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
