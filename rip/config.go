package rip

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strconv"
)

type RouterConfigEntry struct {
	Device string `yaml:"device"`
	Ip     string `yaml:"ip"`
	Mask   string `yaml:"mask"`
}

func (entry RouterConfigEntry) String() string {
	return fmt.Sprintf("Device: %s\nIp: %s\nMask: %s\n", entry.Device, entry.Ip, entry.Mask)
}

func ReadConfigAsBytes(configFile string) ([]byte, error) {
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	return io.ReadAll(file)

}

func ReadConfig(configFile string) []RouterEntry {
	data, err := ReadConfigAsBytes(configFile)
	if err != nil {
		fmt.Println(err)
	}

	var routerConfigEntry []map[string]RouterConfigEntry

	err = yaml.Unmarshal(data, &routerConfigEntry)
	if err != nil {
		fmt.Println(err)
	}

	var routerEntry []RouterEntry

	for _, entry := range routerConfigEntry {
		for _, value := range entry {
			routerEntry = append(routerEntry, RouterEntry{
				AddressFamilyIdentifier: 2, // IPv4
				RouteTag:                0, // Ignored
				IpAddress:               stringToBytes(value.Ip),
				SubMask:                 maskToBytes(value.Mask),
				NextHop:                 [4]byte{0, 0, 0, 0},
				Metric:                  1,
			})
		}

	}

	return routerEntry
}

func stringToBytes(s string) [4]byte {
	var bytes [4]byte
	fmt.Sscanf(s, "%d.%d.%d.%d", &bytes[0], &bytes[1], &bytes[2], &bytes[3])
	return bytes
}

func maskToBytes(s string) [4]byte {
	var bytes [4]byte
	mask, _ := strconv.Atoi(s)
	for i := 0; i < mask; i++ {
		bytes[i/8] |= 1 << uint(7-i%8)
	}
	return bytes
}
