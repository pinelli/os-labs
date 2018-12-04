package main

import (
	"github.com/davecgh/go-spew/spew"
)

type keystr string

func (key keystr) LessThan(key2 interface{}) bool{
	keyStr := keystr(key)
	key2Str := key2.(keystr)

	return string(keyStr) < string(key2Str)
}

func main(){
	tree := NewTree()
	var a keystr = "a"
	var b keystr = "b"
	var c keystr = "c"
	var d keystr = "d"
	var e keystr = "e"
	var g keystr = "g"
	var h keystr = "h"

	var i keystr = "i"

	tree.Insert(c, 3)
	tree.Insert(b, 1 )
	tree.Insert(a, 2)
	tree.Insert(d, 2)
	tree.Insert(e, 2)
	tree.Insert(g, 2)
	tree.Insert(h, 2)
	tree.Insert(i, 2)


	//fmt.Print("Found: ", tree.Find(keystr("c")))

	//k := keystr("e")
	//tree.Delete(k)

	spew.Dump(tree)


}
//fmt.Println(a.LessThan(b))
