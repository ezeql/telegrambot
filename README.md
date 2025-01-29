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

   To get your `TELEGRAM_GROUP_CHAT_ID` (for group chats):
   1. Add the bot to your group
   2. Send a message in the group
   3. Access <https://api.telegram.org/bot><YourBOTToken>/getUpdates
   4. Look for the "chat" -> "id" field in the group message response
   Note: Group chat IDs are typically negative numbers

2. Set the environment variables:

   ```bash
   export TELEGRAM_BOT_TOKEN=replace-token-here
   export TELEGRAM_GROUP_CHAT_ID=replace-group-chat-id-here
   ```

3. run

  ```bash
  go run main.go
  ```
