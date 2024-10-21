import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { Round } from '@/types/types'

const errorShowTime = 3000

export const useComponentsStore = defineStore('components', () => {
    const animateBonuses = ref<boolean>(false)
    const bonusesLabelPos = ref([0,0])
    const errors = ref<string[]>([])
    const round = ref<Round>({
        id: 470,
        endTime: 1729530000,
        element: "desert",
        battles: []
    })

    const addError = (error: string) => {
        errors.value.push(error)
        setTimeout(() => {
            errors.value.shift()
        }, errorShowTime)
    }

    return { animateBonuses, bonusesLabelPos, errors, round, addError }
})
