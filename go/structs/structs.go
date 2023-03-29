package structs

type PgDumpOutputStruct struct {
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

type NodeStruct struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	TypeID      int    `json:"type_id"`
	Children    []int  `json:"children,omitempty"`
	PoolWeights struct {
	} `json:"pool_weights,omitempty"`
	DeviceClass     string  `json:"device_class,omitempty"`
	CrushWeight     float64 `json:"crush_weight,omitempty"`
	Depth           int     `json:"depth,omitempty"`
	Exists          int     `json:"exists,omitempty"`
	Status          string  `json:"status,omitempty"`
	Reweight        int     `json:"reweight,omitempty"`
	PrimaryAffinity int     `json:"primary_affinity,omitempty"`
}

// type OsdDumpOutputStruct struct {
// 	Nodes []struct {
// 		ID          int    `json:"id"`
// 		Name        string `json:"name"`
// 		Type        string `json:"type"`
// 		TypeID      int    `json:"type_id"`
// 		Children    []int  `json:"children,omitempty"`
// 		PoolWeights struct {
// 		} `json:"pool_weights,omitempty"`
// 		DeviceClass     string  `json:"device_class,omitempty"`
// 		CrushWeight     float64 `json:"crush_weight,omitempty"`
// 		Depth           int     `json:"depth,omitempty"`
// 		Exists          int     `json:"exists,omitempty"`
// 		Status          string  `json:"status,omitempty"`
// 		Reweight        int     `json:"reweight,omitempty"`
// 		PrimaryAffinity int     `json:"primary_affinity,omitempty"`
// 	} `json:"nodes"`
// 	Stray []any `json:"stray"`
// }

type OsdTreeOutputStruct struct {
	Nodes []NodeStruct `json:"nodes"`
	Stray []any        `json:"stray"`
}

type BucketDistribution struct {
	Root       interface{}
	Region     interface{}
	Zone       interface{}
	Datacenter interface{}
	Rack       interface{}
	Chassis    interface{}
	Host       interface{}
	Osd        interface{}
}

type PoolStruct struct {
	Pool                 int    `json:"pool,omitempty"`
	PoolName             string `json:"pool_name,omitempty"`
	CreateTime           string `json:"create_time,omitempty"`
	Flags                int    `json:"flags,omitempty"`
	FlagsNames           string `json:"flags_names,omitempty"`
	Type                 int    `json:"type,omitempty"`
	Size                 int    `json:"size,omitempty"`
	MinSize              int    `json:"min_size,omitempty"`
	CrushRule            int    `json:"crush_rule,omitempty"`
	ObjectHash           int    `json:"object_hash,omitempty"`
	PgAutoscaleMode      string `json:"pg_autoscale_mode,omitempty"`
	PgNum                int    `json:"pg_num,omitempty"`
	PgPlacementNum       int    `json:"pg_placement_num,omitempty"`
	PgPlacementNumTarget int    `json:"pg_placement_num_target,omitempty"`
	PgNumTarget          int    `json:"pg_num_target,omitempty"`
	PgNumPending         int    `json:"pg_num_pending,omitempty"`
	LastPgMergeMeta      struct {
		SourcePgid       string `json:"source_pgid,omitempty"`
		ReadyEpoch       int    `json:"ready_epoch,omitempty"`
		LastEpochStarted int    `json:"last_epoch_started,omitempty"`
		LastEpochClean   int    `json:"last_epoch_clean,omitempty"`
		SourceVersion    string `json:"source_version,omitempty"`
		TargetVersion    string `json:"target_version,omitempty"`
	} `json:"last_pg_merge_meta,omitempty"`
	LastChange                     string `json:"last_change,omitempty"`
	LastForceOpResend              string `json:"last_force_op_resend,omitempty"`
	LastForceOpResendPrenautilus   string `json:"last_force_op_resend_prenautilus,omitempty"`
	LastForceOpResendPreluminous   string `json:"last_force_op_resend_preluminous,omitempty"`
	Auid                           int    `json:"auid,omitempty"`
	SnapMode                       string `json:"snap_mode,omitempty"`
	SnapSeq                        int    `json:"snap_seq,omitempty"`
	SnapEpoch                      int    `json:"snap_epoch,omitempty"`
	PoolSnaps                      []any  `json:"pool_snaps,omitempty"`
	RemovedSnaps                   string `json:"removed_snaps,omitempty"`
	QuotaMaxBytes                  int    `json:"quota_max_bytes,omitempty"`
	QuotaMaxObjects                int    `json:"quota_max_objects,omitempty"`
	Tiers                          []any  `json:"tiers,omitempty"`
	TierOf                         int    `json:"tier_of,omitempty"`
	ReadTier                       int    `json:"read_tier,omitempty"`
	WriteTier                      int    `json:"write_tier,omitempty"`
	CacheMode                      string `json:"cache_mode,omitempty"`
	TargetMaxBytes                 int    `json:"target_max_bytes,omitempty"`
	TargetMaxObjects               int    `json:"target_max_objects,omitempty"`
	CacheTargetDirtyRatioMicro     int    `json:"cache_target_dirty_ratio_micro,omitempty"`
	CacheTargetDirtyHighRatioMicro int    `json:"cache_target_dirty_high_ratio_micro,omitempty"`
	CacheTargetFullRatioMicro      int    `json:"cache_target_full_ratio_micro,omitempty"`
	CacheMinFlushAge               int    `json:"cache_min_flush_age,omitempty"`
	CacheMinEvictAge               int    `json:"cache_min_evict_age,omitempty"`
	ErasureCodeProfile             string `json:"erasure_code_profile,omitempty"`
	HitSetParams                   struct {
		Type string `json:"type,omitempty"`
	} `json:"hit_set_params,omitempty"`
	HitSetPeriod              int   `json:"hit_set_period,omitempty"`
	HitSetCount               int   `json:"hit_set_count,omitempty"`
	UseGmtHitset              bool  `json:"use_gmt_hitset,omitempty"`
	MinReadRecencyForPromote  int   `json:"min_read_recency_for_promote,omitempty"`
	MinWriteRecencyForPromote int   `json:"min_write_recency_for_promote,omitempty"`
	HitSetGradeDecayRate      int   `json:"hit_set_grade_decay_rate,omitempty"`
	HitSetSearchLastN         int   `json:"hit_set_search_last_n,omitempty"`
	GradeTable                []any `json:"grade_table,omitempty"`
	StripeWidth               int   `json:"stripe_width,omitempty"`
	ExpectedNumObjects        int   `json:"expected_num_objects,omitempty"`
	FastRead                  bool  `json:"fast_read,omitempty"`
	Options                   struct {
		PgAutoscaleBias    int `json:"pg_autoscale_bias,omitempty"`
		PgNumMin           int `json:"pg_num_min,omitempty"`
		RecoveryOpPriority int `json:"recovery_op_priority,omitempty"`
		RecoveryPriority   int `json:"recovery_priority,omitempty"`
	} `json:"options,omitempty"`
	ApplicationMetadata struct {
		Cephfs struct {
			Data     string `json:"data,omitempty"`
			Metadata string `json:"metadata,omitempty"`
		} `json:"cephfs,omitempty"`
		Rgw struct {
		} `json:"rgw,omitempty"`
		Rbd struct {
		} `json:"rbd,omitempty"`
		MgrDevicehealth struct {
		} `json:"mgr_devicehealth,omitempty"`
	} `json:"application_metadata,omitempty"`
}

type OsdDumpOutputStruct struct {
	Epoch                  int          `json:"epoch,omitempty"`
	Fsid                   string       `json:"fsid,omitempty"`
	Created                string       `json:"created,omitempty"`
	Modified               string       `json:"modified,omitempty"`
	LastUpChange           string       `json:"last_up_change,omitempty"`
	LastInChange           string       `json:"last_in_change,omitempty"`
	Flags                  string       `json:"flags,omitempty"`
	FlagsNum               int          `json:"flags_num,omitempty"`
	FlagsSet               []string     `json:"flags_set,omitempty"`
	CrushVersion           int          `json:"crush_version,omitempty"`
	FullRatio              float64      `json:"full_ratio,omitempty"`
	BackfillfullRatio      float64      `json:"backfillfull_ratio,omitempty"`
	NearfullRatio          float64      `json:"nearfull_ratio,omitempty"`
	ClusterSnapshot        string       `json:"cluster_snapshot,omitempty"`
	PoolMax                int          `json:"pool_max,omitempty"`
	MaxOsd                 int          `json:"max_osd,omitempty"`
	RequireMinCompatClient string       `json:"require_min_compat_client,omitempty"`
	MinCompatClient        string       `json:"min_compat_client,omitempty"`
	RequireOsdRelease      string       `json:"require_osd_release,omitempty"`
	Pools                  []PoolStruct `json:"pools,omitempty"`
	Osds                   []struct {
		Osd             int    `json:"osd,omitempty"`
		UUID            string `json:"uuid,omitempty"`
		Up              int    `json:"up,omitempty"`
		In              int    `json:"in,omitempty"`
		Weight          int    `json:"weight,omitempty"`
		PrimaryAffinity int    `json:"primary_affinity,omitempty"`
		LastCleanBegin  int    `json:"last_clean_begin,omitempty"`
		LastCleanEnd    int    `json:"last_clean_end,omitempty"`
		UpFrom          int    `json:"up_from,omitempty"`
		UpThru          int    `json:"up_thru,omitempty"`
		DownAt          int    `json:"down_at,omitempty"`
		LostAt          int    `json:"lost_at,omitempty"`
		PublicAddrs     struct {
			Addrvec []struct {
				Type  string `json:"type,omitempty"`
				Addr  string `json:"addr,omitempty"`
				Nonce int    `json:"nonce,omitempty"`
			} `json:"addrvec,omitempty"`
		} `json:"public_addrs,omitempty"`
		ClusterAddrs struct {
			Addrvec []struct {
				Type  string `json:"type,omitempty"`
				Addr  string `json:"addr,omitempty"`
				Nonce int    `json:"nonce,omitempty"`
			} `json:"addrvec,omitempty"`
		} `json:"cluster_addrs,omitempty"`
		HeartbeatBackAddrs struct {
			Addrvec []struct {
				Type  string `json:"type,omitempty"`
				Addr  string `json:"addr,omitempty"`
				Nonce int    `json:"nonce,omitempty"`
			} `json:"addrvec,omitempty"`
		} `json:"heartbeat_back_addrs,omitempty"`
		HeartbeatFrontAddrs struct {
			Addrvec []struct {
				Type  string `json:"type,omitempty"`
				Addr  string `json:"addr,omitempty"`
				Nonce int    `json:"nonce,omitempty"`
			} `json:"addrvec,omitempty"`
		} `json:"heartbeat_front_addrs,omitempty"`
		PublicAddr         string   `json:"public_addr,omitempty"`
		ClusterAddr        string   `json:"cluster_addr,omitempty"`
		HeartbeatBackAddr  string   `json:"heartbeat_back_addr,omitempty"`
		HeartbeatFrontAddr string   `json:"heartbeat_front_addr,omitempty"`
		State              []string `json:"state,omitempty"`
	} `json:"osds,omitempty"`
	OsdXinfo []struct {
		Osd                  int    `json:"osd,omitempty"`
		DownStamp            string `json:"down_stamp,omitempty"`
		LaggyProbability     int    `json:"laggy_probability,omitempty"`
		LaggyInterval        int    `json:"laggy_interval,omitempty"`
		Features             int64  `json:"features,omitempty"`
		OldWeight            int    `json:"old_weight,omitempty"`
		LastPurgedSnapsScrub string `json:"last_purged_snaps_scrub,omitempty"`
		DeadEpoch            int    `json:"dead_epoch,omitempty"`
	} `json:"osd_xinfo,omitempty"`
	PgUpmap      []any `json:"pg_upmap,omitempty"`
	PgUpmapItems []any `json:"pg_upmap_items,omitempty"`
	PgTemp       []struct {
		Pgid string `json:"pgid,omitempty"`
		Osds []int  `json:"osds,omitempty"`
	} `json:"pg_temp,omitempty"`
	PrimaryTemp []any `json:"primary_temp,omitempty"`
	Blacklist   struct {
	} `json:"blacklist,omitempty"`
	ErasureCodeProfiles struct {
		Default struct {
			K         string `json:"k,omitempty"`
			M         string `json:"m,omitempty"`
			Plugin    string `json:"plugin,omitempty"`
			Technique string `json:"technique,omitempty"`
		} `json:"default,omitempty"`
		EcPool2EcProfile struct {
			CrushDeviceClass          string `json:"crush-device-class,omitempty"`
			CrushFailureDomain        string `json:"crush-failure-domain,omitempty"`
			CrushRoot                 string `json:"crush-root,omitempty"`
			CrushDeviceClass0         string `json:"crush_device_class,omitempty"`
			JerasurePerChunkAlignment string `json:"jerasure-per-chunk-alignment,omitempty"`
			K                         string `json:"k,omitempty"`
			M                         string `json:"m,omitempty"`
			Plugin                    string `json:"plugin,omitempty"`
			StripeUnit                string `json:"stripe_unit,omitempty"`
			Technique                 string `json:"technique,omitempty"`
			W                         string `json:"w,omitempty"`
		} `json:"ec-pool2_ec_profile,omitempty"`
		Pool3EcProfile struct {
			CrushDeviceClass          string `json:"crush-device-class,omitempty"`
			CrushFailureDomain        string `json:"crush-failure-domain,omitempty"`
			CrushRoot                 string `json:"crush-root,omitempty"`
			CrushDeviceClass0         string `json:"crush_device_class,omitempty"`
			JerasurePerChunkAlignment string `json:"jerasure-per-chunk-alignment,omitempty"`
			K                         string `json:"k,omitempty"`
			M                         string `json:"m,omitempty"`
			Plugin                    string `json:"plugin,omitempty"`
			StripeUnit                string `json:"stripe_unit,omitempty"`
			Technique                 string `json:"technique,omitempty"`
			W                         string `json:"w,omitempty"`
		} `json:"pool3_ec_profile,omitempty"`
	} `json:"erasure_code_profiles,omitempty"`
	RemovedSnapsQueue []any `json:"removed_snaps_queue,omitempty"`
	NewRemovedSnaps   []any `json:"new_removed_snaps,omitempty"`
	NewPurgedSnaps    []any `json:"new_purged_snaps,omitempty"`
	CrushNodeFlags    struct {
	} `json:"crush_node_flags,omitempty"`
	DeviceClassFlags struct {
	} `json:"device_class_flags,omitempty"`
}
