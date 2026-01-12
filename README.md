# Go AI Concurrency Challenge

This is a technical challenge for evaluating Go (Golang) development skills, focusing on concurrency, API integration, and clean code practices.

## Objective

Develop a small Go service that integrates with an AI API (OpenAI) and handles concurrent requests efficiently using goroutines and channels.

## Features

- HTTP endpoint that accepts POST requests with text content
- Sends the text to OpenAI API for processing
- Returns the AI-generated response
- Concurrent handling of multiple requests using a worker pool pattern
- Queuing system using Go channels

## Requirements

- Go 1.19 or later
- OpenAI API key

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/acosta9/go-ai-concurrency-challenge
   cd go-ai-concurrency-challenge
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Create a `.env` file in the root directory and add your OpenAI API key:
   ```
   OPENAI_API_KEY=your_openai_api_key_here
   WORKER_POOL_SIZE=5
   ```

## Usage

1. Run the server:
   ```bash
   go run main.go
   ```

2. The server will start on `http://localhost:8080`

3. API documentation via swagger is available at `http://localhost:8080/swagger/`

## Testing

Send a POST request to `/generate` with text in the body:

```bash
curl -X POST http://localhost:8080/generate \
  -H "Content-Type: text/plain" \
  -d "Hello, how are you?"
```

The response will be the AI-generated text.

## Architecture

- **Concurrency**: Uses goroutines and channels for request queuing and processing
- **Worker Pool**: 5 worker as default goroutines handle AI API calls concurrently
- **Queuing**: Buffered channel queues incoming requests
- **API Integration**: Integrates with OpenAI Chat Completion API