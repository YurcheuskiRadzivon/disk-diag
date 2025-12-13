package base

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
)

type Service interface {
	GetPhysicalDisks() ([]models.DiskInfo, error)
	GetExtendedPartitions(diskIndex int) ([]models.PartitionInfo, error)
	GetCDiskInfo(diskIndex int) (models.CDiskInfo, error)
}

type service struct {
	ctx context.Context
}

func NewService(ctx context.Context) (*service, error) {
	srv := service{
		ctx: ctx,
	}

	return &srv, nil
}

func (s *service) runPowerShell(script string) (string, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	str := strings.TrimSpace(out.String())
	if strings.HasPrefix(str, "{") {
		str = "[" + str + "]"
	}
	if str == "" {
		str = "[]"
	}
	return str, err
}
