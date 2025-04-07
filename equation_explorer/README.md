# Equation Explorer

Equation Explorer is an educational game that combines maze navigation with mathematical problem-solving. Players navigate through a maze while solving various types of equations to progress and earn points.

## Features

- Procedurally generated mazes
- Multiple types of equations based on difficulty levels:
  - Linear equations (Level 1)
  - Multi-step equations (Level 2)
  - Literal equations (Level 3)
  - Inequalities (Level 4)
- Real-time timer tracking
- Score system
- Hints and solution steps
- Immediate feedback on answers

## Controls

- **SPACE**: Start game / Restart game
- **Arrow Keys**: Navigate through the maze
- **H**: Show hint for current equation
- **S**: Show solution steps
- **ESC**: Skip current equation
- **ENTER**: Submit answer

## Gameplay

1. Press SPACE to start the game
2. Navigate through the maze using arrow keys
3. When you reach a green cell, you'll encounter an equation
4. Solve the equation by typing your answer
5. Get hints or solution steps if needed
6. Earn points for correct answers
7. Try to solve as many equations as possible while navigating the maze

## Installation

1. Make sure you have Go installed on your system
2. Clone the repository:
   ```bash
   git clone https://github.com/hsrakri/hd-gocode.git
   ```
3. Navigate to the game directory:
   ```bash
   cd hd-gocode/equation_explorer
   ```
4. Run the game:
   ```bash
   go run main.go
   ```

## Dependencies

- [raylib-go](https://github.com/gen2brain/raylib-go) - Game development library

## Development

The game is structured into three main packages:
- `main.go`: Main game loop and rendering
- `equation/`: Equation generation and solving logic
- `maze/`: Maze generation and navigation

## License

This project is licensed under the MIT License - see the LICENSE file for details. 