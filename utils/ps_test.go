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
	list, err := pl.SortByMem()
	if err != nil {
		t.Fatalf("Error:%v", err)
	}

	if pl.Number() != len(pl.Processes) {
		t.Fatal("Error:Number() return is not correct")
	}
}
