package main

import "fmt"

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

func main() {
	memSize = 15

	p1 := procCreate("1").
		addSeg(&Seg{name: "1", len: 2}).
		addSeg(&Seg{name: "2", len: 3})
	p1.allocate()
	memory.print()

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

	p4 := procCreate("4").
		addSeg(&Seg{name: "1", len: 3}) //3 - 2
	p4.allocate()
	memory.print()

	//s1 := &Seg{
	//	name: "A",
	//	len:  5,
	//}
	//
	//s2 := &Seg{
	//	name: "B",
	//	len:  12,
	//}
	//
	//s3 := &Seg{
	//	name: "C",
	//	len:  2,
	//}
	//
	//s4 := &Seg{
	//	name: "D",
	//	len:  2,
	//}
	//
	//s5 := &Seg{
	//	name: "E",
	//	len:  3,
	//}
	//
	//s6 := &Seg{
	//	name: "F",
	//	len:  1,
	//}
	//
	//s7 := &Seg{
	//	name: "G",
	//	len:  1,
	//}
	//
	//s8 := &Seg{
	//	name: "H",
	//	len:  1,
	//}
	//
	////
	////memory.add(&MmuEntry{s1, 0})
	////memory.add(&MmuEntry{s2, 20})
	////memory.add(&MmuEntry{s3, 15})
	//
	//fmt.Println(alloc(s1))
	//fmt.Println(alloc(s2))
	//fmt.Println(alloc(s3))
	//
	////fmt.Println("FREE:")
	//fmt.Println(free("A"))
	////memory.print()
	//fmt.Println(alloc(s4))
	//fmt.Println(alloc(s5))
	//fmt.Println(alloc(s1))
	//fmt.Println(alloc(s6))
	//fmt.Println(alloc(s7))
	//fmt.Println(alloc(s8))
	//
	//memory.print()
}
