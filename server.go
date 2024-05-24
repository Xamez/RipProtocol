package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
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

		//data, err := UnmarshalRipPacket(buf[:n])
		//if err != nil {
		//	fmt.Println("Error unmarshalling RipPacket:", err)
		//	continue
		//}

		var routerConfigEntry []map[string]RouterConfigEntry

		err = yaml.Unmarshal(buf[:n], &routerConfigEntry)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Received from client: \n", routerConfigEntry)
		connection.WriteToUDP([]byte("Hello from server\n"), addr)
	}

}
