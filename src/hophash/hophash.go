package hophash

import (
	"errors"
	"math"
)

type HashMap interface {
	Add(key string,value string) error
	Contains(key string) bool
	Get(key string) (string,bool)
}

type HashCoder interface {
	GetCode(key string) uint64
}


type entry struct {	
	value string
}

type item struct {
	hashKey uint64	
	entry *entry
}


type Hopscotch struct {
	arr []item	
	size uint64
	sizeForModulo uint64
	hopRange uint64
	coder HashCoder
}


func NewHophash(coder HashCoder,hopRange int,power int) *Hopscotch {
	h := new(Hopscotch)
	h.hopRange = uint64(hopRange)
	h.size = uint64(math.Pow(2,float64(power)))
	h.sizeForModulo = h.size - 1
	h.arr = make([]item,(h.size))
	h.coder = coder
	return h
}




func (h *Hopscotch) Add(key string,value string) error {
	// compute the hash value
	e := new(entry)
	e.value = value
	hashKey := h.coder.GetCode(key)
	ok := h.addToHopscoth(e,hashKey)

	if !ok {		
		return errors.New("Could not insert to map")
	} else {
		return nil
	}

	
}

func (h *Hopscotch) addToHopscoth(e *entry,hashKey uint64) bool {
	//base := e.hashKey % h.size

	base := hashKey & h.sizeForModulo

	var slot uint64
	
	slotFound := false
	for i:= base;i<h.size;i++ {
		if h.arr[i].entry == nil {
			// found an empty space
			slot  = i
			slotFound = true
			break
		}
	}

	if !slotFound {
		return false
	}

	// check if this within the hopRange of base
	if slot <= (base+h.hopRange-1) {
		//within hop-range
		h.arr[slot].entry = e	
		h.arr[slot].hashKey = hashKey
		return true
	} else {
		// not within hop range
		// adjust the table
		nextslot,found := h.adjustTable(slot,base)

		if !found {
			// resize
			// currently we do not resize, just return that
			// we cannot insert the key
			
			return false
		} else {
			h.arr[nextslot].entry = e			
			return true
		}
	}
	
}

func (h *Hopscotch) adjustTable(emptyIndex uint64,base uint64) (uint64,bool) {
	

	baseRange := (base + h.hopRange-1)
	// we can try adjusting until we reach the base
	for emptyIndex > base {
		
		slot,ok := h.getSlotToSwap(emptyIndex)
		
		if !ok {
			return 0,false // we need to resize the table
		}

		// swap the emptyIndex with slot
		hold := h.arr[slot].entry
		h.arr[slot].entry = nil
		h.arr[emptyIndex].entry = hold
		h.arr[emptyIndex].hashKey = h.arr[slot].hashKey
		emptyIndex = slot
		h.arr[emptyIndex].hashKey = 0

		// now check if the newly created empty slot is within the range of (H -1) of base
		if emptyIndex <= baseRange {
			// ok we have successfully created an empty slot in H-1 range of base, return it
			return emptyIndex,true
		}

		// else we continue looking, in the next set of previous buckets
	}
	return 0,false
}

func (h *Hopscotch) getSlotToSwap(emptyIndex uint64) (uint64,bool) {

	index := (emptyIndex - h.hopRange -1)

	var base uint64
	for index < emptyIndex {
	
		// get the base of the key in this index
		//base := h.arr[index].hashKey % h.size
		base = h.arr[index].hashKey & h.sizeForModulo
		
		// check if this base is too far from the emptyIndex to swap
		if emptyIndex <= (base + h.hopRange -1) {
			// within range, so we can swap this "index" with "emptyIndex"
			return index,true
		} 

		// else we just continue looking
		index++
	}
	// we did not get anything that we could replace in the previous n-1 entries
	// this means we need to resize the table
	return 0,false
}


func (h *Hopscotch) Contains(key string) bool {

	ok,_ := h.indexOf(key)	
	return ok
}

func (h *Hopscotch) getHopEnd(base uint64) uint64 {

	end := (base+h.hopRange-1)
	if end >= h.size {
		end = h.size -1
	}
	return end
}

func (h *Hopscotch) indexOf(key string) (bool,uint64) {

	hashKey := h.coder.GetCode(key)
	//base := hashKey % h.size
	base :=  h.coder.GetCode(key) & h.sizeForModulo 
	
	end := h.getHopEnd(base)

	for i:=base;i<=end;i++ {
		if h.arr[i].entry != nil &&  h.arr[i].hashKey == hashKey {
			return true,i
		}
	}	
	return false,0
}


func (h *Hopscotch) Get(key string) (string,bool) {

	ok,index := h.indexOf(key)
	
	if ok {
		return h.arr[index].entry.value,true
	}
	return "",false
}

