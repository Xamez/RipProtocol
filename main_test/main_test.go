package main_test

import (
	"ProtocolRIP/rip"
	"testing"
)

func TestServer(t *testing.T) {
	go rip.UdpServer("localhost", 521, []rip.RouterEntry{})
	ripPacket := rip.CreateRipPacket("")
	dataToSend, err := rip.MarshalRipPacket(ripPacket)
	if err != nil {
		t.Error(err)
	}
	rip.UdpClient("localhost", 521, dataToSend)
	// Not a real test, just to see the output
}

func TestMarshallUnMarshall(t *testing.T) {
	routingTable := rip.ReadConfig("../config/routeur-r1.yaml")
	ripPacket := rip.RipPacket{Command: 1, Version: 2, Unused: [2]byte{0, 0}, RoutingTable: routingTable}
	dataToSend, err := rip.MarshalRipPacket(ripPacket)
	if err != nil {
		t.Error(err)
	}
	data, err := rip.UnmarshalRipPacket(dataToSend)
	if err != nil {
		t.Error(err)
	}
	if data.String() != ripPacket.String() {
		t.Errorf("Expected %s but got %s", ripPacket.String(), data.String())
	}
}

func TestSendRoutingTableToAnotherRouter(t *testing.T) {
	go rip.UdpServer("localhost", 521, rip.ReadConfig("../config/routeur-r1.yaml"))
	ripPacket := rip.CreateRipPacket("../config/routeur-r2.yaml")
	dataToSend, err := rip.MarshalRipPacket(ripPacket)
	if err != nil {
		t.Error(err)
	}
	rip.UdpClient("localhost", 521, dataToSend)
	// Not a real test, just to see the output
}

func TestSendR2ToR1(t *testing.T) {
	router1 := rip.ReadConfig("../config/routeur-r1.yaml")
	router2 := rip.ReadConfig("../config/routeur-r2.yaml")
	router1 = rip.MergeRoutingTable(router1, router2)
	expectedTable := []rip.RouterEntry{
		{IpAddress: "192.168.1.0", SubMask: "255.255.255.0", HasNextHop: false, Interface: "192.168.1.254", Metric: 1},
		{IpAddress: "10.1.1.0", SubMask: "255.255.255.252", HasNextHop: false, Interface: "10.1.1.1", Metric: 1},
		{IpAddress: "10.1.2.0", SubMask: "255.255.255.252", NextHop: "10.1.1.2", HasNextHop: true, Interface: "10.1.1.1", Metric: 2},
		{IpAddress: "10.1.4.0", SubMask: "255.255.255.252", NextHop: "10.1.1.2", HasNextHop: true, Interface: "10.1.1.1", Metric: 2},
		{IpAddress: "10.1.3.0", SubMask: "255.255.255.252", NextHop: "10.1.1.2", HasNextHop: true, Interface: "10.1.1.1", Metric: 2},
	}
	checkRoutingTable(t, router1, expectedTable)
}

func TestMergeEverything(t *testing.T) {
	router1 := rip.ReadConfig("../config/routeur-r1.yaml")
	router2 := rip.ReadConfig("../config/routeur-r2.yaml")
	router4 := rip.ReadConfig("../config/routeur-r4.yaml")
	router5 := rip.ReadConfig("../config/routeur-r5.yaml")
	router6 := rip.ReadConfig("../config/routeur-r6.yaml")
	router5 = rip.MergeRoutingTable(router5, router6)
	router4 = rip.MergeRoutingTable(router4, router5)
	router2 = rip.MergeRoutingTable(router2, router4)
	router2 = rip.MergeRoutingTable(router2, router5)
	router1 = rip.MergeRoutingTable(router1, router2)
	expectedTable := []rip.RouterEntry{
		{IpAddress: "10.1.1.0", SubMask: "255.255.255.252", HasNextHop: false, Interface: "10.1.1.1", Metric: 1},
		{IpAddress: "192.168.1.0", SubMask: "255.255.255.0", HasNextHop: false, Interface: "192.168.1.254", Metric: 1},
		{IpAddress: "10.1.2.0", SubMask: "255.255.255.252", NextHop: "10.1.1.2", HasNextHop: true, Interface: "10.1.1.1", Metric: 2},
		{IpAddress: "10.1.3.0", SubMask: "255.255.255.252", NextHop: "10.1.1.2", HasNextHop: true, Interface: "10.1.1.1", Metric: 2},
		{IpAddress: "10.1.4.0", SubMask: "255.255.255.252", NextHop: "10.1.1.2", HasNextHop: true, Interface: "10.1.1.1", Metric: 2},
		{IpAddress: "10.1.5.0", SubMask: "255.255.255.252", NextHop: "10.1.1.2", HasNextHop: true, Interface: "10.1.1.1", Metric: 3},
		{IpAddress: "10.1.6.0", SubMask: "255.255.255.252", NextHop: "10.1.1.2", HasNextHop: true, Interface: "10.1.1.1", Metric: 3},
		{IpAddress: "10.1.7.0", SubMask: "255.255.255.252", NextHop: "10.1.1.2", HasNextHop: true, Interface: "10.1.1.1", Metric: 3},
		{IpAddress: "172.16.180.0", SubMask: "255.255.255.0", NextHop: "10.1.1.2", HasNextHop: true, Interface: "10.1.1.1", Metric: 4},
	}
	checkRoutingTable(t, router1, expectedTable)

}

func areRoutingTableEqual(routingTable1 []rip.RouterEntry, routingTable2 []rip.RouterEntry) bool {
	if len(routingTable1) != len(routingTable2) {
		return false
	}
	for _, entry1 := range routingTable1 {
		found := false
		for _, entry2 := range routingTable2 {
			if entry1.String() == entry2.String() {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func checkRoutingTable(t *testing.T, routingTable []rip.RouterEntry, expectedTable []rip.RouterEntry) {
	if !areRoutingTableEqual(routingTable, expectedTable) {
		t.Error("Table mismatch")
		println("Expected")
		rip.PrintRoutingTable(expectedTable)
		println("But got")
		rip.PrintRoutingTable(routingTable)
	}
}
