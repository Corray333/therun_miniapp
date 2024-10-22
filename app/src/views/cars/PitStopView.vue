<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import { computed, ref, onBeforeMount } from 'vue'
import raceIcon from '@/components/icons/race-icon.vue';

import { useAccountStore } from '@/stores/account'
import { useComponentsStore } from '@/stores/components'
import { useCarsStore } from '@/stores/cars'
import { useRoute, useRouter } from 'vue-router'
import { Car, RaceState } from '@/types/types'

const { t } = useI18n()
const accStore = useAccountStore()
const componentsStore = useComponentsStore()
const carsStore = useCarsStore()

onBeforeMount(async () => {
    setInterval(() => {
        countMiles(carsStore.mainCar.speed)
    }, 1000)
})

const miles = ref<number>(0)

const countMiles = (speed: number) => {
    let timeSpent = carsStore.raceState.startTime == 0 ? 0 : (Math.floor(new Date().getTime() / 1000) - carsStore.raceState.startTime)
    miles.value = speed * timeSpent/60/60 + carsStore.raceState.currentMiles
}

const fuel = computed(() => {
    const wasting = carsStore.mainCar.fuel-miles.value/carsStore.mainCar.fuelWasting
    return wasting > 0 ? wasting : 0
})

const health = computed(() => {
    const wasting = carsStore.mainCar.health-miles.value/carsStore.mainCar.healthWasting
    return wasting > 0 ? wasting : 0
})

</script>

<template>
    <section class="pb-20 p-4 pt-10 relative text-center flex flex-col gap-4">
        <router-link to="/cars"> <i class=" pi fixed pi-times bg-dark text-white aspect-square p-1 rounded-full top-4 right-4"></i></router-link>

        <img class="rounded-2xl w-full" src="../../assets/images/cars/pit-stop/pit-stop-banner.png" alt="">
        <div>
            <h3 class=" font-bold text-xl">{{ t('screens.cars.pitStop.header') }}</h3>
            <p>{{ t('screens.cars.pitStop.description') }}</p>
        </div>

        <div class="flex gap-4">
            <div class="card">
                <div class="circular-progress-bar flex relative items-center justify-center w-20">
                    <div class="icon-container z-10 absolute p-3 bg-white rounded-full">
                        <img class="block w-6" src="../../assets/images/cars/main/fuel-icon.png" alt="">
                    </div>
                    <svg style="transform: rotate(18deg); position: relative;"  width="64" height="64" viewBox="0 0 150 150">
                        <g transform="translate(75,75)" stroke="white" stroke-width="2">
                            <path d="M0 0 70 0A99 99 0 0 1 56.6 41.4Z" :fill="fuel>=4 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 56.6 41.4A99 99 0 0 1 21.2 67.9Z" :fill="fuel>=5 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 21.2 67.9A99 99 0 0 1 -21.2 67.9Z" :fill="fuel>=6 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 -21.2 67.9A99 99 0 0 1 -56.6 41.4Z" :fill="fuel>=7 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 -56.6 41.4A99 99 0 0 1 -70 0Z" :fill="fuel>=8 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 -70 0A99 99 0 0 1 -56.6 -41.4Z" :fill="fuel>=9 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 -56.6 -41.4A99 99 0 0 1 -21.2 -67.9Z" :fill="fuel>=10 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 -21.2 -67.9A99 99 0 0 1 21.2 -67.9Z" :fill="fuel>=1 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 21.2 -67.9A99 99 0 0 1 56.6 -41.4Z" :fill="fuel>=2 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 56.6 -41.4A99 99 0 0 1 70 0Z" :fill="fuel>=3 ? 'var(--green)' : 'var(--half-dark)'"/>
                        </g>
                    </svg>
                </div>
                <span>
                    <raceIcon class=" text-xl"/> 1 / {{ t('screens.cars.pitStop.slot') }}
                </span>
                <button>{{ t('screens.cars.pitStop.refuel') }}</button>
            </div>

            <div class="card">
                <div class="circular-progress-bar flex relative items-center justify-center w-20">
                    <div class="icon-container z-10 absolute p-3 bg-white rounded-full">
                        <img class="block w-6" src="../../assets/images/cars/main/health-icon.png" alt="">
                    </div>
                    <svg style="transform: rotate(18deg); position: relative;"  width="64" height="64" viewBox="0 0 150 150">
                        <g transform="translate(75,75)" stroke="white" stroke-width="2">
                            <path d="M0 0 70 0A99 99 0 0 1 56.6 41.4Z" :fill="health>=4 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 56.6 41.4A99 99 0 0 1 21.2 67.9Z" :fill="health>=5 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 21.2 67.9A99 99 0 0 1 -21.2 67.9Z" :fill="health>=6 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 -21.2 67.9A99 99 0 0 1 -56.6 41.4Z" :fill="health>=7 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 -56.6 41.4A99 99 0 0 1 -70 0Z" :fill="health>=8 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 -70 0A99 99 0 0 1 -56.6 -41.4Z" :fill="health>=9 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 -56.6 -41.4A99 99 0 0 1 -21.2 -67.9Z" :fill="health>=10 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 -21.2 -67.9A99 99 0 0 1 21.2 -67.9Z" :fill="health>=1 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 21.2 -67.9A99 99 0 0 1 56.6 -41.4Z" :fill="health>=2 ? 'var(--green)' : 'var(--half-dark)'"/>
                            <path d="M0 0 56.6 -41.4A99 99 0 0 1 70 0Z" :fill="health>=3 ? 'var(--green)' : 'var(--half-dark)'"/>
                        </g>
                    </svg>
                </div>
                <span>
                    <raceIcon class=" text-xl"/> 1 / {{ t('screens.cars.pitStop.slot') }}
                </span>
                <button>{{ t('screens.cars.pitStop.repair') }}</button>
            </div>
        </div>

        <div class="chibi-banner">
            <span>
                <p>{{ t('screens.cars.pitStop.chibiBanner.description1') }}</p>
                <raceIcon class="inline text-xl"/>
                <p>{{ t('screens.cars.pitStop.chibiBanner.description2') }}</p>
            </span>
            <router-link to="/chibi/city"><button>{{ t('screens.cars.pitStop.chibiBanner.goToChibiCityBtn') }}</button></router-link>
        </div>
    </section>
</template>


<style scoped>

.card{
    @apply w-full flex flex-col items-center gap-2 bg-half_dark p-4 rounded-2xl
}

.card span{
    @apply flex items-center gap-2 font-bold
}
.card span>*{
    @apply inline
}

.chibi-banner{
    @apply bg-half_dark p-4 rounded-2xl text-left flex flex-col gap-2
}

.chibi-banner>span *{
    @apply inline
}

</style>