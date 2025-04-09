# TL;DR Telegram Bot

## Overview
The TL;DR Telegram Bot is a Go-based application that integrates with the Telegram messaging platform to provide summarization of chat messages using a local LLM (Ollama). The bot listens for specific trigger words in replies and collects messages for summarization, which it then sends to the Ollama server for processing.

## Features
- Responds to specific trigger words in Telegram group chats.
- Collects and summarizes messages from the chat.
- Integrates with a local Ollama LLM server for summarization.
- Logs all received messages to a database.
- Configurable via environment variables.

## Project Structure
```
tldr-telegram-bot
├── cmd
│   └── bot
│       └── main.go
├── internal
│   ├── config
│   ├── db
│   ├── llm
│   ├── telegram
│   └── utils
├── .dockerignoreI
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

## Requirements
- Go 1.18 or later
- Docker and Docker Compose

## Environment Variables
Create a `.env` file in the root directory based on the provided `.env.example` file. The following environment variables are required:
- `TELEGRAM_BOT_TOKEN`: Your Telegram bot token.
- `DEFAULT_LANG`: Default language for summarization (`pt`, `en`, or `es`).
- `OLLAMA_MODEL`: The model name to be used by the Ollama API.
- `OLLAMA_MODELS`: Comma-separated list of models available for summarization.
- `AUTHORIZED_GROUPS`: Comma-separated list of authorized group IDs.

## Running the Project Locally
1. Clone the repository:
   ```
   git clone https://tldr-telegram-bot/.git
   cd tldr-telegram-bot
   ```

2. Create a `.env` file based on `.env.example` and fill in the required values.

3. Install dependencies:
   ```
   go mod tidy
   ```

4. Run the bot:
   ```
   go run cmd/bot/main.go
   ```

## Running the Project with Docker
1. Ensure Docker and Docker Compose are installed.

2. Create a `.env` file based on `.env.example` and fill in the required values.

3. Build and run the services:
   ```
   docker-compose up --build
   ```

4. Access the bot in your authorized Telegram group.

## Logging
The bot uses structured logging to track critical events and errors. Ensure to monitor the logs for any issues during operation.

## Testing
Unit tests are included for critical functionalities. Run tests using:
```
go test ./...
```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.