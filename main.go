package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	binanceAPI   = "https://api.binance.com/api/v3/ticker/24hr?symbol=BTCUSDT"
	maxRetries   = 3
	retryDelay   = 5 * time.Second
)

// Update struct for Binance response
type BinanceResponse struct {
	LastPrice      string `json:"lastPrice"`
	PriceChange    string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
}

func fetchBitcoinPrice() (price float64, change24h float64, err error) {
	resp, err := http.Get(binanceAPI)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var data BinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, err
	}

	price, err = strconv.ParseFloat(data.LastPrice, 64)
	if err != nil {
		return 0, 0, err
	}

	change24h, err = strconv.ParseFloat(data.PriceChangePercent, 64)
	if err != nil {
		return 0, 0, err
	}

	log.Printf("Bitcoin price: $%.2f", price)

	return price, change24h, nil
}

func formatPriceMessage(priceUSD float64, change24h float64) string {
	var symbol string
	if change24h < 0 {
		symbol = ""
		// Add salt for negative changes
		saltCount := int(-change24h / 2.5)
		if saltCount > 0 {
			if saltCount > 3 {
				saltCount = 3
			}
			symbol += strings.Repeat("🧂", saltCount)
		}
	} else {
		symbol = ""
		// Add rockets for positive changes
		rockets := int(change24h / 2.5)
		if rockets > 0 {
			if rockets > 3 {
				rockets = 3
			}
			symbol += strings.Repeat("🚀", rockets)
		}
	}
	return fmt.Sprintf("$%.2f (%+.2f%%) %s", priceUSD, change24h, symbol)
}

func main() {
	
	// Get token and group chat ID from environment variables
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
	}

	groupChatIDStr := os.Getenv("TELEGRAM_GROUP_CHAT_ID")
	if groupChatIDStr == "" {
		log.Fatal("TELEGRAM_GROUP_CHAT_ID environment variable is not set")
	}

	groupChatID, err := strconv.ParseInt(groupChatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid TELEGRAM_GROUP_CHAT_ID value: %v", err)
	}

	b, err := bot.New(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	price, change24h, err := fetchBitcoinPrice()
	if err != nil {
		log.Printf("Failed to get initial price: %v", err)
		return
	}

	initialMessage := formatPriceMessage(price, change24h)
	log.Printf("Sending initial price: %s", initialMessage)

	var lastMessageID int
	message := &bot.SendMessageParams{
		ChatID:    groupChatID,
		Text:      initialMessage,
		ParseMode: models.ParseModeHTML,
	}
	resp, err := b.SendMessage(context.Background(), message)
	if err != nil {
		log.Printf("Failed to send initial price message: %v", err)
		return
	}

	lastMessageID = resp.ID
	log.Printf("Message ID for updates: %d", lastMessageID)

	// Keep track of last price
	lastPrice := price
	lastChange := change24h

	// Update prices periodically
	for range ticker.C {
		if err := updatePriceMessage(b, groupChatID, &lastMessageID, &lastPrice, &lastChange); err != nil {
			log.Printf("Failed to update price message: %v", err)
		}
	}
}

func updatePriceMessage(b *bot.Bot, groupChatID int64, lastMessageID *int, lastPrice, lastChange *float64) error {
	priceUSD, change24h, err := fetchBitcoinPrice()
	if err != nil {
		return fmt.Errorf("failed to get Bitcoin price: %v", err)
	}

	// Skip update if price and change haven't changed
	if priceUSD == *lastPrice && change24h == *lastChange {
		log.Printf("Price unchanged, skipping update")
		return nil
	}

	message := formatPriceMessage(priceUSD, change24h)
	log.Printf("Price changed. Sending new message: %s", message)

	// Send new message
	resp, err := b.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID:    groupChatID,
		Text:      message,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		return fmt.Errorf("failed to send new message: %v", err)
	}

	// Delete previous message
	if *lastMessageID != 0 {
		_, err = b.DeleteMessage(context.Background(), &bot.DeleteMessageParams{
			ChatID:    groupChatID,
			MessageID: *lastMessageID,
		})
		if err != nil {
			log.Printf("Failed to delete previous message: %v", err)
		}
	}

	// Update message ID and last known values
	*lastMessageID = resp.ID
	*lastPrice = priceUSD
	*lastChange = change24h
	return nil
}
