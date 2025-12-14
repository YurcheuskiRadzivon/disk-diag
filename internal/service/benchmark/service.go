package benchmark

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
	sigar "github.com/cloudfoundry/gosigar"
)

const Blocksize = 0x1 << 16

type Service interface {
	WriteTest(diskIndex int, cfg *models.TestConfig) (*models.WriteTestResult, error)
	ReadTest(diskIndex int, cfg *models.TestConfig) (*models.ReadTestResult, error)
	IOPSTest(diskIndex int, cfg *models.TestConfig) (*models.IOPSTestResult, error)
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

func (s *service) physMemGet() (uint64, error) {
	mem := sigar.Mem{}
	err := mem.Get()
	if err != nil {
		return 0, err
	}
	return mem.Total, nil
}

func (s *service) SetTestDir(cfg *models.TestConfig) error {
	err := createDirIfNeeded(cfg.Dir)
	if err != nil {
		return err
	}
	cfg.Dir = path.Join(cfg.Dir, "diskdiag")
	return createDirIfNeeded(cfg.Dir)
}

func createDirIfNeeded(dir string) error {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	} else if !fileInfo.IsDir() {
		return errors.New(fmt.Sprintf("'%s' is not a directory!", dir))
	}
	return nil
}

func (s *service) CreateRandomBlock(cfg *models.TestConfig) error {
	cfg.RandomBlock = make([]byte, Blocksize)
	lenRandom, err := rand.Read(cfg.RandomBlock)
	if err != nil {
		return err
	}
	if len(cfg.RandomBlock) != lenRandom {
		return fmt.Errorf("CreateRandomBlock(): RandomBlock didn't get the correct number of bytes, %d != %d",
			len(cfg.RandomBlock), lenRandom)
	}
	return nil
}

func (s *service) MegaBytesPerSecond(bytes int, duration time.Duration) float64 {
	return float64(bytes) / float64(duration.Seconds()) / 1000000
}

func IOPS(operations int, duration time.Duration) float64 {
	return float64(operations) / float64(duration.Seconds())
}
