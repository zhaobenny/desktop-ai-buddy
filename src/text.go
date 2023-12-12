package main

import (
	"image"
	"strings"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/basicfont"
)

func addLineBreaks(text string, maxLineLength int) string {
	var lines []string

	words := strings.Fields(text)
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+len(word)+1 <= maxLineLength {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	lines = append(lines, currentLine)

	return strings.Join(lines, "\n")
}

func create_text_bubble(text string) image.Image {
	dc := gg.NewContext(1, 1) // Start with a small context, will be resized later

	dc.SetFontFace(basicfont.Face7x13)

	// add a new line if the text is too long and doesn't have any already
	if !strings.Contains(text, "\n") {
		text = addLineBreaks(text, 20)
	}

	lines := strings.Split(text, "\n")

	var textWidth, textHeight float64
	for _, line := range lines {
		width, height := dc.MeasureString(line)
		if width > textWidth {
			textWidth = width
		}
		textHeight += height
	}

	bubbleWidth := textWidth + 35
	if bubbleWidth < 80 {
		bubbleWidth = 80
	}
	bubbleHeight := textHeight + 20
	totalHeight := int(bubbleHeight + 20)
	totalWidth := int(bubbleWidth + 10)

	dc = gg.NewContext(totalWidth, totalHeight)

	dc.SetRGBA(255, 255, 255, 0)
	dc.Clear()

	x, y, r := 5, 0, 20
	dc.SetRGB(1, 1, 1)
	dc.DrawRoundedRectangle(float64(x), float64(y), bubbleWidth, bubbleHeight, float64(r))
	dc.Fill()

	dc.SetRGB(1, 1, 1)
	dc.MoveTo(bubbleWidth/5, bubbleHeight-5)
	if totalHeight > 100 {
		dc.LineTo(50, float64(totalHeight/2))
	} else {
		dc.LineTo(0, float64(totalHeight))

	}
	dc.SetLineWidth(5)
	dc.Stroke()

	dc.SetRGB(0, 0, 0)

	currentY := float64(y + 20)
	for _, line := range lines {
		dc.DrawStringAnchored(line, float64(x+20), currentY, 0.0, 0.0)
		currentY += dc.FontHeight() // Move to the next line
	}

	return dc.Image()
}
