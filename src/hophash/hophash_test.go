package hophash

import (
	"testing"
	"hash/fnv"
	"math"
	"strconv"
)

type testHash struct {
	m map[string]uint64
}



func (t *testHash) GetCode(key string) uint64 {
	return t.m[key]
}




func TestFitWithinRangeSameBase(t *testing.T) {

	var coder HashCoder
	tHash := new(testHash)
	coder = tHash
	tHash.m = make(map[string]uint64)
	tHash.m["a"] = 1
	tHash.m["b"] = 1
	tHash.m["c"] = 1
	tHash.m["d"] = 1
	
	h := NewHophash(coder,4,16)
	// scenario to be tested:
	// fits within the hop range

	h.Add("a","a")
	h.Add("b","b")
	h.Add("c","c")
	h.Add("d","d")

	if h.arr[4].value != "d" {
		t.Fail()
	}

}


func TestFitWithinRangeDifBase(t *testing.T) {

	var coder HashCoder
	tHash := new(testHash)
	coder = tHash
	tHash.m = make(map[string]uint64)
	tHash.m["a"] = 1
	tHash.m["b"] = 2
	tHash.m["c"] = 2
	tHash.m["d"] = 4
	
	h := NewHophash(coder,4,16)
	// scenario to be tested:
	// fits within the hop range

	h.Add("a","a")
	h.Add("b","b")
	h.Add("c","c")
	h.Add("d","d")

	if h.arr[4].value != "d" {
		t.Fail()
	}

	if h.arr[3].value != "c" {
		t.Fail()
	}

	if h.arr[1].value != "a" {
		t.Fail()
	}

}


func TestOutOfRange(t *testing.T) {

	var coder HashCoder
	tHash := new(testHash)
	coder = tHash
	tHash.m = make(map[string]uint64)
	tHash.m["a"] = 1
	tHash.m["b"] = 2
	tHash.m["c"] = 2
	tHash.m["d"] = 2
	
	h := NewHophash(coder,4,16)
	// scenario to be tested:
	// fits within the hop range

	h.Add("a","a")
	h.Add("b","b")
	h.Add("c","c")
	h.Add("d","d")

	tHash.m["e"] = 1
	// now we expect b in slot 2 to be moved to 5 and e to put in 2
	h.Add("e","e")


	if h.arr[2].value != "e" {
		t.Fail()
	}

	if h.arr[5].value != "b" {
		t.Fail()
	}

}



func TestOutOfRangeMultipleIterationsToFindEmptySlot(t *testing.T) {

	var coder HashCoder
	tHash := new(testHash)
	coder = tHash
	tHash.m = make(map[string]uint64)
	tHash.m["a"] = 1
	tHash.m["b"] = 2
	tHash.m["c"] = 2
	tHash.m["d"] = 3
	tHash.m["e"] = 3
	tHash.m["f"] = 4
	tHash.m["g"] = 6
	tHash.m["h"] = 6
	tHash.m["i"] = 8

	h := NewHophash(coder,4,16)
	// scenario to be tested:
	// fits within the hop range

	h.Add("a","a")
	h.Add("b","b")
	h.Add("c","c")
	h.Add("d","d")
	h.Add("e","e")
	h.Add("f","f")
	h.Add("g","g")
	h.Add("h","h")
	h.Add("i","i")


	tHash.m["j"] = 4
	// now we expect in multiple iterations
	// to place i(base=4) in slot at 7, and then first i to moved to 10 and g to be moved to 9
	h.Add("j","j")

	if h.arr[10].value != "i" {
		t.Fail()
	}

	if h.arr[9].value != "g" {
		t.Fail()
	}

	if h.arr[7].value != "j" {
		t.Fail()
	}

}


type fnvHash struct {
}


func (f *fnvHash) GetCode(key string) uint64 {
	h := fnv.New64a()
        h.Write([]byte(key))
        return h.Sum64()
}


func TestAddGetHighLoadLimited(t *testing.T) {

	var coder HashCoder
	coder = new(fnvHash)
	h := NewHophash(coder,2,uint64(16))

        for i:=1;i<15;i++ {	
		str := strconv.Itoa(i)
		key := "key"+str
		
		err := h.Add(key,str)
		
		
		if err != nil {		
			t.Log("Warning: could not add key: "+key)		
		
		} else {
			value,ok := h.Get(key)
			if !ok {
				t.Fatal(i)
				t.Fatal("value not found")			
			}

			if value != str {
				t.Fatal("value not same as what was added")			
			}
		}
	}
}


func TestAddGetHighLoad(t *testing.T) {

	var coder HashCoder
	coder = new(fnvHash)
	h := NewHophash(coder,32,uint64(math.Pow(2,28)))

	for i:=0;i<(int(math.Pow(2,16))-2000);i++ {	
		str := strconv.Itoa(i)
		key := "key"+str
		err := h.Add(key,str)
		
		if err != nil {		
			t.Log("Warning: could not add key: "+key)
		} else {

			value,ok := h.Get(key)
			if !ok {
				t.Fatal(i)
				t.Fatal("value not found")			
			}

			if value != str {
				t.Fatal("value not same as what was added")			
			}
		}
	}
}

