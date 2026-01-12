package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "go-ai-concurrency-challenge/docs"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go AI Concurrency Challenge API
// @version 1.0
// @description A simple API for AI text generation with concurrency handling.
// @host localhost:8080
// @BasePath /

type AIRequest struct {
	Text     string
	Response chan string
}

type Server struct {
	requests chan AIRequest
}

// GenerateAIResponse handles the AI generation request
// @Summary Generate AI response
// @Description Send text to AI and get response
// @Accept plain
// @Produce plain
// @Param text body string true "Text to process"
// @Success 200 {string} string "AI response"
// @Failure 400 {string} string "Bad request"
// @Failure 405 {string} string "Method not allowed"
// @Router /generate [post]
func (s *Server) GenerateAIResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	text := string(body)

	// Create request
	req := AIRequest{Text: text, Response: make(chan string, 1)}

	// Send to channel
	s.requests <- req

	// Wait for response
	response := <-req.Response

	w.Write([]byte(response))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY not set")
	}

	client := openai.NewClient(apiKey)

	// Channel for requests
	requests := make(chan AIRequest, 100) // buffer for queuing

	// Start workers
	numWorkers := os.Getenv("WORKER_POOL_SIZE")
	workers, err := strconv.Atoi(numWorkers)
	if err != nil {
		workers = 5
	}
	log.Println("Starting worker pool with size:", workers)
	for i := 0; i < workers; i++ {
		go worker(client, requests)
	}

	server := &Server{requests: requests}

	// HTTP handler
	http.HandleFunc("/generate", server.GenerateAIResponse)

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func worker(client *openai.Client, requests <-chan AIRequest) {
	for req := range requests {
		// Call OpenAI
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT5Nano,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: req.Text,
					},
				},
			},
		)
		if err != nil {
			if apiErr, ok := err.(*openai.APIError); ok && apiErr.HTTPStatusCode == 429 {
				req.Response <- "Quota exceeded. Please check your OpenAI account and billing details."
			} else {
				req.Response <- "Error: " + err.Error()
			}
			continue
		}

		if len(resp.Choices) > 0 {
			req.Response <- resp.Choices[0].Message.Content
		} else {
			req.Response <- "No response"
		}
	}
}
