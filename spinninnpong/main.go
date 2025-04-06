package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth      = 800
	screenHeight     = 450
	paddleWidth      = 10 // Reduced width back to original size
	paddleHeight     = 45 // Reduced height from 60 to 45
	ballSize         = 10
	winScore         = 15  // Score needed to win
	initialBallSpeed = 3.0 // Slower initial ball speed
)

type GameState int

const (
	Instructions GameState = iota
	Playing
	GameOver
)

type Paddle struct {
	X, Y     float32
	Speed    float32
	Angle    float32
	SpinRate float32
}

type Ball struct {
	X, Y    float32
	SpeedX  float32
	SpeedY  float32
	LastHit string // "left" or "right" to track who hit last
}

type Game struct {
	State       GameState
	LeftPaddle  Paddle
	RightPaddle Paddle
	Ball        Ball
	ScoreLeft   int
	ScoreRight  int
}

func NewGame() *Game {
	return &Game{
		State: Instructions,
		LeftPaddle: Paddle{
			X:        50,
			Y:        screenHeight / 2,
			Speed:    5,
			SpinRate: 5.0,
		},
		RightPaddle: Paddle{
			X:        screenWidth - 50,
			Y:        screenHeight / 2,
			Speed:    5,
			SpinRate: 5.0,
		},
		Ball: Ball{
			X:      screenWidth / 2,
			Y:      screenHeight / 2,
			SpeedX: initialBallSpeed,
			SpeedY: initialBallSpeed,
		},
	}
}

func (g *Game) Update() {
	if g.State == Instructions {
		if rl.IsKeyPressed(rl.KeySpace) {
			g.State = Playing
		}
		return
	}

	if g.State == GameOver {
		if rl.IsKeyPressed(rl.KeySpace) {
			// Reset game
			*g = *NewGame()
		}
		return
	}

	g.updatePaddles()
	g.updateBall()

	// Check for win condition
	if g.ScoreLeft >= winScore || g.ScoreRight >= winScore {
		g.State = GameOver
	}
}

func (g *Game) updatePaddles() {
	// Left paddle controls (W/S for vertical, A/D for horizontal)
	if rl.IsKeyDown(rl.KeyW) {
		g.LeftPaddle.Y -= g.LeftPaddle.Speed
	}
	if rl.IsKeyDown(rl.KeyS) {
		g.LeftPaddle.Y += g.LeftPaddle.Speed
	}
	if rl.IsKeyDown(rl.KeyA) {
		g.LeftPaddle.X -= g.LeftPaddle.Speed
	}
	if rl.IsKeyDown(rl.KeyD) {
		g.LeftPaddle.X += g.LeftPaddle.Speed
	}

	// Right paddle controls (Up/Down for vertical, Left/Right for horizontal)
	if rl.IsKeyDown(rl.KeyUp) {
		g.RightPaddle.Y -= g.RightPaddle.Speed
	}
	if rl.IsKeyDown(rl.KeyDown) {
		g.RightPaddle.Y += g.RightPaddle.Speed
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		g.RightPaddle.X -= g.RightPaddle.Speed
	}
	if rl.IsKeyDown(rl.KeyRight) {
		g.RightPaddle.X += g.RightPaddle.Speed
	}

	// Keep paddles within screen bounds
	g.LeftPaddle.Y = clamp(g.LeftPaddle.Y, paddleHeight, screenHeight-paddleHeight)
	g.RightPaddle.Y = clamp(g.RightPaddle.Y, paddleHeight, screenHeight-paddleHeight)

	// Keep paddles within horizontal bounds (up to center line)
	g.LeftPaddle.X = clamp(g.LeftPaddle.X, 50, screenWidth/2-50)
	g.RightPaddle.X = clamp(g.RightPaddle.X, screenWidth/2+50, screenWidth-50)

	// Update paddle rotation
	g.LeftPaddle.Angle += g.LeftPaddle.SpinRate * rl.GetFrameTime()
	g.RightPaddle.Angle += g.RightPaddle.SpinRate * rl.GetFrameTime()
}

func (g *Game) updateBall() {
	g.Ball.X += g.Ball.SpeedX
	g.Ball.Y += g.Ball.SpeedY

	// Ball collision with top and bottom walls
	if g.Ball.Y < 0 || g.Ball.Y > screenHeight {
		g.Ball.SpeedY *= -1
	}

	// Ball collision with paddles
	if g.checkPaddleCollision(&g.LeftPaddle, "left") || g.checkPaddleCollision(&g.RightPaddle, "right") {
		g.Ball.SpeedX *= -1.1                       // Increase speed slightly
		g.Ball.SpeedY += (rand.Float32() - 0.5) * 2 // Add some randomness
	}

	// Score points and reset ball
	if g.Ball.X < 0 {
		g.ScoreRight++
		g.resetBall("right")
	} else if g.Ball.X > screenWidth {
		g.ScoreLeft++
		g.resetBall("left")
	}
}

func (g *Game) checkPaddleCollision(paddle *Paddle, side string) bool {
	// Calculate distance from ball to paddle center
	dx := g.Ball.X - paddle.X
	dy := g.Ball.Y - paddle.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// Check if ball is within paddle's length
	if distance > paddleHeight {
		return false
	}

	// Calculate angle between ball and paddle
	angle := float32(math.Atan2(float64(dy), float64(dx)))
	relativeAngle := angle - paddle.Angle

	// Normalize angle to -π to π
	for relativeAngle > float32(math.Pi) {
		relativeAngle -= 2 * float32(math.Pi)
	}
	for relativeAngle < -float32(math.Pi) {
		relativeAngle += 2 * float32(math.Pi)
	}

	// Check if ball is within paddle's width
	if math.Abs(float64(relativeAngle)) < float64(paddleWidth/2) {
		// Add a small offset to prevent the ball from getting stuck
		if side == "left" {
			g.Ball.X = paddle.X + paddleWidth/2
		} else {
			g.Ball.X = paddle.X - paddleWidth/2
		}
		g.Ball.LastHit = side
		return true
	}
	return false
}

func (g *Game) resetBall(lastScored string) {
	g.Ball.X = screenWidth / 2
	g.Ball.Y = screenHeight / 2
	g.Ball.SpeedX = initialBallSpeed
	if lastScored == "right" {
		g.Ball.SpeedX = -initialBallSpeed
	}
	g.Ball.SpeedY = (rand.Float32() - 0.5) * initialBallSpeed
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	if g.State == Instructions {
		g.drawInstructions()
	} else if g.State == Playing {
		g.drawGame()
	} else if g.State == GameOver {
		g.drawGame()
		g.drawGameOver()
	}

	rl.EndDrawing()
}

func (g *Game) drawInstructions() {
	titleText := "SPINNINN PONG"
	controlsText := []string{
		"Controls:",
		"Left Paddle:  W/S keys (up/down)",
		"Left Paddle:  A/D keys (left/right)",
		"Right Paddle: UP/DOWN arrow keys (up/down)",
		"Right Paddle: LEFT/RIGHT arrow keys (left/right)",
		"",
		fmt.Sprintf("First to %d points wins!", winScore),
		"",
		"Press SPACE to start",
	}

	// Draw title
	titleSize := 40
	titleWidth := rl.MeasureText(titleText, int32(titleSize))
	rl.DrawText(titleText, int32(screenWidth/2-titleWidth/2), 100, int32(titleSize), rl.Green)

	// Draw controls
	fontSize := 20
	startY := 200
	for i, text := range controlsText {
		width := rl.MeasureText(text, int32(fontSize))
		rl.DrawText(text, int32(screenWidth/2-width/2), int32(startY+i*30), int32(fontSize), rl.Green)
	}
}

func (g *Game) drawGame() {
	// Draw paddles
	g.drawPaddle(g.LeftPaddle)
	g.drawPaddle(g.RightPaddle)

	// Draw ball
	rl.DrawCircle(int32(g.Ball.X), int32(g.Ball.Y), ballSize, rl.Green)

	// Draw scores
	scoreSize := int32(40)
	scoreY := int32(20)
	// Left score
	leftScore := fmt.Sprintf("%d", g.ScoreLeft)
	leftScoreWidth := rl.MeasureText(leftScore, scoreSize)
	rl.DrawText(leftScore, int32(screenWidth/4-leftScoreWidth/2), scoreY, scoreSize, rl.Green)

	// Right score
	rightScore := fmt.Sprintf("%d", g.ScoreRight)
	rightScoreWidth := rl.MeasureText(rightScore, scoreSize)
	rl.DrawText(rightScore, int32(3*screenWidth/4-rightScoreWidth/2), scoreY, scoreSize, rl.Green)

	// Draw center line
	for i := 0; i < screenHeight; i += 20 {
		rl.DrawRectangle(int32(screenWidth/2-2), int32(i), 4, 10, rl.Green)
	}
}

func (g *Game) drawPaddle(p Paddle) {
	// Draw rotating paddle
	rl.DrawLineEx(
		rl.Vector2{X: p.X - float32(math.Cos(float64(p.Angle)))*paddleHeight,
			Y: p.Y - float32(math.Sin(float64(p.Angle)))*paddleHeight},
		rl.Vector2{X: p.X + float32(math.Cos(float64(p.Angle)))*paddleHeight,
			Y: p.Y + float32(math.Sin(float64(p.Angle)))*paddleHeight},
		float32(paddleWidth),
		rl.Green,
	)
}

func (g *Game) drawGameOver() {
	winner := "Left Player"
	if g.ScoreRight > g.ScoreLeft {
		winner = "Right Player"
	}

	gameOverText := fmt.Sprintf("%s Wins!", winner)
	restartText := "Press SPACE to play again"

	// Draw game over message
	fontSize := 40
	textWidth := rl.MeasureText(gameOverText, int32(fontSize))
	rl.DrawText(gameOverText, int32(screenWidth/2-textWidth/2), screenHeight/2-30, int32(fontSize), rl.Green)

	// Draw restart message
	smallFontSize := 20
	textWidth = rl.MeasureText(restartText, int32(smallFontSize))
	rl.DrawText(restartText, int32(screenWidth/2-textWidth/2), screenHeight/2+30, int32(smallFontSize), rl.Green)
}

func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func main() {
	rl.InitWindow(int32(screenWidth), int32(screenHeight), "Spinninnpong")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	rand.Seed(time.Now().UnixNano())
	game := NewGame()

	for !rl.WindowShouldClose() {
		game.Update()
		game.Draw()
	}
}
