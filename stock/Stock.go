package stock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"go-stock/config"
)

// Configurable list of Indian stocks (expandable via config)
var defaultStocks = []string{
	// Large Cap Stocks (Market Cap > ₹50,000 Cr)
	"RELIANCE.NS",  // Reliance Industries - Energy & Retail
	"TCS.NS",       // Tata Consultancy Services - IT
	"HDFCBANK.NS",  // HDFC Bank - Banking
	"INFY.NS",      // Infosys - IT
	"ICICIBANK.NS", // ICICI Bank - Banking

	// Mid Cap Stocks (Market Cap ₹10,000-50,000 Cr)
	"TATAMOTORS.NS", // Tata Motors - Auto
	"ADANIENT.NS",   // Adani Enterprises - Infrastructure
	"BAJFINANCE.NS", // Bajaj Finance - NBFC
	"TITAN.NS",      // Titan Company - Consumer Goods
	"MARICO.NS",     // Marico - FMCG

	// Small Cap Stocks (Market Cap < ₹10,000 Cr)
	"JUBLFOOD.NS",   // Jubilant FoodWorks - Food Services
	"FORTIS.NS",     // Fortis Healthcare - Healthcare
	"KALYANKJIL.NS", // Kalyan Jewellers - Retail
	"SUPREMEIND.NS", // Supreme Industries - Plastics
	"VBL.NS",        // Varun Beverages - Beverages
}

type StockData struct {
	Symbol        string  `json:"symbol"`
	Price         float64 `json:"regularMarketPrice"`
	PreviousClose float64 `json:"previousClose"`
	High          float64 `json:"regularMarketDayHigh"`
	Low           float64 `json:"regularMarketDayLow"`
	Volume        int64   `json:"regularMarketVolume"`
}

type StockMetrics struct {
	Symbol       string
	Price        float64
	PriceChange  float64
	DailyRange   float64
	Volatility   float64
	Volume       int64
	VolumeChange float64
	MA5          float64 // 5-day moving average
	MA20         float64 // 20-day moving average
	PriceVsMA5   float64 // Price vs 5-day MA
	PriceVsMA20  float64 // Price vs 20-day MA
	RSI          float64 // 14-day RSI
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Send Telegram notification with improved formatting
func sendStockTelegramNotification(message string) {
	cfg := config.GetConfig()
	if cfg.TelegramBotToken == "" || len(cfg.TelegramChatIDs) == 0 {
		fmt.Println("Warning: Telegram credentials not set")
		return
	}

	// Escape special characters in the message
	message = strings.ReplaceAll(message, "_", "\\_")
	message = strings.ReplaceAll(message, "*", "\\*")
	message = strings.ReplaceAll(message, "[", "\\[")
	message = strings.ReplaceAll(message, "]", "\\]")
	message = strings.ReplaceAll(message, "(", "\\(")
	message = strings.ReplaceAll(message, ")", "\\)")
	message = strings.ReplaceAll(message, "~", "\\~")
	message = strings.ReplaceAll(message, "`", "\\`")
	message = strings.ReplaceAll(message, ">", "\\>")
	message = strings.ReplaceAll(message, "#", "\\#")
	message = strings.ReplaceAll(message, "+", "\\+")
	message = strings.ReplaceAll(message, "-", "\\-")
	message = strings.ReplaceAll(message, "=", "\\=")
	message = strings.ReplaceAll(message, "|", "\\|")
	message = strings.ReplaceAll(message, "{", "\\{")
	message = strings.ReplaceAll(message, "}", "\\}")
	message = strings.ReplaceAll(message, ".", "\\.")
	message = strings.ReplaceAll(message, "!", "\\!")

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.TelegramBotToken)

	// Send to each chat ID
	for _, chatID := range cfg.TelegramChatIDs {
		chatID = strings.TrimSpace(chatID)
		if chatID == "" {
			continue
		}

		payload := map[string]interface{}{
			"chat_id":    chatID,
			"text":       message,
			"parse_mode": "MarkdownV2",
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("Error marshaling payload for chat %s: %v\n", chatID, err)
			continue
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Printf("Error sending Telegram notification to chat %s: %v\n", chatID, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("Notification sent successfully to chat %s!\n", chatID)
		} else {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Failed to send notification to chat %s. Status code: %d, Response: %s\n",
				chatID, resp.StatusCode, string(body))
		}
	}
}

// Fetch current stock data from Yahoo Finance
func fetchYahooFinanceData(symbol string) (StockData, error) {
	client := resty.New().
		SetRetryCount(3).                     // Retry on failure
		SetRetryWaitTime(2 * time.Second).    // Initial wait
		SetRetryMaxWaitTime(10 * time.Second) // Max wait

	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", symbol)
	resp, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0").
		Get(url)

	if err != nil {
		return StockData{}, fmt.Errorf("failed to fetch data for %s: %v", symbol, err)
	}

	var result struct {
		Chart struct {
			Result []struct {
				Meta struct {
					RegularMarketPrice   float64 `json:"regularMarketPrice"`
					PreviousClose        float64 `json:"previousClose"`
					RegularMarketDayHigh float64 `json:"regularMarketDayHigh"`
					RegularMarketDayLow  float64 `json:"regularMarketDayLow"`
					RegularMarketVolume  int64   `json:"regularMarketVolume"`
				} `json:"meta"`
			} `json:"result"`
		} `json:"chart"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return StockData{}, fmt.Errorf("failed to parse response for %s: %v", symbol, err)
	}

	if len(result.Chart.Result) == 0 {
		return StockData{}, fmt.Errorf("no data found for symbol %s", symbol)
	}

	meta := result.Chart.Result[0].Meta

	// Validate critical fields
	if meta.RegularMarketPrice <= 0 || meta.PreviousClose <= 0 {
		return StockData{}, fmt.Errorf("invalid price data for %s: price=%.2f, previousClose=%.2f", symbol, meta.RegularMarketPrice, meta.PreviousClose)
	}

	return StockData{
		Symbol:        symbol,
		Price:         meta.RegularMarketPrice,
		PreviousClose: meta.PreviousClose,
		High:          meta.RegularMarketDayHigh,
		Low:           meta.RegularMarketDayLow,
		Volume:        meta.RegularMarketVolume,
	}, nil
}

// Fetch historical data (30 days for better MA and RSI)
func fetchHistoricalData(symbol string) ([]StockData, error) {
	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -30) // Fetch 30 days to cover holidays

	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", symbol)
	resp, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0").
		SetQueryParams(map[string]string{
			"period1":  fmt.Sprintf("%d", startTime.Unix()),
			"period2":  fmt.Sprintf("%d", endTime.Unix()),
			"interval": "1d",
		}).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical data for %s: %v", symbol, err)
	}

	var result struct {
		Chart struct {
			Result []struct {
				Timestamp  []int64 `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Close  []float64 `json:"close"`
						High   []float64 `json:"high"`
						Low    []float64 `json:"low"`
						Volume []int64   `json:"volume"`
					} `json:"quote"`
				} `json:"indicators"`
			} `json:"result"`
		} `json:"chart"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse historical response for %s: %v", symbol, err)
	}

	if len(result.Chart.Result) == 0 || len(result.Chart.Result[0].Indicators.Quote) == 0 {
		return nil, fmt.Errorf("no historical data found for symbol %s", symbol)
	}

	quote := result.Chart.Result[0].Indicators.Quote[0]
	var historicalData []StockData

	for i := 0; i < len(quote.Close); i++ {
		historicalData = append(historicalData, StockData{
			Symbol:        symbol,
			Price:         quote.Close[i],
			High:          quote.High[i],
			Low:           quote.Low[i],
			Volume:        quote.Volume[i],
			PreviousClose: quote.Close[i], // For simplicity, reuse close as previous close
		})
	}

	return historicalData, nil
}

// Calculate moving averages with sufficient data
func calculateMovingAverages(historicalData []StockData) (float64, float64) {
	if len(historicalData) == 0 {
		return 0, 0
	}

	// 5-day MA
	ma5Period := min(5, len(historicalData))
	var ma5 float64
	for i := 0; i < ma5Period; i++ {
		ma5 += historicalData[i].Price
	}
	ma5 /= float64(ma5Period)

	// 20-day MA
	ma20Period := min(20, len(historicalData))
	var ma20 float64
	for i := 0; i < ma20Period; i++ {
		ma20 += historicalData[i].Price
	}
	ma20 /= float64(ma20Period)

	return ma5, ma20
}

// Fixed RSI calculation (14-day period, correct sequence)
func calculateRSI(prices []float64) float64 {
	if len(prices) < 15 { // Need 14 days + 1 for change
		return 50 // Neutral if insufficient data
	}

	var gains, losses float64
	// Use the most recent 14 days (newest first in historical data)
	for i := 1; i <= 14; i++ {
		change := prices[i-1] - prices[i] // Newer - Older
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / 14
	avgLoss := losses / 14

	if avgLoss == 0 {
		return 100 // Perfect uptrend
	}
	if avgGain == 0 {
		return 0 // Perfect downtrend
	}

	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

// Calculate stock metrics
func calculateMetrics(data StockData, avgVolume int64, historicalData []StockData) StockMetrics {
	priceChange := ((data.Price - data.PreviousClose) / data.PreviousClose) * 100
	dailyRange := data.High - data.Low
	volatility := (dailyRange / data.Price) * 100
	volumeChange := ((float64(data.Volume) - float64(avgVolume)) / float64(avgVolume)) * 100

	ma5, ma20 := calculateMovingAverages(historicalData)
	priceVsMA5 := ((data.Price - ma5) / ma5) * 100
	priceVsMA20 := ((data.Price - ma20) / ma20) * 100

	var prices []float64
	for _, hData := range historicalData {
		prices = append(prices, hData.Price)
	}
	rsi := calculateRSI(prices)

	return StockMetrics{
		Symbol:       data.Symbol,
		Price:        data.Price,
		PriceChange:  priceChange,
		DailyRange:   dailyRange,
		Volatility:   volatility,
		Volume:       data.Volume,
		VolumeChange: volumeChange,
		MA5:          ma5,
		MA20:         ma20,
		PriceVsMA5:   priceVsMA5,
		PriceVsMA20:  priceVsMA20,
		RSI:          rsi,
	}
}

// Get AI insights from Gemini API
func getGeminiInsights(metrics StockMetrics, historicalData []StockData) (string, error) {
	client := resty.New()
	cfg := config.GetConfig()
	if cfg.GeminiAPIKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not set")
	}

	// 5-day trend
	var trend string
	if len(historicalData) >= 5 {
		firstPrice := historicalData[4].Price // 5th day back
		lastPrice := historicalData[0].Price  // Most recent
		if firstPrice > 0 {
			trendChange := ((lastPrice - firstPrice) / firstPrice) * 100
			if trendChange > 0 {
				trend = fmt.Sprintf("5-day trend: 📈 +%.2f%%", trendChange)
			} else {
				trend = fmt.Sprintf("5-day trend: 📉 %.2f%%", trendChange)
			}
		}
	}

	rsiSignal := "Neutral"
	if metrics.RSI > 70 {
		rsiSignal = "Overbought"
	} else if metrics.RSI < 30 {
		rsiSignal = "Oversold"
	}

	maSignal := "Sideways"
	if metrics.PriceVsMA5 > 1 && metrics.PriceVsMA20 > 1 {
		maSignal = "Strong Uptrend"
	} else if metrics.PriceVsMA5 > 0 && metrics.PriceVsMA20 > 0 {
		maSignal = "Uptrend"
	} else if metrics.PriceVsMA5 < -1 && metrics.PriceVsMA20 < -1 {
		maSignal = "Strong Downtrend"
	} else if metrics.PriceVsMA5 < 0 && metrics.PriceVsMA20 < 0 {
		maSignal = "Downtrend"
	}

	volumeSignal := "Normal Volume"
	if metrics.VolumeChange > 50 {
		volumeSignal = "Very High Volume"
	} else if metrics.VolumeChange > 20 {
		volumeSignal = "High Volume"
	} else if metrics.VolumeChange < -50 {
		volumeSignal = "Very Low Volume"
	} else if metrics.VolumeChange < -20 {
		volumeSignal = "Low Volume"
	}

	prompt := fmt.Sprintf(
		"Analyze %s stock and provide a clear trading recommendation:\n"+
			"Current Price: ₹%.2f (%.2f%% today)\n"+
			"%s\n"+
			"Technical Indicators:\n"+
			"- Price vs 5-day MA: %.2f%%\n"+
			"- Price vs 20-day MA: %.2f%%\n"+
			"- RSI (14): %.2f (%s)\n"+
			"- Moving Average Trend: %s\n"+
			"- Volume Analysis: %s\n"+
			"- Volatility: %.2f%%\n"+
			"Provide a clear, actionable recommendation in 2-3 lines. Include:\n"+
			"1. Entry price and stop-loss levels\n"+
			"2. Key resistance/support levels\n"+
			"3. Risk level (Low/Medium/High)\n"+
			"Be direct and decisive.",
		metrics.Symbol, metrics.Price, metrics.PriceChange, trend,
		metrics.PriceVsMA5, metrics.PriceVsMA20, metrics.RSI, rsiSignal,
		maSignal, volumeSignal, metrics.Volatility,
	)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("key", cfg.GeminiAPIKey).
		SetBody(payload).
		Post("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent")

	if err != nil {
		return "", fmt.Errorf("Gemini API request failed: %v", err)
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(resp.Body(), &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse Gemini response: %v", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no insights returned from Gemini API")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

// Process stocks and generate report
func processStocks() {
	cfg := config.GetConfig()
	stocks := strings.Split(cfg.StockList, ",")

	// Group stocks by market cap
	largeCap := []string{}
	midCap := []string{}
	smallCap := []string{}

	for _, symbol := range stocks {
		symbol = strings.TrimSpace(symbol)
		if symbol == "" {
			continue
		}

		// Categorize stocks based on the defaultStocks list
		switch {
		case contains(defaultStocks[:5], symbol):
			largeCap = append(largeCap, symbol)
		case contains(defaultStocks[5:10], symbol):
			midCap = append(midCap, symbol)
		case contains(defaultStocks[10:], symbol):
			smallCap = append(smallCap, symbol)
		default:
			// If not in default list, assume it's a large cap
			largeCap = append(largeCap, symbol)
		}
	}

	// Process each group separately
	processStockGroup("Large Cap Stocks", largeCap)
	processStockGroup("Mid Cap Stocks", midCap)
	processStockGroup("Small Cap Stocks", smallCap)
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// Process a group of stocks and send notification
func processStockGroup(groupName string, stocks []string) {
	if len(stocks) == 0 {
		return
	}

	cfg := config.GetConfig()
	message := fmt.Sprintf("📊 *%s* - %s\n\n", groupName, time.Now().Format("02-Jan-2006"))

	// Process stocks in groups of 5
	for i := 0; i < len(stocks); i += 5 {
		end := i + 5
		if end > len(stocks) {
			end = len(stocks)
		}

		currentGroup := stocks[i:end]
		var messages []string

		for _, symbol := range currentGroup {
			data, err := fetchYahooFinanceData(symbol)
			if err != nil {
				fmt.Printf("Error fetching %s: %v\n", symbol, err)
				continue
			}

			historicalData, err := fetchHistoricalData(symbol)
			if err != nil {
				fmt.Printf("Error fetching historical data for %s: %v\n", symbol, err)
				continue
			}

			var avgVolume int64
			if len(historicalData) > 0 {
				totalVolume := int64(0)
				for _, hData := range historicalData {
					totalVolume += hData.Volume
				}
				avgVolume = totalVolume / int64(len(historicalData))
			}

			metrics := calculateMetrics(data, avgVolume, historicalData)
			insights, err := getGeminiInsights(metrics, historicalData)
			if err != nil {
				fmt.Printf("Error getting insights for %s: %v\n", symbol, err)
				continue
			}

			priceChangeEmoji := "📈"
			if metrics.PriceChange < 0 {
				priceChangeEmoji = "📉"
			}

			volumeChangeEmoji := "📊"
			if metrics.VolumeChange > 20 {
				volumeChangeEmoji = "🚀"
			} else if metrics.VolumeChange < -20 {
				volumeChangeEmoji = "📉"
			}

			companyName := strings.TrimSuffix(symbol, ".NS")
			technicalIndicators := fmt.Sprintf(
				"Price vs 5-day MA: %.2f%%\n"+
					"Price vs 20-day MA: %.2f%%\n"+
					"RSI (14): %.2f\n"+
					"Volatility: %.2f%%",
				metrics.PriceVsMA5, metrics.PriceVsMA20, metrics.RSI, metrics.Volatility,
			)

			stockMessage := fmt.Sprintf(
				"*%s* (%s)\n"+
					"💰 *Price*: ₹%.2f %s (*%.2f%%*)\n"+
					"📈 *Volume*: %s %.2f%% vs avg\n"+
					"📊 *Technical Indicators*:\n```\n%s\n```\n"+
					"🤖 *AI Insights*:\n%s\n",
				companyName, symbol, data.Price, priceChangeEmoji, metrics.PriceChange,
				volumeChangeEmoji, metrics.VolumeChange, technicalIndicators, insights,
			)
			messages = append(messages, stockMessage)
		}

		if len(messages) > 0 {
			groupMessage := message + strings.Join(messages, "\n---\n")
			fmt.Println(groupMessage)

			if cfg.TelegramBotToken != "" && len(cfg.TelegramChatIDs) > 0 {
				sendStockTelegramNotification(groupMessage)
			}
		}
	}
}

// RunStockAnalysis executes the analysis
func RunStockAnalysis() {
	fmt.Println("Running stock analysis...")
	processStocks()
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
