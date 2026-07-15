package ai

import (
	"context"
	"fmt"
	"sync"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var client openai.Client
var userHistories = make(map[string][]openai.ChatCompletionMessageParamUnion)
var mu sync.Mutex

func InitAi() {
	client = openai.NewClient(
		option.WithBaseURL("http://localhost:8080/v1"),
		option.WithAPIKey("dummy"),
	)

	fmt.Println("AI Engine berhasil diinisialisasi.")
}

func TanyaAi(userID string, userInput string) string {
	ctx := context.Background()

	instruksiSistem := `Anda adalah FikomBot, asisten virtual resmi Fakultas Ilmu Komputer UDB Surakarta.
Jawablah dengan bahasa Indonesia yang sopan, singkat, dan mudah dipahami.`

	// Ambil histori user
	mu.Lock()
	chatHistory := userHistories[userID]
	mu.Unlock()

	// Susun pesan
	var currentPayload []openai.ChatCompletionMessageParamUnion

	currentPayload = append(currentPayload,
		openai.SystemMessage(instruksiSistem),
	)

	currentPayload = append(currentPayload, chatHistory...)

	currentPayload = append(currentPayload,
		openai.UserMessage(userInput),
	)

	// Request ke LLaMA Server
	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel("/Users/agiestanajwaputri/Downloads/llama-b9632/qwen2.5-0.5b-instruct-q4_0.gguf"),
		Messages: currentPayload,
	})

	if err != nil {
		fmt.Println("AI ERROR:", err)
		return err.Error()
	}
	var jawabanAi string

	if len(resp.Choices) > 0 {
		jawabanAi = resp.Choices[0].Message.Content
	} else {
		jawabanAi = "Maaf, AI tidak memberikan jawaban."
	}

	// Simpan histori maksimal 10 percakapan
	mu.Lock()

	if len(chatHistory) >= 10 {
		chatHistory = chatHistory[2:]
	}

	chatHistory = append(chatHistory,
		openai.UserMessage(userInput),
		openai.AssistantMessage(jawabanAi),
	)

	userHistories[userID] = chatHistory

	mu.Unlock()

	return jawabanAi
}
