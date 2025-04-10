# Multiplication Tumbler

A fast-paced multiplication game where players need to guess the missing factor in a multiplication problem before time runs out!

## Game Overview

Multiplication Tumbler is an educational game that challenges players to quickly solve multiplication problems. The game shows one number, hides the other (replacing it with a "?"), and displays the result. Players must figure out the hidden number within a 60-second time limit.

## Features

- Random multiplication problems with numbers 1-12
- 60-second countdown timer
- Score tracking
- Immediate feedback on answers
- Dynamic problem generation
- Visual feedback for correct and incorrect answers

## How to Play

1. Press ENTER to start the game
2. You'll see a problem in the format: `n × ? = result` or `? × n = result`
3. Type your answer for the missing number
4. Press ENTER to submit your answer
5. Get as many correct answers as possible before time runs out!

## Controls

- Number keys (0-9): Enter your answer
- Backspace: Delete last digit entered
- Enter: Submit answer/Start game
- Esc: Exit game

## Installation

1. Ensure you have Go installed on your system
2. Clone the repository
3. Navigate to the multiplicationtumbler directory
4. Run the following commands:
   ```bash
   go mod tidy
   go run main.go
   ```

## Dependencies

- Go 1.21 or later
- Raylib-go library for graphics and input handling

## Learning Outcomes

- Practice multiplication facts
- Develop quick mental math skills
- Improve response time under pressure
- Build confidence with multiplication tables

## Development

The game is built using Go and the Raylib library for graphics. It features:
- Clean, modular code structure
- Efficient random problem generation
- Smooth graphics and text rendering
- Real-time input handling

## License

This project is licensed under the MIT License - see the LICENSE file for details. 