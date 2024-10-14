package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Define structures for the requests

type IndexRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// Function to get current date and date one week ago in required format
func getDates() (string, string) {
	endDate := time.Now().Format("02-Jan-2006")
	startDate := time.Now().AddDate(0, 0, -7).Format("02-Jan-2006")
	return startDate, endDate
}

// Function to send Telegram notification
func sendTelegramNotification(message string, botToken string, chatID string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}

	payloadBytes, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error sending Telegram notification:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Notification sent successfully!")
	} else {
		fmt.Printf("Failed to send notification. Status code: %d\n", resp.StatusCode)
	}
}

// Function to fetch data from the API
func fetchIndexData(index IndexRequest) (string, error) {
	url := "https://www.niftyindices.com/Backpage.aspx/getHistoricaldataDBtoString"

	headers := map[string]string{
		"Content-Type":     "application/json; charset=utf-8",
		"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:131.0) Gecko/20100101 Firefox/131.0",
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"X-Requested-With": "XMLHttpRequest",
		"Origin":           "https://www.niftyindices.com",
	}

	reqBody, _ := json.Marshal(index)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	output := strings.Split(result["d"], "[")[0]
	return strings.TrimSpace(output), nil
}

func main() {
	// Define the start and end dates
	startDate, endDate := getDates()

	// Define the index data
	indices := []IndexRequest{
		{Name: "Nifty 50", StartDate: startDate, EndDate: endDate},
		{Name: "NIFTY100 LOWVOL30", StartDate: startDate, EndDate: endDate},
		{Name: "Nifty200Momentm30", StartDate: startDate, EndDate: endDate},
		{Name: "Nifty500 Momentum 50", StartDate: startDate, EndDate: endDate},
		{Name: "Nifty Midcap150 Momentum 50", StartDate: startDate, EndDate: endDate},
	}

	// Telegram bot token and chat ID
	botToken := "7541372157:AAFxgVbgOVH2uk04yMD0dR78rwVYnTgKo6c"
	chatID := "754100903"

	var returns []float64
	var messages []string

	for _, index := range indices {
		// Fetch data from the API
		data, err := fetchIndexData(index)
		if err != nil {
			fmt.Println("Error fetching data for", index.Name, ":", err)
			continue
		}

		fmt.Printf("Data: %s\n", data)

		// Convert the return value to float64
		returnValue, err := strconv.ParseFloat(data, 64)
		if err != nil {
			fmt.Println("Error parsing return value for", index.Name, ":", err)
			continue
		}

		returns = append(returns, returnValue)
		messages = append(messages, fmt.Sprintf("%s: %f", index.Name, returnValue))
	}

	// Check if all returns are negative and send a Telegram notification if so
	allNegative := true
	for _, ret := range returns {
		if ret >= 0 {
			allNegative = false
			break
		}
	}

	if allNegative {
		message := fmt.Sprintf("Indices are negative, from: %s to: %s\n%s", startDate, endDate, strings.Join(messages, "\n"))
		fmt.Println(message)
		sendTelegramNotification(message, botToken, chatID)
	} else {
		fmt.Println("Not all returns are negative.")
		message := "Stock Update : Not a big gap"
		sendTelegramNotification(message, botToken, chatID)
	}
}
