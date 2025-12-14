package models

type SmartInfo struct {
	Device struct {
		InfoName string `json:"info_name"`
		Name     string `json:"name"`
		Protocol string `json:"protocol"`
		Type     string `json:"type"`
	} `json:"device"`

	EnduranceUsed struct {
		CurrentPercent int `json:"current_percent"`
	} `json:"endurance_used"`

	FirmwareVersion string `json:"firmware_version"`

	JSONFormatVersion []int `json:"json_format_version"`

	LocalTime struct {
		Asctime string `json:"asctime"`
		TimeT   int64  `json:"time_t"`
	} `json:"local_time"`

	LogicalBlockSize int    `json:"logical_block_size"`
	ModelName        string `json:"model_name"`

	NVMeCompositeTemperatureThreshold struct {
		Critical int `json:"critical"`
		Warning  int `json:"warning"`
	} `json:"nvme_composite_temperature_threshold"`

	NVMeControllerID int `json:"nvme_controller_id"`

	NVMeErrorInformationLog struct {
		Read   int `json:"read"`
		Unread int `json:"unread"`
		Size   int `json:"size"`
		Table  []struct {
			CommandID  int `json:"command_id"`
			ErrorCount int `json:"error_count"`
			LBA        struct {
				Value uint64 `json:"value"`
			} `json:"lba"`
			NSID              int  `json:"nsid"`
			ParmErrorLocation int  `json:"parm_error_location"`
			PhaseTag          bool `json:"phase_tag"`
			StatusField       struct {
				DoNotRetry     bool   `json:"do_not_retry"`
				StatusCode     int    `json:"status_code"`
				StatusCodeType int    `json:"status_code_type"`
				String         string `json:"string"`
				Value          int    `json:"value"`
			} `json:"status_field"`
			SubmissionQueueID int `json:"submission_queue_id"`
		} `json:"table"`
	} `json:"nvme_error_information_log"`

	NVMeFirmwareUpdateCapabilities struct {
		ActiviationWithoutReset bool `json:"activiation_without_reset"`
		FirstSlotIsReadOnly     bool `json:"first_slot_is_read_only"`
		MultipleUpdateDetection bool `json:"multiple_update_detection"`
		Other                   int  `json:"other"`
		Slots                   int  `json:"slots"`
		Value                   int  `json:"value"`
	} `json:"nvme_firmware_update_capabilities"`

	NVMeIEEEOuiIdentifier int `json:"nvme_ieee_oui_identifier"`

	NVMeLogPageAttributes struct {
		CommandsEffectsLog      bool `json:"commands_effects_log"`
		ExtendedGetLogPageCmd   bool `json:"extended_get_log_page_cmd"`
		PersistentEventLog      bool `json:"persistent_event_log"`
		SmartHealthPerNamespace bool `json:"smart_health_per_namespace"`
		SupportedLogPagesLog    bool `json:"supported_log_pages_log"`
		TelemetryDataArea4      bool `json:"telemetry_data_area_4"`
		TelemetryLog            bool `json:"telemetry_log"`
		Other                   int  `json:"other"`
		Value                   int  `json:"value"`
	} `json:"nvme_log_page_attributes"`

	NVMeMaximumDataTransferPages int `json:"nvme_maximum_data_transfer_pages"`

	NVMeNamespaces []struct {
		ID int `json:"id"`

		Capacity struct {
			Blocks uint64 `json:"blocks"`
			Bytes  uint64 `json:"bytes"`
		} `json:"capacity"`

		Size struct {
			Blocks uint64 `json:"blocks"`
			Bytes  uint64 `json:"bytes"`
		} `json:"size"`

		Utilization struct {
			Blocks uint64 `json:"blocks"`
			Bytes  uint64 `json:"bytes"`
		} `json:"utilization"`

		EUI64 struct {
			ExtID uint64 `json:"ext_id"`
			OUI   int    `json:"oui"`
		} `json:"eui64"`

		Features struct {
			DeallocOrUnwrittenBlockError bool `json:"dealloc_or_unwritten_block_error"`
			NAFields                     bool `json:"na_fields"`
			NPFields                     bool `json:"np_fields"`
			ThinProvisioning             bool `json:"thin_provisioning"`
			UIDReuse                     bool `json:"uid_reuse"`
			Other                        int  `json:"other"`
			Value                        int  `json:"value"`
		} `json:"features"`

		FormattedLBASize int `json:"formatted_lba_size"`

		LBAFormats []struct {
			DataBytes           int  `json:"data_bytes"`
			Formatted           bool `json:"formatted"`
			MetadataBytes       int  `json:"metadata_bytes"`
			RelativePerformance int  `json:"relative_performance"`
		} `json:"lba_formats"`
	} `json:"nvme_namespaces"`

	NVMeNumberOfNamespaces int `json:"nvme_number_of_namespaces"`

	NVMeOptionalAdminCommands struct {
		CommandAndFeatureLockdown bool `json:"command_and_feature_lockdown"`
		Directives                bool `json:"directives"`
		DoorbellBufferConfig      bool `json:"doorbell_buffer_config"`
		FirmwareDownload          bool `json:"firmware_download"`
		FormatNVM                 bool `json:"format_nvm"`
		GetLBAStatus              bool `json:"get_lba_status"`
		MISendReceive             bool `json:"mi_send_receive"`
		NamespaceManagement       bool `json:"namespace_management"`
		SecuritySendReceive       bool `json:"security_send_receive"`
		SelfTest                  bool `json:"self_test"`
		VirtualizationManagement  bool `json:"virtualization_management"`
		Other                     int  `json:"other"`
		Value                     int  `json:"value"`
	} `json:"nvme_optional_admin_commands"`

	NVMeOptionalNVMCommands struct {
		Compare                  bool `json:"compare"`
		Copy                     bool `json:"copy"`
		DatasetManagement        bool `json:"dataset_management"`
		Reservations             bool `json:"reservations"`
		SaveSelectFeatureNonZero bool `json:"save_select_feature_nonzero"`
		Timestamp                bool `json:"timestamp"`
		Verify                   bool `json:"verify"`
		WriteUncorrectable       bool `json:"write_uncorrectable"`
		WriteZeroes              bool `json:"write_zeroes"`
		Other                    int  `json:"other"`
		Value                    int  `json:"value"`
	} `json:"nvme_optional_nvm_commands"`

	NVMePCIVendor struct {
		ID          int `json:"id"`
		SubsystemID int `json:"subsystem_id"`
	} `json:"nvme_pci_vendor"`

	NVMePowerStates []struct {
		EntryLatencyUS int `json:"entry_latency_us"`
		ExitLatencyUS  int `json:"exit_latency_us"`
		MaxPower       struct {
			Scale        int `json:"scale"`
			UnitsPerWatt int `json:"units_per_watt"`
			Value        int `json:"value"`
		} `json:"max_power"`
		NonOperationalState     bool `json:"non_operational_state"`
		RelativeReadLatency     int  `json:"relative_read_latency"`
		RelativeReadThroughput  int  `json:"relative_read_throughput"`
		RelativeWriteLatency    int  `json:"relative_write_latency"`
		RelativeWriteThroughput int  `json:"relative_write_throughput"`
	} `json:"nvme_power_states"`

	NVMeSelfTestLog struct {
		NSID                     int `json:"nsid"`
		CurrentSelfTestOperation struct {
			String string `json:"string"`
			Value  int    `json:"value"`
		} `json:"current_self_test_operation"`
		Table []struct {
			PowerOnHours int `json:"power_on_hours"`
			SelfTestCode struct {
				String string `json:"string"`
				Value  int    `json:"value"`
			} `json:"self_test_code"`
			SelfTestResult struct {
				String string `json:"string"`
				Value  int    `json:"value"`
			} `json:"self_test_result"`
		} `json:"table"`
	} `json:"nvme_self_test_log"`

	NVMeSmartHealthInformationLog struct {
		AvailableSpare          int    `json:"available_spare"`
		AvailableSpareThreshold int    `json:"available_spare_threshold"`
		ControllerBusyTime      int    `json:"controller_busy_time"`
		CriticalCompTime        int    `json:"critical_comp_time"`
		CriticalWarning         int    `json:"critical_warning"`
		DataUnitsRead           uint64 `json:"data_units_read"`
		DataUnitsWritten        uint64 `json:"data_units_written"`
		HostReads               uint64 `json:"host_reads"`
		HostWrites              uint64 `json:"host_writes"`
		MediaErrors             int    `json:"media_errors"`
		NSID                    int    `json:"nsid"`
		NumErrLogEntries        int    `json:"num_err_log_entries"`
		PercentageUsed          int    `json:"percentage_used"`
		PowerCycles             int    `json:"power_cycles"`
		PowerOnHours            int    `json:"power_on_hours"`
		Temperature             int    `json:"temperature"`
		UnsafeShutdowns         int    `json:"unsafe_shutdowns"`
		WarningTempTime         int    `json:"warning_temp_time"`
	} `json:"nvme_smart_health_information_log"`

	NVMeTotalCapacity       uint64 `json:"nvme_total_capacity"`
	NVMeUnallocatedCapacity uint64 `json:"nvme_unallocated_capacity"`

	NVMeVersion struct {
		String string `json:"string"`
		Value  int    `json:"value"`
	} `json:"nvme_version"`

	PowerCycleCount int `json:"power_cycle_count"`

	PowerOnTime struct {
		Hours int `json:"hours"`
	} `json:"power_on_time"`

	SerialNumber string `json:"serial_number"`

	SmartStatus struct {
		Passed bool `json:"passed"`
		NVMe   struct {
			Value int `json:"value"`
		} `json:"nvme"`
	} `json:"smart_status"`

	SmartSupport struct {
		Available bool `json:"available"`
		Enabled   bool `json:"enabled"`
	} `json:"smart_support"`

	Smartctl struct {
		Argv         []string `json:"argv"`
		BuildInfo    string   `json:"build_info"`
		ExitStatus   int      `json:"exit_status"`
		PlatformInfo string   `json:"platform_info"`
		PreRelease   bool     `json:"pre_release"`
		SVNRevision  string   `json:"svn_revision"`
		Version      []int    `json:"version"`
	} `json:"smartctl"`

	SpareAvailable struct {
		CurrentPercent   int `json:"current_percent"`
		ThresholdPercent int `json:"threshold_percent"`
	} `json:"spare_available"`

	Temperature struct {
		CriticalLimitMax int `json:"critical_limit_max"`
		Current          int `json:"current"`
		OpLimitMax       int `json:"op_limit_max"`
	} `json:"temperature"`

	UserCapacity struct {
		Blocks uint64 `json:"blocks"`
		Bytes  uint64 `json:"bytes"`
	} `json:"user_capacity"`
}
