package main

func parse(bytecode []byte) [20][]byte {
	x := [20][]byte{}
	var j = 0
	var pc = 0

	for i := 0; i < len(bytecode); i++ {
		if bytecode[i] == 0x56 || bytecode[i] == 0x57 {
			x[j] = bytecode[pc : i+1]
			pc = i + 1
			j++
		}
	}

	return x
}
