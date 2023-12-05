package main

import (
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type DashboardBase struct {
	// WebsiteNumber  int `json:"websiteNumber"`
	// DatabaseNumber int `json:"databaseNumber"`
	// CronjobNumber  int `json:"cronjobNumber"`
	// AppInstalledNumber int `json:"appInstalledNumber"`

	// Hostname        string `json:"hostname"`
	// OS              string `json:"os"`
	// Platform        string `json:"platform"`
	// PlatformFamily  string `json:"platformFamily"`
	// PlatformVersion string `json:"platformVersion"`
	// KernelArch      string `json:"kernelArch"`
	// KernelVersion   string `json:"kernelVersion"`
	// VirtualizationSystem string         `json:"virtualizationSystem"`
	HostInfo        *host.InfoStat `json:"hostInfo"`
	CPUCores        int            `json:"cpuCores"`
	CPULogicalCores int            `json:"cpuLogicalCores"`
	CPUModelName    string         `json:"cpuModelName"`

	CurrentInfo DashboardCurrent `json:"currentInfo"`
	BootTime    string           `json:"bootTime"`
}

type DashboardCurrent struct {
	Uptime          uint64 `json:"uptime"`
	TimeSinceUptime string `json:"timeSinceUptime"`

	Procs uint64 `json:"procs"`

	Load1            float64 `json:"load1"`
	Load5            float64 `json:"load5"`
	Load15           float64 `json:"load15"`
	LoadUsagePercent float64 `json:"loadUsagePercent"`

	CPUPercent     []float64 `json:"cpuPercent"`
	CPUUsedPercent float64   `json:"cpuUsedPercent"`
	CPUUsed        float64   `json:"cpuUsed"`
	CPUTotal       int       `json:"cpuTotal"`

	MemoryTotal      uint64 `json:"memoryTotal"`
	MemoryAvailable  uint64 `json:"memoryAvailable"`
	MemoryUsed       uint64 `json:"memoryUsed"`
	MemoryTotals     string `json:"memoryTotals"`
	MemoryAvailables string `json:"memoryAvailables"`
	MemoryUseds      string `json:"memoryUseds"`

	MemoryUsedPercent float64 `json:"memoryUsedPercent"`

	SwapMemoryTotal       uint64  `json:"swapMemoryTotal"`
	SwapMemoryAvailable   uint64  `json:"swapMemoryAvailable"`
	SwapMemoryUsed        uint64  `json:"swapMemoryUsed"`
	SwapMemoryTotals      string  `json:"swapMemoryTotals"`
	SwapMemoryAvailables  string  `json:"swapMemoryAvailables"`
	SwapMemoryUseds       string  `json:"swapMemoryUseds"`
	SwapMemoryUsedPercent float64 `json:"swapMemoryUsedPercent"`

	IOReadBytes  uint64 `json:"ioReadBytes"`
	IOWriteBytes uint64 `json:"ioWriteBytes"`
	IOCount      uint64 `json:"ioCount"`
	IOReadTime   uint64 `json:"ioReadTime"`
	IOWriteTime  uint64 `json:"ioWriteTime"`

	DiskData []DiskInfo `json:"diskData"`

	NetBytesSent uint64 `json:"netBytesSent"`
	NetBytesRecv uint64 `json:"netBytesRecv"`

	ShotTime time.Time `json:"shotTime"`
}

type DiskInfo struct {
	Path        string  `json:"path"`
	Type        string  `json:"type"`
	Device      string  `json:"device"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`

	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

type DashboardService struct{}

type IDashboardService interface {
	LoadBaseInfo(ioOption string, netOption string) (*DashboardBase, error)
	LoadCurrentInfo(ioOption string, netOption string) *DashboardCurrent
}

func NewIDashboardService() IDashboardService {
	return &DashboardService{}
}

func (u *DashboardService) LoadBaseInfo(ioOption string, netOption string) (*DashboardBase, error) {
	var baseInfo DashboardBase
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	// baseInfo.Hostname = hostInfo.Hostname
	// baseInfo.OS = hostInfo.OS
	// baseInfo.Platform = hostInfo.Platform
	// baseInfo.PlatformFamily = hostInfo.PlatformFamily
	// baseInfo.PlatformVersion = hostInfo.PlatformVersion
	// baseInfo.KernelArch = hostInfo.KernelArch
	// baseInfo.KernelVersion = hostInfo.KernelVersion
	// ss, _ := json.Marshal(hostInfo)
	// baseInfo.VirtualizationSystem = string(ss)
	baseInfo.HostInfo = hostInfo
	goTime := time.Unix(int64(hostInfo.BootTime), 0)
	baseInfo.BootTime = goTime.Format("2006-01-02 15:04:05")
	cpuInfo, err := cpu.Info()
	if err == nil {
		baseInfo.CPUModelName = cpuInfo[0].ModelName
	}

	baseInfo.CPUCores, _ = cpu.Counts(false)
	baseInfo.CPULogicalCores, _ = cpu.Counts(true)

	baseInfo.CurrentInfo = *u.LoadCurrentInfo(ioOption, netOption)
	return &baseInfo, nil
}

func (u *DashboardService) LoadCurrentInfo(ioOption string, netOption string) *DashboardCurrent {
	var currentInfo DashboardCurrent
	hostInfo, _ := host.Info()
	currentInfo.Uptime = hostInfo.Uptime
	currentInfo.TimeSinceUptime = time.Now().Add(-time.Duration(hostInfo.Uptime) * time.Second).Format("2006-01-02 15:04:05")
	currentInfo.Procs = hostInfo.Procs

	currentInfo.CPUTotal, _ = cpu.Counts(true)
	totalPercent, _ := cpu.Percent(0, false)
	if len(totalPercent) == 1 {
		currentInfo.CPUUsedPercent = totalPercent[0]
		currentInfo.CPUUsed = currentInfo.CPUUsedPercent * 0.01 * float64(currentInfo.CPUTotal)
	}
	currentInfo.CPUPercent, _ = cpu.Percent(0, true)

	loadInfo, _ := load.Avg()
	currentInfo.Load1 = loadInfo.Load1
	currentInfo.Load5 = loadInfo.Load5
	currentInfo.Load15 = loadInfo.Load15
	currentInfo.LoadUsagePercent = loadInfo.Load1 / (float64(currentInfo.CPUTotal*2) * 0.75) * 100

	memoryInfo, _ := mem.VirtualMemory()
	currentInfo.MemoryTotal = memoryInfo.Total
	currentInfo.MemoryAvailable = memoryInfo.Available
	currentInfo.MemoryUsed = memoryInfo.Used

	currentInfo.MemoryTotals = formatBytes(memoryInfo.Total)
	currentInfo.MemoryAvailables = formatBytes(memoryInfo.Available)
	currentInfo.MemoryUseds = formatBytes(memoryInfo.Used)

	currentInfo.MemoryUsedPercent = memoryInfo.UsedPercent

	swapInfo, _ := mem.SwapMemory()
	currentInfo.SwapMemoryTotal = swapInfo.Total
	currentInfo.SwapMemoryAvailable = swapInfo.Free
	currentInfo.SwapMemoryUsed = swapInfo.Used
	currentInfo.SwapMemoryTotals = formatBytes(swapInfo.Total)
	currentInfo.SwapMemoryAvailables = formatBytes(swapInfo.Free)
	currentInfo.SwapMemoryUseds = formatBytes(swapInfo.Used)

	currentInfo.SwapMemoryUsedPercent = swapInfo.UsedPercent

	currentInfo.DiskData = loadDiskInfo()

	if ioOption == "all" {
		diskInfo, _ := disk.IOCounters()
		for _, state := range diskInfo {
			currentInfo.IOReadBytes += state.ReadBytes
			currentInfo.IOWriteBytes += state.WriteBytes
			currentInfo.IOCount += (state.ReadCount + state.WriteCount)
			currentInfo.IOReadTime += state.ReadTime
			currentInfo.IOWriteTime += state.WriteTime
		}
	} else {
		diskInfo, _ := disk.IOCounters(ioOption)
		for _, state := range diskInfo {
			currentInfo.IOReadBytes += state.ReadBytes
			currentInfo.IOWriteBytes += state.WriteBytes
			currentInfo.IOCount += (state.ReadCount + state.WriteCount)
			currentInfo.IOReadTime += state.ReadTime
			currentInfo.IOWriteTime += state.WriteTime
		}
	}

	if netOption == "all" {
		netInfo, _ := net.IOCounters(false)
		if len(netInfo) != 0 {
			currentInfo.NetBytesSent = netInfo[0].BytesSent
			currentInfo.NetBytesRecv = netInfo[0].BytesRecv
		}
	} else {
		netInfo, _ := net.IOCounters(true)
		for _, state := range netInfo {
			if state.Name == netOption {
				currentInfo.NetBytesSent = state.BytesSent
				currentInfo.NetBytesRecv = state.BytesRecv
			}
		}
	}

	currentInfo.ShotTime = time.Now()
	return &currentInfo
}

type diskInfo struct {
	Type   string
	Mount  string
	Device string
}

func loadDiskInfo() []DiskInfo {
	var datas []DiskInfo

	partitions, _ := disk.Partitions(false)
	// stdout, err := cmd.ExecWithTimeOut("df -hT -P|grep '/'|grep -v tmpfs|grep -v 'snap/core'|grep -v udev", 2*time.Second)
	// if err != nil {
	// 	return datas
	// }
	// lines := strings.Split(stdout, "\n")
	lines := make([]string, 0)
	for _, p := range partitions {
		lines = append(lines, p.Device)
	}

	var mounts []diskInfo
	var excludes = []string{"/mnt/cdrom", "/boot", "/boot/efi", "/dev", "/dev/shm", "/run/lock", "/run", "/run/shm", "/run/user"}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}
		if fields[1] == "tmpfs" {
			continue
		}
		if strings.Contains(fields[2], "M") || strings.Contains(fields[2], "K") {
			continue
		}
		if strings.Contains(fields[6], "docker") {
			continue
		}
		isExclude := false
		if runtime.GOOS != "windows" {
			for _, exclude := range excludes {
				if exclude == fields[6] {
					isExclude = true
				}
			}
		}
		if isExclude {
			continue
		}
		mounts = append(mounts, diskInfo{Type: fields[1], Device: fields[0], Mount: fields[6]})
	}

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	wg.Add(len(mounts))
	for i := 0; i < len(mounts); i++ {
		go func(timeoutCh <-chan time.Time, mount diskInfo) {
			defer wg.Done()

			var itemData DiskInfo
			itemData.Path = mount.Mount
			itemData.Type = mount.Type
			itemData.Device = mount.Device
			select {
			case <-timeoutCh:
				mu.Lock()
				datas = append(datas, itemData)
				mu.Unlock()
				// global.LOG.Errorf("load disk info from %s failed, err: timeout", mount.Mount)
			default:
				state, err := disk.Usage(mount.Mount)
				if err != nil {
					mu.Lock()
					datas = append(datas, itemData)
					mu.Unlock()
					// global.LOG.Errorf("load disk info from %s failed, err: %v", mount.Mount, err)
					return
					// return nil, errors.New(fmt.Sprintf("load disk info from %s failed, err: %v", mount.Mount, err))
				}
				itemData.Total = state.Total
				itemData.Free = state.Free
				itemData.Used = state.Used
				itemData.UsedPercent = state.UsedPercent
				itemData.InodesTotal = state.InodesTotal
				itemData.InodesUsed = state.InodesUsed
				itemData.InodesFree = state.InodesFree
				itemData.InodesUsedPercent = state.InodesUsedPercent
				mu.Lock()
				datas = append(datas, itemData)
				mu.Unlock()
			}
		}(time.After(5*time.Second), mounts[i])
	}
	wg.Wait()

	sort.Slice(datas, func(i, j int) bool {
		return datas[i].Path < datas[j].Path
	})
	return datas
}
