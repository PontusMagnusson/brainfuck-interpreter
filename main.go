package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// 30kb recommended for program memory space
var data [30000]int
var ptr = 0

func main() {
	filePath := os.Args[1]
	fmt.Println(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	instructions, err := ioutil.ReadAll(file)
	fmt.Println(instructions)
	fmt.Println(len(instructions))

	for _, ins := range instructions {
		switch ins {
		case 43: // +
			data[ptr]++
		case 45: // -
			data[ptr]--
		case 62: // >
			ptr++
		case 60: // <
			ptr--
		case 91: // [
			{
				if data[ptr] == 0 {
					// set ptr to bit after matching ]
				}
			}
		case 93:
			{
				if data[ptr] != 0 {
					// set ptr to previous [ in the instruction queue
				}
			}
		}
	}

	for _, bit := range data[:100] {
		fmt.Printf("%v ", bit)
	}
}
