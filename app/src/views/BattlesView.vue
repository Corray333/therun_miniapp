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

const remainingTime = ref<string>('00:00:00');

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
    let secondsLeft = round.value.endTime - now;

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
        <Transition>
            <section v-if="round === undefined"
                class=" fixed z-40 top-0 left-0 w-full h-screen bg-white flex justify-center items-center">
                <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
            </section>
        </Transition>
        <section class=" p-4 flex flex-col gap-4">
            <Balances/>

            <div class="timer flex flex-col items-center justify-center p-4 rounded-2xl bg-secondary">
                <p class=" round-timer text-4xl font-bold">{{ remainingTime }}</p>
                <p class="font-bold">{{ t('screens.battles.worldRound') }} #{{ round?.id }}</p>
            </div>

            <div class="banner rounded-2xl w-full p-4 bg-cover text-white ">
                <p class="font-bold text-xl">{{ t('screens.battles.banner.header') }}</p>
                <p class="">{{ t('screens.battles.banner.description') }}</p>
            </div>

            <section class=" flex flex-col gap-4">
                <Battle v-for="(battle, i) of round?.battles" :key="i" :battle="battle" />
            </section>

            <span class="flex gap-2 text-dark">
                <i class=" pi pi-info-circle mt-1" style="font-size: 1.25rem;"></i>
                <p>{{ t('screens.battles.info') }}</p>
            </span>
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