package main

//插件模板
import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

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
	VirtualMemory *mem.VirtualMemoryStat `json:"vertual_memory"`
	SwapMemory    *mem.SwapMemoryStat    `json:"swap_memory"`
}

type Memory2 struct {
	VirtualMemory *VirtualMemoryStat `json:"vertual_memory"`
	SwapMemory    *SwapMemoryStat    `json:"swap_memory"`
}
type User struct {
	Info        *host.InfoStat         `jons:"info"`
	User        []host.UserStat        `json:"user"`
	Temperature []host.TemperatureStat `json:"temperature"`
}

type Cpu struct {
	Info []cpu.InfoStat `jons:"info"`
}
type Disk struct {
	Partitions []disk.PartitionStat           `json:"partitions"`
	IOCounters map[string]disk.IOCountersStat `json:"io_counters"`
	Usage      []*disk.UsageStat              `json:"usage"`
	// SerialNumbers map[string]string              `json:"serialNumbers"`
}
type Load struct {
	Misc *load.MiscStat `json:"misc"`
	Avg  *load.AvgStat  `jon:"avg"`
}
type Net struct {
	IOCounters     []net.IOCountersStat            `json:"io_counters"`
	ProtoCounters  []net.ProtoCountersStat         `json:"proto_counters"`
	FilterCounters []net.FilterStat                `json:"filter_counters"`
	ConntrackStats []net.ConntrackStat             `json:"conntrack_stats"`
	Connections    map[string][]net.ConnectionStat `json:"connections"`
	Interfaces     net.InterfaceStatList           `json:"interfaces"`
}
type Process struct {
	Processes []*process.Process `json:"processes"`
}

type Docker struct {
	GroupStat  []docker.CgroupDockerStat        `json:"docker_stat"`
	DockerList []string                         `json:"docker_list"`
	CpuStat    map[string]*docker.CgroupCPUStat `json:"cpu"`
	MemStat    map[string]*docker.CgroupMemStat `json:"mem"`
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
		cpuInfo.Info, _ = cpu.Info()
		data = cpuInfo
	case "disk":
		diskInfo := Disk{}
		diskInfo.Partitions, _ = disk.Partitions(true)
		diskInfo.IOCounters, _ = disk.IOCounters()
		// diskInfo.SerialNumbers = make(map[string]string, 0)
		for _, p := range diskInfo.Partitions {
			usage, err := disk.Usage(p.Mountpoint)
			if err == nil && usage != nil {
				diskInfo.Usage = append(diskInfo.Usage, usage)
			}
			// serial, err := disk.SerialNumber(p.Device)
			// if err == nil {
			// 	diskInfo.SerialNumbers[p.Device] = serial
			// }
		}

		data = diskInfo
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
		user := User{}
		user.User, _ = host.Users()
		user.Info, _ = host.Info()

		// if err != nil {
		// 	return nil, err
		// }
		data = user
	case "load":
		loadInfo := Load{}
		loadInfo.Avg, _ = load.Avg()
		loadInfo.Misc, _ = load.Misc()
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
		}
	default:
		isOk = false
		v = map[string]interface{}{"code": 500, "message": fmt.Sprintf("方法不支持：%s", name)}
	}
	if data != nil && isOk {
		v = map[string]interface{}{"code": 200, "result": data}
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
