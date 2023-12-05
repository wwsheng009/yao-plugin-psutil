package main

//插件模板
import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/docker"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"

	// "github.com/shirou/gopsutil/v3/winservices"
	"github.com/yaoapp/kun/grpc"
)

// https://github.com/shirou/gopsutil
// 定义插件类型，包含grpc.Plugin
type PsUtilPlugin struct{ grpc.Plugin }

// 设置插件日志到单独的文件
func (plugin *PsUtilPlugin) setLogFile() {
	var output io.Writer = os.Stdout
	//开启日志
	logroot := os.Getenv("GOU_TEST_PLG_LOG")
	if logroot == "" {
		logroot = "./logs"
	}
	if logroot != "" {
		logfile, err := os.OpenFile(path.Join(logroot, "psutil.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err == nil {
			output = logfile
		}
	}
	plugin.Plugin.SetLogger(output, grpc.Trace)
}

type Memory struct {
	VirtualMemory *mem.VirtualMemoryStat `json:"vertualMemory"`
	SwapMemory    *mem.SwapMemoryStat    `json:"swapMemory"`
}

type Memory2 struct {
	VirtualMemory *VirtualMemoryStat `json:"vertualMemory"`
	SwapMemory    *SwapMemoryStat    `json:"swapMemory"`
}
type Host struct {
	BootTime    string                 `json:"bootTime"`
	HostInfo    *host.InfoStat         `json:"hostInfo"`
	UserStat    []host.UserStat        `json:"userStat"`
	Temperature []host.TemperatureStat `json:"temperature"`
}

type Cpu struct {
	CpuInfo []cpu.InfoStat `jons:"cpuInfo"`
}

type DiskUsageStat struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             string  `json:"total"`
	Free              string  `json:"free"`
	Used              string  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       string  `json:"inodesTotal"`
	InodesUsed        string  `json:"inodesUsed"`
	InodesFree        string  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

type Disk struct {
	Partitions []disk.PartitionStat           `json:"diskPartitions"`
	IOCounters map[string]disk.IOCountersStat `json:"ioCounters"`
	Usage      []*disk.UsageStat              `json:"diskUsage"`
	Usages     []DiskUsageStat                `json:"diskUsages"`
	// SerialNumbers map[string]string              `json:"serialNumbers"`
}
type Load struct {
	LoadMisc *load.MiscStat `json:"loadMisc"`
	LoadAvg  *load.AvgStat  `jon:"loadAvg"`
}
type Net struct {
	IOCounters     []net.IOCountersStat            `json:"ioCounters"`
	ProtoCounters  []net.ProtoCountersStat         `json:"protoCounters"`
	FilterCounters []net.FilterStat                `json:"filterCounters"`
	ConntrackStats []net.ConntrackStat             `json:"conntrackStats"`
	Connections    map[string][]net.ConnectionStat `json:"connections"`
	Interfaces     net.InterfaceStatList           `json:"interfaces"`
}
type Process struct {
	Processes []*process.Process `json:"processes"`
}

type Docker struct {
	GroupStat  []docker.CgroupDockerStat        `json:"dockerStat"`
	DockerList []string                         `json:"dockerList"`
	CpuStat    map[string]*docker.CgroupCPUStat `json:"dockerCpu"`
	MemStat    map[string]*docker.CgroupMemStat `json:"dockerMem"`
}

// 插件执行需要实现的方法
// 参数name是在调用插件时的方法名，比如调用插件demo的Hello方法是的规则是plugins.demo.Hello时。
//
// 注意：name会自动的变成小写
//
// args参数是一个数组，需要在插件中自行解析。判断它的长度与类型，再转入具体的go类型。
//
// Exec 插件入口函数

func (plugin *PsUtilPlugin) Exec(name string, args ...interface{}) (*grpc.Response, error) {
	// plugin.Logger.Log(hclog.Trace, "plugin method called", name)
	// plugin.Logger.Log(hclog.Trace, "args", args)
	// isOk := true
	// var v = map[string]interface{}{"code": 200, "message": ""}
	isOk := true
	var v = map[string]interface{}{"code": 500, "message": "信息读取失败"}

	var data any
	var bytes []byte

	switch name {
	case "cpu":
		cpuInfo := Cpu{}
		cpuInfo.CpuInfo, _ = cpu.Info()
		data = cpuInfo
	case "disk":
		diskInfo := Disk{}
		readAll := false
		if len(args) > 0 {
			readAll = convertToBoolean(args[0])
		}
		diskInfo.Partitions, _ = disk.Partitions(readAll)
		diskInfo.IOCounters, _ = disk.IOCounters()
		// diskInfo.SerialNumbers = make(map[string]string, 0)
		for _, p := range diskInfo.Partitions {
			usage, err := disk.Usage(p.Mountpoint)
			if err == nil && usage != nil {
				diskInfo.Usage = append(diskInfo.Usage, usage)

				diskUsage := DiskUsageStat{}
				diskUsage.Path = usage.Path
				diskUsage.Fstype = usage.Fstype
				diskUsage.Total = formatBytes(usage.Total)
				diskUsage.Free = formatBytes(usage.Free)
				diskUsage.Used = formatBytes(usage.Used)
				diskUsage.UsedPercent = usage.UsedPercent
				diskUsage.InodesTotal = formatBytes(usage.InodesTotal)
				diskUsage.InodesUsed = formatBytes(usage.InodesUsed)
				diskUsage.InodesFree = formatBytes(usage.InodesFree)
				diskUsage.InodesUsedPercent = usage.InodesUsedPercent
				diskInfo.Usages = append(diskInfo.Usages, diskUsage)
			}

		}

		data = diskInfo
	case "dashboard":
		ds := DashboardService{}
		ioOption := "all"
		netOption := "all"

		if len(args) > 0 {
			ioOption = args[0].(string)
		}
		if len(args) > 1 {
			netOption = args[1].(string)
		}

		dsInfo, _ := ds.LoadBaseInfo(netOption, ioOption)
		data = dsInfo
	case "docker":
		dockerInfo := Docker{}
		dockerInfo.GroupStat, _ = docker.GetDockerStat()
		dockerInfo.DockerList, _ = docker.GetDockerIDList()
		dockerInfo.CpuStat = make(map[string]*docker.CgroupCPUStat)
		for _, id := range dockerInfo.DockerList {
			dockerInfo.CpuStat[id], _ = docker.CgroupCPUDocker(id)
			dockerInfo.MemStat[id], _ = docker.CgroupMemDocker(id)
		}
		data = dockerInfo
	case "host":
		hostInfo := Host{}
		hostInfo.UserStat, _ = host.Users()
		hostInfo.HostInfo, _ = host.Info()
		goTime := time.Unix(int64(hostInfo.HostInfo.BootTime), 0)
		hostInfo.BootTime = goTime.Format("2006-01-02 15:04:05")

		data = hostInfo
	case "load":
		loadInfo := Load{}
		loadInfo.LoadAvg, _ = load.Avg()
		loadInfo.LoadMisc, _ = load.Misc()
		data = loadInfo
	case "mem":
		mem1 := Memory{}
		mem1.VirtualMemory, _ = mem.VirtualMemory()
		mem1.SwapMemory, _ = mem.SwapMemory()
		data = mem1
	case "mem2":
		mem1 := Memory{}
		mem1.VirtualMemory, _ = mem.VirtualMemory()
		mem1.SwapMemory, _ = mem.SwapMemory()
		// 转换成mb/gb 显示
		mem2 := Memory2{}
		convertToFormattedString(&mem1, &mem2)
		data = mem2
	case "net2":
		data, _ = getNetConnections(NetConfig{})
	case "net":
		netInfo := Net{}
		netInfo.Connections = make(map[string][]net.ConnectionStat)
		list := []string{
			"all",
			"tcp",
			"tcp4",
			"tcp6",
			"udp",
			"udp4",
			"udp6",
			"inet",
			"inet4",
			"inet6"}
		for _, T := range list {
			netInfo.Connections[T], _ = net.Connections(T)
		}
		netInfo.ConntrackStats, _ = net.ConntrackStats(false)
		netInfo.FilterCounters, _ = net.FilterCounters()
		netInfo.IOCounters, _ = net.IOCounters(true)
		netInfo.Interfaces, _ = net.Interfaces()
		netInfo.ProtoCounters, _ = net.ProtoCounters([]string{})
		data = netInfo
	case "process":
		data, _ = getProcessData(PsProcessConfig{})
	case "ssh_session":
		data, _ = getSSHSessions(SSHSessionConfig{})
	case "winservices":
		if runtime.GOOS == "windows" {
			data = Winservice()
		} else {
			isOk = false
			v = map[string]interface{}{"code": 500, "message": fmt.Sprintf("方法只支持windows操作系统：%s", name)}
		}
	default:
		isOk = false
		v = map[string]interface{}{"code": 500, "message": fmt.Sprintf("方法不支持：%s", name)}
	}
	if data != nil && isOk {
		v = map[string]interface{}{"system": data}
		//输出前需要转换成字节
	}
	bytes1, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	bytes = bytes1

	//设置输出数据的类型
	//支持的类型：map/interface/string/integer,int/float,double/array,slice
	return &grpc.Response{Bytes: bytes, Type: "map"}, nil
}

// 生成插件时函数名修改成main
func main() {
	plugin := &PsUtilPlugin{}
	plugin.setLogFile()
	grpc.Serve(plugin)
}
