package main

import (
	"fmt"
	"time"
	"unsafe"
)

type waitq struct {
	first uintptr
	last  uintptr
}

type channelStruct struct {
	qcount   uint      // total data in the queue
	dataqsiz uint      // size of the circular queue
	buf      *[4]int32 // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype uintptr // element type
	sendx    uint    // send index
	recvx    uint    // receive index
	recvq    waitq   // list of recv waiters
	sendq    waitq   // list of send waiters
	lock     uintptr
}

func Scalpel(channel *(chan int32)) *channelStruct {
	cs := unsafe.Pointer(*(*uintptr)(unsafe.Pointer(channel)))
	return (*channelStruct)(cs)
}

func Microscope(cs *channelStruct) {
	fmt.Printf("Total data in queue: %d\n", cs.qcount)
	fmt.Printf("Size of the queue: %d\n", cs.dataqsiz)
	fmt.Printf("Buffer address: %p\n", cs.buf)
	fmt.Printf("Element size: %d\n", cs.elemsize)
	fmt.Printf("Queued elements: %v\n", *cs.buf)
	fmt.Printf("Closed: %d\n", cs.closed)
	fmt.Printf("Element Type Address: %d\n", cs.elemtype)
	fmt.Printf("Send Index: %d\n", cs.sendx)
	fmt.Printf("Receive Index: %d\n", cs.recvx)
	fmt.Printf("Receive Wait list first address: 0x%x\n", cs.recvq.first)
	fmt.Printf("Receive Wait list last address: 0x%x\n", cs.recvq.last)
	fmt.Printf("Send Wait list first address: 0x%x\n", cs.sendq.first)
	fmt.Printf("Send Wait list last address: 0x%x\n", cs.sendq.last)
	fmt.Println("-------------------------------")
}

func main() {
	c := make(chan int32, 4)
	cs := Scalpel(&c)
	Microscope(cs)

	c <- 5
	Microscope(cs)
	go func() {
		c <- 4
		c <- 3
		c <- 2
		c <- 1
	}()
	time.Sleep(2 * time.Millisecond)
	Microscope(cs)

	<-c
	Microscope(cs)

	<-c
	Microscope(cs)

	go func() {
		<-c
		<-c
		<-c
		<-c
	}()
	time.Sleep(2 * time.Millisecond)
	Microscope(cs)

	close(c)

	Microscope(cs)
}
