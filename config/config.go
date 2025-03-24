package config

import (
	"os"
	"strings"
)

// Config holds all configuration values
type Config struct {
	TelegramBotToken string
	TelegramChatID   string
	GeminiAPIKey     string
	StockList        string
}

// Default stock lists by market cap
const (
	// Large Cap Stocks (Market Cap > ₹50,000 Cr)
	DefaultLargeCapStocks = "RELIANCE.NS,TCS.NS,HDFCBANK.NS,INFY.NS,ICICIBANK.NS"

	// Mid Cap Stocks (Market Cap ₹10,000-50,000 Cr)
	DefaultMidCapStocks = "TATAMOTORS.NS,ADANIENT.NS,BAJFINANCE.NS,TITAN.NS,MARICO.NS"

	// Small Cap Stocks (Market Cap < ₹10,000 Cr)
	DefaultSmallCapStocks = "JUBLFOOD.NS,FORTIS.NS,KALYANKJIL.NS,SUPREMEIND.NS,VBL.NS"

	// Combined default stock list
	DefaultStockList = DefaultLargeCapStocks + "," + DefaultMidCapStocks + "," + DefaultSmallCapStocks
)

// GetConfig retrieves configuration values from environment variables
func GetConfig() *Config {
	// Ensure Telegram credentials are properly formatted
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	// Remove any whitespace from credentials
	if botToken != "" {
		botToken = strings.TrimSpace(botToken)
	}
	if chatID != "" {
		chatID = strings.TrimSpace(chatID)
	}

	return &Config{
		TelegramBotToken: botToken,
		TelegramChatID:   chatID,
		GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
		StockList:        getStockList(),
	}
}

// getStockList returns the stock list from environment variable or default list
func getStockList() string {
	if stockList := os.Getenv("STOCK_LIST"); stockList != "" {
		return stockList
	}
	return DefaultStockList
}
