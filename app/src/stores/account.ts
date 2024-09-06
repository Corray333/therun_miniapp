import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import {User} from '@/types/types'

export const useAccountStore = defineStore('account', () => {
    const user = ref<User>(new User())
    const token = ref<string>("")
    // function increment() {
    //   count.value++
    // }

    return { user, token }
})
