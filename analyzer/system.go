package analyzer

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type processInfo struct {
	Pid         int32                `json:"pid"`
	Name        string               `json:"name"`
	ExecPath    string               `json:"exec_path"`
	Threads     int32                `json:"threads_number"`
	User        string               `json:"username"`
	Status      string               `json:"status"`
	MemoryUsage float32              `json:"memory_percent"`
	CPUUsage    float64              `json:"cpu_percent"`
	IOCounter   []net.IOCountersStat `json:"io"`
}

// SystemStatus define the system status
type SystemStatus struct {
	TimeStamp int64 `json:"timestamp"`
	// 主机基本信息
	HostInfoState host.InfoStat `json:"host_info"`
	// 内存信息
	MemoryState mem.VirtualMemoryStat `json:"virtual_memory"`
	// cpu信息
	CPUTimes []cpu.TimesStat `json:"cpu_times"`
	// 磁盘信息
	DiskInfo []disk.PartitionStat `json:"disk_info"`
	// 磁盘使用信息
	DiskUsage []disk.UsageStat `json:"disk_usage"`
	// 系统负载信息
	LoadAvg load.AvgStat `json:"load_avg"`
	// 网络IO操作信息
	IOState []net.IOCountersStat `json:"network_io"`
	// 系统进程信息
	Processes []processInfo `json:"processes"`
}

// GetSystemStatus will get some system runtime data
func GetSystemStatus() *SystemStatus {
	// 初始化
	systemStatus := new(SystemStatus)
	systemStatus.TimeStamp = time.Now().Unix()
	// 获取
	hostInfo, err := host.Info()
	if err == nil {
		systemStatus.HostInfoState = *hostInfo
	}

	// 获取内存
	v, err := mem.VirtualMemory()
	if err == nil {
		systemStatus.MemoryState = *v
	}
	// 获取cpu的信息
	cpuTimesState, err := cpu.Times(true)
	if err == nil {
		systemStatus.CPUTimes = cpuTimesState
	}
	// 获取disk的信息
	diskInfo, err := disk.Partitions(false)
	if err == nil {
		var diskUsage []disk.UsageStat
		for _, d := range diskInfo {
			usage, err := disk.Usage(d.Mountpoint)
			if err == nil {
				diskUsage = append(diskUsage, *usage)
			}

		}
		systemStatus.DiskInfo = diskInfo
		systemStatus.DiskUsage = diskUsage
	}
	// 获取系统负载情况
	avgStat, err := load.Avg()
	if err == nil {
		systemStatus.LoadAvg = *avgStat
	}

	// 获取网络网卡信息
	ioStat, err := net.IOCounters(true)
	if err == nil {
		systemStatus.IOState = ioStat
	}
	//获取系统进程信息

	processes, err := process.Processes()
	pInfos := []processInfo{}
	if err == nil {

		for _, p := range processes {
			if p.Pid == 0 {
				continue
			}
			name, err := p.Name()
			if err != nil || name == "" {
				continue
			}
			info := processInfo{}
			info.Pid = p.Pid
			info.Name, _ = p.Name()
			info.ExecPath, _ = p.Exe()
			info.Status, _ = p.Status()
			info.CPUUsage, _ = p.CPUPercent()
			info.IOCounter, _ = p.NetIOCounters(false)
			info.Threads, _ = p.NumThreads()
			info.MemoryUsage, _ = p.MemoryPercent()
			pInfos = append(pInfos, info)
		}
		systemStatus.Processes = pInfos
	}
	return systemStatus
}
