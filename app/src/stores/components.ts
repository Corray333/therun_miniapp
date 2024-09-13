import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useComponentsStore = defineStore('components', () => {
    const animateBonuses = ref<boolean>(false)
    const bonusesLabelPos = ref([0,0])

    return { animateBonuses, bonusesLabelPos }
})
