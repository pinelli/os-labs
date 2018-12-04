package main

import "fmt"

func free(name string) bool {
	for i, e := range memory {
		if e.seg.name == name {
			memory = append(memory[:i], memory[i+1:]...)
			return true
		}
	}
	return false
}

func alloc(seg *Seg) bool {
	fmt.Print("Allocating ", seg.name, ", size:", seg.len, "->>")
	size := seg.len
	beg := 0
	for _, e := range memory {
		if e.offset-beg >= size {
			memory.add(&MmuEntry{
				seg:    seg,
				offset: beg,
			})
			return true
		}
		beg = e.offset + e.seg.len
	}

	if memSize-beg >= size {
		memory.add(&MmuEntry{
			seg:    seg,
			offset: beg,
		})
		return true
	}

	return false
}

//func alloc(seg *Seg) bool {
//	size := seg.len
//	beg := 0
//	for _, e := range memory {
//		if e.offset+e.seg.len-beg >= size {
//			memory.add(&MmuEntry{
//				seg:    seg,
//				offset: beg,
//			})
//			return true
//		}
//		beg = e.offset + e.seg.len
//	}
//
//	if memSize-beg >= size {
//		memory.add(&MmuEntry{
//			seg:    seg,
//			offset: beg,
//		})
//		return true
//	}
//
//	return false
//}

//func (e MmuEntry) String() string {
//	return e.seg.name + ":" + strconv.Itoa(e.offset)
//}
