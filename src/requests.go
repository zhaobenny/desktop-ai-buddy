package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/jpeg"
	"log"

	openai "github.com/sashabaranov/go-openai"
)

const (
	chatModel       = "gryphe/mythomist-7b"
	chatTemperature = 1.25
	visionModel     = "openai/gpt-4-vision-preview"
)

func createClient(key string) *openai.Client {
	config := openai.DefaultConfig(key)
	config.BaseURL = "https://openrouter.ai/api/v1/"

	client := openai.NewClientWithConfig(config)

	return client
}

func requestChat(client openai.Client, prompt string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: chatModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
			},
			MaxTokens:   50,
			Temperature: chatTemperature,
		},
	)
	if err != nil {
		log.Println("Chat error:", err)

		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func requestImageComment(client openai.Client, image image.Image) (string, error) {
	base64image, err := imageToBase64(image)
	if err != nil {
		return "", err
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "openai/gpt-4-vision-preview",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Do not act like an assistant. Cheerfully give a short comment on the activity in the image. You are not allowed to fully describe the activity.",
				},
				{
					Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{
							Type: openai.ChatMessagePartTypeText,
							Text: "This is what I am doing right now:",
						},
						{
							Type: openai.ChatMessagePartTypeImageURL,
							ImageURL: &openai.ChatMessageImageURL{
								Detail: openai.ImageURLDetailAuto,
								URL:    "data:image/png;base64," + base64image,
							},
						},
					},
				},
			},
			MaxTokens: 200,
		},
	)

	if err != nil {
		log.Println("Image comment error:", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func imageToBase64(img image.Image) (string, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	if err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64String, nil
}
