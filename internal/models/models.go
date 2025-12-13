package models

type DiskInfo struct {
	Index           int    `json:"index"`
	DeviceID        string `json:"deviceId"`
	Model           string `json:"model"`
	MediaType       string `json:"mediaType"`
	BusType         string `json:"busType"`
	Size            int64  `json:"size"`
	SerialNumber    string `json:"serialNumber"`
	HealthStatus    string `json:"healthStatus"`
	FirmwareVersion string `json:"firmwareVersion"`
	SectorSize      int    `json:"sectorSize"`
	TotalHeads      int    `json:"totalHeads"`
	TotalCylinders  int64  `json:"totalCylinders"`
	PartitionsCount int    `json:"partitionsCount"`
}

type PartitionInfo struct {
	DiskIndex          int    `json:"diskIndex"`
	PartitionID        string `json:"partitionId"`
	DriveLetter        string `json:"driveLetter"`
	Size               int64  `json:"size"`
	FreeSpace          int64  `json:"freeSpace"`
	FileSystem         string `json:"fileSystem"`
	Label              string `json:"label"`
	Type               string `json:"type"` // GPT/MBR
	Offset             int64  `json:"offset"`
	IsBoot             bool   `json:"isBoot"`
	IsReadOnly         bool   `json:"isReadOnly"`
	IsHidden           bool   `json:"isHidden"`
	BitLockerStatus    string `json:"bitLockerStatus"`
	AllocationUnitSize int    `json:"allocationUnitSize"`
	GptType            string `json:"gptType"`
}
