package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type RouterConfigEntry struct {
	Device string `yaml:"device"`
	Ip     string `yaml:"ip"`
	Mask   string `yaml:"mask"`
}

type RouterEntry struct {
	AddressFamilyIdentifier uint16
	RouteTag                uint16
	IpAddress               [4]byte
	SubMask                 [4]byte
	NextHop                 [4]byte
	Metric                  uint32
}

type RipPacket struct {
	Command     byte
	Version     byte
	Unused      [2]byte
	RouterEntry []RouterEntry
}

func (packet RipPacket) String() string {
	str := fmt.Sprintf("Command: %d\n", packet.Command)
	str += fmt.Sprintf("Version: %d\n", packet.Version)
	str += fmt.Sprintf("Unused: %d\n", binary.BigEndian.Uint16(packet.Unused[:]))
	for _, entry := range packet.RouterEntry {
		str += fmt.Sprintf("AddressFamilyIdentifier: %d\n", entry.AddressFamilyIdentifier)
		str += fmt.Sprintf("RouteTag: %d\n", entry.RouteTag)
		str += fmt.Sprintf("IpAddress: %d.%d.%d.%d\n", entry.IpAddress[0], entry.IpAddress[1], entry.IpAddress[2], entry.IpAddress[3])
		str += fmt.Sprintf("SubMask: %d.%d.%d.%d\n", entry.SubMask[0], entry.SubMask[1], entry.SubMask[2], entry.SubMask[3])
		str += fmt.Sprintf("NextHop: %d.%d.%d.%d\n", entry.NextHop[0], entry.NextHop[1], entry.NextHop[2], entry.NextHop[3])
		str += fmt.Sprintf("Metric: %d\n", entry.Metric)
	}
	return str
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

	for _, entry := range packet.RouterEntry {
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
		packet.RouterEntry = append(packet.RouterEntry, entry)
	}

	return packet, nil
}
