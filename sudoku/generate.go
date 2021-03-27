package sudoku

import (
	"fmt"
	"math/rand"
	"time"
)

// difficulties
var hard = [5]int{23, 24, 25, 26, 27}
var medium = [5]int{28, 29, 30, 31, 32}
var easy = [5]int{33, 34, 35, 36, 37}

func Generate(difficulty string) (*Grid, error) {
	initialGrid := &Grid{}

	initialGrid.RandomFill()

	solvedGrid, err := Solve(initialGrid)
	if err != nil {
		fmt.Println(initialGrid.Draw())
		return nil, fmt.Errorf("Failed to generate Sudoku")
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5)

	squaresToKeep := easy[n]

	if difficulty == "hard" {
		squaresToKeep = hard[n]
	} else if difficulty == "medium" {
		squaresToKeep = medium[n]
	}

	solvedGrid.RemoveSquares(81 - squaresToKeep)

	return solvedGrid, nil
}
