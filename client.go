package main

import (
	"bufio"
	"fmt"
	"net"
)

func UdpClient(host string, port int, dataToSend RipPacket) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	CheckForError(err)

	connection, err := net.DialUDP("udp", nil, udpAddr)
	CheckForError(err)

	// Send message to server
	serializedPacket, err := MarshalRipPacket(dataToSend)
	_, err = connection.Write(serializedPacket)
	CheckForError(err)

	// Receive message from server
	data, err := bufio.NewReader(connection).ReadString('\n')
	CheckForError(err)

	fmt.Println("Received from server: ", data)

}
