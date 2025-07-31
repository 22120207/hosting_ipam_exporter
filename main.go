package main

import (
	"hosting_ipam_exporter/internal/helper"
	"log"
	"math/rand"
	"strings"
	"time"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 45 minutes
	minDuration := 45 * time.Minute

	// 1 hour
	maxDuration := 1 * time.Hour

	for {
		randomDuration := minDuration + time.Duration(r.Int63n(int64(maxDuration-minDuration+1)))

		// Send IP to Webhook Server
		RunTask()

		time.Sleep(randomDuration)
	}
}

func RunTask() {
	cmd := `/usr/sbin/ip -brief address show | awk '{ if ($1 == "bond0") { for (i = 3; i <= NF; i++) print $i } }'`

	output, err := helper.RunCommand(cmd)
	if err != nil {
		log.Println(err)
	}

	ipList := strings.Split(output, "\n")

	filteredIpList := []string{}
	for _, ip := range ipList {
		fields := strings.Split(ip, "/")

		if len(fields) > 0 {
			ip = fields[0]
		} else {
			continue
		}

		ip = strings.TrimSpace(ip)

		if ip == "" {
			continue
		}

		if isPublicIP, err := helper.IsPublicIPv4(ip); isPublicIP && err == nil {
			filteredIpList = append(filteredIpList, ip)
		} else if err != nil {
			log.Printf("Error in check IPv4: %v", err)
		}
	}

	err = helper.SendToWebhook(filteredIpList)
	if err != nil {
		log.Println(err)
	}
}
