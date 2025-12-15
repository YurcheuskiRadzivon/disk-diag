package benchmark

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
)

func (s *service) WriteTest(diskIndex int, cfg *models.TestConfig) (*models.WriteTestResult, error) {
	if diskIndex != 0 {
		return nil, fmt.Errorf("No avaliable write test for device with id - %d", diskIndex)
	}

	defaultTestDir, err := ioutil.TempDir("", "diskdiag")
	if err != nil {
		return nil, fmt.Errorf("tempDir failed - %w", err)
	}
	if !cfg.Retry {
		defer os.RemoveAll(defaultTestDir)
	}
	log.Printf("%v", defaultTestDir)

	if cfg.Threads <= 0 || cfg.Threads > runtime.NumCPU() {
		cfg.Threads = 4
	}

	physMem, err := s.physMemGet()
	if err != nil {
		return nil, fmt.Errorf("physMemGet failed - %w", err)
	}

	if cfg.SizeGiB == 0 || cfg.SizeGiB > math.Floor(float64(2*int(physMem>>20)))/1024 {
		cfg.SizeGiB = math.Floor(float64(2*int((physMem/10)>>20))) / 1024
	}

	if cfg.Runs == 0 || cfg.Runs > models.MaxRuns {
		cfg.Runs = 1
	}

	if cfg.Dir == "" {
		cfg.Dir = defaultTestDir
	}

	if cfg.Seconds == 0 {
		cfg.Seconds = 300
	}

	err = s.SetTestDir(cfg)
	if err != nil {
		return nil, fmt.Errorf("setTestDir failed - %w", err)
	}
	if !cfg.Retry {
		defer os.RemoveAll(cfg.Dir)
	}

	log.Printf("%v", cfg)

	cfg.Balloon = false

	if cfg.Balloon {
		err = s.Balloon()
		if err != nil {
			return nil, fmt.Errorf("balloon failed - %w", err)
		}
	}

	wtr := models.WriteTestResult{}

	wtr.CPUNum = runtime.NumCPU()
	wtr.TotalSystemRAM = physMem >> 20
	wtr.WorkingDir = cfg.Dir
	wtr.Runs = cfg.Runs
	err = s.CreateRandomBlock(cfg)
	if err != nil {
		return nil, fmt.Errorf("createRandomBlock failed - %w", err)
	}

	finishTime := wtr.Start.Add(time.Duration(cfg.Seconds) * time.Second)
	for i := 0; (i < cfg.Runs) || time.Now().Before(finishTime); i++ {
		log.Printf("%v", finishTime)
		if err = s.RunSequentialWriteTest(cfg, &wtr); err != nil {
			return nil, fmt.Errorf("runSequentialWriteTest failed - %w", err)
		}
		wtr.Results[i].WrittenMiB = wtr.Results[i].WrittenBytes >> 20
		wtr.Results[i].WrittenMB = float64(wtr.Results[i].WrittenBytes) / 1000000
		wtr.Results[i].MBps = s.MegaBytesPerSecond(wtr.Results[i].WrittenBytes, wtr.Results[i].Duration)
		log.Printf("%v", wtr)
	}

	return &wtr, nil
}

func (s *service) RunSequentialWriteTest(cfg *models.TestConfig, wtr *models.WriteTestResult) error {
	cfg.FileSize = (int(cfg.SizeGiB*(1<<10)) << 20) / cfg.Threads

	log.Printf("%v", cfg.FileSize)
	for i := 0; i < cfg.Threads; i++ {
		err := os.RemoveAll(path.Join(cfg.Dir, fmt.Sprintf("diskdiag.%d", i)))
		if err != nil {
			return err
		}
	}

	log.Printf("%v", cfg)

	wtr.Results = append(wtr.Results, models.SingleWriteResult{Start: time.Now()})
	newResult := &wtr.Results[len(wtr.Results)-1]

	bytesWritten := make(chan models.ThreadResult)
	start := time.Now()
	log.Printf("%v", start)
	for i := 0; i < cfg.Threads; i++ {
		go s.singleThreadWriteTest(cfg, wtr, path.Join(cfg.Dir, fmt.Sprintf("diskdiag.%d", i)), bytesWritten)
	}

	log.Printf("%v", wtr)
	newResult.WrittenBytes = 0
	for i := 0; i < cfg.Threads; i++ {
		result := <-bytesWritten
		if result.Error != nil {
			return result.Error
		}
		newResult.WrittenBytes += result.Result
	}

	log.Printf("%v", wtr)
	newResult.Duration = time.Now().Sub(start)
	return nil
}

func (s *service) singleThreadWriteTest(cfg *models.TestConfig, wtr *models.WriteTestResult, filename string, bytesWrittenChannel chan<- models.ThreadResult) {
	f, err := os.Create(filename)
	if err != nil {
		bytesWrittenChannel <- models.ThreadResult{
			Result: 0, Error: err,
		}
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	bytesWritten := 0
	for i := 0; i < cfg.FileSize; i += len(cfg.RandomBlock) {
		n, err := w.Write(cfg.RandomBlock)
		if err != nil {
			bytesWrittenChannel <- models.ThreadResult{
				Result: 0, Error: err,
			}
			return
		}
		bytesWritten += n
	}

	err = w.Flush()
	if err != nil {
		bytesWrittenChannel <- models.ThreadResult{
			Result: 0, Error: err,
		}
		return
	}
	err = f.Close()
	if err != nil {
		bytesWrittenChannel <- models.ThreadResult{
			Result: 0, Error: err,
		}
		return
	}
	bytesWrittenChannel <- models.ThreadResult{Result: bytesWritten, Error: nil}
}
