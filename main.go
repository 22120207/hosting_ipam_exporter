package main

import (
	"hosting_ipam_exporter/internal/helper"
	"log"
	"strings"
)

func main() {
	cmd := `ip -brief address show | awk '{ if ($1 == "bond0") { for (i = 3; i <= NF; i++) print $i } }'`

	output, err := helper.RunCommand(cmd)
	if err != nil {
		log.Println(err)
	}

	ipList := strings.Split(output, "\n")

	filteredIpList := []string{}
	for _, ip := range ipList {
		if isPublicIP, err := helper.IsPublicIPv4(ip); isPublicIP && err == nil {
			filteredIpList = append(filteredIpList, ip)
		} else if err != nil {
			log.Printf("Error in check IPv4: %v", err)
		}
	}

	helper.SendToWebhook(filteredIpList)
}
