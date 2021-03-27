<template>
  <div id="app">
    <div v-if="loading" class="fixed z-50 inset-0 overflow-y-auto">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">

        <!-- This element is to trick the browser into centering the modal contents. -->
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

        <span class="text-purple-500 opacity-75">
          <i class="fas fa-circle-notch fa-spin fa-5x"></i>
        </span>

      </div>
    </div>

    <div v-if="!initializing">
      <div class="flex flex-col items-center justify-center mt-5">
        <p class="text-4xl ...">Sudoku</p>
        <p class="text-lg ...">{{ prettyDifficulty }}</p>
      </div>

      <div class="flex flex-grow items-center justify-center py-5 sm:py-8">
        <div class="grid grid-flow-row sm:grid-flow-col justify-center gap-8">
          <SudokuGrid :cells="cells" :on-update="updateGrid"/>
        </div>
      </div>

      <div class="flex flex-row items-center justify-center mb-10">
        <div :class="['flex inline-flex', canMakeChanges ? 'mr-10' : '']" role="group">
          <button type="button"
                  class="focus:outline-none text-white text-md py-2.5 px-5 bg-gray-500 rounded-l-md hover:bg-gray-600 hover:shadow-lg"
                  @click="setDifficulty('easy')">
            Easy
          </button>
          <button type="button"
                  class="focus:outline-none text-white text-md py-2.5 px-5 bg-gray-500 hover:bg-gray-600 hover:shadow-lg"
                  @click="setDifficulty('medium')">
            Medium
          </button>
          <button type="button"
                  class="focus:outline-none text-white text-md py-2.5 px-5 bg-gray-500 rounded-r-md hover:bg-gray-600 hover:shadow-lg"
                  @click="setDifficulty('hard')">
            Hard
          </button>
        </div>

        <div v-if="canMakeChanges" class="flex ml-10">
          <button v-if="completed" type="button"
                  class="focus:outline-none text-white m-1 text-md py-2.5 px-5 rounded-md bg-pink-500 hover:bg-pink-600 hover:shadow-lg"
                  @click="validate">
            Validate
          </button>
          <button v-if="hasMadeChanges" type="button"
                  class="focus:outline-none text-white m-1 text-md py-2.5 px-5 rounded-md bg-indigo-500 hover:bg-indigo-600 hover:shadow-lg"
                  @click="clear">
            Clear
          </button>
          <button type="button"
                  class="focus:outline-none text-white m-1 text-md py-2.5 px-5 rounded-md bg-purple-500 hover:bg-purple-600 hover:shadow-lg"
                  @click="solve">
            Solve
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import Swal from 'sweetalert2'
import SudokuGrid from "./components/SudokuGrid"

export default {
  name: 'App',
  components: {
    SudokuGrid
  },
  data: () => {
    return {
      initializing: true,
      loading: true,
      difficulty: "easy",
      cells: [],
    }
  },
  computed: {
    completed() {
      for (let i = 0; i < this.cells.length; i++) {
        if (this.cells[i].data === 0) {
          return false
        }
      }

      return true
    },
    canMakeChanges() {
      for (let i = 0; i < this.cells.length; i++) {
        if (this.cells[i].isConfirmed === false) {
          return true
        }
      }

      return false
    },
    hasMadeChanges() {
      for (let i = 0; i < this.cells.length; i++) {
        if (this.cells[i].isConfirmed === false && this.cells[i].data !== 0) {
          return true
        }
      }

      return false
    },
    prettyDifficulty() {
      return this.difficulty.charAt(0).toUpperCase() + this.difficulty.slice(1)
    }
  },
  mounted() {
    this.generateBoard()
  },
  methods: {
    generateBoard() {
      axios.get(`https://api.sudoku.raymondwilkinson.com/generate?difficulty=${this.difficulty}`)
        .then(({data}) => {
          this.updateCellsFromString(data.grid)

          this.initializing = false
        })
        .catch(() => {
          Swal.fire(
              'An error occurred.',
              "Please try again later.",
              'error'
          )
        })
        .finally(() => {
          this.loading = false
        })
    },

    updateCellsFromString(boardString) {
      let cells = []

      for (let i = 0; i < boardString.length; i++) {
        let data = parseInt(boardString[i])

        const cell = {
          data: data,
          isConfirmed: data > 0 && data < 10
        }

        cells.push(cell)
      }

      this.cells = cells
    },

    updateGrid(row, column, value) {
      let currentCell = this.cells[row * 9 + column]

      if (currentCell.isConfirmed) {
        return
      }

      if (currentCell.data !== value) {
        currentCell.data = value
      }
    },

    solve() {
      this.loading = true

      let grid = ""
      this.cells.forEach(cell => {
        grid += `${cell.data === "" ? 0 : cell.data}`
      })

      axios.get(`https://api.sudoku.raymondwilkinson.com/solve?grid=${grid}`)
        .then(({data}) => {
          this.updateCellsFromString(data.grid)
        })
        .catch(({response}) => {
          Swal.fire(
              'Uh oh.',
              response.data.message,
              'error'
          )
        })
        .finally(() => {
          this.loading = false
        })
    },

    validate() {
      this.loading = true

      let grid = ""
      this.cells.forEach(cell => {
        grid += `${cell.data === "" ? 0 : cell.data}`
      })

      axios.get(`https://api.sudoku.raymondwilkinson.com/solve?grid=${grid}`)
        .then(() => {
          Swal.fire(
              'Success!',
              'You solved the Sudoku puzzle.',
              'success'
          )
        })
        .catch(() => {
          Swal.fire(
              'Uh oh.',
              'Looks like the Sudoku has not been solved.',
              'error'
          )
        })
        .finally(() => {
          this.loading = false
        })
    },

    setDifficulty(difficulty) {
      this.difficulty = difficulty
      this.loading = true
      this.generateBoard()
    },

    clear() {
      for (let i = 0; i < this.cells.length; i++) {
        let cell = this.cells[i]
        if (cell.isConfirmed === false && cell.data !== 0) {
          cell.data = 0
          this.cells[i] = cell
        }
      }
    }
  }
}
</script>
