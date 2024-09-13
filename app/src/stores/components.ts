import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useComponentsStore = defineStore('components', () => {
    const animateBonuses = ref<boolean>(false)

    return { animateBonuses }
})
