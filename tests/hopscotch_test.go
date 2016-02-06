package tests

import "testing"
import "hophash"
import "strconv"
import "fmt"
import "math"
import "hash/fnv"

type fnvHash struct {
}


func (f *fnvHash) GetCode(key string) uint64 {
	h := fnv.New64a()
        h.Write([]byte(key))
        return h.Sum64()
}



func BenchmarkAddGet(b *testing.B) {

	var h hophash.HashMap
	var coder hophash.HashCoder
	coder = new(fnvHash)
	h = hophash.NewHophash(coder,64,28)
	for i:=0;i<b.N;i++ {
		testAddGet(h,uint64(i),b)
	}
}


func BenchmarkAddFirstGetLaterSmallSize(b *testing.B) {

	for i:=0;i<b.N;i++ {
		testAddFirstGetLater(b,10)
	}
}

func BenchmarkAddFirstGetLaterBuiltInSmallSize(b *testing.B) {


	for i:=0;i<b.N;i++ {
		testAddFirstGetLaterBuiltIn(b,10)	
	}
}

func BenchmarkAddFirstGetLaterLargeSize(b *testing.B) {

	for i:=0;i<b.N;i++ {
		testAddFirstGetLater(b,28)
	}
}

func BenchmarkAddFirstGetLaterBuiltInLargeSize(b *testing.B) {


	for i:=0;i<b.N;i++ {
		testAddFirstGetLaterBuiltIn(b,28)	
	}
}


func BenchmarkAddFirstGetLater(b *testing.B) {


	for i:=0;i<b.N;i++ {
		testAddFirstGetLater(b,16)
	}
}

func testAddFirstGetLater(b *testing.B,power int) {
	var h hophash.HashMap
	var coder hophash.HashCoder
	coder = new(fnvHash)
	h = hophash.NewHophash(coder,32,power)

	n := uint64(math.Pow(2,float64(power)))
	firstPart := uint64(n / 10)

	// do 10% inserts first
	
	for i:=uint64(0);i<firstPart;i++ {
		str := strconv.FormatUint(i,10)
		h.Add(str,str)
	}

	// now do contains with some misses and so hits
	// if "i" is even we should try successfull hit
	// if not we will try a miss
	for i:=(n-firstPart);i<n;i++ {
		if i % 2 == 0 {
			str := strconv.FormatUint((i/100),10)
			value,ok := h.Get(str)
			if !ok {				
				b.Fatal("value not found")
				return
			}
			if value != str {
				b.Fatal("value not same as what was added")
				return
			}
		} else {
			str := strconv.FormatUint(i,10)
			_,ok := h.Get(str)	
			if ok {
				b.Fatal("This key should not have been found")
			}
		}
	}
}


func BenchmarkAddFirstGetLaterBuiltIn(b *testing.B) {


	for i:=0;i<b.N;i++ {
		testAddFirstGetLaterBuiltIn(b,16)	
	}
}

func testAddFirstGetLaterBuiltIn(b *testing.B,power int) {
	
	m := make(map[string]string)

	n := uint64(math.Pow(2,float64(power)))
	firstPart := uint64(n / 10)

	// do 10% inserts first
	
	for i:=uint64(0);i<firstPart;i++ {
		str := strconv.FormatUint(i,10)
		m[str] = str
	}

	// now do contains with some misses and so hits
	// if "i" is even we should try successfull hit
	// if not we will try a miss
	for i:=(n-firstPart);i<n;i++ {
		if i % 2 == 0 {
			str := strconv.FormatUint((i/100),10)
			value,ok := m[str]
			if !ok {				
				b.Fatal("value not found")
				return
			}
			if value != str {
				b.Fatal("value not same as what was added")
				return
			}
		} else {
			str := strconv.FormatUint(i,10)
			_,ok := m[str]	
			if ok {
				b.Fatal("This key should not have been found")
			}
		}
	}
}


func BenchmarkAddGetHighLoad(b *testing.B) {

	var h hophash.HashMap
	var coder hophash.HashCoder
	coder = new(fnvHash)
	h = hophash.NewHophash(coder,16,28)
	for i:=uint64(0);i<uint64(math.Pow(2,20));i++ {
		testAddGet(h,i,b)
	}
}

func BenchmarkBuiltInMapHighLoad(b *testing.B) {

	var m map[string]string
	m = make(map[string]string)
	for i:=uint64(0);i<uint64(math.Pow(2,20));i++ {
		testMap(m,i,b)
	}
}




func testAddGet(h hophash.HashMap,index uint64, b *testing.B) {
	str := strconv.FormatUint(index,10)
	key := "key"+str
	err := h.Add(key,str)
	
	if err != nil {		
		b.Log("Warning could not add key: "+key)
		

		value,ok := h.Get(key)
		if !ok {
			fmt.Println(index)
			b.Fatal("value not found")
			return
		}

		if value != str {
			b.Fatal("value not same as what was added")
			return
		}
	}

}


func BenchmarkBuiltInMap(b *testing.B) {

	var m map[string]string
	m = make(map[string]string)
	for i:=0;i<b.N;i++ {
		testMap(m,uint64(i),b)
	}
}

func testMap(m map[string]string,index uint64, b *testing.B) {
	str := strconv.FormatUint(index,10)
	key := "key"+str
	
	_,contains := m[key]
	
	if !contains {
		m[key] = str
		
		value,ok := m[key]
		if !ok {
			fmt.Println(index)
			b.Fatal("value not found")
			return
		}

		if value != str {
			b.Fatal("value not same as what was added")
			return
		}
	}
	
}
