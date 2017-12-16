package analyzer

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/zhangmingkai4315/dns-dashboard/utils"
)

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
	// 系统进程信息Top10内存
	ProcessesMemory []utils.Process `json:"processes_memory"`
	// 系统进程信息Top10CPU
	ProcessesCPU []utils.Process `json:"processes_cpu"`
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
	processes, err := utils.NewProcessList()
	if err == nil {
		// 默认仅仅获取前10条 内存占用最大的进程
		systemStatus.ProcessesMemory, _ = processes.SortByMem(10, true)
		systemStatus.ProcessesCPU, _ = processes.SortByCPU(10, true)
	}
	return systemStatus
}
