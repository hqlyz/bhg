package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket/pcap"
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}
	for _, device := range devices {
		fmt.Println(device.Name)
		for _, add := range device.Addresses {
			fmt.Printf("IP: %s\n", add.IP)
			fmt.Printf("Netmask: %s\n", add.Netmask)
		}
	}
}
