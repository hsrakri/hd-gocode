package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 600
	cellSize     = 40
)

type Shape struct {
	Points []rl.Vector2
	Angle  float32
}

type Question struct {
	Shape     Shape
	Question  string
	Answer    float32
	Attempts  int
	StartTime time.Time
}

type Player struct {
	Name     string
	Score    int
	Time     time.Duration
	Finished bool
}

type Game struct {
	State       int
	Player      Player
	Questions   []Question
	CurrentQ    int
	Input       string
	Message     string
	MessageTime float32
	Leaderboard []Player
}

const (
	StateMenu      = 0
	StateNameInput = 1
	StatePlaying   = 2
	StateAnswer    = 3
	StateGameOver  = 4
)

func NewGame() *Game {
	game := &Game{
		State:  StateMenu,
		Player: Player{},
		Questions: []Question{
			{
				Shape: Shape{
					Points: []rl.Vector2{
						{X: 400, Y: 300},
						{X: 500, Y: 300},
						{X: 450, Y: 200},
					},
					Angle: 45,
				},
				Question: "What is the angle at the top vertex?",
				Answer:   45,
			},
			{
				Shape: Shape{
					Points: []rl.Vector2{
						{X: 400, Y: 300},
						{X: 500, Y: 300},
						{X: 450, Y: 400},
					},
					Angle: 135,
				},
				Question: "What is the angle at the bottom vertex?",
				Answer:   135,
			},
			{
				Shape: Shape{
					Points: []rl.Vector2{
						{X: 400, Y: 300},
						{X: 500, Y: 300},
						{X: 450, Y: 250},
					},
					Angle: 30,
				},
				Question: "What is the angle at the top vertex?",
				Answer:   30,
			},
			{
				Shape: Shape{
					Points: []rl.Vector2{
						{X: 400, Y: 300},
						{X: 500, Y: 300},
						{X: 450, Y: 350},
					},
					Angle: 150,
				},
				Question: "What is the angle at the bottom vertex?",
				Answer:   150,
			},
			{
				Shape: Shape{
					Points: []rl.Vector2{
						{X: 400, Y: 300},
						{X: 500, Y: 300},
						{X: 450, Y: 280},
					},
					Angle: 15,
				},
				Question: "What is the angle at the top vertex?",
				Answer:   15,
			},
			{
				Shape: Shape{
					Points: []rl.Vector2{
						{X: 400, Y: 300},
						{X: 500, Y: 300},
						{X: 450, Y: 320},
					},
					Angle: 165,
				},
				Question: "What is the angle at the bottom vertex?",
				Answer:   165,
			},
			{
				Shape: Shape{
					Points: []rl.Vector2{
						{X: 400, Y: 300},
						{X: 500, Y: 300},
						{X: 450, Y: 270},
					},
					Angle: 20,
				},
				Question: "What is the angle at the top vertex?",
				Answer:   20,
			},
			{
				Shape: Shape{
					Points: []rl.Vector2{
						{X: 400, Y: 300},
						{X: 500, Y: 300},
						{X: 450, Y: 330},
					},
					Angle: 160,
				},
				Question: "What is the angle at the bottom vertex?",
				Answer:   160,
			},
			{
				Shape: Shape{
					Points: []rl.Vector2{
						{X: 400, Y: 300},
						{X: 500, Y: 300},
						{X: 450, Y: 260},
					},
					Angle: 25,
				},
				Question: "What is the angle at the top vertex?",
				Answer:   25,
			},
			{
				Shape: Shape{
					Points: []rl.Vector2{
						{X: 400, Y: 300},
						{X: 500, Y: 300},
						{X: 450, Y: 340},
					},
					Angle: 155,
				},
				Question: "What is the angle at the bottom vertex?",
				Answer:   155,
			},
		},
		CurrentQ:    0,
		Input:       "",
		Message:     "",
		MessageTime: 0,
	}

	// Load leaderboard
	game.LoadLeaderboard()

	return game
}

func (g *Game) LoadLeaderboard() {
	// TODO: Load from file
	g.Leaderboard = []Player{
		{Name: "Player 1", Score: 8, Time: time.Minute * 2},
		{Name: "Player 2", Score: 7, Time: time.Minute * 3},
	}
}

func (g *Game) SaveLeaderboard() {
	// TODO: Save to file
}

func (g *Game) Update() {
	switch g.State {
	case StateMenu:
		g.UpdateMenu()
	case StateNameInput:
		g.UpdateNameInput()
	case StatePlaying:
		g.UpdatePlaying()
	case StateAnswer:
		g.UpdateAnswer()
	case StateGameOver:
		g.UpdateGameOver()
	}

	// Update message timer
	if g.MessageTime > 0 {
		g.MessageTime -= rl.GetFrameTime()
		if g.MessageTime <= 0 {
			g.Message = ""
		}
	}
}

func (g *Game) UpdateMenu() {
	if rl.IsKeyPressed(rl.KeySpace) {
		g.State = StateNameInput
		g.Input = ""
	}
}

func (g *Game) UpdateNameInput() {
	key := rl.GetCharPressed()
	for key > 0 {
		if key >= 32 && key <= 125 {
			g.Input += string(key)
		}
		key = rl.GetCharPressed()
	}

	if rl.IsKeyPressed(rl.KeyBackspace) {
		if len(g.Input) > 0 {
			g.Input = g.Input[:len(g.Input)-1]
		}
	}

	if rl.IsKeyPressed(rl.KeyEnter) && len(g.Input) > 0 {
		g.Player.Name = g.Input
		g.State = StatePlaying
		g.Player.Time = 0
		g.CurrentQ = 0
		g.Questions[0].StartTime = time.Now()
	}
}

func (g *Game) UpdatePlaying() {
	// Update time
	g.Player.Time = time.Since(g.Questions[0].StartTime)

	// Check if all questions are answered
	if g.CurrentQ >= len(g.Questions) {
		g.State = StateGameOver
		g.Player.Finished = true
		return
	}

	// Get current question
	q := &g.Questions[g.CurrentQ]

	// Handle input
	key := rl.GetCharPressed()
	for key > 0 {
		if key >= 48 && key <= 57 || key == 46 {
			g.Input += string(key)
		}
		key = rl.GetCharPressed()
	}

	if rl.IsKeyPressed(rl.KeyBackspace) {
		if len(g.Input) > 0 {
			g.Input = g.Input[:len(g.Input)-1]
		}
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		if len(g.Input) > 0 {
			var answer float32
			fmt.Sscanf(g.Input, "%f", &answer)
			q.Attempts++

			if math.Abs(float64(answer-float32(q.Answer))) <= 5 {
				g.Player.Score++
				g.Message = "Correct! Within 5 degrees."
				g.MessageTime = 2
				g.Input = ""
				g.CurrentQ++
				if g.CurrentQ < len(g.Questions) {
					g.Questions[g.CurrentQ].StartTime = time.Now()
				}
			} else if q.Attempts >= 3 {
				g.Message = fmt.Sprintf("The correct answer was %.1f degrees", q.Answer)
				g.MessageTime = 2
				g.Input = ""
				g.CurrentQ++
				if g.CurrentQ < len(g.Questions) {
					g.Questions[g.CurrentQ].StartTime = time.Now()
				}
			} else {
				g.Message = fmt.Sprintf("Try again! (%d/3 attempts)", q.Attempts)
				g.MessageTime = 1
				g.Input = ""
			}
		}
	}
}

func (g *Game) UpdateAnswer() {
	if rl.IsKeyPressed(rl.KeyEnter) {
		g.State = StatePlaying
		g.Input = ""
	}
}

func (g *Game) UpdateGameOver() {
	if rl.IsKeyPressed(rl.KeySpace) {
		// Add player to leaderboard
		g.Leaderboard = append(g.Leaderboard, g.Player)
		// Sort by score (descending) and time (ascending)
		sort.Slice(g.Leaderboard, func(i, j int) bool {
			if g.Leaderboard[i].Score == g.Leaderboard[j].Score {
				return g.Leaderboard[i].Time < g.Leaderboard[j].Time
			}
			return g.Leaderboard[i].Score > g.Leaderboard[j].Score
		})
		// Keep only top 10
		if len(g.Leaderboard) > 10 {
			g.Leaderboard = g.Leaderboard[:10]
		}
		// Save leaderboard
		g.SaveLeaderboard()
		// Reset game
		g.State = StateMenu
		g.Player = Player{}
		g.Input = ""
	}
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	switch g.State {
	case StateMenu:
		g.DrawMenu()
	case StateNameInput:
		g.DrawNameInput()
	case StatePlaying:
		g.DrawPlaying()
	case StateAnswer:
		g.DrawAnswer()
	case StateGameOver:
		g.DrawGameOver()
	}

	rl.EndDrawing()
}

func (g *Game) DrawMenu() {
	title := "Guess the Angle"
	subtitle := "Press SPACE to Start"

	titleWidth := rl.MeasureText(title, 40)
	subtitleWidth := rl.MeasureText(subtitle, 20)

	rl.DrawText(title,
		int32(screenWidth/2-int(titleWidth/2)),
		int32(screenHeight/2-50),
		40,
		rl.Black,
	)

	rl.DrawText(subtitle,
		int32(screenWidth/2-int(subtitleWidth/2)),
		int32(screenHeight/2+50),
		20,
		rl.Gray,
	)
}

func (g *Game) DrawNameInput() {
	prompt := "Enter your name:"
	promptWidth := rl.MeasureText(prompt, 20)

	rl.DrawText(prompt,
		int32(screenWidth/2-int(promptWidth/2)),
		int32(screenHeight/2-50),
		20,
		rl.Black,
	)

	rl.DrawText(g.Input+"_",
		int32(screenWidth/2-int(rl.MeasureText(g.Input+"_", 20)/2)),
		int32(screenHeight/2),
		20,
		rl.Black,
	)
}

func (g *Game) DrawPlaying() {
	// Draw shape
	q := g.Questions[g.CurrentQ]
	rl.DrawLineEx(q.Shape.Points[0], q.Shape.Points[1], 2, rl.Black)
	rl.DrawLineEx(q.Shape.Points[1], q.Shape.Points[2], 2, rl.Black)
	rl.DrawLineEx(q.Shape.Points[2], q.Shape.Points[0], 2, rl.Black)

	// Draw angle arc and highlight vertex
	vertex := q.Shape.Points[2]       // The vertex we want to measure
	rl.DrawCircleV(vertex, 5, rl.Red) // Red dot at the vertex

	// Calculate vectors for the angle
	v1 := rl.Vector2{
		X: q.Shape.Points[0].X - vertex.X,
		Y: q.Shape.Points[0].Y - vertex.Y,
	}
	v2 := rl.Vector2{
		X: q.Shape.Points[1].X - vertex.X,
		Y: q.Shape.Points[1].Y - vertex.Y,
	}

	// Calculate angle between vectors
	angle := float32(math.Atan2(float64(v2.Y), float64(v2.X)) - math.Atan2(float64(v1.Y), float64(v1.X)))
	if angle < 0 {
		angle += 2 * float32(math.Pi)
	}

	// Draw arc
	radius := float32(30)
	rl.DrawCircleLines(int32(vertex.X), int32(vertex.Y), radius, rl.Red)
	rl.DrawLineEx(vertex, rl.Vector2{
		X: vertex.X + radius*float32(math.Cos(float64(angle))),
		Y: vertex.Y + radius*float32(math.Sin(float64(angle))),
	}, 2, rl.Red)

	// Draw question
	questionWidth := rl.MeasureText(q.Question, 20)
	rl.DrawText(q.Question,
		int32(screenWidth/2-int(questionWidth/2)),
		50,
		20,
		rl.Black,
	)

	// Draw input
	rl.DrawText(g.Input+"_",
		int32(screenWidth/2-int(rl.MeasureText(g.Input+"_", 20)/2)),
		int32(screenHeight-100),
		20,
		rl.Black,
	)

	// Draw message
	if g.Message != "" {
		messageWidth := rl.MeasureText(g.Message, 20)
		rl.DrawText(g.Message,
			int32(screenWidth/2-int(messageWidth/2)),
			100,
			20,
			rl.Red,
		)
	}

	// Draw score and time
	scoreText := fmt.Sprintf("Score: %d/%d", g.Player.Score, g.CurrentQ+1)
	timeText := fmt.Sprintf("Time: %.1f s", g.Player.Time.Seconds())
	rl.DrawText(scoreText, 10, 10, 20, rl.Black)
	rl.DrawText(timeText, 10, 40, 20, rl.Black)
}

func (g *Game) DrawAnswer() {
	message := "Press ENTER to continue"
	messageWidth := rl.MeasureText(message, 20)
	rl.DrawText(message,
		int32(screenWidth/2-int(messageWidth/2)),
		int32(screenHeight/2),
		20,
		rl.Black,
	)
}

func (g *Game) DrawGameOver() {
	title := "Game Over"
	score := fmt.Sprintf("Final Score: %d/%d", g.Player.Score, len(g.Questions))
	time := fmt.Sprintf("Time: %.1f seconds", g.Player.Time.Seconds())
	restart := "Press SPACE to see leaderboard"

	titleWidth := rl.MeasureText(title, 40)
	scoreWidth := rl.MeasureText(score, 30)
	timeWidth := rl.MeasureText(time, 30)
	restartWidth := rl.MeasureText(restart, 20)

	rl.DrawText(title,
		int32(screenWidth/2-int(titleWidth/2)),
		int32(screenHeight/2-100),
		40,
		rl.Black,
	)

	rl.DrawText(score,
		int32(screenWidth/2-int(scoreWidth/2)),
		int32(screenHeight/2-30),
		30,
		rl.Black,
	)

	rl.DrawText(time,
		int32(screenWidth/2-int(timeWidth/2)),
		int32(screenHeight/2+30),
		30,
		rl.Black,
	)

	rl.DrawText(restart,
		int32(screenWidth/2-int(restartWidth/2)),
		int32(screenHeight/2+100),
		20,
		rl.Gray,
	)

	// Draw leaderboard
	y := 150
	rl.DrawText("Top Players:", 10, int32(y), 20, rl.Black)
	y += 30
	for i, p := range g.Leaderboard {
		if i >= 5 {
			break
		}
		entry := fmt.Sprintf("%d. %s - Score: %d, Time: %.1f s", i+1, p.Name, p.Score, p.Time.Seconds())
		rl.DrawText(entry, 10, int32(y), 20, rl.Black)
		y += 30
	}
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Guess the Angle")
	defer rl.CloseWindow()

	game := NewGame()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		game.Update()
		game.Draw()
	}
}
