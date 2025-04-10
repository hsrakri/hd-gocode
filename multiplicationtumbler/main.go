package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 600
	fontSize     = 40
)

type Player struct {
	name  string
	score int32
	time  float32
}

type GameState int

const (
	StateNameInput GameState = iota
	StatePlaying
	StateGameOver
	StateLeaderboard
)

type Game struct {
	number1     int32
	number2     int32
	result      int32
	userInput   string
	score       int32
	timeLeft    float32
	state       GameState
	message     string
	messageTime float32
	playerName  string
	leaderboard []Player
}

func NewGame() *Game {
	return &Game{
		timeLeft:    60.0,
		message:     "Enter your name and press ENTER to start!",
		state:       StateNameInput,
		leaderboard: make([]Player, 0),
	}
}

func (g *Game) generateNewProblem() {
	g.number1 = int32(rand.Intn(12) + 1)
	g.number2 = int32(rand.Intn(12) + 1)
	g.result = g.number1 * g.number2

	// Store the hidden number before swapping
	hiddenNumber := g.number2

	// Randomly decide which number to hide (1 or 2)
	if rand.Float32() < 0.5 {
		// If we're hiding number1, swap the numbers
		g.number1, g.number2 = g.number2, g.number1
		hiddenNumber = g.number1
	}

	// Store the visible number for comparison
	visibleNumber := g.number1
	if rand.Float32() < 0.5 {
		visibleNumber = g.number2
		g.number1, g.number2 = g.number2, g.number1
	}

	// Set up the problem
	g.number2 = hiddenNumber  // This is the number player needs to guess
	g.number1 = visibleNumber // This is the number shown
	g.result = g.number1 * hiddenNumber
}

func (g *Game) handleInput() {
	key := rl.GetCharPressed()

	switch g.state {
	case StateNameInput:
		// Handle name input
		if key >= 32 && key <= 125 && len(g.playerName) < 15 { // Allow most printable characters
			g.playerName += string(key)
		}
		if rl.IsKeyPressed(rl.KeyBackspace) && len(g.playerName) > 0 {
			g.playerName = g.playerName[:len(g.playerName)-1]
		}
		if rl.IsKeyPressed(rl.KeyEnter) && len(g.playerName) > 0 {
			g.state = StatePlaying
			g.generateNewProblem()
		}

	case StatePlaying:
		// Handle game input
		if key >= '0' && key <= '9' && len(g.userInput) < 3 {
			g.userInput += string(key)
		}
		if rl.IsKeyPressed(rl.KeyBackspace) && len(g.userInput) > 0 {
			g.userInput = g.userInput[:len(g.userInput)-1]
		}
		if rl.IsKeyPressed(rl.KeyEnter) && len(g.userInput) > 0 {
			guess := 0
			fmt.Sscanf(g.userInput, "%d", &guess)
			if int32(guess) == g.number2 {
				g.score++
				g.message = "Correct!"
				g.messageTime = 1.0
			} else {
				g.message = fmt.Sprintf("Wrong! It was %d (%d × %d = %d)",
					g.number2, g.number1, g.number2, g.result)
				g.messageTime = 1.0
			}
			g.userInput = ""
			g.generateNewProblem()
		}

	case StateGameOver:
		if rl.IsKeyPressed(rl.KeySpace) {
			// Reset for new game
			g.state = StateNameInput
			g.playerName = ""
			g.score = 0
			g.timeLeft = 60.0
			g.userInput = ""
			g.message = "Enter your name and press ENTER to start!"
		}

	case StateLeaderboard:
		if rl.IsKeyPressed(rl.KeySpace) {
			g.state = StateNameInput
			g.playerName = ""
			g.score = 0
			g.timeLeft = 60.0
			g.userInput = ""
			g.message = "Enter your name and press ENTER to start!"
		}
	}
}

func (g *Game) update() {
	if g.state != StatePlaying {
		return
	}

	g.timeLeft -= rl.GetFrameTime()
	if g.messageTime > 0 {
		g.messageTime -= rl.GetFrameTime()
	}

	if g.timeLeft <= 0 {
		// Add player to leaderboard
		g.leaderboard = append(g.leaderboard, Player{
			name:  g.playerName,
			score: g.score,
			time:  60.0 - g.timeLeft,
		})

		// Sort leaderboard by score (higher is better)
		sort.Slice(g.leaderboard, func(i, j int) bool {
			return g.leaderboard[i].score > g.leaderboard[j].score
		})

		// Keep only top 10 scores
		if len(g.leaderboard) > 10 {
			g.leaderboard = g.leaderboard[:10]
		}

		g.state = StateLeaderboard
		g.message = "Game Over! Press SPACE to play again"
	}
}

func (g *Game) draw() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	switch g.state {
	case StateNameInput:
		namePrompt := "Enter your name:"
		rl.DrawText(namePrompt, screenWidth/2-rl.MeasureText(namePrompt, fontSize)/2, screenHeight/3, fontSize, rl.Black)

		nameText := g.playerName + "_"
		rl.DrawText(nameText, screenWidth/2-rl.MeasureText(nameText, fontSize)/2, screenHeight/2, fontSize, rl.DarkGray)

		instructions := "Press ENTER when ready"
		rl.DrawText(instructions, screenWidth/2-rl.MeasureText(instructions, 20)/2, screenHeight*2/3, 20, rl.DarkGray)

	case StatePlaying:
		// Draw timer
		timerText := fmt.Sprintf("Time: %.1f", g.timeLeft)
		rl.DrawText(timerText, 10, 10, 30, rl.DarkGray)

		// Draw score
		scoreText := fmt.Sprintf("Score: %d", g.score)
		rl.DrawText(scoreText, screenWidth-200, 10, 30, rl.DarkGray)

		// Draw player name
		rl.DrawText(g.playerName, 10, 50, 20, rl.DarkGray)

		// Draw problem
		problemText := fmt.Sprintf("%d × ? = %d", g.number1, g.result)
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

	case StateLeaderboard:
		title := "LEADERBOARD"
		rl.DrawText(title, screenWidth/2-rl.MeasureText(title, fontSize)/2, 50, fontSize, rl.Black)

		// Draw the leaderboard entries
		startY := 150
		for i, player := range g.leaderboard {
			entry := fmt.Sprintf("%d. %s - Score: %d", i+1, player.name, player.score)
			rl.DrawText(entry, screenWidth/4, int32(startY+i*40), 30, rl.DarkGray)
		}

		// Draw restart instruction
		restart := "Press SPACE to play again"
		rl.DrawText(restart, screenWidth/2-rl.MeasureText(restart, 30)/2, screenHeight-100, 30, rl.DarkGray)
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
