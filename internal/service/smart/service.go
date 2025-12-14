package smart

import (
	"context"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
)

type Service interface {
	GetNVMeSmart(diskIndex int) (*models.SmartInfo, error)
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
