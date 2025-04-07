package equation

import (
	"fmt"
	"math/rand"
)

// Equation represents a mathematical equation
type Equation struct {
	Text       string
	Answer     float64
	Type       string
	Difficulty int
	Hint       string
	Steps      []string
}

// GenerateEquation creates a new equation based on difficulty level
func GenerateEquation(level int) *Equation {
	switch level {
	case 1:
		return generateLinearEquation()
	case 2:
		return generateMultiStepEquation()
	case 3:
		return generateLiteralEquation()
	case 4:
		return generateInequality()
	default:
		return generateLinearEquation()
	}
}

// generateLinearEquation creates a simple linear equation
func generateLinearEquation() *Equation {
	// Generate coefficients and constants
	a := rand.Intn(10) + 1  // coefficient of x
	b := rand.Intn(20) - 10 // constant term
	c := rand.Intn(20) - 10 // right side constant

	// Calculate the answer
	answer := float64(c-b) / float64(a)

	// Generate the equation string
	eq := fmt.Sprintf("%dx + %d = %d", a, b, c)

	// Generate solution steps
	steps := []string{
		fmt.Sprintf("Original equation: %s", eq),
		fmt.Sprintf("Subtract %d from both sides: %dx = %d", b, a, c-b),
		fmt.Sprintf("Divide both sides by %d: x = %.2f", a, answer),
	}

	return &Equation{
		Text:       eq,
		Answer:     answer,
		Type:       "Linear",
		Difficulty: 1,
		Hint:       "Try to isolate x by performing inverse operations",
		Steps:      steps,
	}
}

// generateMultiStepEquation creates a multi-step equation
func generateMultiStepEquation() *Equation {
	// Generate coefficients and constants
	a := rand.Intn(5) + 1   // coefficient of x
	b := rand.Intn(10) + 1  // constant term
	c := rand.Intn(5) + 1   // multiplier
	d := rand.Intn(20) - 10 // right side constant

	// Calculate the answer
	answer := float64(d-b) / float64(a*c)

	// Generate the equation string
	eq := fmt.Sprintf("%d(%dx + %d) = %d", c, a, b, d)

	// Generate solution steps
	steps := []string{
		fmt.Sprintf("Original equation: %s", eq),
		fmt.Sprintf("Distribute %d: %dx + %d = %d", c, a*c, b*c, d),
		fmt.Sprintf("Subtract %d from both sides: %dx = %d", b*c, a*c, d-b*c),
		fmt.Sprintf("Divide both sides by %d: x = %.2f", a*c, answer),
	}

	return &Equation{
		Text:       eq,
		Answer:     answer,
		Type:       "Multi-Step",
		Difficulty: 2,
		Hint:       "First distribute the multiplier, then solve like a linear equation",
		Steps:      steps,
	}
}

// generateLiteralEquation creates a literal equation
func generateLiteralEquation() *Equation {
	// Generate coefficients and constants
	a := rand.Intn(5) + 1  // coefficient of x
	b := rand.Intn(5) + 1  // coefficient of y
	c := rand.Intn(10) + 1 // constant term

	// Generate the equation string
	eq := fmt.Sprintf("%dx + %dy = %d", a, b, c)

	// Generate solution steps for solving for x
	steps := []string{
		fmt.Sprintf("Original equation: %s", eq),
		fmt.Sprintf("Subtract %dy from both sides: %dx = %d - %dy", b, a, c, b),
		fmt.Sprintf("Divide both sides by %d: x = (%.2f - %.2fy)/%.2f", a, float64(c), float64(b), float64(a)),
	}

	return &Equation{
		Text:       eq,
		Answer:     0, // Not applicable for literal equations
		Type:       "Literal",
		Difficulty: 3,
		Hint:       "Treat y as a constant and solve for x",
		Steps:      steps,
	}
}

// generateInequality creates an inequality
func generateInequality() *Equation {
	// Generate coefficients and constants
	a := rand.Intn(5) + 1  // coefficient of x
	b := rand.Intn(10) - 5 // constant term
	c := rand.Intn(10) + 1 // right side constant

	// Calculate the answer
	answer := float64(c-b) / float64(a)

	// Generate the equation string
	eq := fmt.Sprintf("%dx + %d > %d", a, b, c)

	// Generate solution steps
	steps := []string{
		fmt.Sprintf("Original inequality: %s", eq),
		fmt.Sprintf("Subtract %d from both sides: %dx > %d", b, a, c-b),
		fmt.Sprintf("Divide both sides by %d: x > %.2f", a, answer),
	}

	return &Equation{
		Text:       eq,
		Answer:     answer,
		Type:       "Inequality",
		Difficulty: 4,
		Hint:       "Solve like a linear equation, but remember to flip the inequality if you multiply or divide by a negative number",
		Steps:      steps,
	}
}
