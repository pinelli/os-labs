package main

import (
	"fmt"
	"sort"
)

var memory Memory
var memSize int // = 26

type MmuEntry struct {
	seg    *Seg
	offset int
}
type Memory []*MmuEntry

func (mem Memory) Len() int {
	return len(mem)
}

func (mem Memory) Swap(i, j int) {
	mem[i], mem[j] = mem[j], mem[i]
}

func (mem Memory) Less(i, j int) bool {
	return mem[i].offset < mem[j].offset
}

func (mem *Memory) add(entry *MmuEntry) {
	memory = append(memory, entry)
	sort.Sort(Memory(memory))
}

func (mem *Memory) print() {
	//fmt.Println(([]*MmuEntry)(memory))
	//spew.Dump(memory)
	cnt := 0
	for _, e := range *mem {
		if e.seg.len == 1 { //
			fmt.Println(cnt, "--------", " ", e.seg.name, "&&", " END ", e.seg.name)
			cnt++
			continue
		}

		for ; cnt < e.offset; cnt++ {
			fmt.Println(cnt)
		}

		fmt.Println(cnt, "--------", " ", e.seg.name)
		cnt++
		cntFrom := cnt
		for ; cnt < cntFrom+e.seg.len-2; cnt++ {
			fmt.Println(cnt)
		}

		fmt.Println(cnt, "--------", " 	END ", e.seg.name)
		cnt++

	}

	for ; cnt < memSize; cnt++ {
		fmt.Println(cnt)
	}
}
