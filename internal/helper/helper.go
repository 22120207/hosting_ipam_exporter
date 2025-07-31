package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

type ipType string

var (
	normalIP   ipType = "Normal IP"
	outgoingIP ipType = "Outgoing IP"
)

type IPv4 struct {
	Value string `json:"value"`
	Type  ipType `json:"type"`
}

type Request struct {
	Hostname  string `json:"hostname"`
	HostIP    []IPv4 `json:"ipv4"`
	AuthenKey string `json:"auth"`
}

func SendToWebhook(ipList []string) error {
	var req Request
	req.Hostname = GetHostname()
	req.AuthenKey = "Ts3GkAzpAx1xG7Q"

	for index, ip := range ipList {
		ipv4 := IPv4{
			Value: ip,
			Type:  normalIP,
		}

		if index == 0 {
			ipv4.Type = outgoingIP
		}

		req.HostIP = append(req.HostIP, ipv4)
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	jsonStr := string(jsonData)

	cmd := `curl \
	-X POST \
	-d '` + jsonStr + `' \
	"http://14.225.204.41:5555/v1/hosting/ipam"`

	log.Printf("Curl command: %s", cmd)

	output, err := RunCommand(cmd)
	if err != nil {
		return fmt.Errorf("error %s in send to webhook with output %s", err, output)
	} else {
		log.Printf("Webhook Output: %s", output)
	}

	return nil
}

func RunCommand(cmd string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	output, err := exec.CommandContext(ctx, "bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error in run command %s: %v", cmd, err)
	}

	return string(output), nil
}

func GetHostname() string {
	cmd := `cat /etc/hostname`

	output, err := RunCommand(cmd)
	if err != nil {
		log.Println(err)
	}

	return strings.TrimSpace(string(output))
}

func IsPublicIPv4(ipStr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address format: %s", ipStr)
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return false, nil // Not an IPv4 address
	}

	// Check if it's a loopback address (127.0.0.0/8)
	if ip4.IsLoopback() {
		return false, nil
	}

	// Check if it's a link-local address (169.254.0.0/16)
	if ip4.IsLinkLocalUnicast() {
		return false, nil
	}

	// Check if it's an unspecified address (0.0.0.0)
	if ip4.IsUnspecified() {
		return false, nil
	}

	// Check if it's a multicast address (224.0.0.0/4)
	if ip4.IsMulticast() {
		return false, nil
	}

	// Check if it's a private address (RFC 1918)
	if ip4.IsPrivate() {
		return false, nil
	}

	// If it is public IP address
	if ip4.IsGlobalUnicast() {
		return true, nil
	}

	return false, nil
}
