package main

import (
	"fmt"
)

type Proc struct {
	id   string
	segs []*Seg
}

func procCreate(id string) *Proc {
	var p = &Proc{}
	p.id = id
	p.segs = nil
	return p
}

func (p *Proc) addSeg(seg *Seg) *Proc {
	seg.name = "P_" + p.id + "_" + seg.name
	p.segs = append(p.segs, seg)
	return p
}

func (p *Proc) allocate() {
	var allocated []*Seg = nil
	fmt.Println("Allocate memory for pid #", p.id)

	for _, e := range p.segs {

		if !alloc(e) {
			fmt.Println("Cannot allocate memory for segment", e.name, "revert changes")
			for _, s := range allocated {
				free(s.name)
			}
			return
		}
		fmt.Println("	Allocated seg ", e.name, ", size:", e.len)
		allocated = append(allocated, e)
	}
	fmt.Println("Successfully allocated memory for pid #", p.id)

}

func (p *Proc) free() {
	fmt.Println("Free memory for pid #", p.id)
	for _, e := range p.segs {
		fmt.Println("	Free seg ", e.name, ", size:", e.len)
		free(e.name)
	}
	fmt.Println("Successfully released memory for pid #", p.id)

}

type Seg struct {
	name string
	len  int
}

func acess(segName string, offset int) {
	fmt.Println("Access segment ", segName)
	for _, e := range memory {
		if e.seg.name == segName { //
			if offset > e.seg.len-1 {
				fmt.Println("Segmentation fault ")
				return
			}
			fmt.Println("	Seg base: ", e.offset)
			fmt.Println("	Offset: ", offset)
			fmt.Println("Accessed. Seg base + offset =", e.offset+offset)

			//fmt.Println(cnt, "--------", " ", e.seg.name, "&&", " END ", e.seg.name)
			//cnt++
			//continue
		}
	}
}

func main() {
	memSize = 15

	p1 := procCreate("1").
		addSeg(&Seg{name: "1", len: 2}).
		addSeg(&Seg{name: "2", len: 3})
	p1.allocate()
	memory.print()

	//acess("P_1_1", 1)

	p2 := procCreate("2").
		addSeg(&Seg{name: "1", len: 1}).
		addSeg(&Seg{name: "2", len: 1})
	p2.allocate()
	memory.print()

	p3 := procCreate("3").
		addSeg(&Seg{name: "1", len: 2}).
		addSeg(&Seg{name: "2", len: 3})
	p3.allocate()
	memory.print()

	p2.free()
	memory.print()

	//acess("P_1_2", 2) //3 -segfault

	p4 := procCreate("4").
		addSeg(&Seg{name: "1", len: 3}) //3 - 2
	p4.allocate()
	memory.print()
}
