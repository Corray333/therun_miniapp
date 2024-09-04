import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useBalanceStore = defineStore('balance', () => {
  const race = ref(0)
  const bonuses = ref(0)
  const keys = ref(0)
  // function increment() {
  //   count.value++
  // }

  return { race, bonuses, keys }
})
