package main

import (
	"fmt"
	"os"
)

func main() {
	go UdpServer("localhost", 521)
	//routerEntries := ReadConfig()
	//ripPacket := RipPacket{Command: 1, Version: 2, Unused: [2]byte{0, 0}, RouterEntries: routerEntries}
	routerConfigEntry, err := ReadConfigAsBytes("config/routeur-r1.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	UdpClient("localhost", 521, routerConfigEntry)
}

func CheckForError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
