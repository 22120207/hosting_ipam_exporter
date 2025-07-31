package models

// DiscordEmbed represents a Discord embed message
type DiscordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
}

// A struct of Discord Payload Data
type DiscordPayload struct {
	Content string         `json:"content"`
	Embeds  []DiscordEmbed `json:"embeds"`
}

// Embed colors
const (
	ColorDefault = 0
	ColorRed     = 15158332
	ColorGreen   = 3066993
	ColorYellow  = 15844367
	ColorBlue    = 3447003
)
