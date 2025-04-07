# Equation Explorer

An educational game that combines algebra problem-solving with maze navigation, designed for 9th-grade students.

## Overview

Equation Explorer is an interactive game where players navigate through a maze by solving algebraic equations. The game progressively introduces more complex mathematical concepts while maintaining an engaging and fun experience.

## Features

### Core Gameplay
- Maze navigation with equation-solving obstacles
- Progressive difficulty levels
- Immediate feedback on solutions
- Scoring system with bonuses for speed and accuracy

### Mathematical Content
- Linear Equations
  - Single-variable equations
  - Equations with integers, fractions, and decimals
- Multi-Step Equations
  - Multiple operations
  - Distributive property
- Literal Equations
  - Solving for specific variables
  - Formula rearrangement
- Inequalities
  - Solving inequalities
  - Understanding inequality direction
- Systems of Linear Equations (Advanced Level)
  - Two equations with two variables

### Learning Features
- Clear equation presentation
- Optional hints system
- Step-by-step solution guides
- Progress tracking
- Difficulty level selection

## Technical Details

- Written in Go
- Uses the Raylib library for graphics
- Cross-platform support (Windows, macOS, Linux)

## Getting Started

1. Ensure you have Go installed on your system
2. Clone the repository
3. Navigate to the equation_explorer directory
4. Run `go run main.go`

## Controls

- Arrow keys: Navigate through the maze
- Enter: Submit equation solution
- H: Request a hint
- S: Show solution steps
- ESC: Exit game

## Development

The project is structured as follows:
- `main.go`: Main game loop and initialization
- `equation/`: Equation generation and validation
- `maze/`: Maze generation and navigation
- `ui/`: User interface components
- `game/`: Core game logic

## License

This project is licensed under the MIT License - see the LICENSE file for details. 