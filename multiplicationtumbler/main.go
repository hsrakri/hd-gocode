package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 600
	fontSize     = 40
)

type Game struct {
	number1     int32
	number2     int32
	result      int32
	userInput   string
	score       int32
	timeLeft    float32
	gameStarted bool
	message     string
	messageTime float32
}

func NewGame() *Game {
	return &Game{
		timeLeft: 60.0,
		message:  "Press ENTER to start!",
	}
}

func (g *Game) generateNewProblem() {
	g.number1 = int32(rand.Intn(12) + 1)
	g.number2 = int32(rand.Intn(12) + 1)
	g.result = g.number1 * g.number2
	// Randomly decide which number to hide (1 or 2)
	if rand.Float32() < 0.5 {
		g.number1, g.number2 = g.number2, g.number1 // Swap numbers
	}
	g.number2 = -1 // Hide the second number
}

func (g *Game) handleInput() {
	key := rl.GetCharPressed()
	if key >= '0' && key <= '9' && len(g.userInput) < 3 {
		g.userInput += string(key)
	}
	if rl.IsKeyPressed(rl.KeyBackspace) && len(g.userInput) > 0 {
		g.userInput = g.userInput[:len(g.userInput)-1]
	}
	if rl.IsKeyPressed(rl.KeyEnter) {
		if !g.gameStarted {
			g.gameStarted = true
			g.generateNewProblem()
			return
		}
		if len(g.userInput) > 0 {
			guess := 0
			fmt.Sscanf(g.userInput, "%d", &guess)
			if int32(guess) == g.number2 {
				g.score++
				g.message = "Correct!"
				g.messageTime = 1.0
			} else {
				g.message = fmt.Sprintf("Wrong! It was %d", g.number2)
				g.messageTime = 1.0
			}
			g.userInput = ""
			g.generateNewProblem()
		}
	}
}

func (g *Game) update() {
	if !g.gameStarted {
		return
	}

	g.timeLeft -= rl.GetFrameTime()
	if g.messageTime > 0 {
		g.messageTime -= rl.GetFrameTime()
	}

	if g.timeLeft <= 0 {
		g.gameStarted = false
		g.message = fmt.Sprintf("Game Over! Score: %d - Press ENTER to restart", g.score)
		g.timeLeft = 60.0
		g.score = 0
		g.userInput = ""
	}
}

func (g *Game) draw() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	if !g.gameStarted {
		rl.DrawText(g.message, screenWidth/2-rl.MeasureText(g.message, fontSize)/2, screenHeight/2, fontSize, rl.Black)
		return
	}

	// Draw timer
	timerText := fmt.Sprintf("Time: %.1f", g.timeLeft)
	rl.DrawText(timerText, 10, 10, 30, rl.DarkGray)

	// Draw score
	scoreText := fmt.Sprintf("Score: %d", g.score)
	rl.DrawText(scoreText, screenWidth-200, 10, 30, rl.DarkGray)

	// Draw problem
	var problemText string
	if g.number2 == -1 {
		problemText = fmt.Sprintf("%d × ? = %d", g.number1, g.result)
	} else {
		problemText = fmt.Sprintf("? × %d = %d", g.number2, g.result)
	}
	rl.DrawText(problemText, screenWidth/2-rl.MeasureText(problemText, fontSize)/2, screenHeight/2-50, fontSize, rl.Black)

	// Draw user input
	inputText := fmt.Sprintf("Your answer: %s", g.userInput)
	rl.DrawText(inputText, screenWidth/2-rl.MeasureText(inputText, 30)/2, screenHeight/2+50, 30, rl.DarkGray)

	// Draw message
	if g.messageTime > 0 {
		messageColor := rl.DarkGray
		if g.message == "Correct!" {
			messageColor = rl.Green
		} else if g.message[:5] == "Wrong!" {
			messageColor = rl.Red
		}
		rl.DrawText(g.message, screenWidth/2-rl.MeasureText(g.message, 30)/2, screenHeight/2+100, 30, messageColor)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	rl.InitWindow(screenWidth, screenHeight, "Multiplication Tumbler")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	game := NewGame()

	for !rl.WindowShouldClose() {
		game.handleInput()
		game.update()
		game.draw()
	}
}
