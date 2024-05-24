package main

import (
	"fmt"
	"os"
)

func main() {
	go UdpServer("localhost", 521)
	ripPacket := RipPacket{
		Command: 2,
		Version: 2,
		Unused:  [2]byte{0, 0},
		RouterEntry: []RouterEntry{
			{
				AddressFamilyIdentifier: 2,
				RouteTag:                0,
				IpAddress:               [4]byte{192, 168, 1, 1},
				SubMask:                 [4]byte{255, 255, 255, 0},
				NextHop:                 [4]byte{0, 0, 0, 0},
				Metric:                  1,
			},
		},
	}
	UdpClient("localhost", 521, ripPacket)
	ReadConfig()
}

func CheckForError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
