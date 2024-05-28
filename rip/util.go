package rip

import (
	"fmt"
	"os"
)

const TerminatedChar = 0x00

func CheckForError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func applyMask(ip [4]byte, mask [4]byte) [4]byte {
	var result [4]byte
	for i := 0; i < 4; i++ {
		result[i] = ip[i] & mask[i]
	}
	return result
}
