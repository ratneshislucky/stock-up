# Stocking Up

This project contains a Go program that runs daily using GitHub Actions.

## Project Structure

- `MarketFall.go`: The main Go program file.
- `.github/workflows/daily-schedule.yml`: The GitHub Actions workflow to run the Go file daily at 2 PM IST.

## How It Works

- The GitHub Actions workflow (`.github/workflows/run-go-file.yml`) is scheduled to run every day at **2:00 PM IST**.
- It sets up the Go environment and runs the Go program (`MarketFall.go`).
