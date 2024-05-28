package rip

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type RouterEntry struct {
	IpAddress  [4]byte
	SubMask    [4]byte
	NextHop    [4]byte
	Interface  [4]byte
	Metric     uint32
	HasNextHop bool
}

func (entry RouterEntry) String() string {
	nextHop := "nil"
	if entry.HasNextHop {
		nextHop = fmt.Sprintf("%d.%d.%d.%d", entry.NextHop[0], entry.NextHop[1], entry.NextHop[2], entry.NextHop[3])
	}
	return fmt.Sprintf("%-15s %-15s %-15s %-15s %-5d",
		fmt.Sprintf("%d.%d.%d.%d", entry.IpAddress[0], entry.IpAddress[1], entry.IpAddress[2], entry.IpAddress[3]),
		fmt.Sprintf("%d.%d.%d.%d", entry.SubMask[0], entry.SubMask[1], entry.SubMask[2], entry.SubMask[3]),
		nextHop,
		fmt.Sprintf("%d.%d.%d.%d", entry.Interface[0], entry.Interface[1], entry.Interface[2], entry.Interface[3]),
		entry.Metric)
}

type RipPacket struct {
	Command      byte
	Version      byte
	Unused       [2]byte
	RoutingTable []RouterEntry
}

func (packet RipPacket) String() string {
	str := fmt.Sprintf("Command: %d\n", packet.Command)
	str += fmt.Sprintf("Version: %d\n", packet.Version)
	str += fmt.Sprintf("Unused: %d\n", binary.BigEndian.Uint16(packet.Unused[:]))
	str += "Ip Address      Mask            Next Hop        Interface       Metric\n"
	for _, entry := range packet.RoutingTable {
		str += fmt.Sprintf("%s\n", entry.String())
	}
	return str
}

func CreateRipPacket(configFile string) RipPacket {
	var s []RouterEntry
	if configFile != "" {
		s = ReadConfig(configFile)
	}
	return RipPacket{Command: 1, Version: 2, Unused: [2]byte{0, 0}, RoutingTable: s}
}

func CreateRipPacketFromRoutingTable(routingTable []RouterEntry) RipPacket {
	return RipPacket{Command: 1, Version: 2, Unused: [2]byte{0, 0}, RoutingTable: routingTable}
}

func MarshalRipPacket(packet RipPacket) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, packet.Command); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, packet.Version); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, packet.Unused); err != nil {
		return nil, err
	}

	for _, entry := range packet.RoutingTable {
		if err := binary.Write(buf, binary.BigEndian, entry); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func UnmarshalRipPacket(data []byte) (RipPacket, error) {
	packet := RipPacket{}
	buf := bytes.NewReader(data)

	if err := binary.Read(buf, binary.BigEndian, &packet.Command); err != nil {
		return packet, err
	}
	if err := binary.Read(buf, binary.BigEndian, &packet.Version); err != nil {
		return packet, err
	}

	if err := binary.Read(buf, binary.BigEndian, &packet.Unused); err != nil {
		return packet, err
	}

	for buf.Len() > 0 {
		entry := RouterEntry{}
		if err := binary.Read(buf, binary.BigEndian, &entry); err != nil {
			return packet, err
		}
		packet.RoutingTable = append(packet.RoutingTable, entry)
	}

	return packet, nil
}
