package rip

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type RouterEntry struct {
	IpAddress  string
	SubMask    string
	NextHop    string
	Interface  string
	Metric     uint32
	HasNextHop bool
}

func (entry RouterEntry) String() string {
	nextHop := ""
	if entry.HasNextHop {
		nextHop = entry.NextHop
	}
	return fmt.Sprintf("%-15s %-15s %-15s %-15s %-5d",
		entry.IpAddress,
		entry.SubMask,
		nextHop,
		entry.Interface,
		entry.Metric)
}

func PrintRoutingTable(routingTable []RouterEntry) {
	fmt.Println("Ip Address      Mask            Next Hop        Interface       Metric")
	for i := 0; i < len(routingTable); i++ {
		for j := i + 1; j < len(routingTable); j++ {
			if routingTable[i].Metric > routingTable[j].Metric || (routingTable[i].Metric == routingTable[j].Metric && routingTable[i].IpAddress > routingTable[j].IpAddress) {
				routingTable[i], routingTable[j] = routingTable[j], routingTable[i]
			}
		}
	}

	for _, entry := range routingTable {
		fmt.Println(entry.String())
	}
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
		ipBytes := parseIpToBytes(entry.IpAddress)
		if err := binary.Write(buf, binary.BigEndian, ipBytes); err != nil {
			return nil, err
		}
		subMaskBytes := parseIpToBytes(entry.SubMask)
		if err := binary.Write(buf, binary.BigEndian, subMaskBytes); err != nil {
			return nil, err
		}
		nextHopBytes := parseIpToBytes(entry.NextHop)
		if err := binary.Write(buf, binary.BigEndian, nextHopBytes); err != nil {
			return nil, err
		}
		interfaceBytes := parseIpToBytes(entry.Interface)
		if err := binary.Write(buf, binary.BigEndian, interfaceBytes); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, entry.Metric); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, entry.HasNextHop); err != nil {
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
		var ipBytes [4]byte
		if err := binary.Read(buf, binary.BigEndian, &ipBytes); err != nil {
			return packet, err
		}
		entry.IpAddress = parseBytesToIp(ipBytes)
		var subMaskBytes [4]byte
		if err := binary.Read(buf, binary.BigEndian, &subMaskBytes); err != nil {
			return packet, err
		}
		entry.SubMask = parseBytesToIp(subMaskBytes)
		var nextHopBytes [4]byte
		if err := binary.Read(buf, binary.BigEndian, &nextHopBytes); err != nil {
			return packet, err
		}
		entry.NextHop = parseBytesToIp(nextHopBytes)
		var interfaceBytes [4]byte
		if err := binary.Read(buf, binary.BigEndian, &interfaceBytes); err != nil {
			return packet, err
		}
		entry.Interface = parseBytesToIp(interfaceBytes)
		if err := binary.Read(buf, binary.BigEndian, &entry.Metric); err != nil {
			return packet, err
		}
		if err := binary.Read(buf, binary.BigEndian, &entry.HasNextHop); err != nil {
			return packet, err
		}
		packet.RoutingTable = append(packet.RoutingTable, entry)
	}

	return packet, nil
}
