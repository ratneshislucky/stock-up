package config

import (
	"os"
	"strings"
)

// Config holds all configuration values
type Config struct {
	TelegramBotToken string
	TelegramChatIDs  []string
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
	chatIDsStr := os.Getenv("TELEGRAM_CHAT_IDS")

	// Remove any whitespace from credentials
	if botToken != "" {
		botToken = strings.TrimSpace(botToken)
	}

	// Split and clean chat IDs
	var chatIDs []string
	if chatIDsStr != "" {
		chatIDsStr = strings.TrimSpace(chatIDsStr)
		rawChatIDs := strings.Split(chatIDsStr, ",")
		for _, id := range rawChatIDs {
			if trimmedID := strings.TrimSpace(id); trimmedID != "" {
				chatIDs = append(chatIDs, trimmedID)
			}
		}
	}

	return &Config{
		TelegramBotToken: botToken,
		TelegramChatIDs:  chatIDs,
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

// getEnvVar retrieves an environment variable and ensures it's not empty
func getEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Environment variable " + key + " is not set")
	}
	return value
}
