package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

func ReadConfig() {
	file, err := os.Open("config/routeur-r1.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var interfaces []map[string]RouterConfigEntry

	// Unmarshal the YAML
	err = yaml.Unmarshal(data, &interfaces)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Iterate over the interfaces and print them
	for _, iface := range interfaces {
		for _, intf := range iface {
			fmt.Printf("Device: %s, IP: %s, Mask: %s\n", intf.Device, stringToBytes(intf.Ip), stringToBytes(intf.Mask))
		}
	}
}

func stringToBytes(s string) [4]byte {
	var bytes [4]byte
	fmt.Sscanf(s, "%d.%d.%d.%d", &bytes[0], &bytes[1], &bytes[2], &bytes[3])
	return bytes
}
