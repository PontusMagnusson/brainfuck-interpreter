package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang-collections/collections/stack"
)

// 30kB recommended for program memory space
var data [30000]byte
var instructions []byte
var ptr = 0
var loopStarts = stack.New()

func main() {
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	dirtyIns, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	instructions = cleanInstructions(dirtyIns)

	//fmt.Printf("%v\n", string(instructions))

	for insPtr := 0; insPtr < len(instructions); insPtr++ {
		//fmt.Printf("[%v : %v]\n", insPtr, string(instructions[insPtr]))
		//fmt.Printf("Value at current address (%v): %v\n", ptr, data[ptr])
		switch instructions[insPtr] {
		case 43: // +
			data[ptr]++
		case 44: // ,
			{
				fmt.Printf("Reading byte from io, [%d]\n", insPtr)
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadByte()
				data[ptr] = input
			}
		case 45: // -
			data[ptr]--
		case 46: // .
			fmt.Printf("%v", string(data[ptr]))
		case 62: // >
			ptr++
		case 60: // <
			ptr--
		case 91: // [
			{
				if data[ptr] == 0 {
					// set insPtr to bit after matching ]
					loopRange := getLoopEndIndex(insPtr)
					insPtr = insPtr + loopRange + 1
				} else {
					loopStarts.Push(insPtr) // Save location of loop start
				}
			}
		case 93: // ]
			{
				if data[ptr] != 0 {
					insPtr = loopStarts.Peek().(int)
				} else {
					loopStarts.Pop()
				}
			}
		}
		//fmt.Printf("Stack len: %d\n", loopStarts.Len())
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

func cleanInstructions(dirtyInstructions []byte) []byte {
	var buffer bytes.Buffer
	validOperands := []byte{'>', '<', '+', '-', '.', ',', '[', ']'}

	for _, char := range dirtyInstructions {
		if contains(validOperands, char) {
			buffer.WriteByte(char)
		}
	}

	return buffer.Bytes()
}

func getLoopEndIndex(loopStart int) int {
	relIndex := 0
	relDepth := 0

	for i, c := range instructions[loopStart+1:] {
		if c == 91 { // [
			relDepth++
		}
		if c == 93 { // ]
			if relDepth > 0 {
				relDepth--
			} else {
				relIndex = i
				break
			}
		}
	}
	return relIndex
}
