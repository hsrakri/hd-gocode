# Guess the Angle

A fun educational game that tests your ability to estimate angles in geometric shapes. The game presents various triangles and asks you to guess the measure of specific angles.

## Features

- 10 different angle estimation challenges
- Visual representation of geometric shapes
- Tolerance of ±5 degrees for correct answers
- Three attempts per question
- Real-time feedback
- Score tracking
- Time tracking
- Top 10 leaderboard
- Player name registration

## Controls

- **SPACE**: Start game / View leaderboard
- **ENTER**: Submit answer / Continue
- **BACKSPACE**: Delete character
- **Numbers and Decimal Point**: Enter angle measurement

## Gameplay

1. Press SPACE to start the game
2. Enter your name
3. For each question:
   - Look at the displayed triangle
   - Estimate the angle at the specified vertex
   - Type your answer (in degrees)
   - Press ENTER to submit
   - You have three attempts per question
   - If you're wrong after three attempts, the correct answer will be shown
4. After completing all questions:
   - View your final score and time
   - Press SPACE to see the leaderboard
   - Your score will be added to the top 10 if it's high enough

## Installation

1. Make sure you have Go installed on your system
2. Clone the repository:
   ```bash
   git clone https://github.com/hsrakri/hd-gocode.git
   ```
3. Navigate to the game directory:
   ```bash
   cd hd-gocode/guesstheangle
   ```
4. Run the game:
   ```bash
   go run main.go
   ```

## Dependencies

- [raylib-go](https://github.com/gen2brain/raylib-go) - Game development library

## Development

The game is written in Go and uses the Raylib library for graphics and input handling. The code is organized into a single main file that handles all game logic, rendering, and state management.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 