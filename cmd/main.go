package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/lsshawn/go-todo/internal/dto"
	"github.com/lsshawn/go-todo/views"
	"github.com/lsshawn/go-todo/views/components"
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
		msgs := []*dto.MessageDto{}
		component := views.Index(msgs)
		return component.Render(context.Background(), c.Response().Writer)
	})

	app.POST("/user-message", func(c echo.Context) error {
		userInput := c.FormValue("text") // Retrieve the user's input from the form

		msgDto := &dto.MessageDto{
			Role: "user",
			Text: userInput,
		}

		component := components.ChatMessage(msgDto)
		return component.Render(context.Background(), c.Response().Writer)
	})

	app.POST("/ask", func(c echo.Context) error {
		question := c.FormValue("text")
		fmt.Printf("LS -> cmd/main.go:50 -> question: %+v\n", question)
		msg, err := aiAssistant(question)
		if err != nil {
			return err
		}
		msgDto := &dto.MessageDto{
			Role:  msg.Role,
			Text:  msg.Content[0].Text.Value,
			RunID: msg.RunID,
		}
		component := components.ChatMessage(msgDto)
		return component.Render(context.Background(), c.Response().Writer)
	})

	app.Static("/css", "css")
	app.Static("/static", "static")

	app.Logger.Fatal(app.Start(":1323"))
}

func checkRun(client *openai.Client, threadID, runID string) (*openai.Run, error) {
	ctx := context.Background()
	for {
		run, err := client.RetrieveRun(ctx, threadID, runID)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving run: %v\n", err)
		}

		switch run.Status {
		case "completed":
			fmt.Println("Run completed")
			return &run, nil
		case "expired":
			fmt.Println("Run expired")
			return &run, nil
		case "requires_action":
			// Handle required action here (not implemented in this example)
			fmt.Println("Run requires action")
			fmt.Printf("Run: %+v\n", run)
			return &run, nil
		default:
			fmt.Println("Running...")
			time.Sleep(3 * time.Second)
		}
	}
}

func aiAssistant(question string) (*openai.Message, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()

	req := openai.CreateThreadAndRunRequest{
		RunRequest: openai.RunRequest{
			AssistantID: "asst_L2KwWjgK0JXDdgixLwzKkJYm",
		},
		Thread: openai.ThreadRequest{
			Messages: []openai.ThreadMessage{
				{
					Role:    openai.ThreadMessageRoleUser,
					Content: question,
				},
			},
		},
	}

	resp, err := client.CreateThreadAndRun(ctx, req)
	if err != nil {
		fmt.Printf("CreateThreadAndRunRequest error: %v\n", err)
		return nil, err
	}
	fmt.Printf("Thread ID: %s\n", resp.ThreadID)
	fmt.Printf("Run ID: %s\n", resp.ID)

	run, runErr := checkRun(client, resp.ThreadID, resp.ID)
	if runErr != nil {
		fmt.Printf("checkRun error: %v\n", runErr)
		return nil, runErr
	}
	fmt.Printf("Run: %+v\n", run)

	limit := 3
	res, listErr := client.ListMessage(ctx, resp.ThreadID, &limit, nil, nil, nil)
	for _, msg := range res.Messages {
		fmt.Printf("msg %+v\n", msg)
		fmt.Printf("Message: %s: %s\n", msg.Role, msg.Content[0].Text.Value)
	}
	// Check if there are any messages
	if len(res.Messages) > 0 {
		lastMsg := res.Messages[0]
		fmt.Printf("Last Message: %s: %s\n", lastMsg.Role, lastMsg.Content[0].Text.Value)
		return &lastMsg, nil
	} else {
		fmt.Println("No messages found.")
		return nil, listErr
	}
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
