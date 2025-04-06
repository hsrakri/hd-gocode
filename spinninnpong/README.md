# Spinninnpong

A unique twist on the classic Pong game where the paddles spin vertically, adding an interesting spin effect to the ball when hit.

## Features

- Spinning paddles that add variable spin to the ball
- Smooth paddle movement with keyboard controls
- Progressive difficulty (ball speeds up with each hit)
- Dynamic ball spin effects
- Automatic ball reset when out of bounds
- Clean, minimalist graphics

## Controls

### Left Paddle
- W: Move Up
- S: Move Down

### Right Paddle
- Up Arrow: Move Up
- Down Arrow: Move Down

## Game Mechanics

- Paddles spin continuously, creating a wave-like motion
- Ball gains spin effect based on where it hits the paddle
- Ball speed increases slightly with each paddle hit
- Ball resets to center with random direction when out of bounds
- Paddles are constrained to screen boundaries

## Requirements

- Go 1.16 or later
- Raylib-go library

## Installation

1. Install Go if not already installed:
   ```bash
   # macOS with Homebrew
   brew install go

   # Linux
   sudo apt-get install golang
   ```

2. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/hd-gocode.git
   cd hd-gocode/spinninnpong
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the game:
   ```bash
   go run main.go
   ```

## Development

The game is built using:
- Go programming language
- Raylib-go for graphics and input handling
- Standard Go math libraries for physics calculations

### Project Structure
```
spinninnpong/
├── main.go      # Main game logic
├── go.mod       # Go module file
├── go.sum       # Go module checksums
└── README.md    # This file
```

## Contributing

Feel free to submit issues, fork the repository, and create pull requests for any improvements.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 