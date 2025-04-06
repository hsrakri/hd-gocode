package main

import (
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
	paddleWidth  = 20
	paddleHeight = 100
	ballSize     = 10
)

type Paddle struct {
	X, Y   float32
	Speed  float32
	Angle  float32
	Radius float32
}

type Ball struct {
	X, Y   float32
	SpeedX float32
	SpeedY float32
}

func main() {
	rl.InitWindow(int32(screenWidth), int32(screenHeight), "Spinninnpong")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	leftPaddle := Paddle{X: 50, Y: screenHeight / 2, Speed: 5, Radius: paddleHeight / 2}
	rightPaddle := Paddle{X: screenWidth - 50, Y: screenHeight / 2, Speed: 5, Radius: paddleHeight / 2}

	ball := Ball{X: screenWidth / 2, Y: screenHeight / 2, SpeedX: 4, SpeedY: 4}

	rand.Seed(time.Now().UnixNano())

	for !rl.WindowShouldClose() {
		// Update
		updatePaddle(&leftPaddle)
		updatePaddle(&rightPaddle)
		updateBall(&ball, &leftPaddle, &rightPaddle)

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		drawPaddle(leftPaddle)
		drawPaddle(rightPaddle)
		drawBall(ball)

		rl.EndDrawing()
	}
}

func updatePaddle(p *Paddle) {
	// Add keyboard controls for paddles
	if p.X < screenWidth/2 { // Left paddle
		if rl.IsKeyDown(rl.KeyW) {
			p.Y -= p.Speed
		}
		if rl.IsKeyDown(rl.KeyS) {
			p.Y += p.Speed
		}
	} else { // Right paddle
		if rl.IsKeyDown(rl.KeyUp) {
			p.Y -= p.Speed
		}
		if rl.IsKeyDown(rl.KeyDown) {
			p.Y += p.Speed
		}
	}

	// Keep paddles within screen bounds
	if p.Y < p.Radius {
		p.Y = p.Radius
	}
	if p.Y > screenHeight-p.Radius {
		p.Y = screenHeight - p.Radius
	}

	// Update spinning angle
	p.Angle += 0.1
	if p.Angle > 2*math.Pi {
		p.Angle -= 2 * math.Pi
	}
}

func updateBall(b *Ball, leftPaddle, rightPaddle *Paddle) {
	b.X += b.SpeedX
	b.Y += b.SpeedY

	// Ball collision with top and bottom walls
	if b.Y < 0 || b.Y > screenHeight {
		b.SpeedY *= -1
	}

	// Ball collision with paddles
	if checkCollision(b, leftPaddle) || checkCollision(b, rightPaddle) {
		b.SpeedX *= -1.1                 // Increase speed slightly on each hit
		b.SpeedY += rand.Float32()*2 - 1 // Add spin effect
	}

	// Reset ball if it goes out of bounds
	if b.X < 0 || b.X > screenWidth {
		b.X = screenWidth / 2
		b.Y = screenHeight / 2
		b.SpeedX = 4 * float32(map[bool]int{true: 1, false: -1}[rand.Float32() > 0.5])
		b.SpeedY = 4 * float32(map[bool]int{true: 1, false: -1}[rand.Float32() > 0.5])
	}
}

func checkCollision(b *Ball, p *Paddle) bool {
	paddleTop := p.Y - p.Radius
	paddleBottom := p.Y + p.Radius
	paddleLeft := p.X - paddleWidth/2
	paddleRight := p.X + paddleWidth/2

	return b.X-ballSize/2 < paddleRight && b.X+ballSize/2 > paddleLeft && b.Y > paddleTop && b.Y < paddleBottom
}

func drawPaddle(p Paddle) {
	paddleTop := p.Y - p.Radius
	paddleBottom := p.Y + p.Radius

	// Draw spinning paddle segments
	for y := paddleTop; y < paddleBottom; y += 5 {
		xOffset := 5 * math.Sin(float64(p.Angle+(y-paddleTop)/p.Radius))
		rl.DrawRectangle(
			int32(p.X-paddleWidth/2+float32(xOffset)),
			int32(y),
			paddleWidth,
			5,
			rl.Black,
		)
	}
}

func drawBall(b Ball) {
	rl.DrawCircle(int32(b.X), int32(b.Y), ballSize, rl.Red)
}
