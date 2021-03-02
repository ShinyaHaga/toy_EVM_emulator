package parser

import "fmt"

func parse(bytecode []byte) [20][]byte {
	x := [20][]byte{}
	var j = 0
	var pc = 0

	for i := 0; i < len(bytecode); i++ {
		if bytecode[i] == 0x56 {
			x[j] = bytecode[pc : i+1]
			pc = i + 1
			j++
		}
	}

	return x
}

func main() {
	fmt.Println(parse([]byte{0x11, 0x12, 0x32, 0x34, 0x56, 0x11, 0x32, 0x56}))
}
