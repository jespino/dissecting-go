package main

import (
	"fmt"
	"math"
	"unsafe"
)

type mapStruct struct {
	count     int
	flags     uint8
	B         uint8
	noverflow uint16
	hash0     uint32

	buckets    unsafe.Pointer
	oldbuckets unsafe.Pointer
	nevacuate  uintptr

	extra *struct {
		overflow     []*bucketStruct
		oldoverflow  []*bucketStruct
		nextOverflow *bucketStruct
	}
}

type bucketStruct struct {
	topHash     uint64
	keys        [8]int
	values      [8]int
	overflowPtr uintptr
}

func main() {
	m := map[int]int{}
	ms := Scalpel(&m)
	Microscope(ms)

	m[1] = 10
	Microscope(ms)
	m[2] = 20
	m[3] = 30
	m[4] = 40
	m[5] = 50
	m[6] = 60
	m[7] = 70
	m[8] = 80
	m[9] = 90
	Microscope(ms)
	m[10] = 100
	m[11] = 110
	m[12] = 120
	m[13] = 130
	Microscope(ms)
}

func Scalpel(mapValue *map[int]int) *mapStruct {
	ms := unsafe.Pointer(*(*uintptr)(unsafe.Pointer(mapValue)))
	return (*mapStruct)(ms)
}

func Microscope(ms *mapStruct) {
	totalBuckets := int(math.Pow(2, float64(ms.B)))
	oldTotalBuckets := int(math.Pow(2, float64(ms.B-1)))
	fmt.Printf("Map size: %d\n", ms.count)
	fmt.Printf("Map flags: %d\n", ms.flags)
	fmt.Printf("Map B: %d\n", ms.B)
	fmt.Printf("Map number of overflow buckets (aprox): %d\n", ms.noverflow)
	fmt.Printf("Map hash seed: %d\n", ms.hash0)
	fmt.Printf("Map buckets: %v\n", ms.buckets)
	for x := 0; x < totalBuckets; x++ {
		bucket := uintptr(ms.buckets) + unsafe.Sizeof(bucketStruct{})*uintptr(x)
		data := (*bucketStruct)(unsafe.Pointer(bucket))
		fmt.Printf("  Bucket %d:\n", x)
		fmt.Printf("    Tophash: %v\n", data.topHash)
		fmt.Printf("    Keys: %v\n", data.keys)
		fmt.Printf("    Values: %v\n", data.values)
		fmt.Printf("    OverflowPtr: %v\n", data.overflowPtr)
		if data.overflowPtr != 0 {
			ovfBucket := data.overflowPtr
			ovfData := (*bucketStruct)(unsafe.Pointer(ovfBucket))
			fmt.Printf("      Overflow, Tophash: %v, Keys: %v, Values: %v, OverflowPtr: %v\n", ovfData.topHash, ovfData.keys, ovfData.values, ovfData.overflowPtr)
		}
	}
	fmt.Printf("Map old buckets: %v\n", ms.oldbuckets)
	if ms.oldbuckets != nil {
		for x := 0; x < oldTotalBuckets; x++ {
			bucket := uintptr(ms.oldbuckets) + unsafe.Sizeof(bucketStruct{})*uintptr(x)
			data := (*bucketStruct)(unsafe.Pointer(bucket))
			fmt.Printf("  Bucket %d:\n", x)
			fmt.Printf("    Tophash: %v\n", data.topHash)
			fmt.Printf("    Keys: %v\n", data.keys)
			fmt.Printf("    Values: %v\n", data.values)
			fmt.Printf("    OverflowPtr: %v\n", data.overflowPtr)
			if data.overflowPtr != 0 {
				ovfBucket := data.overflowPtr
				ovfData := (*bucketStruct)(unsafe.Pointer(ovfBucket))
				fmt.Printf("      Overflow:\n")
				fmt.Printf("        Tophash: %v\n", ovfData.topHash)
				fmt.Printf("        Keys: %v\n", ovfData.keys)
				fmt.Printf("        Values: %v\n", ovfData.values)
				fmt.Printf("        OverflowPtr: %v\n", ovfData.overflowPtr)
			}
		}
	}
	fmt.Printf("Map number of evacuated buckets: %d\n", ms.nevacuate)
}
