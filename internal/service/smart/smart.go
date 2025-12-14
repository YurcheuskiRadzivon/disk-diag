package smart

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
)

func (s *service) GetNVMeSmart(diskIndex int) (*models.SmartInfo, error) {
	dev := "/dev/sda"

	if diskIndex != 0 {
		return &models.SmartInfo{}, fmt.Errorf("No smart info for device with id - %d", diskIndex)
	}

	out, err := exec.Command("smartctl", "-a", "-d", "nvme", "-j", dev).Output()
	if err != nil {
		return nil, fmt.Errorf("smartctl error: %w", err)
	}

	var raw models.SmartInfo
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return &raw, nil
}
