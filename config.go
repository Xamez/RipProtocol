package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

func ReadConfigAsBytes(configFile string) ([]byte, error) {
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	return io.ReadAll(file)

}

func ReadConfig(configFile string) []map[string]RouterConfigEntry {
	data, err := ReadConfigAsBytes(configFile)
	if err != nil {
		fmt.Println(err)
	}

	var routerConfigEntry []map[string]RouterConfigEntry

	err = yaml.Unmarshal(data, &routerConfigEntry)
	if err != nil {
		fmt.Println(err)
	}

	return routerConfigEntry

	//var routerEntry []RouterEntry
	//
	//for _, entry := range routerConfigEntry {
	//	for _, value := range entry {
	//		routerEntry = append(routerEntry, RouterEntry{
	//			AddressFamilyIdentifier: 1, // IPv4
	//			RouteTag:                0, // Ignored
	//			IpAddress:               stringToBytes(value.Ip),
	//			SubMask:                 stringToBytes(value.Mask),
	//			NextHop:                 [4]byte{0, 0, 0, 0},
	//			Metric:                  1,
	//		})
	//	}
	//
	//}
	//
	//return routerEntry
}

func stringToBytes(s string) [4]byte {
	var bytes [4]byte
	fmt.Sscanf(s, "%d.%d.%d.%d", &bytes[0], &bytes[1], &bytes[2], &bytes[3])
	return bytes
}
