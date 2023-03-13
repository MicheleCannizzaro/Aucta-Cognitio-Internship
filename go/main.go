package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

//----------------STRUCTS----------------

type CephDumpOutputStruct struct {
	PgReady bool `json:"pg_ready"`
	PgMap   struct {
		Version         int    `json:"version"`
		Stamp           string `json:"stamp"`
		LastOsdmapEpoch int    `json:"last_osdmap_epoch"`
		LastPgScan      int    `json:"last_pg_scan"`
		PgStatsSum      struct {
			StatSum struct {
				NumBytes                   int64 `json:"num_bytes"`
				NumObjects                 int   `json:"num_objects"`
				NumObjectClones            int   `json:"num_object_clones"`
				NumObjectCopies            int   `json:"num_object_copies"`
				NumObjectsMissingOnPrimary int   `json:"num_objects_missing_on_primary"`
				NumObjectsMissing          int   `json:"num_objects_missing"`
				NumObjectsDegraded         int   `json:"num_objects_degraded"`
				NumObjectsMisplaced        int   `json:"num_objects_misplaced"`
				NumObjectsUnfound          int   `json:"num_objects_unfound"`
				NumObjectsDirty            int   `json:"num_objects_dirty"`
				NumWhiteouts               int   `json:"num_whiteouts"`
				NumRead                    int   `json:"num_read"`
				NumReadKb                  int   `json:"num_read_kb"`
				NumWrite                   int   `json:"num_write"`
				NumWriteKb                 int   `json:"num_write_kb"`
				NumScrubErrors             int   `json:"num_scrub_errors"`
				NumShallowScrubErrors      int   `json:"num_shallow_scrub_errors"`
				NumDeepScrubErrors         int   `json:"num_deep_scrub_errors"`
				NumObjectsRecovered        int   `json:"num_objects_recovered"`
				NumBytesRecovered          int64 `json:"num_bytes_recovered"`
				NumKeysRecovered           int   `json:"num_keys_recovered"`
				NumObjectsOmap             int   `json:"num_objects_omap"`
				NumObjectsHitSetArchive    int   `json:"num_objects_hit_set_archive"`
				NumBytesHitSetArchive      int   `json:"num_bytes_hit_set_archive"`
				NumFlush                   int   `json:"num_flush"`
				NumFlushKb                 int   `json:"num_flush_kb"`
				NumEvict                   int   `json:"num_evict"`
				NumEvictKb                 int   `json:"num_evict_kb"`
				NumPromote                 int   `json:"num_promote"`
				NumFlushModeHigh           int   `json:"num_flush_mode_high"`
				NumFlushModeLow            int   `json:"num_flush_mode_low"`
				NumEvictModeSome           int   `json:"num_evict_mode_some"`
				NumEvictModeFull           int   `json:"num_evict_mode_full"`
				NumObjectsPinned           int   `json:"num_objects_pinned"`
				NumLegacySnapsets          int   `json:"num_legacy_snapsets"`
				NumLargeOmapObjects        int   `json:"num_large_omap_objects"`
				NumObjectsManifest         int   `json:"num_objects_manifest"`
				NumOmapBytes               int   `json:"num_omap_bytes"`
				NumOmapKeys                int   `json:"num_omap_keys"`
				NumObjectsRepaired         int   `json:"num_objects_repaired"`
			} `json:"stat_sum"`
			StoreStats struct {
				Total                   int `json:"total"`
				Available               int `json:"available"`
				InternallyReserved      int `json:"internally_reserved"`
				Allocated               int `json:"allocated"`
				DataStored              int `json:"data_stored"`
				DataCompressed          int `json:"data_compressed"`
				DataCompressedAllocated int `json:"data_compressed_allocated"`
				DataCompressedOriginal  int `json:"data_compressed_original"`
				OmapAllocated           int `json:"omap_allocated"`
				InternalMetadata        int `json:"internal_metadata"`
			} `json:"store_stats"`
			LogSize       int `json:"log_size"`
			OndiskLogSize int `json:"ondisk_log_size"`
			Up            int `json:"up"`
			Acting        int `json:"acting"`
			NumStoreStats int `json:"num_store_stats"`
		} `json:"pg_stats_sum"`
		OsdStatsSum struct {
			UpFrom             int   `json:"up_from"`
			Seq                int   `json:"seq"`
			NumPgs             int   `json:"num_pgs"`
			NumOsds            int   `json:"num_osds"`
			NumPerPoolOsds     int   `json:"num_per_pool_osds"`
			NumPerPoolOmapOsds int   `json:"num_per_pool_omap_osds"`
			Kb                 int64 `json:"kb"`
			KbUsed             int   `json:"kb_used"`
			KbUsedData         int   `json:"kb_used_data"`
			KbUsedOmap         int   `json:"kb_used_omap"`
			KbUsedMeta         int   `json:"kb_used_meta"`
			KbAvail            int64 `json:"kb_avail"`
			Statfs             struct {
				Total                   int64 `json:"total"`
				Available               int64 `json:"available"`
				InternallyReserved      int64 `json:"internally_reserved"`
				Allocated               int64 `json:"allocated"`
				DataStored              int64 `json:"data_stored"`
				DataCompressed          int   `json:"data_compressed"`
				DataCompressedAllocated int   `json:"data_compressed_allocated"`
				DataCompressedOriginal  int   `json:"data_compressed_original"`
				OmapAllocated           int   `json:"omap_allocated"`
				InternalMetadata        int64 `json:"internal_metadata"`
			} `json:"statfs"`
			HbPeers           []any `json:"hb_peers"`
			SnapTrimQueueLen  int   `json:"snap_trim_queue_len"`
			NumSnapTrimming   int   `json:"num_snap_trimming"`
			NumShardsRepaired int   `json:"num_shards_repaired"`
			OpQueueAgeHist    struct {
				Histogram  []int `json:"histogram"`
				UpperBound int   `json:"upper_bound"`
			} `json:"op_queue_age_hist"`
			PerfStat struct {
				CommitLatencyMs int `json:"commit_latency_ms"`
				ApplyLatencyMs  int `json:"apply_latency_ms"`
				CommitLatencyNs int `json:"commit_latency_ns"`
				ApplyLatencyNs  int `json:"apply_latency_ns"`
			} `json:"perf_stat"`
			Alerts           []any `json:"alerts"`
			NetworkPingTimes []any `json:"network_ping_times"`
		} `json:"osd_stats_sum"`
		PgStatsDelta struct {
			StatSum struct {
				NumBytes                   int `json:"num_bytes"`
				NumObjects                 int `json:"num_objects"`
				NumObjectClones            int `json:"num_object_clones"`
				NumObjectCopies            int `json:"num_object_copies"`
				NumObjectsMissingOnPrimary int `json:"num_objects_missing_on_primary"`
				NumObjectsMissing          int `json:"num_objects_missing"`
				NumObjectsDegraded         int `json:"num_objects_degraded"`
				NumObjectsMisplaced        int `json:"num_objects_misplaced"`
				NumObjectsUnfound          int `json:"num_objects_unfound"`
				NumObjectsDirty            int `json:"num_objects_dirty"`
				NumWhiteouts               int `json:"num_whiteouts"`
				NumRead                    int `json:"num_read"`
				NumReadKb                  int `json:"num_read_kb"`
				NumWrite                   int `json:"num_write"`
				NumWriteKb                 int `json:"num_write_kb"`
				NumScrubErrors             int `json:"num_scrub_errors"`
				NumShallowScrubErrors      int `json:"num_shallow_scrub_errors"`
				NumDeepScrubErrors         int `json:"num_deep_scrub_errors"`
				NumObjectsRecovered        int `json:"num_objects_recovered"`
				NumBytesRecovered          int `json:"num_bytes_recovered"`
				NumKeysRecovered           int `json:"num_keys_recovered"`
				NumObjectsOmap             int `json:"num_objects_omap"`
				NumObjectsHitSetArchive    int `json:"num_objects_hit_set_archive"`
				NumBytesHitSetArchive      int `json:"num_bytes_hit_set_archive"`
				NumFlush                   int `json:"num_flush"`
				NumFlushKb                 int `json:"num_flush_kb"`
				NumEvict                   int `json:"num_evict"`
				NumEvictKb                 int `json:"num_evict_kb"`
				NumPromote                 int `json:"num_promote"`
				NumFlushModeHigh           int `json:"num_flush_mode_high"`
				NumFlushModeLow            int `json:"num_flush_mode_low"`
				NumEvictModeSome           int `json:"num_evict_mode_some"`
				NumEvictModeFull           int `json:"num_evict_mode_full"`
				NumObjectsPinned           int `json:"num_objects_pinned"`
				NumLegacySnapsets          int `json:"num_legacy_snapsets"`
				NumLargeOmapObjects        int `json:"num_large_omap_objects"`
				NumObjectsManifest         int `json:"num_objects_manifest"`
				NumOmapBytes               int `json:"num_omap_bytes"`
				NumOmapKeys                int `json:"num_omap_keys"`
				NumObjectsRepaired         int `json:"num_objects_repaired"`
			} `json:"stat_sum"`
			StoreStats struct {
				Total                   int `json:"total"`
				Available               int `json:"available"`
				InternallyReserved      int `json:"internally_reserved"`
				Allocated               int `json:"allocated"`
				DataStored              int `json:"data_stored"`
				DataCompressed          int `json:"data_compressed"`
				DataCompressedAllocated int `json:"data_compressed_allocated"`
				DataCompressedOriginal  int `json:"data_compressed_original"`
				OmapAllocated           int `json:"omap_allocated"`
				InternalMetadata        int `json:"internal_metadata"`
			} `json:"store_stats"`
			LogSize       int    `json:"log_size"`
			OndiskLogSize int    `json:"ondisk_log_size"`
			Up            int    `json:"up"`
			Acting        int    `json:"acting"`
			NumStoreStats int    `json:"num_store_stats"`
			StampDelta    string `json:"stamp_delta"`
		} `json:"pg_stats_delta"`
		PgStats []struct {
			Pgid                    string `json:"pgid"`
			Version                 string `json:"version"`
			ReportedSeq             int    `json:"reported_seq"`
			ReportedEpoch           int    `json:"reported_epoch"`
			State                   string `json:"state"`
			LastFresh               string `json:"last_fresh"`
			LastChange              string `json:"last_change"`
			LastActive              string `json:"last_active"`
			LastPeered              string `json:"last_peered"`
			LastClean               string `json:"last_clean"`
			LastBecameActive        string `json:"last_became_active"`
			LastBecamePeered        string `json:"last_became_peered"`
			LastUnstale             string `json:"last_unstale"`
			LastUndegraded          string `json:"last_undegraded"`
			LastFullsized           string `json:"last_fullsized"`
			MappingEpoch            int    `json:"mapping_epoch"`
			LogStart                string `json:"log_start"`
			OndiskLogStart          string `json:"ondisk_log_start"`
			Created                 int    `json:"created"`
			LastEpochClean          int    `json:"last_epoch_clean"`
			Parent                  string `json:"parent"`
			ParentSplitBits         int    `json:"parent_split_bits"`
			LastScrub               string `json:"last_scrub"`
			LastScrubStamp          string `json:"last_scrub_stamp"`
			LastDeepScrub           string `json:"last_deep_scrub"`
			LastDeepScrubStamp      string `json:"last_deep_scrub_stamp"`
			LastCleanScrubStamp     string `json:"last_clean_scrub_stamp"`
			LogSize                 int    `json:"log_size"`
			OndiskLogSize           int    `json:"ondisk_log_size"`
			StatsInvalid            bool   `json:"stats_invalid"`
			DirtyStatsInvalid       bool   `json:"dirty_stats_invalid"`
			OmapStatsInvalid        bool   `json:"omap_stats_invalid"`
			HitsetStatsInvalid      bool   `json:"hitset_stats_invalid"`
			HitsetBytesStatsInvalid bool   `json:"hitset_bytes_stats_invalid"`
			PinStatsInvalid         bool   `json:"pin_stats_invalid"`
			ManifestStatsInvalid    bool   `json:"manifest_stats_invalid"`
			SnaptrimqLen            int    `json:"snaptrimq_len"`
			StatSum                 struct {
				NumBytes                   int `json:"num_bytes"`
				NumObjects                 int `json:"num_objects"`
				NumObjectClones            int `json:"num_object_clones"`
				NumObjectCopies            int `json:"num_object_copies"`
				NumObjectsMissingOnPrimary int `json:"num_objects_missing_on_primary"`
				NumObjectsMissing          int `json:"num_objects_missing"`
				NumObjectsDegraded         int `json:"num_objects_degraded"`
				NumObjectsMisplaced        int `json:"num_objects_misplaced"`
				NumObjectsUnfound          int `json:"num_objects_unfound"`
				NumObjectsDirty            int `json:"num_objects_dirty"`
				NumWhiteouts               int `json:"num_whiteouts"`
				NumRead                    int `json:"num_read"`
				NumReadKb                  int `json:"num_read_kb"`
				NumWrite                   int `json:"num_write"`
				NumWriteKb                 int `json:"num_write_kb"`
				NumScrubErrors             int `json:"num_scrub_errors"`
				NumShallowScrubErrors      int `json:"num_shallow_scrub_errors"`
				NumDeepScrubErrors         int `json:"num_deep_scrub_errors"`
				NumObjectsRecovered        int `json:"num_objects_recovered"`
				NumBytesRecovered          int `json:"num_bytes_recovered"`
				NumKeysRecovered           int `json:"num_keys_recovered"`
				NumObjectsOmap             int `json:"num_objects_omap"`
				NumObjectsHitSetArchive    int `json:"num_objects_hit_set_archive"`
				NumBytesHitSetArchive      int `json:"num_bytes_hit_set_archive"`
				NumFlush                   int `json:"num_flush"`
				NumFlushKb                 int `json:"num_flush_kb"`
				NumEvict                   int `json:"num_evict"`
				NumEvictKb                 int `json:"num_evict_kb"`
				NumPromote                 int `json:"num_promote"`
				NumFlushModeHigh           int `json:"num_flush_mode_high"`
				NumFlushModeLow            int `json:"num_flush_mode_low"`
				NumEvictModeSome           int `json:"num_evict_mode_some"`
				NumEvictModeFull           int `json:"num_evict_mode_full"`
				NumObjectsPinned           int `json:"num_objects_pinned"`
				NumLegacySnapsets          int `json:"num_legacy_snapsets"`
				NumLargeOmapObjects        int `json:"num_large_omap_objects"`
				NumObjectsManifest         int `json:"num_objects_manifest"`
				NumOmapBytes               int `json:"num_omap_bytes"`
				NumOmapKeys                int `json:"num_omap_keys"`
				NumObjectsRepaired         int `json:"num_objects_repaired"`
			} `json:"stat_sum"`
			Up                   []int `json:"up"`
			Acting               []int `json:"acting"`
			AvailNoMissing       []any `json:"avail_no_missing"`
			ObjectLocationCounts []any `json:"object_location_counts"`
			BlockedBy            []any `json:"blocked_by"`
			UpPrimary            int   `json:"up_primary"`
			ActingPrimary        int   `json:"acting_primary"`
			PurgedSnaps          []any `json:"purged_snaps"`
		} `json:"pg_stats"`
		PoolStats []struct {
			Poolid  int `json:"poolid"`
			NumPg   int `json:"num_pg"`
			StatSum struct {
				NumBytes                   int64 `json:"num_bytes"`
				NumObjects                 int   `json:"num_objects"`
				NumObjectClones            int   `json:"num_object_clones"`
				NumObjectCopies            int   `json:"num_object_copies"`
				NumObjectsMissingOnPrimary int   `json:"num_objects_missing_on_primary"`
				NumObjectsMissing          int   `json:"num_objects_missing"`
				NumObjectsDegraded         int   `json:"num_objects_degraded"`
				NumObjectsMisplaced        int   `json:"num_objects_misplaced"`
				NumObjectsUnfound          int   `json:"num_objects_unfound"`
				NumObjectsDirty            int   `json:"num_objects_dirty"`
				NumWhiteouts               int   `json:"num_whiteouts"`
				NumRead                    int   `json:"num_read"`
				NumReadKb                  int   `json:"num_read_kb"`
				NumWrite                   int   `json:"num_write"`
				NumWriteKb                 int   `json:"num_write_kb"`
				NumScrubErrors             int   `json:"num_scrub_errors"`
				NumShallowScrubErrors      int   `json:"num_shallow_scrub_errors"`
				NumDeepScrubErrors         int   `json:"num_deep_scrub_errors"`
				NumObjectsRecovered        int   `json:"num_objects_recovered"`
				NumBytesRecovered          int   `json:"num_bytes_recovered"`
				NumKeysRecovered           int   `json:"num_keys_recovered"`
				NumObjectsOmap             int   `json:"num_objects_omap"`
				NumObjectsHitSetArchive    int   `json:"num_objects_hit_set_archive"`
				NumBytesHitSetArchive      int   `json:"num_bytes_hit_set_archive"`
				NumFlush                   int   `json:"num_flush"`
				NumFlushKb                 int   `json:"num_flush_kb"`
				NumEvict                   int   `json:"num_evict"`
				NumEvictKb                 int   `json:"num_evict_kb"`
				NumPromote                 int   `json:"num_promote"`
				NumFlushModeHigh           int   `json:"num_flush_mode_high"`
				NumFlushModeLow            int   `json:"num_flush_mode_low"`
				NumEvictModeSome           int   `json:"num_evict_mode_some"`
				NumEvictModeFull           int   `json:"num_evict_mode_full"`
				NumObjectsPinned           int   `json:"num_objects_pinned"`
				NumLegacySnapsets          int   `json:"num_legacy_snapsets"`
				NumLargeOmapObjects        int   `json:"num_large_omap_objects"`
				NumObjectsManifest         int   `json:"num_objects_manifest"`
				NumOmapBytes               int   `json:"num_omap_bytes"`
				NumOmapKeys                int   `json:"num_omap_keys"`
				NumObjectsRepaired         int   `json:"num_objects_repaired"`
			} `json:"stat_sum"`
			StoreStats struct {
				Total                   int   `json:"total"`
				Available               int   `json:"available"`
				InternallyReserved      int   `json:"internally_reserved"`
				Allocated               int64 `json:"allocated"`
				DataStored              int64 `json:"data_stored"`
				DataCompressed          int   `json:"data_compressed"`
				DataCompressedAllocated int   `json:"data_compressed_allocated"`
				DataCompressedOriginal  int   `json:"data_compressed_original"`
				OmapAllocated           int   `json:"omap_allocated"`
				InternalMetadata        int   `json:"internal_metadata"`
			} `json:"store_stats"`
			LogSize       int `json:"log_size"`
			OndiskLogSize int `json:"ondisk_log_size"`
			Up            int `json:"up"`
			Acting        int `json:"acting"`
			NumStoreStats int `json:"num_store_stats"`
		} `json:"pool_stats"`
		OsdStats []struct {
			Osd                int   `json:"osd"`
			UpFrom             int   `json:"up_from"`
			Seq                int64 `json:"seq"`
			NumPgs             int   `json:"num_pgs"`
			NumOsds            int   `json:"num_osds"`
			NumPerPoolOsds     int   `json:"num_per_pool_osds"`
			NumPerPoolOmapOsds int   `json:"num_per_pool_omap_osds"`
			Kb                 int64 `json:"kb"`
			KbUsed             int   `json:"kb_used"`
			KbUsedData         int   `json:"kb_used_data"`
			KbUsedOmap         int   `json:"kb_used_omap"`
			KbUsedMeta         int   `json:"kb_used_meta"`
			KbAvail            int64 `json:"kb_avail"`
			Statfs             struct {
				Total                   int64 `json:"total"`
				Available               int64 `json:"available"`
				InternallyReserved      int   `json:"internally_reserved"`
				Allocated               int64 `json:"allocated"`
				DataStored              int64 `json:"data_stored"`
				DataCompressed          int   `json:"data_compressed"`
				DataCompressedAllocated int   `json:"data_compressed_allocated"`
				DataCompressedOriginal  int   `json:"data_compressed_original"`
				OmapAllocated           int   `json:"omap_allocated"`
				InternalMetadata        int   `json:"internal_metadata"`
			} `json:"statfs"`
			HbPeers           []int `json:"hb_peers"`
			SnapTrimQueueLen  int   `json:"snap_trim_queue_len"`
			NumSnapTrimming   int   `json:"num_snap_trimming"`
			NumShardsRepaired int   `json:"num_shards_repaired"`
			OpQueueAgeHist    struct {
				Histogram  []any `json:"histogram"`
				UpperBound int   `json:"upper_bound"`
			} `json:"op_queue_age_hist"`
			PerfStat struct {
				CommitLatencyMs int `json:"commit_latency_ms"`
				ApplyLatencyMs  int `json:"apply_latency_ms"`
				CommitLatencyNs int `json:"commit_latency_ns"`
				ApplyLatencyNs  int `json:"apply_latency_ns"`
			} `json:"perf_stat"`
			Alerts           []any `json:"alerts"`
			NetworkPingTimes []struct {
				Osd        int    `json:"osd"`
				LastUpdate string `json:"last update"`
				Interfaces []struct {
					Interface string `json:"interface"`
					Average   struct {
						OneMin  float64 `json:"1min"`
						FiveMin float64 `json:"5min"`
						One5Min float64 `json:"15min"`
					} `json:"average"`
					Min struct {
						OneMin  float64 `json:"1min"`
						FiveMin float64 `json:"5min"`
						One5Min float64 `json:"15min"`
					} `json:"min"`
					Max struct {
						OneMin  float64 `json:"1min"`
						FiveMin float64 `json:"5min"`
						One5Min float64 `json:"15min"`
					} `json:"max"`
					Last float64 `json:"last"`
				} `json:"interfaces"`
			} `json:"network_ping_times"`
		} `json:"osd_stats"`
	} `json:"pg_map"`
}

type OsdPoolTuple struct {
	osd, pool interface{}
}

// ----------------UTILITY----------------
func readJson() CephDumpOutputStruct {
	jsonFile, err := os.Open("ceph_dump_beautified.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	var pgStruc CephDumpOutputStruct

	json.Unmarshal(byteValue, &pgStruc)

	return pgStruc
}

func removeDuplicateStr(strSlice []string) []string {
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

func removeDuplicateInt(intSlice []int) []int {
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

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
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

// ----------------INFORMATIONS GATHERING----------------
func getPgIds(CephDumpOutput CephDumpOutputStruct) []string {
	defer duration(track("getPgIds"))

	pgStats := CephDumpOutput.PgMap.PgStats
	pgSlice := []string{}

	for _, s := range pgStats {
		pgSlice = append(pgSlice, s.Pgid)
	}

	return pgSlice
}

func getPools(CephDumpOutput CephDumpOutputStruct) []string {
	defer duration(track("getPools"))

	pgSlice := getPgIds(CephDumpOutput)
	poolSlice := []string{}

	for _, s := range pgSlice {
		poolId := strings.Split(s, ".")[0]

		poolSlice = append(poolSlice, poolId)
	}

	poolSlice = removeDuplicateStr(poolSlice)

	return poolSlice
}

func getPgIdOsdMap(CephDumpOutput CephDumpOutputStruct) map[string][]int {
	defer duration(track("getPgIdOsdMap"))

	pgStats := CephDumpOutput.PgMap.PgStats
	pgIdOsdMap := make(map[string][]int)

	for _, s := range pgStats {
		pgOsdSlice := s.Up
		pgIdOsdMap[s.Pgid] = pgOsdSlice
	}

	return pgIdOsdMap
}

func getOsdPgMap(CephDumpOutput CephDumpOutputStruct) map[int][]string {
	defer duration(track("getOsdPgMap"))

	pgIdOsdMap := getPgIdOsdMap(CephDumpOutput)
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

func getOsds(CephDumpOutput CephDumpOutputStruct) []int {
	defer duration(track("getOsds"))

	osdPgMap := getOsdPgMap(CephDumpOutput)
	osds := []int{}

	for osd := range osdPgMap {
		osds = append(osds, osd)
	}

	return osds
}

func getNumberOfAssociatedPgsPerOsdMap(osdPgMap map[int][]string) map[int]int {
	defer duration(track("getNumberOfAssociatedPgsPerOsdMap"))

	osdPgNumberMap := make(map[int]int)

	for key, value := range osdPgMap {
		osdPgNumberMap[key] = len(value)
	}

	return osdPgNumberMap
}

func getTotalNumberOfPgs(osdPgNumberMap map[int]int) int {
	defer duration(track("getTotalNumberOfPgs"))

	totalPgs := 0

	for _, osdPgCount := range osdPgNumberMap {
		totalPgs += osdPgCount
	}

	return totalPgs
}

func getOsdPoolPgMap(CephDumpOutput CephDumpOutputStruct) map[OsdPoolTuple][]string {
	defer duration(track("getOsdPoolPgMap"))

	osdPgMap := getOsdPgMap(CephDumpOutput)
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

func getOsdPoolNumberPerPgsMap(CephDumpOutput CephDumpOutputStruct) map[OsdPoolTuple]map[string]int {
	defer duration(track("getOsdPoolNumberPerPgsMap"))

	osdPoolPgMap := getOsdPoolPgMap(CephDumpOutput)
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

func getOsdsContainingPool(pool string, CephDumpOutput CephDumpOutputStruct) []int {
	defer duration(track("getOsdsContainingPool"))

	osdPoolPgMap := getOsdPoolPgMap(CephDumpOutput)

	poolOsdsSlice := []int{}

	for key := range osdPoolPgMap {
		if key.pool == pool {
			poolOsdsSlice = append(poolOsdsSlice, key.osd.(int))
		}

	}

	return removeDuplicateInt(poolOsdsSlice)
}

func getOsdsContainingPg(givenPg string, CephDumpOutput CephDumpOutputStruct) []int {
	defer duration(track("getOsdsContainingPg"))

	osdPoolPgMap := getOsdPoolPgMap(CephDumpOutput)

	pgOsdsSlice := []int{}

	for key, pgSlice := range osdPoolPgMap {
		for _, pg := range pgSlice {
			if pg == givenPg {
				pgOsdsSlice = append(pgOsdsSlice, key.osd.(int))
			}
		}
	}

	return pgOsdsSlice
}

func getAffectedPools(faultOsd int, CephDumpOutput CephDumpOutputStruct) []string {
	defer duration(track("getAffectedPools"))

	osdPgMap := getOsdPgMap(CephDumpOutput)

	affectedPoolSlice := []string{}

	for _, pg := range osdPgMap[faultOsd] {
		poolId := strings.Split(pg, ".")[0]

		affectedPoolSlice = append(affectedPoolSlice, poolId)
	}

	return removeDuplicateStr(affectedPoolSlice)
}

func getAffectedPgs(faultOsd int, CephDumpOutput CephDumpOutputStruct) []string {
	defer duration(track("getAffectedPgs"))

	osdPgMap := getOsdPgMap(CephDumpOutput)

	affectedPgSlice := []string{}

	affectedPgSlice = append(affectedPgSlice, osdPgMap[faultOsd]...)

	return affectedPgSlice
}

func getNumAffectedReplicasPgItem(faultOsd int, givenPg string, CephDumpOutput CephDumpOutputStruct) float64 {
	defer duration(track("getNumAffectedReplicasPgItem"))

	affectedPgs := getAffectedPgs(faultOsd, CephDumpOutput)

	count := 0

	for _, pg := range affectedPgs {
		if pg == givenPg {
			count += 1
		}
	}

	return float64(count)
}

func getAllPgInAllOsds(CephDumpOutput CephDumpOutputStruct) []string {
	defer duration(track("getAllPgInAllOsds"))

	osdPgMap := getOsdPgMap(CephDumpOutput)

	allPgs := []string{}

	for _, pgSlice := range osdPgMap {
		allPgs = append(allPgs, pgSlice...)
	}

	return allPgs
}

func getPgTotalReplicas(givenPg string, CephDumpOutput CephDumpOutputStruct) float64 {
	defer duration(track("getPgTotalReplicas"))

	allPg := getAllPgInAllOsds(CephDumpOutput)

	//how many replicas of givenPg in allPg
	count := 0
	for _, pg := range allPg {
		if pg == givenPg {
			count += 1
		}
	}

	return float64(count)
}

func percentageCalculationAffectedReplicasPg(faultOsd int, givenPg string, CephDumpOutput CephDumpOutputStruct) float64 {
	defer duration(track("percentageCalculationAffectedReplicasPg"))

	numAffectedReplicaPgItem := getNumAffectedReplicasPgItem(faultOsd, givenPg, CephDumpOutput)

	pgTotalReplicas := getPgTotalReplicas(givenPg, CephDumpOutput)

	percentage := (numAffectedReplicaPgItem / pgTotalReplicas) * 100

	return roundFloat(percentage, 2)
}

func warningCheck(percentage float64, faultOsd int, givenPg string, CephDumpOutput CephDumpOutputStruct) {
	defer duration(track("warningCheck"))
	//log.SetOutput(io.Discard) //disable log

	osds := getOsds(CephDumpOutput)
	warningOsdSlice := make(map[int]float64)

	if percentage > 0.0 {
		for osd := range osds {
			percentage = percentageCalculationAffectedReplicasPg(osd, givenPg, CephDumpOutput)
			if percentage > 0.0 && osd != faultOsd {
				warningOsdSlice[osd] = percentage
			}
		}
	}

	//log.SetOutput(os.Stdout) //renable log
	fmt.Printf("Warning: check these other OSDs and percentages for %s ->%v\n", givenPg, warningOsdSlice)
}

func getTotalAffectedPgsAndPools(faultOsdSlice []int, CephDumpOutput CephDumpOutputStruct) ([]string, []string) {
	defer duration(track("getTotalAffectedPgsAndPools"))

	totalAffectedPgs := []string{}
	totalAffectedPools := []string{}

	for _, osd := range faultOsdSlice {
		affectedPgs := getAffectedPgs(osd, CephDumpOutput)
		affectedPools := getAffectedPools(osd, CephDumpOutput)

		totalAffectedPgs = append(totalAffectedPgs, affectedPgs...)
		totalAffectedPools = append(totalAffectedPools, affectedPools...)
	}

	totalAffectedPools = removeDuplicateStr(totalAffectedPools)

	return totalAffectedPgs, totalAffectedPools
}

func getPgNumberOfAffectedReplicaMap(faultOsdSlice []int, CephDumpOutput CephDumpOutputStruct) map[string]int {
	defer duration(track("getPgNumberOfAffectedReplicaMap"))

	totalAffectedPgs, _ := getTotalAffectedPgsAndPools(faultOsdSlice, CephDumpOutput)
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

func getPgsWithHighProbabilityOfLosingData(faultOsdSlice []int, CephDumpOutput CephDumpOutputStruct) ([]string, []string, []string, []string) {
	defer duration(track("getPgsWithHighProbabilityOfLosingData"))

	pgNumberOfAffectedReplicaMap := getPgNumberOfAffectedReplicaMap(faultOsdSlice, CephDumpOutput)

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

	inHelthPgs := Difference(removeDuplicateStr(getAllPgInAllOsds(CephDumpOutput)), allCompromisedPgs)

	return inHelthPgs, goodPgs, warningPgs, lostPgs
}

func extractPoolsFromPgSlice(pgSlice []string) []string {
	defer duration(track("extractPoolsFromPgSlice"))

	poolSlice := []string{}

	for _, pg := range pgSlice {
		poolId := strings.Split(pg, ".")[0]
		poolSlice = append(poolSlice, poolId)
	}
	return removeDuplicateStr(poolSlice)
}

func extractOsdsFromPoolSlice(poolSlice []string, CephDumpOutput CephDumpOutputStruct) []int {
	defer duration(track("extractOsdsFromPoolSlice"))
	osdSlice := []int{}

	for _, pool := range poolSlice {
		osdsContainingPool := getOsdsContainingPool(pool, CephDumpOutput)
		osdSlice = append(osdSlice, osdsContainingPool...)
	}

	return removeDuplicateInt(osdSlice)
}

func main() {
	log.SetOutput(io.Discard) //disable log

	//UnMarshalling Json
	CephDumpOutput := readJson()
	//fmt.Println(CephDumpOutput.PgMap.PgStats[0].Pgid)

	//---getting different pg_ids---
	fmt.Println("------------------------------------------")
	pgSlice := getPgIds(CephDumpOutput)
	fmt.Printf("#Different pgids  -> %d\n\n", len(pgSlice))
	fmt.Println("------------------------------------------")

	//---getting pools---
	poolSlice := getPools(CephDumpOutput)
	fmt.Printf("Pools -> %s\n", poolSlice)
	fmt.Printf("#Pools -> %d\n\n", len(poolSlice))
	fmt.Println("------------------------------------------")

	//---getting_osds---
	osds := getOsds(CephDumpOutput)
	fmt.Printf("OSDs ->%d\n", osds)
	fmt.Printf("#OSDs ->%d\n\n", len(osds))
	fmt.Println("------------------------------------------")

	//---getting_pgs_details---
	osdPgMap := getOsdPgMap(CephDumpOutput)
	osdPgNumberMap := getNumberOfAssociatedPgsPerOsdMap(osdPgMap)
	fmt.Printf("Number of mapped PGs per OSD ->%d\n\n", osdPgNumberMap)
	fmt.Printf("#PGs ->%d\n\n", getTotalNumberOfPgs(osdPgNumberMap))
	fmt.Println("------------------------------------------")

	//---getting_pgIdOsdMap---
	fmt.Println("-> pgIdOsdMap")
	pgIdOsdMap := getPgIdOsdMap(CephDumpOutput)
	fmt.Printf("pgid 31.0 is on these osds -> %d\n\n", pgIdOsdMap["31.0"])
	fmt.Println("------------------------------------------")

	//---osd_pgs_mapping---
	fmt.Println("-> osd_pgs_mapping")
	numberOfPg := getNumberOfAssociatedPgsPerOsdMap(osdPgMap)[3]
	fmt.Printf("OSD 3 contains these #%d PGs-> %s\n\n", numberOfPg, osdPgMap[3])
	fmt.Println("------------------------------------------")

	//------osd_pool_pgs_map------
	//	   {(osd,pool):[pgs]}
	fmt.Println("-> osd_pool_pgs_map   {(osd,pool):[pgs]}")
	osdPoolPgMap := getOsdPoolPgMap(CephDumpOutput)
	fmt.Printf("OSD 6 pool 31 contains these PGs -> %s\n\n", osdPoolPgMap[OsdPoolTuple{6, "31"}])

	//	   {(osd,pool):{pg:number_replicas},...}
	fmt.Println("-> osd_pool_pgs_map   {(osd,pool):{pg:number_replicas},...}")
	osdPollNumberPerPgsMap := getOsdPoolNumberPerPgsMap(CephDumpOutput)
	fmt.Printf("OSD 6 pool 31 contains these PGs -> %v\n\n", osdPollNumberPerPgsMap[OsdPoolTuple{6, "31"}])
	fmt.Println("------------------------------------------")

	//---osds_containing_pool---
	fmt.Println("->   Pool->OSDs")
	pool := "3"
	osdsContaingPool := getOsdsContainingPool(pool, CephDumpOutput)
	fmt.Printf("Pool: %s is spread between these OSDs-> %v\n\n", pool, osdsContaingPool)
	fmt.Println("------------------------------------------")

	//---osds_containing_pg---
	fmt.Println("->   PG->OSDs")
	givenPg := "19.1f"
	osdsContainingPg := getOsdsContainingPg(givenPg, CephDumpOutput)
	fmt.Printf("PG: %s is spread between these OSDs-> %d\n\n", givenPg, osdsContainingPg)
	fmt.Println("------------------------------------------")

	//---affected pool for osd crush---
	fmt.Println("-> Affected pools")
	faultOsd := 2
	affectedPools := getAffectedPools(faultOsd, CephDumpOutput)
	fmt.Printf("Affected pools for crush on OSD: %d  -> #%d ->  %v\n\n", faultOsd, len(affectedPools), affectedPools)
	fmt.Println("------------------------------------------")

	//---affected pg for osd crush---
	fmt.Println("-> Affected PGs")
	affectedPgs := getAffectedPgs(faultOsd, CephDumpOutput)
	fmt.Printf("Affected pgid for crush on OSD: %d -> %s\n\n", faultOsd, affectedPgs)
	fmt.Println("------------------------------------------")

	//---percentage of lost replicas for single pgid_item in fault_osd--- (OSD DEGRADATION)
	fmt.Printf("-> lost OSD: %d -> percentage%% of lost replicas for %s---(OSD DEGRADATION)\n", faultOsd, givenPg)
	percentage := percentageCalculationAffectedReplicasPg(faultOsd, givenPg, CephDumpOutput)
	fmt.Printf("\nPercentage of %s replicas lost -> %.2f%%\n\n", givenPg, percentage)
	fmt.Println("------------------------------------------")

	warningCheck(percentage, faultOsd, givenPg, CephDumpOutput)
	fmt.Printf("---------------------------------------------------------------\n\n")

	//---if these osds crush which are the affected Pgs?---
	faultOsdSlice := []int{2, 4, 3}
	fmt.Printf("-> If these osds:%d crush which Pgs and Pools are affected?\n", faultOsdSlice)
	totalAffectedPgs, totalAffectedPools := getTotalAffectedPgsAndPools(faultOsdSlice, CephDumpOutput)
	fmt.Printf("\n#%d Total affected pgs for osds:%d  -> (map is hidden, uncomment to view)", len(totalAffectedPgs), faultOsdSlice)
	//fmt.Printf("%s\n\n", totalAffectedPgs)
	fmt.Printf("\n#%d Total affected pools for osds:%d  -> %s\n\n", len(totalAffectedPools), faultOsdSlice, totalAffectedPools)

	//    {pg: number_affected_replicas}
	fmt.Printf("->   {pg: number_affected_replicas}\n")
	pgNumberOfAffectedReplicaMap := getPgNumberOfAffectedReplicaMap(faultOsdSlice, CephDumpOutput)
	fmt.Printf("%v\n\n", pgNumberOfAffectedReplicaMap)

	//count := 0
	//for _, value := range pgNumberOfAffectedReplicaMap {
	//	count += value
	//}
	//fmt.Printf("-->%d", count)

	inHelthPgs, goodPgs, warningPgs, lostPgs := getPgsWithHighProbabilityOfLosingData(faultOsdSlice, CephDumpOutput)

	fmt.Printf("---------------------------------------------------------------\n")

	fmt.Printf("inHelthPgs: replicaLost=0%%\n\n%s\n", inHelthPgs)
	fmt.Printf("---------------------------------------------------------------\n")

	fmt.Printf("goodPgs (replicaLost<50%%):\n\n%s\n\n", goodPgs)
	goodPgsPoolSlice := extractPoolsFromPgSlice(goodPgs)
	fmt.Printf("pools ->%s\n", goodPgsPoolSlice)
	fmt.Printf("previous pools are spread on these osds ->%d\n", extractOsdsFromPoolSlice(goodPgsPoolSlice, CephDumpOutput))
	fmt.Printf("---------------------------------------------------------------\n")

	fmt.Printf("warningPgs (replicaLost>=50%% and replicaLost<100%%):\n\n%s\n\n", warningPgs)
	warningPgsPoolSlice := extractPoolsFromPgSlice(warningPgs)
	fmt.Printf("pools ->%s\n", warningPgsPoolSlice)
	fmt.Printf("previous pools are spread on these osds ->%d\n", extractOsdsFromPoolSlice(warningPgsPoolSlice, CephDumpOutput))
	fmt.Printf("---------------------------------------------------------------\n")

	fmt.Printf("lostPgs  (replicaLost=100%%):\n\n%s\n\n", lostPgs)
	lostPgsPoolSlice := extractPoolsFromPgSlice(lostPgs)
	fmt.Printf("pools ->%s\n", lostPgsPoolSlice)
	fmt.Printf("previous pools are spread on these osds ->%d\n", extractOsdsFromPoolSlice(lostPgsPoolSlice, CephDumpOutput))
	fmt.Printf("---------------------------------------------------------------\n")

	//checkTotal := len(inHelthPgs) + len(goodPgs) + len(warningPgs) + len(lostPgs)
	//fmt.Printf("\ncheck total: %d\n", checkTotal)

}
