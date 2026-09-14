package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Message represents a single chat turn for the Hugging Face API
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// HfRequest matches the Hugging Face chat completion payload structure
type HfRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// HfStreamResponse handles the server-sent events chunk structure
type HfStreamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

const systemPrompt = `You are Light, a sharp, tech-savvy digital friend and versatile terminal companion. You have a laid-back, peer-to-peer vibe—like an experienced co-developer hanging out in the shell with the user. You are concise, deeply technical when needed, and always ready to help code, debug, or just talk tech. Never sound like a corporate customer service bot.`

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "chat":
		runChatLoop()
	case "ask":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing prompt for 'ask'. Usage: light ask \"your question here\"")
			os.Exit(1)
		}
		prompt := strings.Join(os.Args[2:], " ")
		runSingleShot(prompt)
	case "clear":
		clearHistory()
	case "help":
		printUsage()
	default:
		// Default behavior: treat all arguments as a single-shot query if they don't match a command
		prompt := strings.Join(os.Args[1:], " ")
		runSingleShot(prompt)
	}
}

func printUsage() {
	fmt.Println("Light - Your terminal companion")
	fmt.Println("\nUsage:")
	fmt.Println("  light <prompt>          Quick single-shot question")
	fmt.Println("  light ask <prompt>      Explicit single-shot question")
	fmt.Println("  light chat              Start an interactive persistent session")
	fmt.Println("  light clear             Clear local chat conversation history")
	fmt.Println("  light help              Show this help menu")
}

func getApiKey() string {
	apiKey := os.Getenv("HUGGINGFACE_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: HUGGINGFACE_API_KEY environment variable is not set.")
		fmt.Println("Run: export HUGGINGFACE_API_KEY=\"your_token_here\"")
		os.Exit(1)
	}
	return apiKey
}

func getHistoryPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".light_history.json"
	}
	configDir := filepath.Join(home, ".config", "light")
	_ = os.MkdirAll(configDir, 0755)
	return filepath.Join(configDir, "history.json")
}

func loadHistory() []Message {
	path := getHistoryPath()
	file, err := os.Open(path)
	if err != nil {
		// Return default system prompt if no history exists yet
		return []Message{
			{Role: "system", Content: systemPrompt},
		}
	}
	defer file.Close()

	var messages []Message
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&messages); err != nil {
		return []Message{
			{Role: "system", Content: systemPrompt},
		}
	}
	return messages
}

func saveHistory(messages []Message) {
	path := getHistoryPath()
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(messages)
}

func clearHistory() {
	path := getHistoryPath()
	if err := os.Remove(path); err == nil {
		fmt.Println("Light: Conversation history cleared.")
	} else {
		fmt.Println("Light: No history found to clear.")
	}
}

func runSingleShot(prompt string) {
	apiKey := getApiKey()
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: prompt},
	}
	streamChatCompletion(apiKey, messages, false)
	fmt.Println()
}

func runChatLoop() {
	apiKey := getApiKey()
	messages := loadHistory()

	fmt.Println("Light: Session started. Type 'exit' or 'quit' to leave.")
	fmt.Println("-----------------------------------------------------")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\033[32m> \033[0m") // Neon green prompt indicator
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" {
			fmt.Println("Light: Catch you later.")
			break
		}

		// Append user message
		messages = append(messages, Message{Role: "user", Content: input})

		// Stream response from API
		fmt.Print("\033[36mLight: \033[0m")
		fullResponse := streamChatCompletion(apiKey, messages, true)
		fmt.Println()

		// Append assistant response to history and persist
		messages = append(messages, Message{Role: "assistant", Content: fullResponse})
		saveHistory(messages)
	}
}

func streamChatCompletion(apiKey string, messages []Message, stream bool) string {
	// Using a reliable, high-performance open model available on the Hugging Face Serverless API
	payload := HfRequest{
		Model:    "meta-llama/Llama-3.3-70B-Instruct",
		Messages: messages,
		Stream:   stream,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error building request: %v\n", err)
		return ""
	}

	req, err := http.NewRequest("POST", "https://router.huggingface.co/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err)
		return ""
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("\nConnection error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("\nAPI Error (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
		return ""
	}

	var fullResponse strings.Builder

	if stream {
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				break
			}

			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}

			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				if data == "[DONE]" {
					break
				}

				var chunk HfStreamResponse
				if err := json.Unmarshal([]byte(data), &chunk); err == nil {
					if len(chunk.Choices) > 0 {
						content := chunk.Choices[0].Delta.Content
						if content != "" {
							fmt.Print(content)
							fullResponse.WriteString(content)
						}
					}
				}
			}
		}
	} else {
		// Handle non-streaming response format fallback
		var nonStreamResp struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(bodyBytes, &nonStreamResp); err == nil && len(nonStreamResp.Choices) > 0 {
			content := nonStreamResp.Choices[0].Message.Content
			fmt.Print(content)
			fullResponse.WriteString(content)
		} else {
			fmt.Println(string(bodyBytes))
		}
	}

	return fullResponse.String()
}
