package sudoku

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Grid struct {
	grid [9][9]uint8
}

func ParseGrid(gridStr string) *Grid {
	if len(gridStr) != 81 {
		return nil
	}

	grid := &Grid{}

	for index, numStr := range gridStr {
		num, err := strconv.ParseUint(string(numStr), 10, 8)
		if err != nil {
			log.Error(err)
			return nil
		}

		if num < 0 || num > 9 {
			return nil
		}

		grid.grid[index/9][index%9] = uint8(num)
	}

	return grid
}

func (g *Grid) String() string {
	output := ""

	for _, row := range g.grid {
		for _, num := range row {
			output += fmt.Sprintf("%d", num)
		}
	}

	return output
}

func (g *Grid) Draw() string {
	output := "\n"

	for _, row := range g.grid {
		rowString := ""
		for _, num := range row {
			rowString = fmt.Sprintf("%s %d", rowString, num)
		}

		output += strings.TrimSpace(rowString) + "\n"
	}

	return output
}

func (g *Grid) EachNumber(f func(row int, column int, num int)) {
	for rowIndex, row := range g.grid {
		for columnIndex, num := range row {
			if num != 0 {
				f(rowIndex+1, columnIndex+1, int(num))
			}
		}
	}
}

func (g *Grid) Fill(row int, column int, number int) {
	g.grid[row-1][column-1] = uint8(number)
}

func (g *Grid) RandomFill() {
	numbers := [9]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(time.Now().UnixNano())
	for i := len(numbers) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	g.grid[0] = numbers

	columnNumbers := []uint8{}
	clashingNumbers := []uint8{}
	for i, number := range numbers {
		if i != 0 {
			if i <= 2 {
				clashingNumbers = append(clashingNumbers, number)
			} else {
				columnNumbers = append(columnNumbers, number)
			}
		}
	}

	for i := 1; i < 9; i++ {
		var chosenIndex int

		// need to ensure there isn't a clash within the box
		if i == 3 {
			columnNumbers = append(clashingNumbers, columnNumbers...)
		}

		chosenIndex = rand.Intn(len(columnNumbers))

		g.grid[i][0] = columnNumbers[chosenIndex]

		columnNumbers = append(columnNumbers[:chosenIndex], columnNumbers[chosenIndex+1:]...)
	}
}

func (g *Grid) RemoveSquares(squaresToRemove int) {
	for i, row := range g.grid {
		for j, num := range row {
			rand.Seed(time.Now().UnixNano())
			r := rand.Intn(2)
			if num != 0 && squaresToRemove > 0 && r == 0 {
				g.grid[i][j] = 0
				squaresToRemove--
			}
		}
	}

	if squaresToRemove > 0 {
		g.RemoveSquares(squaresToRemove)
	}
}

type Pair struct {
	A int
	B int
}

var BlockCenters = []Pair{
	{2, 2},
	{2, 5},
	{2, 8},

	{5, 2},
	{5, 5},
	{5, 8},

	{8, 2},
	{8, 5},
	{8, 8},
}

var BlockCenterOffsets = []Pair{
	{-1, -1},
	{-1, 0},
	{-1, 1},

	{0, -1},
	{0, 0},
	{0, 1},

	{1, -1},
	{1, 0},
	{1, 1},
}
