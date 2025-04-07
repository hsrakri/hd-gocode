package maze

import (
	"math/rand"
)

// Cell represents a single cell in the maze
type Cell struct {
	X, Y     int
	Wall     bool
	Equation bool
	Visited  bool
}

// Maze represents the game maze
type Maze struct {
	Width  int
	Height int
	Cells  [][]Cell
}

// NewMaze creates a new maze with the given dimensions
func NewMaze(width, height int) *Maze {
	maze := &Maze{
		Width:  width,
		Height: height,
		Cells:  make([][]Cell, height),
	}

	// Initialize cells
	for y := range maze.Cells {
		maze.Cells[y] = make([]Cell, width)
		for x := range maze.Cells[y] {
			maze.Cells[y][x] = Cell{
				X:       x,
				Y:       y,
				Wall:    true,
				Visited: false,
			}
		}
	}

	return maze
}

// Generate creates a maze using recursive backtracking
func (m *Maze) Generate() {
	// Start from the top-left corner
	startX, startY := 1, 1
	m.Cells[startY][startX].Wall = false
	m.Cells[startY][startX].Visited = true

	// Generate the maze
	m.carve(startX, startY)

	// Add equations to some cells
	m.addEquations()
}

// carve recursively carves out the maze
func (m *Maze) carve(x, y int) {
	// Define possible directions (up, right, down, left)
	directions := [][2]int{
		{0, -2}, // up
		{2, 0},  // right
		{0, 2},  // down
		{-2, 0}, // left
	}

	// Shuffle directions
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	// Try each direction
	for _, dir := range directions {
		newX, newY := x+dir[0], y+dir[1]

		// Check if the new position is valid and unvisited
		if newX > 0 && newX < m.Width-1 && newY > 0 && newY < m.Height-1 && !m.Cells[newY][newX].Visited {
			// Carve a path
			m.Cells[y+dir[1]/2][x+dir[0]/2].Wall = false
			m.Cells[newY][newX].Wall = false
			m.Cells[newY][newX].Visited = true

			// Recursively carve from the new position
			m.carve(newX, newY)
		}
	}
}

// addEquations adds equations to some cells in the maze
func (m *Maze) addEquations() {
	// Calculate how many equations to add (about 10% of the maze size)
	numEquations := (m.Width * m.Height) / 10

	// Add equations to random non-wall cells
	for i := 0; i < numEquations; i++ {
		x := rand.Intn(m.Width)
		y := rand.Intn(m.Height)

		if !m.Cells[y][x].Wall {
			m.Cells[y][x].Equation = true
		}
	}
}

// IsWall checks if a cell is a wall
func (m *Maze) IsWall(x, y int) bool {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return true
	}
	return m.Cells[y][x].Wall
}

// HasEquation checks if a cell has an equation
func (m *Maze) HasEquation(x, y int) bool {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return false
	}
	return m.Cells[y][x].Equation
}

// GetCell returns a cell at the given coordinates
func (m *Maze) GetCell(x, y int) *Cell {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return nil
	}
	return &m.Cells[y][x]
}
