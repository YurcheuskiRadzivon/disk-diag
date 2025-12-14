package benchmark

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/models"
)

func (s *service) IOPSTest(diskIndex int, cfg *models.TestConfig) (*models.IOPSTestResult, error) {
	if diskIndex != 0 {
		return nil, fmt.Errorf("No avaliable write test for device with id - %d", diskIndex)
	}

	defaultTestDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, fmt.Errorf("tempDir failed - %w", err)
	}

	defer os.RemoveAll(defaultTestDir)
	log.Printf("%v", defaultTestDir)

	if cfg.Threads <= 0 || cfg.Threads > runtime.NumCPU() {
		cfg.Threads = 4
	}

	if cfg.IODuration <= 0 {
		cfg.IODuration = 15
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

	defer os.RemoveAll(cfg.Dir)

	log.Printf("%v", cfg)

	wtr := models.IOPSTestResult{}

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
		if err = s.RunIOPSTest(cfg, &wtr); err != nil {
			return nil, fmt.Errorf("runSequentialIOPSTest failed - %w", err)
		}
		wtr.Results[i].IOPS = IOPS(wtr.Results[i].IOOperations, wtr.Results[i].IODuration)
		log.Printf("%v", wtr)
	}

	return &wtr, nil
}

func (s *service) RunIOPSTest(cfg *models.TestConfig, wtr *models.IOPSTestResult) error {
	wtr.Results = append(wtr.Results, models.SingleIOPSResult{Start: time.Now()})
	newResult := &wtr.Results[len(wtr.Results)-1]

	opsPerformed := make(chan models.ThreadResult)
	start := time.Now()

	for i := 0; i < cfg.Threads; i++ {
		go s.singleThreadIOPSTest(cfg, wtr, path.Join(cfg.Dir, fmt.Sprintf("diskdiag.%d", i)), opsPerformed)
	}
	newResult.IOOperations = 0
	for i := 0; i < cfg.Threads; i++ {
		result := <-opsPerformed
		if result.Error != nil {
			return result.Error
		}
		newResult.IOOperations += result.Result
	}

	newResult.IODuration = time.Now().Sub(start)
	return nil
}

func (s *service) singleThreadIOPSTest(cfg *models.TestConfig, wtr *models.IOPSTestResult, filename string, numOpsChannel chan<- models.ThreadResult) {
	diskBlockSize := 0x1 << 9 // 512 bytes, nostalgia: in the olden (System V) days, disk blocks were 512 bytes
	fileInfo, err := os.Stat(filename)
	if err != nil {
		numOpsChannel <- models.ThreadResult{
			Result: 0, Error: err,
		}
		return
	}
	fileSizeLessOneDiskBlock := fileInfo.Size() - int64(diskBlockSize) // give myself room to not read past EOF
	numOperations := 0

	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		numOpsChannel <- models.ThreadResult{
			Result: 0, Error: err,
		}
		return
	}
	defer f.Close()

	data := make([]byte, diskBlockSize)
	checksum := make([]byte, diskBlockSize)

	start := time.Now()
	for i := 0; time.Now().Sub(start).Seconds() < cfg.IODuration; i++ { // run for xx seconds then blow this taco stand
		f.Seek(rand.Int63n(fileSizeLessOneDiskBlock), 0)
		// TPC-E has a reads:writes ratio of 9.7:1  http://www.cs.cmu.edu/~chensm/papers/TPCE-sigmod-record10.pdf
		// we round to 10:1
		if i%10 != 0 {
			length, err := f.Read(data)
			if err != nil {
				numOpsChannel <- models.ThreadResult{
					Result: 0, Error: err,
				}
				return
			}
			if length != diskBlockSize {
				panic(fmt.Sprintf("I expected to read %d bytes, instead I read %d bytes!", diskBlockSize, length))
			}
			for j := 0; j < diskBlockSize; j++ {
				checksum[j] ^= data[j]
			}
		} else {
			length, err := f.Write(checksum)
			if err != nil {
				numOpsChannel <- models.ThreadResult{
					Result: 0, Error: err,
				}
				return
			}
			if length != diskBlockSize {
				numOpsChannel <- models.ThreadResult{
					Result: 0,
					Error: fmt.Errorf("I expected to write %d bytes, instead I wrote %d bytes!",
						diskBlockSize, length),
				}
				return
			}
		}
		numOperations++
	}
	err = f.Close() // redundant, I know. I want to make sure writes are flushed
	if err != nil {
		numOpsChannel <- models.ThreadResult{
			Result: 0, Error: err,
		}
		return
	}
	numOpsChannel <- models.ThreadResult{Result: int(numOperations), Error: nil}
}
