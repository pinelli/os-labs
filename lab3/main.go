package main

import "fmt"

type Proc struct {
	id   string
	segs []*Seg
}

func (p *Proc) create(id string) {
	p.id = id
	p.segs = nil
}

func (p *Proc) addSeg(seg *Seg) {
	seg.name = "P_" + p.id + seg.name
	p.segs = append(p.segs, seg)
}

func (p *Proc) allocate() {
	var allocated []*Seg = nil

	for _, e := range p.segs {
		if !alloc(e) {
			for _, s := range allocated {
				free(s.name)
			}
		}
		allocated = append(allocated, e)
	}
	fmt.Println("Allocated memory for pid #", p.id)

}

type Seg struct {
	name string
	len  int
}

func main() {

	s1 := &Seg{
		name: "A",
		len:  5,
	}

	s2 := &Seg{
		name: "B",
		len:  12,
	}

	s3 := &Seg{
		name: "C",
		len:  2,
	}

	s4 := &Seg{
		name: "D",
		len:  2,
	}

	s5 := &Seg{
		name: "E",
		len:  3,
	}

	s6 := &Seg{
		name: "F",
		len:  1,
	}

	s7 := &Seg{
		name: "G",
		len:  1,
	}

	s8 := &Seg{
		name: "H",
		len:  1,
	}

	//
	//memory.add(&MmuEntry{s1, 0})
	//memory.add(&MmuEntry{s2, 20})
	//memory.add(&MmuEntry{s3, 15})

	fmt.Println(alloc(s1))
	fmt.Println(alloc(s2))
	fmt.Println(alloc(s3))

	//fmt.Println("FREE:")
	fmt.Println(free("A"))
	//memory.print()
	fmt.Println(alloc(s4))
	fmt.Println(alloc(s5))
	fmt.Println(alloc(s1))
	fmt.Println(alloc(s6))
	fmt.Println(alloc(s7))
	fmt.Println(alloc(s8))

	memory.print()
}
