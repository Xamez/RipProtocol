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
				IpAddress: applyMask(value.Ip, parseToMask(value.Mask)),
				SubMask:   parseToMask(value.Mask),
				Interface: value.Ip,
				Metric:    1,
			})
		}

	}

	return routerEntry
}

func parseToMask(m string) string {
	mask, _ := strconv.Atoi(m)
	var result [4]byte
	for i := 0; i < 4; i++ {
		if mask >= 8 {
			result[i] = 255
			mask -= 8
		} else {
			result[i] = byte(255 << uint(8-mask))
			mask = 0
		}
	}
	return fmt.Sprintf("%d.%d.%d.%d", result[0], result[1], result[2], result[3])
}

func applyMask(ip string, mask string) string {
	ipBytes := parseIpToBytes(ip)
	maskBytes := parseIpToBytes(mask)
	for i := 0; i < 4; i++ {
		ipBytes[i] = ipBytes[i] & maskBytes[i]
	}
	return fmt.Sprintf("%d.%d.%d.%d", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
}
