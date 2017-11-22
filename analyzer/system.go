package analyzer

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// SystemStatus define the system status
type SystemStatus struct {
	Memory_State mem.VirtualMemoryStat `json:"virtual_memory"`
	CPU_Times    []cpu.TimesStat       `json:"cpu_times"`
	Disk_Info    []disk.PartitionStat  `json:"disk_info"`
	Disk_Usage   []disk.UsageStat      `json:"disk_usage"`
}

func GetSystemStatus() *SystemStatus {
	systemStatus := new(SystemStatus)
	v, err := mem.VirtualMemory()
	if err == nil {
		systemStatus.Memory_State = *v
	}
	cpuTimesState, err := cpu.Times(true)
	if err == nil {
		systemStatus.CPU_Times = cpuTimesState
	}
	diskInfo, err := disk.Partitions(false)
	if err == nil {
		var disk_usage []disk.UsageStat
		for _, d := range diskInfo {
			usage, err := disk.Usage(d.Mountpoint)
			if err == nil {
				disk_usage = append(disk_usage, *usage)
			}

		}
		systemStatus.Disk_Info = diskInfo
		systemStatus.Disk_Usage = disk_usage
	}
	return systemStatus
}
