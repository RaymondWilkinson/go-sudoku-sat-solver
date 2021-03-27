<template>
  <td
      :class="[
          'box-content inline-block w-8 h-8 sm:w-12 sm:h-12 p-0',
           column % 3 === 2 && column < 8 ? 'border-r-4 border-gray-500' : 'border-r border-gray-700']"
  >
    <div
        class="relative grid grid-rows-3 grid-cols-3 w-8 h-8 sm:w-12 sm:h-12 focus:outline-none focus:ring focus:border-blue-500 z-10"
        @click="edit"
    >
      <p
          v-show="!editing"
          :class="[
              'row-span-3 col-span-3 self-center text-3xl text-center sm:text-4xl',
              cell.isConfirmed ? 'font-sans font-bold' : 'font-cursive text-gray-400'
          ]"
      >
        {{ cell.data !== 0 ? cell.data : '' }}
      </p>

      <input v-show="editing" @blur="blur" v-model="inputData" type="text" maxlength="1" ref="input" class="row-span-3 col-span-3 self-center text-3xl text-center text-gray-400 font-cursive">
    </div>
  </td>
</template>

<script>
export default {
  props: [
    'row',
    'column',
    'cell',
    'onUpdate'
  ],
  watch: {
    cell: {
      handler() {
        let data = this.cell.data === 0 ? "" : this.cell.data.toString()
        if (data !== this.inputData) {
          this.inputData = data
        }
      },
      deep: true,
      immediate: true
    },
    inputData(newData, oldData) {
      if (newData === "") {
        this.onUpdate(this.row, this.column, 0)
        return
      }

      if (newData.length > 1) {
        this.inputData = oldData
      }

      let num = parseInt(newData)
      if (isNaN(num) || num < 1 || num > 9) {
        this.inputData = oldData
      }

      this.onUpdate(this.row, this.column, num)
    }
  },
  data: () => {
    return {
      editing: false,
      inputData: ""
    }
  },
  mounted() {
    this.inputData = this.cell.data === 0 ? "" : this.cell.data.toString()
  },
  methods: {
    edit() {
      if (this.cell.isConfirmed) {
        return
      }

      this.editing = true

      setTimeout(() => {
        this.$refs.input.focus()
      }, 0)
    },

    blur() {
      this.editing = false
    },

    updateCell(e) {
      console.log(e)
      return false
    }
  }
}
</script>
