<template>
  <table
      role="grid"
      class="flex border-collapse select-none"
  >
    <tbody class="flex-grow-0 flex-shrink-0 border-4 border-gray-500">
    <tr
      v-for="(row, rindex) in numbers"
      :key="rindex"
      :class="[rindex % 3 === 2 && rindex < 8? 'border-b-4 border-gray-500' : 'border-b border-gray-700']"
    >
      <SudokuGridCell
          v-for="(column, cindex) in numbers"
          :key="cindex"
          :row="row"
          :column="column"
          :cell="cells[row * 9 + column]"
          :on-update="onUpdate"
      />
    </tr>
    </tbody>
  </table>
</template>

<script>
import SudokuGridCell from "./SudokuGridCell"

export default {
  props: ['cells', 'onUpdate'],
  components: {
    SudokuGridCell
  },
  data: () => {
    return {
      numbers: [0, 1, 2, 3, 4, 5, 6, 7, 8],
    }
  }
}
</script>
