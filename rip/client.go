package rip

import (
	"bufio"
	"fmt"
	"net"
)

func UdpClient(host string, port int, dataToSend ...[]byte) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	CheckForError(err)

	connection, err := net.DialUDP("udp", nil, udpAddr)
	CheckForError(err)

	// Send message to server
	for _, data := range dataToSend {
		_, err = connection.Write(data)
		CheckForError(err)

		// Receive message from server
		data, err := bufio.NewReader(connection).ReadString(TerminatedChar)
		CheckForError(err)
		data = data[:len(data)-1]

		fmt.Println("Received from server: ", data)
	}

}
