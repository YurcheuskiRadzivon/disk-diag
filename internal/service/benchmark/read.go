package benchmark

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
)

func (s *service) ReadTest(diskIndex int, cfg *models.TestConfig) (*models.ReadTestResult, error) {
	if diskIndex != 0 {
		return nil, fmt.Errorf("No avaliable write test for device with id - %d", diskIndex)
	}

	defaultTestDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, fmt.Errorf("tempDir failed - %w", err)
	}
	defer os.RemoveAll(defaultTestDir)

	log.Printf("%v 11111111111111111111111111", defaultTestDir)

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

	cfg.Balloon = false

	if cfg.Balloon {
		err = s.Balloon()
		if err != nil {
			return nil, fmt.Errorf("balloon failed - %w", err)
		}
	}

	wtr := models.ReadTestResult{}

	wtr.CPUNum = runtime.NumCPU()
	wtr.TotalSystemRAM = physMem >> 20
	wtr.WorkingDir = cfg.Dir
	wtr.Runs = cfg.Runs

	finishTime := wtr.Start.Add(time.Duration(cfg.Seconds) * time.Second)
	for i := 0; (i < cfg.Runs) || time.Now().Before(finishTime); i++ {
		log.Printf("%v", finishTime)
		if err = s.RunSequentialReadTest(cfg, &wtr); err != nil {
			return nil, fmt.Errorf("runSequentialReadTest failed - %w", err)
		}
		wtr.Results[i].ReadenMiB = wtr.Results[i].ReadenBytes >> 20
		wtr.Results[i].ReadenMB = float64(wtr.Results[i].ReadenBytes) / 1000000
		wtr.Results[i].MBps = s.MegaBytesPerSecond(wtr.Results[i].ReadenBytes, wtr.Results[i].Duration)
		log.Printf("%v", wtr)
	}

	return &wtr, nil
}

func (s *service) RunSequentialReadTest(cfg *models.TestConfig, wtr *models.ReadTestResult) error {
	wtr.Results = append(wtr.Results, models.SingleReadResult{Start: time.Now()})
	newResult := &wtr.Results[len(wtr.Results)-1]

	bytesRead := make(chan models.ThreadResult)
	start := time.Now()

	for i := 0; i < cfg.Threads; i++ {
		go s.singleThreadReadTest(cfg, wtr, path.Join(cfg.Dir, fmt.Sprintf("diskdiag.%d", i)), bytesRead)
	}
	newResult.ReadenBytes = 0
	for i := 0; i < cfg.Threads; i++ {
		result := <-bytesRead
		if result.Error != nil {
			return result.Error
		}
		newResult.ReadenBytes += result.Result
	}

	newResult.Duration = time.Now().Sub(start)
	return nil
}

func (s *service) singleThreadReadTest(cfg *models.TestConfig, wtr *models.ReadTestResult, filename string, bytesReadChannel chan<- models.ThreadResult) {
	f, err := os.Open(filename)
	if err != nil {
		bytesReadChannel <- models.ThreadResult{
			Result: 0, Error: err,
		}
		return
	}
	defer f.Close()

	bytesRead := 0
	data := make([]byte, Blocksize)

	for {
		n, err := f.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			bytesReadChannel <- models.ThreadResult{
				Result: 0, Error: err,
			}
			return
		}
		bytesRead += n
		if bytesRead%127 == 0 {
			if !bytes.Equal(cfg.RandomBlock, data) {
				bytesReadChannel <- models.ThreadResult{
					Result: 0, Error: fmt.Errorf(
						"Most recent block didn't match random block, bytes read (includes corruption): %d",
						bytesRead),
				}
				return
			}
		}
	}

	bytesReadChannel <- models.ThreadResult{Result: bytesRead, Error: nil}
}
