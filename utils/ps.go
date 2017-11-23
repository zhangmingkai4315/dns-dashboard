// ps.go 通过执行ps命令获得本地的进程状态
package utils

import (
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Process define the process infomation
type Process struct {
	PID     int     `json:"pid"`
	PPID    int     `json:"ppid"`
	Command string  `json:"command"`
	Memory  float64 `json:"memory_percent"`
	CPU     float64 `json:"cpu_percent"`
}

// ProcessList define list of current list
type ProcessList struct {
	processName string
	args        []string
	Processes   []*Process `json:"processes"`
}

// NewProcessList 创建一个新的ProcessList
func NewProcessList() (*ProcessList, error) {

	pl := &ProcessList{
		processName: "ps",
		args:        []string{"-eo", "pid,ppid,cmd,%mem,%cpu"},
		Processes:   []*Process{},
	}
	err := pl.GetProcesses()
	if err != nil {
		return nil, err
	}
	return pl, nil
}

// GetProcesses get the new ps command output and right to the list
func (pl *ProcessList) GetProcesses() error {
	if len(pl.Processes) > 0 {
		pl.Processes = pl.Processes[:0]
	}
	o, err := exec.Command(pl.processName, pl.args...).Output()
	if err != nil {
		return err
	}
	output := strings.Split(string(o), "\n")
	for i, processItem := range output {
		if i == 0 {
			continue
		}
		processItemList := strings.Fields(processItem)
		if len(processItemList) == 5 {
			p := new(Process)
			p.PID, _ = strconv.Atoi(processItemList[0])
			p.PPID, _ = strconv.Atoi(processItemList[1])
			p.Command = processItemList[2]
			p.Memory, _ = strconv.ParseFloat(string(processItem[3]), 32)
			p.CPU, _ = strconv.ParseFloat(string(processItem[4]), 32)
			pl.Processes = append(pl.Processes, p)
		}

	}
	return nil
}

// Number return the size of the Processes
func (pl *ProcessList) Number() int {
	return len(pl.Processes)
}

// SortByMem define the sort func to sort the process
// with memory usage
func (pl *ProcessList) SortByMem(size int, desc bool) []*Process {
	sort.Slice(pl.Processes, func(i, j int) bool {
		if desc == true {
			return pl.Processes[i].Memory > pl.Processes[j].Memory
		}
		return pl.Processes[i].Memory < pl.Processes[j].Memory
	})
	return pl.Processes[:size]
}

// SortByCPU define the sort func to sort the process
// with memory usage
func (pl *ProcessList) SortByCPU(size int, desc bool) []*Process {
	sort.Slice(pl.Processes, func(i, j int) bool {
		if desc == true {
			return pl.Processes[i].CPU > pl.Processes[j].CPU
		}
		return pl.Processes[i].CPU < pl.Processes[j].CPU
	})
	return pl.Processes[:size]
}
