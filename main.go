package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang-collections/collections/stack"
)

// 30kb recommended for program memory space
var data [30000]int
var instructions []byte
var ptr = 0
var stk = stack.New() // Maybe not needed?

func main() {
	filePath := os.Args[1]
	fmt.Println(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	instructions, err = ioutil.ReadAll(file)
	fmt.Println(instructions)

	for insPtr := 0; insPtr < len(instructions); insPtr++ {
		switch instructions[insPtr] {
		case 43: // +
			data[ptr]++
		case 45: // -
			data[ptr]--
		case 46:
			fmt.Printf("%v", string(data[ptr]))
		case 62: // >
			ptr++
		case 60: // <
			ptr--
		case 91: // [
			{
				if data[ptr] == 0 {
					// set insPtr to bit after matching ]
					fmt.Printf("Skipping loop at %d \n", insPtr)
					loopRange := getLoopEndIndex(insPtr)
					insPtr = insPtr + loopRange
				} else {
					stk.Push(insPtr) // Save location of loop start
				}
			}
		case 93:
			{
				if data[ptr] != 0 {
					insPtr = stk.Peek().(int)
				} else {
					stk.Pop()
				}
			}
		}
	}
	fmt.Println("")
	for _, bit := range data[:100] {

		fmt.Printf("%v ", bit)
	}
}

func contains(arr []byte, a byte) bool {
	for _, b := range arr {
		if a == b {
			return true
		}
	}
	return false
}

func getLoopEndIndex(loopStart int) int {
	relIndex := 0
	relDepth := 0

	for i, c := range instructions[loopStart+1:] {
		if c == 91 { // [
			fmt.Printf("Found nested loop at rel %d \n", i)
			relDepth++
		}
		if c == 93 { // ]
			if relDepth > 0 {
				fmt.Printf("Found end of nested loop %d \n", i)
				relDepth--
			} else {
				fmt.Printf("Found loop end at %d \n", i)
				relIndex = i
				break
			}
		}
	}
	return relIndex
}
