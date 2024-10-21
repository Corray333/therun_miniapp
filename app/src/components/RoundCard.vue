<script lang="ts" setup>
import type { Round } from '@/types/types'
import { defineProps, ref, onBeforeMount } from 'vue'
import { useI18n } from 'vue-i18n'
import raceIcon from './icons/race-icon.vue'
import milesinIcon from './icons/milesin-icon.vue'

const { t } = useI18n()

const props = defineProps<{
    round: Round
}>()


const remainingTime = ref<string>('00:00:00')
const remainingSeconds = ref<number>(0)
const baseURL = import.meta.env.VITE_BASE_URL

const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600).toString().padStart(2, '0');
    const minutes = Math.floor((seconds % 3600) / 60).toString().padStart(2, '0');
    const secondsRemaining = (seconds % 60).toString().padStart(2, '0');
    return `${hours}:${minutes}:${secondsRemaining}`;
};

const calculateRemainingTimeAndPoints = () => {
    const now = Date.now() / 1000
    let secondsLeft = props.round.endTime - now
    remainingSeconds.value = secondsLeft

    if (secondsLeft <= 0) {
        remainingTime.value = '00:00:00'
        return
    }

    remainingTime.value = formatTime(Math.floor(secondsLeft))
}

onBeforeMount(() => {
    calculateRemainingTimeAndPoints()
    setInterval(calculateRemainingTimeAndPoints, 1000)
})

const showModal = ref<boolean>(false)

</script>

<template>

    <Transition name="slide-down">
        <section v-show="showModal" class="z-50 w-full h-screen bg-white fixed top-0 left-0">
            <section class=" relative p-4 flex flex-col gap-4">
                <i @click="showModal = false"
                        class=" pi fixed pi-times bg-dark text-white aspect-square p-1 rounded-full top-4 right-4"></i>
                
                <div class="flex flex-col gap-4 mt-10">
                    <img class=" rounded-2xl" :src="`${baseURL}/static/images/round/${round.element}-banner.png`" alt="">
                    <div class="header flex gap-4">
                        <img class="w-12" :src="`${baseURL}/static/images/round/${round.element}-icon.png`" alt="">
                        <h3 class="font-bold text-lg">{{ t(`screens.cars.round.header.${round.element}`) }}</h3>
                    </div>
                </div>

                <span>
                    <p class="font-bold inline">{{ t('screens.cars.round.description1.short') }}</p>
                    <span>{{ t('screens.cars.round.description1.part1') }}<milesinIcon class="inline scale-90"/></span>
                    <span>{{ t('screens.cars.round.description1.part2') }}<raceIcon class="inline scale-90"/></span>
                </span>

                <span class="">
                    <p class="font-bold inline">{{ t('screens.cars.round.description2.short') }}</p>
                    <p class="inline">{{ t('screens.cars.round.description2.text') }}</p>
                </span>
                
            </section>
        </section>
    </Transition>

    <section class=" round-card" @click="showModal = true">
        <div class="text-center">
            <p class=" font-mono font-bold text-xl">{{ remainingTime }}</p>
            <p class=" text-sm">{{ t('other.worldRound') }} #{{ round.id }}</p>
        </div>
        <div class="flex items-center gap-2">
            <img class=" w-12 object-contain aspect-square"
                :src="`${baseURL}/static/images/round/${round.element}-icon.png`" alt="">
            <i class="pi pi-chevron-right"></i>
        </div>
    </section>
</template>


<style scoped>
.slide-down-enter-active,
.slide-down-leave-active {
    transition: transform 0.5s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
    transform: translateY(100%);
}

.round-card {
    @apply flex w-full p-2 px-4 rounded-2xl justify-between items-center bg-secondary;
}
</style>