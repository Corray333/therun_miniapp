<script lang="ts" setup>

import Balances from '@/components/Balances.vue'

import { ref, onBeforeMount } from 'vue'
import { useI18n } from 'vue-i18n'
import { type Round } from '@/types/types'
import { useAccountStore } from '@/stores/account'
import Battle from '@/components/Battle.vue'
import axios, {isAxiosError} from 'axios'
import { useComponentsStore } from '@/stores/components'
import { auth } from '@/utils/helpers'

const accStore = useAccountStore()
const componentsStore = useComponentsStore()

const { t } = useI18n()

const remainingTime = ref<string>('00:00:00')
const remainingSeconds = ref<number>(0)
const showMiles = ref<boolean>(false)

const showMilesAfter = 12 * 60 * 60
const roundDuration = 26 * 60 * 60

const round = ref<Round>()

const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600).toString().padStart(2, '0');
    const minutes = Math.floor((seconds % 3600) / 60).toString().padStart(2, '0');
    const secondsRemaining = (seconds % 60).toString().padStart(2, '0');
    return `${hours}:${minutes}:${secondsRemaining}`;
};

const calculateRemainingTimeAndPoints = () => {
    if (!round.value) {
        return
    }
    const now = Date.now() / 1000
    let secondsLeft = round.value.endTime - now
    remainingSeconds.value = secondsLeft
    showMiles.value = (roundDuration - secondsLeft) > showMilesAfter

    if (secondsLeft <= 0) {
        remainingTime.value = '00:00:00'
        return
    }

    remainingTime.value = formatTime(Math.floor(secondsLeft))
}

onBeforeMount(async () => {
    await getRound()
    calculateRemainingTimeAndPoints()
    setInterval(calculateRemainingTimeAndPoints, 1000)
})

const getRound = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/round`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        
        round.value = data

        return true
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await getRound()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        }
    }
}



</script>

<template>
    <section class="pb-20">

        <p v-if="round == undefined" class="text-center font-dark mt-4"><i class=" pi pi-spinner pi-spin" style="font-size: 1.5rem; color: var(--dark)"></i></p>

        <section v-else class=" p-4 flex flex-col gap-4">
            <Balances/>

            <div class="timer flex flex-col items-center justify-center p-4 rounded-2xl bg-secondary">
                <p class=" round-timer text-4xl font-bold">{{ remainingTime }}</p>
                <p class="font-bold">{{ t('screens.battles.worldRound') }} #{{ round?.id }}</p>
            </div>

            <div class="banner rounded-2xl w-full p-4 bg-cover text-white ">
                <p class="font-bold text-xl">{{ t('screens.battles.banner.header') }}</p>
                <p class="">{{ t('screens.battles.banner.description') }}</p>
            </div>

            <h2 class="text-center text-2xl" v-html="t('screens.battles.header')"></h2>

            <p v-if="!round?.battles.length" class="font-bold text-dark text-center">{{ t('screens.battles.noBattles') }}</p>
            <div v-else>
                <section class=" flex flex-col gap-4">
                    <Battle v-for="(battle, i) of round?.battles" :key="i" :battle="battle" :show-miles="showMiles" />
                </section>
    
                <span class="flex gap-2 text-dark">
                    <i class=" pi pi-info-circle mt-1" style="font-size: 1.25rem;"></i>
                    <p>{{ t('screens.battles.info') }}</p>
                </span>
            </div>
        </section>
    </section>
</template>


<style scoped>

.round-timer{
    width: 8.5rem;
}

.banner{
    background-image: url(../assets/images/battles/banner-bg.png);
}

.v-enter-active,
.v-leave-active {
    transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
    opacity: 0;
}

</style>