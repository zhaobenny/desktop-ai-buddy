package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sashabaranov/go-openai"
)

type Config struct {
	API_Key string `json:"openrouter_api_key"`
}

type Point struct {
	X, Y int
}

type Animation struct {
	sprite      *ebiten.Image
	frameHeight int
	frameWidth  int
	frameCount  int
	speed       int
}

type Action struct {
	name      string
	animation Animation
}

type Game struct {
	drawTick             int
	updateTick           int
	windowPos            Point
	initalLeftPress      Point
	curAction            Action
	isTalking            bool
	idleRightAction      Action
	idleLeftAction       Action
	walkRightAction      Action
	walkLeftAction       Action
	gettingDraggedAction Action
	speechBubble         *ebiten.Image
	openaiClient         *openai.Client
}
