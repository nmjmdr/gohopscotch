package main

import (
	"hophash"
	"fmt"
	"hash/fnv"
	"math"
)

type fnvHash struct {
}


func (f *fnvHash) GetCode(key string) uint64 {
	h := fnv.New64a()
        h.Write([]byte(key))
        return h.Sum64()
}

func main() {
	var h hophash.HashMap
	
	var coder hophash.HashCoder
	coder = new(fnvHash)

	h = hophash.NewHophash(coder,32,uint64(math.Pow(2,16)))

	h.Add("key1","value1")	
	val,ok := h.Get("key1");
	if ok {
		fmt.Println(val)
	} else {
		fmt.Println("no value present")
	}
}
