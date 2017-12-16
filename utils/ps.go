// ps.go 通过执行ps命令获得本地的进程状态
package utils

import (
	"errors"
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
		args:        []string{psArgsOne, psArgsTwo},
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
		println(err.Error())
		return err
	}
	output := strings.Split(string(o), "\n")
	for i, processItem := range output {
		// 第一条记录是程序输出内容
		if i == 0 {
			continue
		}
		processItemList := strings.Fields(processItem)
		if len(processItemList) == 5 {
			p := new(Process)
			p.PID, _ = strconv.Atoi(processItemList[0])
			p.PPID, _ = strconv.Atoi(processItemList[1])
			p.Command = processItemList[2]
			memory, err := strconv.ParseFloat((processItemList[3]), 64)
			if err == nil {
				p.Memory = memory
			} else {
				p.Memory = 0.0
			}
			cpu, err := strconv.ParseFloat((processItemList[4]), 64)
			if err == nil {
				p.CPU = cpu
			} else {
				p.CPU = 0.0
			}
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
func (pl *ProcessList) SortByMem(size int, desc bool) ([]Process, error) {
	if len(pl.Processes) == 0 {
		return nil, errors.New("empty list")
	}
	tempList := []Process{}
	for _, item := range pl.Processes {
		temp := *item
		tempList = append(tempList, temp)
	}
	sort.Slice(tempList, func(i, j int) bool {
		if desc == true {
			return tempList[i].Memory > tempList[j].Memory
		}
		return tempList[i].Memory < tempList[j].Memory
	})
	if size <= 0 || size >= pl.Number() {
		return tempList, nil
	}
	return tempList[:size], nil
}

// SortByCPU define the sort func to sort the process
// with cpu usage
func (pl *ProcessList) SortByCPU(size int, desc bool) ([]Process, error) {
	if len(pl.Processes) == 0 {
		return nil, errors.New("empty list")
	}
	tempList := []Process{}
	for _, item := range pl.Processes {
		temp := *item
		tempList = append(tempList, temp)
	}
	sort.Slice(tempList, func(i, j int) bool {
		if desc == true {
			return tempList[i].CPU > tempList[j].CPU
		}
		return tempList[i].CPU < tempList[j].CPU
	})
	if size <= 0 || size >= pl.Number() {
		return tempList, nil
	}
	return tempList[:size], nil
}
