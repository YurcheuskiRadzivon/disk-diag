package models

import "time"

const (
	MaxRuns = 20
)

type TestConfig struct {
	Retry       bool
	Balloon     bool
	Runs        int
	Seconds     int
	Threads     int
	SizeGiB     float64
	Dir         string
	IODuration  float64
	FileSize    int
	RandomBlock []byte
}

type WriteTestResult struct {
	Start          time.Time
	Runs           int
	CPUNum         int
	TotalSystemRAM uint64
	WorkingDir     string
	Results        []SingleWriteResult
}

type SingleWriteResult struct {
	Start        time.Time
	Duration     time.Duration
	WrittenBytes int
	WrittenMiB   int
	WrittenMB    float64
	MBps         float64
}

type ReadTestResult struct {
	Start          time.Time
	Runs           int
	CPUNum         int
	TotalSystemRAM uint64
	WorkingDir     string
	Results        []SingleReadResult
}

type SingleReadResult struct {
	Start       time.Time
	Duration    time.Duration
	ReadenBytes int
	ReadenMiB   int
	ReadenMB    float64
	MBps        float64
}

type IOPSTestResult struct {
	Start          time.Time
	Runs           int
	CPUNum         int
	TotalSystemRAM uint64
	WorkingDir     string
	Results        []SingleIOPSResult
}

type SingleIOPSResult struct {
	Start        time.Time
	IOOperations int
	IODuration   time.Duration
	IOPS         float64
}

type ThreadResult struct {
	Result int
	Error  error
}
