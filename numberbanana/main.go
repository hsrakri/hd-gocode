package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 600
	fontSize     = 40
)

type NumberFact struct {
	Found bool   `json:"found"`
	Text  string `json:"text"`
}

type Game struct {
	userInput    string
	currentFact  string
	isSearching  bool
	errorMessage string
	facts        []string
}

func NewGame() *Game {
	return &Game{
		facts: make([]string, 0),
	}
}

func (g *Game) fetchNumberFact(number string) {
	g.isSearching = true
	g.errorMessage = ""

	// Try different APIs for number facts
	apis := []string{
		"http://numbersapi.com/%s/math",
		"http://numbersapi.com/%s/trivia",
		"http://numbersapi.com/%s/year",
	}

	for _, apiURL := range apis {
		url := fmt.Sprintf(apiURL, number)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		fact := string(body)
		if fact != "" && !strings.Contains(fact.ToLower(), "error") {
			g.facts = append(g.facts, fact)
		}
	}

	if len(g.facts) == 0 {
		g.errorMessage = "No interesting facts found for this number"
	}

	g.isSearching = false
}

func (g *Game) handleInput() {
	key := rl.GetCharPressed()

	// Handle number input
	if key >= '0' && key <= '9' && len(g.userInput) < 10 {
		g.userInput += string(key)
	}

	// Handle backspace
	if rl.IsKeyPressed(rl.KeyBackspace) && len(g.userInput) > 0 {
		g.userInput = g.userInput[:len(g.userInput)-1]
	}

	// Handle enter
	if rl.IsKeyPressed(rl.KeyEnter) && len(g.userInput) > 0 {
		go g.fetchNumberFact(g.userInput)
		g.userInput = ""
	}
}

func (g *Game) draw() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	// Draw title
	title := "Number Banana"
	rl.DrawText(title, screenWidth/2-rl.MeasureText(title, fontSize)/2, 50, fontSize, rl.Black)

	// Draw input prompt
	prompt := "Enter a number:"
	rl.DrawText(prompt, screenWidth/2-rl.MeasureText(prompt, 30)/2, 150, 30, rl.DarkGray)

	// Draw user input
	input := g.userInput + "_"
	rl.DrawText(input, screenWidth/2-rl.MeasureText(input, 40)/2, 200, 40, rl.Black)

	// Draw loading message
	if g.isSearching {
		loading := "Searching for interesting facts..."
		rl.DrawText(loading, screenWidth/2-rl.MeasureText(loading, 20)/2, 300, 20, rl.DarkGray)
	}

	// Draw error message
	if g.errorMessage != "" {
		rl.DrawText(g.errorMessage, screenWidth/2-rl.MeasureText(g.errorMessage, 20)/2, 300, 20, rl.Red)
	}

	// Draw facts
	startY := 300
	for i, fact := range g.facts {
		// Word wrap the fact
		words := strings.Split(fact, " ")
		line := ""
		y := startY + i*60

		for _, word := range words {
			testLine := line + " " + word
			if rl.MeasureText(testLine, 20) > screenWidth-100 {
				rl.DrawText(line, 50, int32(y), 20, rl.DarkGray)
				y += 25
				line = word
			} else {
				if line == "" {
					line = word
				} else {
					line += " " + word
				}
			}
		}
		if line != "" {
			rl.DrawText(line, 50, int32(y), 20, rl.DarkGray)
		}
	}

	// Draw instructions
	instructions := "Press ENTER to search for facts"
	rl.DrawText(instructions, screenWidth/2-rl.MeasureText(instructions, 20)/2, screenHeight-40, 20, rl.DarkGray)
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Number Banana")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	game := NewGame()

	for !rl.WindowShouldClose() {
		game.handleInput()
		game.draw()
	}
}
