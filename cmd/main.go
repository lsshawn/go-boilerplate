package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/lsshawn/go-todo/views"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	dotenvErr := godotenv.Load()
	if dotenvErr != nil {
		log.Fatal("Error loading .env file")
	}
	app := echo.New()

	app.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	app.GET("/", func(c echo.Context) error {
		// aiChat()
		component := views.Index()
		return component.Render(context.Background(), c.Response().Writer)
	})

	app.GET("/chat-sse", echo.WrapHandler(http.HandlerFunc(sseHandler)))

	app.Static("/css", "css")
	app.Static("/static", "static")

	app.Logger.Fatal(app.Start(":1323"))
}

func aiChat() {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 20,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Hello world",
			},
		},
		Stream: true,
	}
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Printf("\nStream finished")
			return
		}
		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return
		}

		fmt.Printf(response.Choices[0].Delta.Content)
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Send a new message every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Fprintf(w, "data: %s\n\n", time.Now().Format("15:04:05"))
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
