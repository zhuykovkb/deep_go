package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	fIndex := 0
	for i := 0; i < len(memory); i++ {
		oldPointer := unsafe.Pointer(&memory[i])

		if !existInMemory(oldPointer, pointers) {
			continue
		}

		memory[fIndex] = memory[i]
		newPointer := unsafe.Pointer(&memory[fIndex])

		for j := 0; j < len(pointers); j++ {
			if pointers[j] == oldPointer {
				pointers[j] = newPointer
			}
		}

		fIndex++
	}

	for i := fIndex; i < len(memory); i++ {
		memory[i] = 0x00
	}
}

func existInMemory(pointer unsafe.Pointer, pointers []unsafe.Pointer) bool {
	for _, v := range pointers {
		if pointer == v {
			return true
		}
	}
	return false
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
