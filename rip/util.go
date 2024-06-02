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

func parseIpToBytes(ip string) [4]byte {
	var result [4]byte
	fmt.Sscanf(ip, "%d.%d.%d.%d", &result[0], &result[1], &result[2], &result[3])
	return result
}

func parseBytesToIp(ip [4]byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}
