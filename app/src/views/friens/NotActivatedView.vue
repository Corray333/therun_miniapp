<script lang="ts" setup>

import { ref, onBeforeMount } from 'vue'
import { Referal } from '@/types/types'
import axios, { isAxiosError } from 'axios'
import { auth } from '@/utils/helpers'
import { useAccountStore } from '@/stores/account'
import { useComponentsStore } from '@/stores/components'
import { useI18n } from 'vue-i18n'

const componentsStore = useComponentsStore()
const accStore = useAccountStore()
const infoLoaded = ref<boolean>(false)
const { t } = useI18n()


const friends = ref<Referal[]>([])


// const friends = ref<Referal[]>([{
//     username:"Markovnik",
//     avatar: 'https://cdn-icons-png.freepik.com/256/15707/15707874.png?semt=ais_hybrid'
// }])

const getFriends = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/users/0/referals?activated=false`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        friends.value.push(...data)
        return true
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await getFriends()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        }
    }
}

onBeforeMount(async () => {
    await getFriends()
    infoLoaded.value = true
})

</script>

<template>
    <section class=" pb-20">
        <Transition>
            <section v-if="!infoLoaded"
                class=" fixed z-20 top-0 left-0 w-full h-screen bg-white flex justify-center items-center">
                <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
            </section>
        </Transition>
        <section class="p-4 flex flex-col gap-4">
            <h1>{{ t('screens.friens.waitingForActivationFriends') }}</h1>
            <div class=" friends-list flex flex-col gap-1">
                <div class="flex w-full gap-4 items-center bg-half_dark p-2" v-for="(friend, i) of friends" :key="i">
                    <img v-if="friend.avatar != ''" :src="friend.avatar" class="w-12 h-12 rounded-full" alt="">
                    <div v-else  class="w-12 h-12 min-w-12 flex items-center justify-center rounded-full bg-dark text-white font-bold text-2xl"><p>{{ friend.username != '' ? friend.username[0].toUpperCase() : '?' }}</p></div>
                    <div class=" w-full">
                        <p>{{ friend.username }}</p>
                    </div>
                    <i class="pi pi-angle-right text-dark" style="font-size: 1.25rem"></i>
                </div>
            </div>
        </section>
    </section>
</template>


<style scoped>

.friends-list>div:first-child{
    border-top-left-radius: 1rem;
    border-top-right-radius: 1rem;
}

.friends-list>div:last-child{
    border-bottom-left-radius: 1rem;
    border-bottom-right-radius: 1rem;
}

</style>