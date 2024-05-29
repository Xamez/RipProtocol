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

func TestSendR2ToR1(t *testing.T) {
	go rip.UdpServer("localhost", 521, rip.ReadConfig("../config/routeur-r1.yaml"))
	ripPacket := rip.CreateRipPacket("../config/routeur-r2.yaml")
	dataToSend, err := rip.MarshalRipPacket(ripPacket)
	if err != nil {
		t.Error(err)
	}
	rip.UdpClient("localhost", 521, dataToSend)
}

func TestMergeEverything(t *testing.T) {
	router1 := rip.ReadConfig("../config/routeur-r1.yaml")
	router2 := rip.ReadConfig("../config/routeur-r2.yaml")
	//router3 := rip.ReadConfig("../config/routeur-r3.yaml")
	router4 := rip.ReadConfig("../config/routeur-r4.yaml")
	router5 := rip.ReadConfig("../config/routeur-r5.yaml")
	router6 := rip.ReadConfig("../config/routeur-r6.yaml")
	router5 = rip.MergeRoutingTable(router5, router6)
	router4 = rip.MergeRoutingTable(router4, router5)
	router2 = rip.MergeRoutingTable(router2, router4)
	router2 = rip.MergeRoutingTable(router2, router5)
	router1 = rip.MergeRoutingTable(router1, router2)
	rip.PrintRoutingTable(router1)
	rip.PrintRoutingTable(router2)
}
