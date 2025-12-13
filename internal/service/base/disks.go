package base

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
)

func (s *service) GetPhysicalDisks() ([]models.DiskInfo, error) {
	psScript := `
	$pd = Get-PhysicalDisk | Select-Object DeviceId,FriendlyName,MediaType,BusType,Size,SerialNumber,HealthStatus,FirmwareVersion
	$wmi = Get-CimInstance Win32_DiskDrive | Select-Object Index,BytesPerSector,TotalHeads,TotalCylinders,Partitions
	
	$pd | ForEach-Object {
		$currId = $_.DeviceId
		$wmiMatch = $wmi | Where-Object { $_.Index -eq $currId }
		@{
			DeviceId = $_.DeviceId
			Model = $_.FriendlyName
			MediaType = $_.MediaType
			BusType = $_.BusType
			Size = $_.Size
			SerialNumber = $_.SerialNumber
			HealthStatus = $_.HealthStatus
			Firmware = $_.FirmwareVersion
			SectorSize = if ($wmiMatch) { $wmiMatch.BytesPerSector } else { 512 }
			TotalHeads = if ($wmiMatch) { $wmiMatch.TotalHeads } else { 0 }
			TotalCylinders = if ($wmiMatch) { $wmiMatch.TotalCylinders } else { 0 }
			PartitionsCount = if ($wmiMatch) { $wmiMatch.Partitions } else { 0 }
		}
	} | ConvertTo-Json
	`
	out, err := s.runPowerShell(psScript)
	if err != nil {
		return nil, err
	}

	var rawData []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &rawData); err != nil {
		return nil, err
	}

	var disks []models.DiskInfo
	for _, item := range rawData {
		idxStr := fmt.Sprintf("%v", item["DeviceId"])
		idx, _ := strconv.Atoi(idxStr)
		sizeVal, _ := item["Size"].(float64)
		cylVal, _ := item["TotalCylinders"].(float64)
		headsVal, _ := item["TotalHeads"].(float64)
		sectVal, _ := item["SectorSize"].(float64)
		partsVal, _ := item["PartitionsCount"].(float64)

		disk := models.DiskInfo{
			Index:           idx,
			DeviceID:        fmt.Sprintf("\\\\.\\PhysicalDrive%d", idx),
			Model:           fmt.Sprintf("%v", item["Model"]),
			MediaType:       fmt.Sprintf("%v", item["MediaType"]),
			BusType:         fmt.Sprintf("%v", item["BusType"]),
			Size:            int64(sizeVal),
			SerialNumber:    fmt.Sprintf("%v", item["SerialNumber"]),
			HealthStatus:    fmt.Sprintf("%v", item["HealthStatus"]),
			FirmwareVersion: fmt.Sprintf("%v", item["Firmware"]),
			TotalCylinders:  int64(cylVal),
			TotalHeads:      int(headsVal),
			SectorSize:      int(sectVal),
			PartitionsCount: int(partsVal),
		}
		disks = append(disks, disk)
	}
	return disks, nil
}
