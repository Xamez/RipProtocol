package rip

import (
	"fmt"
	"net"
)

func UdpServer(host string, port int) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	CheckForError(err)

	connection, err := net.ListenUDP("udp", udpAddr)
	CheckForError(err)

	var routingTable []RouterEntry

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
	outputRoutingTable := make([]RouterEntry, len(routingTable))
	for i, entry := range newRoutingTable {
		outputRoutingTable = append(outputRoutingTable, entry)
		outputRoutingTable[i].Metric++
		if !contains(routingTable, entry) { // If the entry is not in the routing table
			routingTable = append(routingTable, entry)
		} else { // If the entry is in the routing table
			for i, oldEntry := range routingTable { // Update the entry if the new entry has a better metric
				if areEqual(entry, oldEntry) && entry.Metric < oldEntry.Metric {
					routingTable[i] = entry
				}
			}
		}
	}
	return routingTable
}

func areEqual(routingTable1 RouterEntry, routingTable2 RouterEntry) bool {
	return routingTable1.IpAddress == routingTable2.IpAddress &&
		routingTable1.NextHop == routingTable2.NextHop &&
		routingTable1.SubMask == routingTable2.SubMask
}

func contains(routingTable []RouterEntry, entry RouterEntry) bool {
	for _, e := range routingTable {
		if areEqual(e, entry) {
			return true
		}
	}
	return false
}
