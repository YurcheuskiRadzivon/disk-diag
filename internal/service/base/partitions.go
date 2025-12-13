package base

import (
	"encoding/json"
	"fmt"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
)

func (s *service) GetExtendedPartitions(diskIndex int) ([]models.PartitionInfo, error) {
	psScript := fmt.Sprintf(`
	Get-Partition -DiskNumber %d | ForEach-Object {
		$vol = $_ | Get-Volume -ErrorAction SilentlyContinue
		$bitlocker = "Unknown"
		if ($vol.DriveLetter) {
			$blObj = Get-BitLockerVolume -MountPoint ($vol.DriveLetter + ":") -ErrorAction SilentlyContinue
			if ($blObj) { $bitlocker = $blObj.VolumeStatus }
		}

		@{
			PartitionNumber = $_.PartitionNumber
			DriveLetter = $_.DriveLetter
			Size = $_.Size
			Type = $_.Type
			GptType = $_.GptType
			Offset = $_.Offset
			IsBoot = $_.IsBoot
			IsReadOnly = $_.IsReadOnly
			IsHidden = $_.IsHidden
			FreeSpace = if ($vol) { $vol.SizeRemaining } else { 0 }
			FileSystem = if ($vol) { $vol.FileSystem } else { "Unknown" }
			Label = if ($vol) { $vol.FileSystemLabel } else { "" }
			AllocationUnitSize = if ($vol) { $vol.AllocationUnitSize } else { 0 }
			BitLockerStatus = $bitlocker
		}
	} | ConvertTo-Json
	`, diskIndex)

	out, err := s.runPowerShell(psScript)
	if err != nil {
		return []models.PartitionInfo{}, nil
	}

	var rawData []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &rawData); err != nil {
		return nil, err
	}

	var parts []models.PartitionInfo
	for _, item := range rawData {
		sizeVal, _ := item["Size"].(float64)
		freeVal, _ := item["FreeSpace"].(float64)
		offsetVal, _ := item["Offset"].(float64)
		allocVal, _ := item["AllocationUnitSize"].(float64)

		driveLetter := fmt.Sprintf("%v", item["DriveLetter"])
		if driveLetter == "<nil>" || driveLetter == "0" {
			driveLetter = ""
		}

		p := models.PartitionInfo{
			DiskIndex:          diskIndex,
			PartitionID:        fmt.Sprintf("%v", item["PartitionNumber"]),
			DriveLetter:        driveLetter,
			Size:               int64(sizeVal),
			FreeSpace:          int64(freeVal),
			FileSystem:         fmt.Sprintf("%v", item["FileSystem"]),
			Label:              fmt.Sprintf("%v", item["Label"]),
			Type:               fmt.Sprintf("%v", item["Type"]),
			GptType:            fmt.Sprintf("%v", item["GptType"]),
			Offset:             int64(offsetVal),
			IsBoot:             item["IsBoot"] == true,
			IsReadOnly:         item["IsReadOnly"] == true,
			IsHidden:           item["IsHidden"] == true,
			BitLockerStatus:    fmt.Sprintf("%v", item["BitLockerStatus"]),
			AllocationUnitSize: int(allocVal),
		}
		parts = append(parts, p)
	}
	return parts, nil
}
