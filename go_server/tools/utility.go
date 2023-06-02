package utility

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"regexp"
	"time"

	s "github.com/MicheleCannizzaro/Aucta-Cognitio-Internship/go_server/structs"
)

func ReadPgDumpJson(jsonFileName string) s.PgDumpOutputStruct {
	jsonFile, err := os.Open(jsonFileName)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	var pgStruc s.PgDumpOutputStruct

	json.Unmarshal(byteValue, &pgStruc)

	return pgStruc
}

func ReadOsdTreeJson(jsonFileName string) s.OsdTreeOutputStruct {
	jsonFile, err := os.Open(jsonFileName)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	var osdStruc s.OsdTreeOutputStruct

	json.Unmarshal(byteValue, &osdStruc)

	return osdStruc
}

func ReadOsdDumpJson(jsonFileName string) s.OsdDumpOutputStruct {
	jsonFile, err := os.Open(jsonFileName)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	var osdStruc s.OsdDumpOutputStruct

	json.Unmarshal(byteValue, &osdStruc)

	return osdStruc
}

func RmvDuplStr(strSlice []string) []string {
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

func IntDifference(a, b []int) []int {
	diff := []int{}
	m := make(map[int]bool)

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

func GetFloatRandomNumber(min int, max int) (randomNumber float64) {
	rand.Seed(time.Now().UTC().UnixNano())

	randomNumber = RoundFloat(float64(rand.Intn(max-min+1)+min), 2)
	return
}

func GetIntRandomNumber(min int, max int) (randomNumber int) {
	rand.Seed(time.Now().UTC().UnixNano())

	randomNumber = rand.Intn(max-min+1) + min
	return
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// works on tag like that (json:"pool3_ec_profile,omitempty") to return pool3_ec_profile
func ExtractTagNameFromString(tag string) string {
	tagRegex, _ := regexp.Compile("^[^,]*")
	tagValue := tagRegex.FindStringSubmatch(tag)
	return tagValue[0]
}
