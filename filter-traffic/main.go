package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)


var (
	iface    = "\\Device\\NPF_{2AF87C94-4D20-4CD4-B9F2-321107C488D9}"
	snaplen  = int32(1600)
	promisc  = false
	timeout  = pcap.BlockForever
	filter   = "tcp and port 3389"
	devFound = false
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}
	for _, device := range devices {
		if device.Name == iface {
			devFound = true
			break
		}
	}
	if !devFound {
		log.Panicf("Device name %s does not exist\n", iface)
	}

	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		log.Panicln(err)
	}
	defer handle.Close()

	if err = handle.SetBPFFilter(filter); err != nil {
		log.Panicln(err)
	}
	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		fmt.Println(packet)
	}
}
