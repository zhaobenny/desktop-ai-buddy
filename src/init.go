package main

import (
	"bytes"
	"embed"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

const (
	startYcoord = 1080 / 3
	startXcoord = 1920 / 2
)

func NewGame(config *Config) *Game {

	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowPosition(startYcoord, startXcoord)

	game := &Game{
		windowPos:    Point{X: startXcoord, Y: startYcoord},
		isTalking:    false,
		speechBubble: nil,
		openaiClient: nil,
	}

	if config.API_Key != defaultKeyMessage {
		game.openaiClient = createClient(config.API_Key)
	}

	loadActions(game)

	return game
}

func loadActions(game *Game) {

	idleRightAction := Action{
		name:      "idleRight",
		animation: loadAnimation("assets/Pink_Monster_Idle_4.png", 4, 50, false),
	}

	idleLeftAction := Action{
		name:      "idleLeft",
		animation: loadAnimation("assets/Pink_Monster_Idle_4.png", 4, 50, true),
	}

	walkRightAction := Action{
		name:      "walkRight",
		animation: loadAnimation("assets/Pink_Monster_Walk_6.png", 6, 25, false),
	}

	walkLeftAction := Action{
		name:      "walkLeft",
		animation: loadAnimation("assets/Pink_Monster_Walk_6.png", 6, 25, true),
	}

	gettingDraggedAction := Action{
		name:      "gettingDragged",
		animation: loadAnimation("assets/Pink_Monster_Dragged_4.png", 4, 35, false),
	}

	game.gettingDraggedAction = gettingDraggedAction

	game.idleRightAction = idleRightAction
	game.idleLeftAction = idleLeftAction

	game.walkRightAction = walkRightAction
	game.walkLeftAction = walkLeftAction

	game.curAction = idleRightAction
}

func loadAnimation(filepath string, frameCount int, speed int, mirror bool) Animation {
	img := loadImageFile(filepath)

	if mirror {
		img = mirrorImage(img)
	}

	return Animation{
		sprite:      ebiten.NewImageFromImage(img),
		frameWidth:  32,
		frameHeight: 32,
		frameCount:  frameCount,
		speed:       speed,
	}
}

func loadImageFile(filepath string) image.Image {
	data, err := assets.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func mirrorImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	mirroredImg := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			colorOriginal := img.At(x, y)
			xMirrored := width - 1 - x
			mirroredImg.Set(xMirrored, y, colorOriginal)
		}
	}

	return mirroredImg
}
