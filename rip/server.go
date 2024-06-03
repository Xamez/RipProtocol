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
	outputTable := make([]RouterEntry, len(routingTable))
	copy(outputTable, routingTable)

	for _, route := range newRoutingTable {
		route.Metric++
		outputTable = append(outputTable, route)
	}

	ipDestMap := make(map[string][]RouterEntry)
	for _, route := range outputTable {
		ipDestMap[route.IpAddress] = append(ipDestMap[route.IpAddress], route)
	}

	ipTable1 := make(map[string]bool)
	for _, route := range routingTable {
		ipTable1[route.Interface] = true
	}

	ipDestMapFilter := make(map[string][]RouterEntry)

	for ipDest, routes := range ipDestMap {
		if len(routes) >= 2 {
			ipDestMapFilter[ipDest] = routes
		}
	}

	if len(ipDestMapFilter) >= 2 {
		minSumMetric := -1
		var minSumMetricIP string
		for ipDest, routes := range ipDestMapFilter {
			sumMetric := 0
			for _, route := range routes {
				sumMetric += int(route.Metric)
			}
			if minSumMetric == -1 || sumMetric < minSumMetric {
				minSumMetric = sumMetric
				minSumMetricIP = ipDest
			}
		}

		filteredIPDestMapFilter := make(map[string][]RouterEntry)

		filteredIPDestMapFilter[minSumMetricIP] = ipDestMapFilter[minSumMetricIP]
		ipDestMapFilter = filteredIPDestMapFilter
	}

	for _, routes := range ipDestMapFilter {
		if len(routes) > 1 {
			for _, route := range routes {
				for i := range outputTable {
					if ipTable1[route.Interface] && !ipTable1[outputTable[i].Interface] {
						outputTable[i].Interface = route.Interface
						outputTable[i].NextHop = ""
						outputTable[i].HasNextHop = false
					} else if !ipTable1[route.Interface] && !outputTable[i].HasNextHop {
						outputTable[i].NextHop = route.Interface
						outputTable[i].HasNextHop = true
					}
				}
			}
		}
	}

	for i, route := range outputTable {
		for _, existingRoute := range routingTable {
			if route.IpAddress == existingRoute.IpAddress {
				outputTable[i].NextHop = existingRoute.NextHop
				outputTable[i].HasNextHop = existingRoute.HasNextHop
			}
		}
	}

	ipDestMap = make(map[string][]RouterEntry)
	outputTable = removeDuplicateDestinations(outputTable)

	return outputTable
}

func removeDuplicateDestinations(routingTable []RouterEntry) []RouterEntry {
	ipDestMap := make(map[string]RouterEntry)
	for _, route := range routingTable {
		if existingRoute, ok := ipDestMap[route.IpAddress]; ok {
			if existingRoute.Metric > route.Metric {
				ipDestMap[route.IpAddress] = route
			}
		} else {
			ipDestMap[route.IpAddress] = route
		}
	}

	outputTable := make([]RouterEntry, 0, len(ipDestMap))
	for _, route := range ipDestMap {
		outputTable = append(outputTable, route)
	}

	return outputTable
}
