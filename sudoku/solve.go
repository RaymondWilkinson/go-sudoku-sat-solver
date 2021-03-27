package sudoku

import (
	"fmt"
	"github.com/mitchellh/go-sat"
	"github.com/mitchellh/go-sat/cnf"
	"math"
)

func SolveFromString(sudoku string) (*Grid, error) {
	grid := ParseGrid(sudoku)
	if grid == nil {
		return nil, fmt.Errorf("Invalid Sudoku provided.")
	}

	return Solve(grid)
}

func Solve(grid *Grid) (*Grid, error) {
	solver := buildSudokuSolver()

	// Add a clause for each number that's already solved
	grid.EachNumber(func(row, column, num int) {
		solver.AddClause(cnf.NewClauseFromInts([]int{encodeVariable(row, column, num)}))
	})

	solved := solver.Solve()

	if !solved {
		return nil, fmt.Errorf("Unsolvable Sudoku provided.")
	}

	resultGrid := &Grid{}

	for variable, assigned := range solver.Assignments() {
		if assigned {
			row, column, number := decodeVariable(variable)
			resultGrid.Fill(row, column, number)
		}
	}

	return resultGrid, nil
}

/*
Build Sudoku SAT solver with all rules
 */
func buildSudokuSolver() *sat.Solver {
	solver := sat.New()

	solver.AddFormula(eachFieldGetsAtLeastOneNumber())
	solver.AddFormula(eachNumberOccursAtMostOncePerRow())
	solver.AddFormula(eachNumberOccursAtMostOncePerColumn())
	solver.AddFormula(eachNumberOccursAtMostOncePerBlock())

	return solver
}

/*
Determine SAT variable name for <r, c, n>

Example:
For row = 1, variables range from 1 to 81
*/
func encodeVariable(row, column, number int) int {
	return ((row - 1) * 81) + ((column - 1) * 9) + number
}

/*
Decode the encoded SAT variable
*/
func decodeVariable(variable int) (int, int, int) {
	row := int(math.Ceil(float64(variable) / 81))

	rowRange := (row - 1) * 81

	column := int(math.Ceil(float64(variable - rowRange) / 9))

	num := variable - rowRange - ((column - 1) * 9)

	return row, column, num
}

/*
First Rule: Each <row, column> field contains at least one number
*/
func eachFieldGetsAtLeastOneNumber() cnf.Formula {
	var clauses [][]int

	for row := 1; row <= 9; row++ {
		for column := 1; column <= 9; column++ {
			var clause []int
			for num := 1; num <= 9; num++ {
				clause = append(clause, encodeVariable(row, column, num))
			}

			clauses = append(clauses, clause)
		}
	}

	return cnf.NewFormulaFromInts(clauses)
}

/*
Second Rule: Every number occurs at most once per row
*/
func eachNumberOccursAtMostOncePerRow() cnf.Formula {
	var clauses [][]int

	for row := 1; row <= 9; row++ {
		for num := 1; num <= 9; num++ {

			for column1 := 1; column1 <= 8; column1++ {
				for column2 := column1 + 1; column2 <= 9; column2++ {
					clauses = append(clauses, []int{
						-encodeVariable(row, column1, num),
						-encodeVariable(row, column2, num),
					})
				}
			}

		}
	}

	return cnf.NewFormulaFromInts(clauses)
}

/*
Third Rule: Every number occurs at most once per column
*/
func eachNumberOccursAtMostOncePerColumn() cnf.Formula {
	var clauses [][]int

	for column := 1; column <= 9; column++ {
		for num := 1; num <= 9; num++ {

			for row1 := 1; row1 <= 8; row1++ {
				for row2 := row1 + 1; row2 <= 9; row2++ {
					clauses = append(clauses, []int{
						-encodeVariable(row1, column, num),
						-encodeVariable(row2, column, num),
					})
				}
			}

		}
	}

	return cnf.NewFormulaFromInts(clauses)
}

/*
Fourth Rule: Every number occurs at most once per block

This one is more complicated
*/
func eachNumberOccursAtMostOncePerBlock() cnf.Formula {
	var clauses [][]int

	for _, center := range BlockCenters {
		for o1Index, offset1 := range BlockCenterOffsets {
			for o2Index, offset2 := range BlockCenterOffsets {
				// Filter out duplicates
				if offset1.A == offset2.A && offset1.B == offset2.B {
					continue
				}

				if o1Index < o2Index {
					for num := 1; num <= 9; num++ {
						clauses = append(clauses, []int{
							-encodeVariable(center.A + offset1.A, center.B + offset1.B, num),
							-encodeVariable(center.A + offset2.A, center.B + offset2.B, num),
						})
					}
				}
			}
		}
	}

	return cnf.NewFormulaFromInts(clauses)
}
