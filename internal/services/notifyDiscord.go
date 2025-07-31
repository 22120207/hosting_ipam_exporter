package services

import (
	"encoding/json"
	"hosting_ipam_exporter/internal/models"
	"log"
	"os/exec"
	"strings"
)

func NotifyDiscord(message string, color int) error {

	// Parse ASCII character
	message = strings.ReplaceAll(message, "%0A", "\n")

	title := ":loudspeaker: **THÔNG BÁO LỖI GỬI IP HOSTING TỚI WEBHOOK** :loudspeaker:"

	description := message

	// Build embed
	embed := models.DiscordEmbed{
		Title:       title,
		Description: description,
		Color:       color,
	}

	// Build the payload
	payload := models.DiscordPayload{
		Embeds: []models.DiscordEmbed{embed},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal Discord payload:", err)
		return err
	}

	jsonStr := string(jsonData)

	// Construct curl command using string concatenation only
	baseCommand := `curl \
	-X POST \
	-H "Content-Type: application/json" \
	-s --connect-timeout 10 \
	-d '` + jsonStr + `' \
	"https://discord.com/api/webhooks/1384379242753818724/_LNnAZAOL55chbhrLj6lKwYxHzEUuYll_aD8pKzSFCpHpeVUIf3ypTEoPkDzpJ1oYYtM?thread_id=1394903623473037382"`

	log.Println("Discord curl command:", baseCommand)

	output, err := exec.Command("bash", "-c", baseCommand).CombinedOutput()
	if err != nil {
		log.Println("Error sending to Discord:", string(output))
		return err
	}

	return nil
}
