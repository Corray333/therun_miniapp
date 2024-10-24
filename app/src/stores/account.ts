import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import {User} from '@/types/types'

export const useAccountStore = defineStore('account', () => {
    const user = ref<User>(new User())
    const token = ref<string>("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk4MDUzMjUsImlkIjozNzc3NDI3NDgsInVzZXJuYW1lIjoiY29ycmF5OSJ9.dmaeS2cCAlvmz8Mtko6OPsTxTRQXXE8PauaJQl3UZnc")
    // function increment() {
    //   count.value++
    // }

    return { user, token }
})
