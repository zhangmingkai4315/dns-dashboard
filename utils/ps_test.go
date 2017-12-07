package utils

import (
	"testing"
)

func TestNewProcessList(t *testing.T) {
	pl, err := NewProcessList()
	if err != nil {
		t.Fatalf("Error:%v", err)
	}
	if pl.Number() == 0 {
		t.Fatal("Error:size of process list is zero")
	}
}

func TestNumber(t *testing.T) {
	pl, err := NewProcessList()
	if err != nil {
		t.Fatalf("Error:%v", err)
	}
	if pl.Number() != len(pl.Processes) {
		t.Fatal("Error:Number() return is not correct")
	}
}

func TestSortByMem(t *testing.T) {
	pl, err := NewProcessList()
	if err != nil {
		t.Fatalf("Error:%v", err)
	}
	list, _ := pl.SortByMem(0, true)

	if pl.Number() != len(list) {
		t.Fatal("Error:Number() return is not correct")
	}
	for i, item := range pl.Processes {
		if i == (pl.Number() - 2) {
			break
		}
		if item.Memory < pl.Processes[i+1].Memory {
			t.Fatal("Error:return not sorted by memory")
		}
	}

	list, _ = pl.SortByMem(0, false)

	if pl.Number() != len(list) {
		t.Fatal("Error:Number() return is not desc correct")
	}
	for i, item := range pl.Processes {
		if i == (pl.Number() - 2) {
			break
		}
		if item.Memory > pl.Processes[i+1].Memory {
			t.Fatal("Error:return not sorted by asc memory")
		}
	}

	list, _ = pl.SortByMem(10, false)

	if 10 < len(list) {
		t.Fatal("Error: size return is not correct")
	}
	for i, item := range pl.Processes {
		if i+2 > len(list) {
			break
		}
		if item.Memory > pl.Processes[i+1].Memory {
			t.Fatal("Error:return not sorted by asc memory")
		}
	}

	list, _ = pl.SortByMem(10, true)

	if 10 < len(list) {
		t.Fatal("Error: size return is not correct")
	}
	for i, item := range pl.Processes {
		if i+2 > len(list) {
			break
		}
		if item.Memory < pl.Processes[i+1].Memory {
			t.Fatal("Error:return not sorted by asc memory")
		}
	}
}

func TestSortByCPU(t *testing.T) {
	pl, err := NewProcessList()
	if err != nil {
		t.Fatalf("Error:%v", err)
	}
	list, _ := pl.SortByCPU(0, true)

	if pl.Number() != len(list) {
		t.Fatalf("Error:Number() return is not correct expect :%d, got %d", pl.Number(), len(list))
	}

	for i, item := range list {
		if i == (pl.Number() - 2) {
			break
		}

		if item.CPU < list[i+1].CPU {
			t.Fatal("Error:return not sorted by CPU")
		}
	}

	list, _ = pl.SortByCPU(0, false)

	if pl.Number() != len(list) {
		t.Fatal("Error:Number() return is not desc correct")
	}
	for i, item := range list {
		if i == (pl.Number() - 2) {
			break
		}
		if item.CPU > list[i+1].CPU {
			t.Fatal("Error:return not sorted by asc CPU")
		}
	}
	list, _ = pl.SortByCPU(10, false)
	if 10 < len(list) {
		t.Fatal("Error: size return is not correct")
	}
	for i, item := range list {
		if i+2 > len(list) {
			break
		}

		if item.CPU > list[i+1].CPU {
			t.Fatal("Error:return not sorted by asc CPU")
		}
	}
	list, _ = pl.SortByCPU(10, true)
	if 10 < len(list) {
		t.Fatal("Error: size return is not correct")
	}
	for i, item := range list {
		if i+2 > len(list) {
			break
		}
		if item.CPU < list[i+1].CPU {
			t.Fatal("Error:return not sorted by asc CPU")
		}
	}
}
