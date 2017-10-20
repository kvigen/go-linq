package main

import (
	"fmt"
	"reflect"
	"time"
)

type input struct {
	A int64
	B int64
}

const numEntries = 1000 * 1000

func testFunc(a, b int64) int64 {
	return a + b
}

func main() {

	total := int64(0)
	start := time.Now()

	for i := 0; i < numEntries; i++ {
		total += testFunc(2, 4)
	}
	fmt.Printf("Total: %d in %dms\n", total, time.Since(start)/1000/1000)

	total = int64(0)
	start = time.Now()
	r := reflect.ValueOf(testFunc)
	values := []reflect.Value{reflect.ValueOf(int64(2)), reflect.ValueOf(int64(4))}
	for i := 0; i < numEntries; i++ {
		values := r.Call(values)
		total += values[0].Int()
	}

	fmt.Printf("Total: %d in %dms\n", total, time.Since(start)/1000/1000)

	total = int64(0)
	start = time.Now()
	giantChannel := make(chan int64, numEntries+10)
	for i := 0; i < numEntries; i++ {
		giantChannel <- 2
	}

	for i := 0; i < numEntries; i++ {
		total += <-giantChannel
	}

	fmt.Printf("Total: %d in %dms\n", total, time.Since(start)/1000/1000)

	/*
		inputs := make([]*input, 0)

		for i := 0; i < numEntries; i++ {
			inputs = append(inputs, &input{2, 4})
		}
		fmt.Printf("Input len: %d\n", len(inputs))

		total := int64(0)
		start := time.Now()
		for _, input := range inputs {
			r := reflect.ValueOf(input)
			val := reflect.Indirect(r).FieldByName("B").Int()
			total += val
		}
		fmt.Printf("Total: %d in %dms\n", total, time.Since(start)/1000/1000)

		total = int64(0)
		start = time.Now()
		for _, input := range inputs {
			total += input.B
		}
		fmt.Printf("Total: %d in %dms\n", total, time.Since(start)/1000/1000)

		total = int64(0)
		start = time.Now()
		input := inputs[0]
		// TODO: Need to handle false here...
		field, _ := reflect.TypeOf(input).FieldByName("B")
		offset := field.Offset
		fmt.Printf("Offset: %d\n", offset)

		for _, input := range inputs {
			ptr := unsafe.Pointer(&input)
			ptr = unsafe.Pointer(uintptr(ptr) + offset)
			v := (*int64)(ptr)
			fmt.Printf("V: %d\n", *v)
			*v = int64(6)
			fmt.Printf("V: %d\n")
		}

		for _, input := range inputs {
			ptr := unsafe.Pointer(&input)
			ptr = unsafe.Pointer(uintptr(ptr) + offset)
			v := *(*int64)(ptr)
			total += v
		}
		fmt.Printf("Total: %d in %dms\n", total, time.Since(start)/1000/1000)

		maps := make([]map[string]interface{}, 0)
		for i := 0; i < numEntries; i++ {
			m := map[string]interface{}{"A": int64(2), "B": int64(4)}
			maps = append(maps, m)
		}
		total = int64(0)
		start = time.Now()
		for _, m := range maps {
			val := m["B"].(int64)
			total += val
		}
		fmt.Printf("Total: %d in %dms\n", total, time.Since(start)/1000/1000)
	*/
}
