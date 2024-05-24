package main

import (
	"fmt"
	"net"
)

func UdpServer(host string, port int) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	CheckForError(err)

	connection, err := net.ListenUDP("udp", udpAddr)
	CheckForError(err)

	for {
		var buf [1024]byte
		n, addr, err := connection.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println(err)
			continue
		}

		packet, err := UnmarshalRipPacket(buf[:n])
		if err != nil {
			fmt.Println("Error unmarshalling RipPacket:", err)
			continue
		}

		fmt.Println("Received from client: \n", packet)
		connection.WriteToUDP([]byte("Hello from server\n"), addr)
	}

}
