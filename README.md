# Spill Telegram Bot Backend

A Go-based backend service for the Spill Telegram bot, providing nickname generation and other features.

## Features

- Nickname generation using word lists (adjectives, colors, nouns)
- Telegram bot integration
- RESTful API endpoints

## Project Structure

```
spill-backend/
├── words/
│   ├── adjectives.txt
│   ├── colors.txt
│   └── nouns.txt
├── .gitignore
└── README.md
```

## Setup

### Prerequisites

- Go 1.21 or higher
- Telegram Bot Token

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd spill-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:
```bash
go run main.go
```

## Development

### Building

```bash
go build -o bin/spill-backend
```

### Running Tests

```bash
go test ./...
```

## License

[Add your license here]
