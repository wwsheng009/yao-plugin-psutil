package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/net"

	"github.com/shirou/gopsutil/v3/process"
)

type PsProcessConfig struct {
	Pid      int32  `json:"pid"`
	Name     string `json:"name"`
	Username string `json:"username"`
}
type PsProcessData struct {
	PID            int32  `json:"PID"`
	Name           string `json:"name"`
	PPID           int32  `json:"PPID"`
	Username       string `json:"username"`
	Status         string `json:"status"`
	StartTime      string `json:"startTime"`
	NumThreads     int32  `json:"numThreads"`
	NumConnections int    `json:"numConnections"`
	CpuPercent     string `json:"cpuPercent"`

	DiskRead  string `json:"diskRead"`
	DiskWrite string `json:"diskWrite"`
	CmdLine   string `json:"cmdLine"`

	Rss    string `json:"rss"`
	VMS    string `json:"vms"`
	HWM    string `json:"hwm"`
	Data   string `json:"data"`
	Stack  string `json:"stack"`
	Locked string `json:"locked"`
	Swap   string `json:"swap"`

	CpuValue float64 `json:"cpuValue"`
	RssValue uint64  `json:"rssValue"`

	Envs []string `json:"envs"`

	OpenFiles []process.OpenFilesStat `json:"openFiles"`
	Connects  []processConnect        `json:"connects"`
}
type processConnect struct {
	Type   string   `json:"type"`
	Status string   `json:"status"`
	Laddr  net.Addr `json:"localaddr"`
	Raddr  net.Addr `json:"remoteaddr"`
	PID    int32    `json:"PID"`
	Name   string   `json:"name"`
}
type ProcessConnects []processConnect

func (p ProcessConnects) Len() int {
	return len(p)
}

func (p ProcessConnects) Less(i, j int) bool {
	return p[i].PID < p[j].PID
}

func (p ProcessConnects) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

const (
	b  = uint64(1)
	kb = 1024 * b
	mb = 1024 * kb
	gb = 1024 * mb
)

func formatBytes(bytes uint64) string {
	switch {
	case bytes < kb:
		return fmt.Sprintf("%dB", bytes)
	case bytes < mb:
		return fmt.Sprintf("%.2fKB", float64(bytes)/float64(kb))
	case bytes < gb:
		return fmt.Sprintf("%.2fMB", float64(bytes)/float64(mb))
	default:
		return fmt.Sprintf("%.2fGB", float64(bytes)/float64(gb))
	}
}

func getProcessData(processConfig PsProcessConfig) (res []PsProcessData, err error) {
	var processes []*process.Process
	processes, err = process.Processes()
	if err != nil {
		return
	}

	var (
		result      []PsProcessData
		resultMutex sync.Mutex
		wg          sync.WaitGroup
		numWorkers  = 4
	)

	handleData := func(proc *process.Process) {
		procData := PsProcessData{
			PID: proc.Pid,
		}
		if processConfig.Pid > 0 && processConfig.Pid != proc.Pid {
			return
		}
		if procName, err := proc.Name(); err == nil {
			procData.Name = procName
		} else {
			procData.Name = "<UNKNOWN>"
		}
		if processConfig.Name != "" && !strings.Contains(procData.Name, processConfig.Name) {
			return
		}
		if username, err := proc.Username(); err == nil {
			procData.Username = username
		}
		if processConfig.Username != "" && !strings.Contains(procData.Username, processConfig.Username) {
			return
		}
		procData.PPID, _ = proc.Ppid()
		statusArray, _ := proc.Status()
		if len(statusArray) > 0 {
			procData.Status = strings.Join(statusArray, ",")
		}
		createTime, procErr := proc.CreateTime()
		if procErr == nil {
			t := time.Unix(createTime/1000, 0)
			procData.StartTime = t.Format("2006-1-2 15:04:05")
		}
		procData.NumThreads, _ = proc.NumThreads()
		connections, procErr := proc.Connections()
		if procErr == nil {
			procData.NumConnections = len(connections)
			for _, conn := range connections {
				if conn.Laddr.IP != "" || conn.Raddr.IP != "" {
					procData.Connects = append(procData.Connects, processConnect{
						Status: conn.Status,
						Laddr:  conn.Laddr,
						Raddr:  conn.Raddr,
					})
				}
			}
		}
		procData.CpuValue, _ = proc.CPUPercent()
		procData.CpuPercent = fmt.Sprintf("%.2f", procData.CpuValue) + "%"
		menInfo, procErr := proc.MemoryInfo()
		if procErr == nil {
			procData.Rss = formatBytes(menInfo.RSS)
			procData.RssValue = menInfo.RSS
			procData.Data = formatBytes(menInfo.Data)
			procData.VMS = formatBytes(menInfo.VMS)
			procData.HWM = formatBytes(menInfo.HWM)
			procData.Stack = formatBytes(menInfo.Stack)
			procData.Locked = formatBytes(menInfo.Locked)
			procData.Swap = formatBytes(menInfo.Swap)
		} else {
			procData.Rss = "--"
			procData.Data = "--"
			procData.VMS = "--"
			procData.HWM = "--"
			procData.Stack = "--"
			procData.Locked = "--"
			procData.Swap = "--"

			procData.RssValue = 0
		}
		ioStat, procErr := proc.IOCounters()
		if procErr == nil {
			procData.DiskWrite = formatBytes(ioStat.WriteBytes)
			procData.DiskRead = formatBytes(ioStat.ReadBytes)
		} else {
			procData.DiskWrite = "--"
			procData.DiskRead = "--"
		}
		procData.CmdLine, _ = proc.Cmdline()
		procData.OpenFiles, _ = proc.OpenFiles()
		procData.Envs, _ = proc.Environ()

		resultMutex.Lock()
		result = append(result, procData)
		resultMutex.Unlock()
	}

	chunkSize := (len(processes) + numWorkers - 1) / numWorkers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(processes) {
			end = len(processes)
		}

		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				handleData(processes[j])
			}
		}(start, end)
	}

	wg.Wait()

	sort.Slice(result, func(i, j int) bool {
		return result[i].PID < result[j].PID
	})
	res = result
	// res, err = json.Marshal(result)
	return
}
