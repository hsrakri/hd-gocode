package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/hsrakri/hd-gocode/equation_explorer/equation"
	"github.com/hsrakri/hd-gocode/equation_explorer/maze"
)

const (
	screenWidth  = 800
	screenHeight = 600
	cellSize     = 40
)

// Game states
const (
	StateMenu = iota
	StatePlaying
	StateSolving
	StateGameOver
)

// Game represents the main game state
type Game struct {
	State       int
	Score       int
	Level       int
	Player      Player
	Maze        *maze.Maze
	Equation    *equation.Equation
	Font        rl.Font
	Input       string
	ShowHint    bool
	ShowSteps   bool
	Message     string
	MessageTime float32
	StartTime   float64
	TotalTime   float64
}

// Player represents the player character
type Player struct {
	X, Y  int
	Speed float32
	Color rl.Color
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Equation Explorer")
	defer rl.CloseWindow()

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create game instance
	game := NewGame()

	// Set target FPS
	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		game.Update()
		game.Draw()
	}
}

// NewGame creates a new game instance
func NewGame() *Game {
	game := &Game{
		State:     StateMenu,
		Score:     0,
		Level:     1,
		StartTime: rl.GetTime(),
		TotalTime: 0,
		Player: Player{
			X:     1,
			Y:     1,
			Speed: 5,
			Color: rl.Blue,
		},
		Input:       "",
		ShowHint:    false,
		ShowSteps:   false,
		Message:     "",
		MessageTime: 0,
	}

	// Generate initial maze
	game.Maze = maze.NewMaze(screenWidth/cellSize, screenHeight/cellSize)
	game.Maze.Generate()

	return game
}

// Update handles game logic updates
func (g *Game) Update() {
	switch g.State {
	case StateMenu:
		g.UpdateMenu()
	case StatePlaying:
		g.UpdatePlaying()
	case StateSolving:
		g.UpdateSolving()
	case StateGameOver:
		g.UpdateGameOver()
	}

	// Update total time
	if g.State == StatePlaying || g.State == StateSolving {
		g.TotalTime = rl.GetTime() - g.StartTime
	}

	// Update message timer
	if g.MessageTime > 0 {
		g.MessageTime -= rl.GetFrameTime()
		if g.MessageTime <= 0 {
			g.Message = ""
		}
	}
}

// Draw handles game rendering
func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	switch g.State {
	case StateMenu:
		g.DrawMenu()
	case StatePlaying:
		g.DrawPlaying()
	case StateSolving:
		g.DrawSolving()
	case StateGameOver:
		g.DrawGameOver()
	}

	// Draw message if any
	if g.Message != "" {
		messageWidth := rl.MeasureTextEx(g.Font, g.Message, 20, 2).X
		rl.DrawTextEx(g.Font, g.Message, rl.Vector2{
			X: float32(screenWidth/2 - int(messageWidth/2)),
			Y: float32(screenHeight - 40),
		}, 20, 2, rl.Red)
	}

	rl.EndDrawing()
}

// GenerateMaze creates a new maze
func (g *Game) GenerateMaze() {
	g.Maze = maze.NewMaze(screenWidth/cellSize, screenHeight/cellSize)
	g.Maze.Generate()
}

// GenerateEquation creates a new equation based on current level
func (g *Game) GenerateEquation() {
	g.Equation = equation.GenerateEquation(g.Level)
}

// UpdateMenu handles menu state updates
func (g *Game) UpdateMenu() {
	if rl.IsKeyPressed(rl.KeySpace) {
		g.State = StatePlaying
		g.Reset()
	}
}

// UpdatePlaying handles playing state updates
func (g *Game) UpdatePlaying() {
	// Handle player movement
	newX, newY := g.Player.X, g.Player.Y

	if rl.IsKeyDown(rl.KeyRight) {
		newX++
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		newX--
	}
	if rl.IsKeyDown(rl.KeyDown) {
		newY++
	}
	if rl.IsKeyDown(rl.KeyUp) {
		newY--
	}

	// Check if the new position is valid
	if !g.Maze.IsWall(newX, newY) {
		g.Player.X, g.Player.Y = newX, newY

		// Check if player reached an equation cell
		if g.Maze.HasEquation(newX, newY) {
			g.State = StateSolving
			g.GenerateEquation()
			g.Input = ""
			g.ShowHint = false
			g.ShowSteps = false
		}
	}
}

// UpdateSolving handles equation solving state updates
func (g *Game) UpdateSolving() {
	// Handle input
	if rl.IsKeyPressed(rl.KeyBackspace) {
		if len(g.Input) > 0 {
			g.Input = g.Input[:len(g.Input)-1]
		}
	} else if rl.IsKeyPressed(rl.KeyEnter) {
		// Check answer
		if answer, err := strconv.ParseFloat(g.Input, 64); err == nil {
			if math.Abs(answer-g.Equation.Answer) < 0.01 {
				g.Score += 100
				g.ShowMessage("Correct! +100 points", 2)
				g.State = StatePlaying
				// Remove equation from maze
				g.Maze.GetCell(g.Player.X, g.Player.Y).Equation = false
			} else {
				g.ShowMessage("Incorrect! Try again", 1)
			}
		}
	} else if rl.IsKeyPressed(rl.KeyH) {
		g.ShowHint = !g.ShowHint
	} else if rl.IsKeyPressed(rl.KeyS) {
		g.ShowSteps = !g.ShowSteps
	} else if rl.IsKeyPressed(rl.KeyEscape) {
		g.State = StatePlaying
	} else {
		// Add character to input if it's printable
		key := rl.GetCharPressed()
		if key >= 32 && key <= 126 {
			g.Input += string(key)
		}
	}
}

// UpdateGameOver handles game over state updates
func (g *Game) UpdateGameOver() {
	if rl.IsKeyPressed(rl.KeySpace) {
		g.State = StateMenu
		g.Reset()
	}
}

// DrawMenu renders the menu screen
func (g *Game) DrawMenu() {
	title := "Equation Explorer"
	subtitle := "Press SPACE to Start"

	titleWidth := rl.MeasureTextEx(g.Font, title, 40, 2).X
	subtitleWidth := rl.MeasureTextEx(g.Font, subtitle, 20, 2).X

	rl.DrawTextEx(g.Font, title, rl.Vector2{
		X: float32(screenWidth/2 - int(titleWidth/2)),
		Y: float32(screenHeight/2 - 50),
	}, 40, 2, rl.Black)

	rl.DrawTextEx(g.Font, subtitle, rl.Vector2{
		X: float32(screenWidth/2 - int(subtitleWidth/2)),
		Y: float32(screenHeight/2 + 50),
	}, 20, 2, rl.Gray)
}

// DrawPlaying renders the game screen
func (g *Game) DrawPlaying() {
	// Draw maze
	for y := range g.Maze.Cells {
		for x := range g.Maze.Cells[y] {
			cell := g.Maze.GetCell(x, y)
			if cell.Wall {
				rl.DrawRectangle(
					int32(x*cellSize),
					int32(y*cellSize),
					int32(cellSize),
					int32(cellSize),
					rl.Gray,
				)
			} else if cell.Equation {
				rl.DrawRectangle(
					int32(x*cellSize),
					int32(y*cellSize),
					int32(cellSize),
					int32(cellSize),
					rl.Green,
				)
			}
		}
	}

	// Draw player
	rl.DrawRectangle(
		int32(g.Player.X*cellSize),
		int32(g.Player.Y*cellSize),
		int32(cellSize),
		int32(cellSize),
		g.Player.Color,
	)

	// Draw score, level, and time
	scoreText := fmt.Sprintf("Score: %d", g.Score)
	levelText := fmt.Sprintf("Level: %d", g.Level)
	timeText := fmt.Sprintf("Time: %.1f s", g.TotalTime)
	rl.DrawTextEx(g.Font, scoreText, rl.Vector2{X: 10, Y: 10}, 20, 2, rl.Black)
	rl.DrawTextEx(g.Font, levelText, rl.Vector2{X: 10, Y: 40}, 20, 2, rl.Black)
	rl.DrawTextEx(g.Font, timeText, rl.Vector2{X: 10, Y: 70}, 20, 2, rl.Black)
}

// DrawSolving renders the equation solving screen
func (g *Game) DrawSolving() {
	// Draw equation
	eqText := fmt.Sprintf("Solve: %s", g.Equation.Text)
	eqWidth := rl.MeasureTextEx(g.Font, eqText, 30, 2).X
	rl.DrawTextEx(g.Font, eqText, rl.Vector2{
		X: float32(screenWidth/2 - int(eqWidth/2)),
		Y: float32(screenHeight/2 - 100),
	}, 30, 2, rl.Black)

	// Draw input
	inputText := fmt.Sprintf("Answer: %s", g.Input)
	inputWidth := rl.MeasureTextEx(g.Font, inputText, 30, 2).X
	rl.DrawTextEx(g.Font, inputText, rl.Vector2{
		X: float32(screenWidth/2 - int(inputWidth/2)),
		Y: float32(screenHeight / 2),
	}, 30, 2, rl.Black)

	// Draw hint if requested
	if g.ShowHint {
		hintWidth := rl.MeasureTextEx(g.Font, g.Equation.Hint, 20, 2).X
		rl.DrawTextEx(g.Font, g.Equation.Hint, rl.Vector2{
			X: float32(screenWidth/2 - int(hintWidth/2)),
			Y: float32(screenHeight/2 + 50),
		}, 20, 2, rl.Gray)
	}

	// Draw steps if requested
	if g.ShowSteps {
		y := screenHeight/2 + 100
		for _, step := range g.Equation.Steps {
			stepWidth := rl.MeasureTextEx(g.Font, step, 20, 2).X
			rl.DrawTextEx(g.Font, step, rl.Vector2{
				X: float32(screenWidth/2 - int(stepWidth/2)),
				Y: float32(y),
			}, 20, 2, rl.Gray)
			y += 30
		}
	}

	// Draw controls
	controls := []string{
		"Press H for hint",
		"Press S for solution steps",
		"Press ESC to skip",
		"Press ENTER to submit answer",
	}

	y := screenHeight - 150
	for _, control := range controls {
		controlWidth := rl.MeasureTextEx(g.Font, control, 20, 2).X
		rl.DrawTextEx(g.Font, control, rl.Vector2{
			X: float32(screenWidth/2 - int(controlWidth/2)),
			Y: float32(y),
		}, 20, 2, rl.Gray)
		y += 30
	}
}

// DrawGameOver renders the game over screen
func (g *Game) DrawGameOver() {
	title := "Game Over"
	score := fmt.Sprintf("Final Score: %d", g.Score)
	time := fmt.Sprintf("Total Time: %.1f seconds", g.TotalTime)
	restart := "Press SPACE to Restart"

	titleWidth := rl.MeasureTextEx(g.Font, title, 40, 2).X
	scoreWidth := rl.MeasureTextEx(g.Font, score, 30, 2).X
	timeWidth := rl.MeasureTextEx(g.Font, time, 30, 2).X
	restartWidth := rl.MeasureTextEx(g.Font, restart, 20, 2).X

	rl.DrawTextEx(g.Font, title, rl.Vector2{
		X: float32(screenWidth/2 - int(titleWidth/2)),
		Y: float32(screenHeight/2 - 100),
	}, 40, 2, rl.Black)

	rl.DrawTextEx(g.Font, score, rl.Vector2{
		X: float32(screenWidth/2 - int(scoreWidth/2)),
		Y: float32(screenHeight/2 - 30),
	}, 30, 2, rl.Black)

	rl.DrawTextEx(g.Font, time, rl.Vector2{
		X: float32(screenWidth/2 - int(timeWidth/2)),
		Y: float32(screenHeight/2 + 30),
	}, 30, 2, rl.Black)

	rl.DrawTextEx(g.Font, restart, rl.Vector2{
		X: float32(screenWidth/2 - int(restartWidth/2)),
		Y: float32(screenHeight/2 + 100),
	}, 20, 2, rl.Gray)
}

// Reset resets the game state
func (g *Game) Reset() {
	g.Score = 0
	g.Level = 1
	g.Player.X = 1
	g.Player.Y = 1
	g.StartTime = rl.GetTime()
	g.TotalTime = 0
	g.GenerateMaze()
	g.Input = ""
	g.ShowHint = false
	g.ShowSteps = false
	g.Message = ""
}

// ShowMessage displays a message for a specified duration
func (g *Game) ShowMessage(text string, duration float32) {
	g.Message = text
	g.MessageTime = duration
}
