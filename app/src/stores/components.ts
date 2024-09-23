import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

let errorShowTime = 3000

export const useComponentsStore = defineStore('components', () => {
    const animateBonuses = ref<boolean>(false)
    const bonusesLabelPos = ref([0,0])
    const errors = ref<string[]>([])

    const addError = (error: string) => {
        errors.value.push(error)
        setTimeout(() => {
            errors.value.shift()
        }, errorShowTime)
    }

    return { animateBonuses, bonusesLabelPos, errors, addError }
})
