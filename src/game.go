package main

import (
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kbinani/screenshot"
)

const (
	windowWidth  = 128
	windowHeight = 128
)

func (g *Game) decideNewAction() Action {
	if g.isTalking {
		if g.windowPos.X > 960 {
			return g.idleLeftAction
		} else {
			return g.idleRightAction
		}
	}

	randomPercent := rand.Float64()
	walkProbability := 0.25
	pointRight := rand.Float64() < 0.5

	switch {
	case randomPercent < walkProbability:
		if pointRight {
			return g.walkRightAction
		}
		return g.walkLeftAction
	default:
		switch g.curAction.name {
		case "walkRight":
			return g.idleRightAction
		case "walkLeft":
			return g.idleLeftAction
		default:
			if pointRight {
				return g.idleRightAction
			}
			return g.idleLeftAction
		}
	}
}

func (g *Game) decideMessage() {
	randomPercent := rand.Float64()
	screenProbability := 0.025
	talkProbability := 0.1

	switch {
	case randomPercent < screenProbability:
		image, error := screenshot.CaptureDisplay(0)
		if error != nil {
			g.isTalking = false
			return
		}
		log.Println("Requesting comment on screenshot")
		msg, err := requestImageComment(*g.openaiClient, image)
		if err != nil {
			g.isTalking = false
			return
		}
		g.speechBubble = ebiten.NewImageFromImage(create_text_bubble(msg))
		g.isTalking = true

	case randomPercent < screenProbability+talkProbability:
		currentTime := time.Now()
		timeString := currentTime.Format("15:04")
		log.Println("Requesting chat")
		msg, err := requestChat(*g.openaiClient, "The time is "+timeString+"."+" Give a cheerful short remark to the user")
		if err != nil {
			g.isTalking = false
			return
		}
		g.speechBubble = ebiten.NewImageFromImage(create_text_bubble(msg))
		g.isTalking = true
	default:
		g.isTalking = false
	}

}

func (g *Game) Update() error {
	g.updateTick++

	isLeftClick := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if isLeftClick {
		g.moveWindow()
	} else if g.curAction.name == "gettingDragged" {
		g.curAction = g.decideNewAction()
	} else if g.updateTick%200 == 0 {
		g.curAction = g.decideNewAction()
	}
	if g.updateTick%500 == 0 && !g.isTalking && g.openaiClient != nil {
		g.speechBubble = nil
		g.isTalking = true
		go g.decideMessage()
	} else if g.isTalking && (ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) || g.updateTick%2000 == 0) {
		g.isTalking = false
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawTick++

	ebiten.SetWindowPosition(g.windowPos.X, g.windowPos.Y)

	x, y := ebiten.WindowSize()
	op := &ebiten.DrawImageOptions{}
	if g.isTalking && g.speechBubble != nil {

		bounds := g.speechBubble.Bounds().Size()
		if (bounds.X+50) > x && y > -1 {
			ebiten.SetWindowSize((windowWidth)+(bounds.X)+25, windowHeight)
		}
		x, y = ebiten.WindowSize()
		if bounds.Y > y {
			ebiten.SetWindowSize(x, (windowHeight/2)+(bounds.Y))
		}

		op.GeoM.Translate(100, 0)
		screen.DrawImage(g.speechBubble, op)

	} else if x != windowWidth || y != windowHeight {
		ebiten.SetWindowSize(windowWidth, windowHeight)
	}

	currSheet := g.curAction.animation.sprite
	frameHeight := g.curAction.animation.frameHeight
	frameWidth := g.curAction.animation.frameWidth
	frameCount := g.curAction.animation.frameCount
	speed := g.curAction.animation.speed

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(18, float64(windowHeight/2)-float64(frameHeight*3/2))
	i := (g.drawTick / speed) % frameCount
	sx := i * frameWidth

	frame := currSheet.SubImage(image.Rect(sx, 0, sx+frameWidth, frameHeight)).(*ebiten.Image)
	screen.DrawImage(frame, op)

	g.updateCharacterPosition()
}

func (g *Game) updateCharacterPosition() {
	x, _ := ebiten.WindowSize()
	if g.curAction.name == "walkRight" {
		if g.windowPos.X+x > 1820 {
			g.curAction = g.walkLeftAction
			g.windowPos.X -= 1
		} else {
			g.windowPos.X += 1
		}
	} else if g.curAction.name == "walkLeft" {
		if g.windowPos.X-x < 100 {
			g.curAction = g.walkRightAction
			g.windowPos.X += 1
		} else {
			g.windowPos.X -= 1
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (windowWidth, windowHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) moveWindow() {
	cursorX, cursorY := ebiten.CursorPosition()
	cursorPos := Point{X: cursorX, Y: cursorY}
	g.curAction = g.gettingDraggedAction

	initalPress := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if initalPress {
		g.initalLeftPress = cursorPos
	}

	currX, currY := ebiten.WindowPosition()
	g.windowPos.X = (cursorPos.X + currX) - g.initalLeftPress.X
	g.windowPos.Y = (cursorPos.Y + currY) - g.initalLeftPress.Y

	if !(g.windowPos.X < 0) {
		return
	}

	monitors := ebiten.AppendMonitors(nil)
	curr := ebiten.Monitor()
	prevIndex := -1

	for i, monitor := range monitors {
		if monitor == curr {
			prevIndex = i - 1
			break
		}
	}

	if prevIndex >= 0 {
		ebiten.SetMonitor(monitors[prevIndex])
	}
}
