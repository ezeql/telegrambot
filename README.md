# Bitcoin Price Telegram Bot 🤖 📈 

A Telegram bot that monitors and reports Bitcoin price changes from Binance, with fun emoji indicators for price movements.

## Features

- Real-time Bitcoin price monitoring from Binance API
- Price updates every 30 seconds
- Price change percentage over 24 hours
- Dynamic emoji indicators:
  - 🚀 for positive price movements (1-3 rockets based on % increase)
  - 🧂 for negative price movements (1-3 salt shakers based on % decrease)

## Setup

1. Create a Telegram bot and get your bot token from [@BotFather](https://t.me/botfather)

2. Set the environment variable:

   ```bash
    TELEGRAM_BOT_TOKEN=replace-token-here go run main.go
    ```
