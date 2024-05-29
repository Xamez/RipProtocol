package rip

import (
	"fmt"
	"net"
)

func UdpServer(host string, port int, defaultRouterConfig []RouterEntry) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	CheckForError(err)

	connection, err := net.ListenUDP("udp", udpAddr)
	CheckForError(err)

	routingTable := defaultRouterConfig

	for {
		var buf [1024]byte
		n, addr, err := connection.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println(err)
			continue
		}

		data, err := UnmarshalRipPacket(buf[:n])
		if err != nil {
			fmt.Println("Error unmarshalling RipPacket:", err)
			continue
		}

		routingTable = MergeRoutingTable(routingTable, data.RoutingTable)

		fmt.Println("Received from client: \n", data.String())

		ripPacket := CreateRipPacketFromRoutingTable(routingTable)
		dataToSend := []byte(ripPacket.String())
		dataToSend = append(dataToSend, byte(TerminatedChar))
		connection.WriteToUDP(dataToSend, addr)
	}

}

func MergeRoutingTable(routingTable []RouterEntry, newRoutingTable []RouterEntry) []RouterEntry {
	var addressInterface [4]byte
	var nextHop [4]byte
	existingRoutes := make(map[[8]byte]int)
	for i, existingRoute := range routingTable {
		existingRoutes[[8]byte(append(existingRoute.IpAddress[:], existingRoute.SubMask[:]...))] = i
	}
	for _, newRoute := range newRoutingTable {
		for i, existingRoute := range routingTable {
			if areAddressesEqual(existingRoute.IpAddress, existingRoute.SubMask, newRoute.IpAddress, newRoute.SubMask) {
				nextHop = newRoute.Interface
				addressInterface = existingRoute.Interface
				if newRoute.Metric < existingRoute.Metric {
					routingTable[i] = newRoute
				}
				break
			}
		}
	}
	for _, newRoute := range newRoutingTable {
		newRoute.NextHop = nextHop
		newRoute.Interface = addressInterface
		newRoute.Metric++
		newRoute.HasNextHop = true
		if _, ok := existingRoutes[[8]byte(append(newRoute.IpAddress[:], newRoute.SubMask[:]...))]; ok {
			continue
		}
		routingTable = append(routingTable, newRoute)
	}
	return routingTable
}

func areAddressesEqual(address [4]byte, mask [4]byte, address2 [4]byte, mask2 [4]byte) bool {
	addressWithMask := applyMask(address, mask)
	address2WithMask := applyMask(address2, mask2)
	for i := 0; i < 4; i++ {
		if addressWithMask[i] != address2WithMask[i] {
			return false
		}
	}
	return true
}
