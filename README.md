# Go Stock Market Insights

This application provides daily stock market insights using Yahoo Finance data and Google's Gemini AI. It runs at 8 AM daily and sends insights via Telegram.

## Features

- Tracks 30 Indian stocks across different market caps (Large, Mid, and Small)
- Fetches real-time data from Yahoo Finance
- Generates AI-powered insights using Google Gemini
- Monitors NIFTY indices for market falls
- Sends daily reports via Telegram
- Runs automatically at 8 AM daily via GitHub Actions

## Prerequisites

- Go 1.21 or higher
- Google Gemini API Key
- Telegram Bot Token and Chat ID

## Setup

### Local Development

1. Clone the repository:
```bash
git clone <repository-url>
cd go-stock
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
export GEMINI_API_KEY="your_gemini_api_key"
export TELEGRAM_BOT_TOKEN="your_telegram_bot_token"
export TELEGRAM_CHAT_ID="your_telegram_chat_id"
```

### GitHub Actions Setup

1. Go to your GitHub repository
2. Click on "Settings"
3. In the left sidebar, click on "Secrets and variables" → "Actions"
4. Click "New repository secret"
5. Add the following secrets:
   - `GEMINI_API_KEY`: Your Google Gemini API key
   - `TELEGRAM_BOT_TOKEN`: Your Telegram bot token
   - `TELEGRAM_CHAT_ID`: Your Telegram chat ID

The GitHub Action will automatically run at 8 AM UTC daily and use these secrets.

## Running the Application

### Local Development
```bash
go run main.go
```

This will start both the stock insights and market fall check jobs:
- Stock Insights: 8 AM UTC daily
- Market Fall Check: 8:30 AM UTC daily

### GitHub Actions
The application runs automatically at 8 AM UTC daily. You can also trigger it manually:
1. Go to the "Actions" tab in your repository
2. Click on "Daily Stock Analysis"
3. Click "Run workflow"

## Stock List

The application tracks the following Indian stocks:

### Large Cap Stocks (NIFTY 50)
- RELIANCE.NS (Reliance Industries)
- TCS.NS (Tata Consultancy Services)
- HDFCBANK.NS (HDFC Bank)
- INFY.NS (Infosys)
- ICICIBANK.NS (ICICI Bank)
- HINDUNILVR.NS (Hindustan Unilever)
- SBIN.NS (State Bank of India)
- BHARTIARTL.NS (Bharti Airtel)
- ITC.NS (ITC Limited)
- KOTAKBANK.NS (Kotak Mahindra Bank)

### Mid Cap Stocks (NIFTY Midcap 100)
- POLYCAB.NS (Polycab India)
- PERSISTENT.NS (Persistent Systems)
- TATAMOTORS.NS (Tata Motors)
- MOTHERSON.NS (Motherson Sumi)
- APOLLOTYRE.NS (Apollo Tyres)
- BAJAJFINSV.NS (Bajaj Finserv)
- BAJAJHLDNG.NS (Bajaj Holdings)
- DIXON.NS (Dixon Technologies)
- ZYDUSLIFE.NS (Zydus Lifesciences)
- ALKEM.NS (Alkem Laboratories)

### Small Cap Stocks (NIFTY Smallcap 100)
- JINDALSAW.NS (Jindal Saw)
- MAHINDRAFORG.NS (Mahindra Forgings)
- BALKRISIND.NS (Balkrishna Industries)
- HUDCO.NS (Housing & Urban Development Corp)
- JINDALSTEL.NS (Jindal Steel & Power)
- GODREJIND.NS (Godrej Industries)
- GODREJPROP.NS (Godrej Properties)
- JUBLFOOD.NS (Jubilant FoodWorks)
- KALYANKJIL.NS (Kalyan Jewellers)
- KARURVYSYA.NS (Karur Vysya Bank)

## Output Format

The daily report includes:

### Stock Insights
- Current stock price (in Indian Rupees)
- Price change percentage
- AI-generated insights including:
  - Market trends
  - Potential risks
  - Investment opportunities
  - Volume analysis
  - Volatility assessment

### Market Fall Check
- NIFTY indices performance
- Weekly return analysis
- Market trend indicators
- Risk assessment
