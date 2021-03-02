package main

import (
	"encoding/hex"
)

func str_to_bytes(str string) []byte {
	decode, _ := hex.DecodeString(str)
	return decode
}
