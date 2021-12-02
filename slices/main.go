package main

import (
	"fmt"
	"unsafe"
)

type sliceStruct struct {
	array unsafe.Pointer
	len   int
	cap   int
}

func Scalpel(slice *[]int) *sliceStruct {
	ss := unsafe.Pointer(slice)
	return (*sliceStruct)(ss)
}

func Microscope(ss *sliceStruct) {
	fmt.Printf("Array Memory address: 0x%x\n", ss.array)
	fmt.Printf("Slice length: %d\n", ss.len)
	fmt.Printf("Slice capacity: %d\n", ss.cap)
	fmt.Printf("Stored data: [")
	for x := 0; x < ss.cap; x++ {
		fmt.Printf("%d,", *(*int)(unsafe.Pointer(uintptr(ss.array) + uintptr(x)*unsafe.Sizeof(int(0)))))
	}
	fmt.Println("]")
}

func main() {
	s := []int{}
	ss := Scalpel(&s)
	Microscope(ss)

	s = append(s, 1)
	Microscope(ss)

	s = append(s, 2)
	s = append(s, 3)
	s = append(s, 4)
	s = append(s, 5)
	Microscope(ss)

	subSlice := s[1:4]
	Microscope(Scalpel(&subSlice))

	subSlice[0] = 0
	Microscope(ss)
	Microscope(Scalpel(&subSlice))

	subSlice = append(subSlice, 6)
	Microscope(ss)
	Microscope(Scalpel(&subSlice))

	s = append(s, 6, 7, 8, 9)
	subSlice[1] = 0
	Microscope(ss)
	Microscope(Scalpel(&subSlice))
}
