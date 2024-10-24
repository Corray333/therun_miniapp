<script lang="ts" setup>
import { ref, onBeforeMount, computed } from 'vue'
import { useAccountStore } from '@/stores/account'
import { useI18n } from 'vue-i18n'
import axios, { isAxiosError } from 'axios'
import { auth } from '@/utils/helpers'
import { useComponentsStore } from '@/stores/components'
import { useCarsStore } from '@/stores/cars'
import { useRoute, useRouter } from 'vue-router'

import milesinIcon from '@/components/icons/milesin-icon.vue'
import raceIcon from '@/components/icons/race-icon.vue'

import type { Car, RaceState } from '@/types/types'

import roundCard from '@/components/RoundCard.vue'

const componentsStore = useComponentsStore()
const carsStore = useCarsStore()


const { t } = useI18n()

const accStore = useAccountStore()
const baseURL = import.meta.env.VITE_BASE_URL

const route = useRoute()
const router = useRouter()
const element = route.query.element


const getMainCar = async (): Promise<boolean> => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/cars/main`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        if (data) {
            carsStore.mainCar = data
            return true
        } else {
            return false
        }
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            try {
                await auth()
                return await getMainCar()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
                return false
            }
        } else {
            return false
        }
    }
}

const getRaceState = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/race`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        if (data) {
            carsStore.raceState = data
        }
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            try {
                await auth()
                await getMainCar()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        }
    }
}

onBeforeMount(async () => {
    if (carsStore.mainCar != null){
        loaded.value = true
    }
    getRound()
    if (!await getMainCar()) {
        router.push('/cars/all')
    }
    await getRaceState()
    countMiles(carsStore.mainCar.speed)
    loaded.value = true
    setInterval(() => {
        countMiles(carsStore.mainCar.speed)
    }, 500)
})


const getRound = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/round`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        componentsStore.round = data
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            try {
                await auth()
                await getRound()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        }
    }
}

const startRace = async () => {
    try {
        const { data } = await axios.post(`${import.meta.env.VITE_API_URL}/start-race`, {}, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        carsStore.raceState = data
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await startRace()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        }
    }
}

const endRace = async () => {
    try {
        const { data } = await axios.post(`${import.meta.env.VITE_API_URL}/end-race`, {}, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        carsStore.raceState = data
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await endRace()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        }
    }
}

const miles = ref<number>(0)

const countMiles = (speed: number) => {
    let timeSpent = carsStore.raceState.startTime == 0 ? 0 : (Math.floor(new Date().getTime() / 500) - carsStore.raceState.startTime*2)
    miles.value = (speed * timeSpent / 60 / 60 / 2 + carsStore.raceState.currentMiles)
}

const fuel = computed(() => {
    const wasting = carsStore.mainCar.fuel - miles.value / carsStore.mainCar.fuelWasting
    return wasting > 0 ? wasting : 0
})

const health = computed(() => {
    const wasting = carsStore.mainCar.health - miles.value / carsStore.mainCar.healthWasting
    return wasting > 0 ? wasting : 0
})

const loaded = ref<boolean>(false)

</script>

<template>
    <div v-if="!loaded" class="w-full flex justify-center p-4">
        <i  class="pi pi-spinner pi-spin"></i>
    </div>
    <section v-else class=" pb-20 flex flex-col h-screen justify-between">
        <div class="p-4">
            <div class=" shields">
                <roundCard :round="componentsStore.round" />
                <div @click="router.push('/cars/pit-stop')" class="shield2">
                    <div>
                        <p class="w-max">{{ t('screens.cars.pitStop.pitStop') }}</p>
                    </div>
                    <img src="../components/icons/pit-stop-icon.png" alt="">
                </div>
            </div>
        </div>

        <div class="flex items-center gap-2 font-mono justify-center">
            <milesinIcon class=" text-3xl" />
            <p class=" text-2xl font-bold">{{ miles.toFixed(2) }}</p>
        </div>

        <div class=" main-shields flex justify-between w-full">
            <div class="shield flex items-center gap-2">
                <div class=" text-center">
                    <div class="flex gap-2 items-center">
                        <raceIcon class=" text-2xl" />
                        <p class="text-xl font-bold">{{ carsStore.raceState.place }}</p>
                    </div>
                    <p>{{ t('screens.cars.main.prize') }}</p>
                </div>
                <i class=" pi pi-chevron-right"></i>
            </div>

            <div class="shield flex items-center  gap-2">
                <div class=" text-center">
                    <div class="flex gap-2 items-center">
                        <p class="text-xl font-bold">x {{ carsStore.raceState.place }}</p>
                    </div>
                </div>
                <i class=" pi pi-chevron-right"></i>
            </div>
        </div>

        <div class="p-4">
            <img id="car" class="p-4" :src="carsStore.mainCar.img" alt="">
        </div>

        <div class="p-4 flex flex-col gap-2">
            <div class="botton-widgets flex gap-2">
                <div class="circular-progress-bar flex relative items-center justify-center w-20">
                    <div class="icon-container z-10 absolute p-3 bg-white rounded-full">
                        <img class="block w-6" src="../assets/images/cars/main/fuel-icon.png" alt="">
                    </div>
                    <svg style="transform: rotate(18deg); position: relative;" width="64" height="64"
                        viewBox="0 0 150 150">
                        <g transform="translate(75,75)" stroke="white" stroke-width="2">
                            <path d="M0 0 70 0A99 99 0 0 1 56.6 41.4Z"
                                :fill="fuel >= 4 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 56.6 41.4A99 99 0 0 1 21.2 67.9Z"
                                :fill="fuel >= 5 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 21.2 67.9A99 99 0 0 1 -21.2 67.9Z"
                                :fill="fuel >= 6 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 -21.2 67.9A99 99 0 0 1 -56.6 41.4Z"
                                :fill="fuel >= 7 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 -56.6 41.4A99 99 0 0 1 -70 0Z"
                                :fill="fuel >= 8 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 -70 0A99 99 0 0 1 -56.6 -41.4Z"
                                :fill="fuel >= 9 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 -56.6 -41.4A99 99 0 0 1 -21.2 -67.9Z"
                                :fill="fuel >= 10 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 -21.2 -67.9A99 99 0 0 1 21.2 -67.9Z"
                                :fill="fuel >= 1 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 21.2 -67.9A99 99 0 0 1 56.6 -41.4Z"
                                :fill="fuel >= 2 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 56.6 -41.4A99 99 0 0 1 70 0Z"
                                :fill="fuel >= 3 ? 'var(--green)' : 'var(--half-dark)'" />
                        </g>
                    </svg>
                </div>

                <div class="circular-progress-bar flex relative items-center justify-center w-20">
                    <div class="icon-container z-10 absolute p-3 bg-white rounded-full">
                        <img class="block w-6" src="../assets/images/cars/main/health-icon.png" alt="">
                    </div>
                    <svg style="transform: rotate(18deg); position: relative;" width="64" height="64"
                        viewBox="0 0 150 150">
                        <g transform="translate(75,75)" stroke="white" stroke-width="2">
                            <path d="M0 0 70 0A99 99 0 0 1 56.6 41.4Z"
                                :fill="health >= 4 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 56.6 41.4A99 99 0 0 1 21.2 67.9Z"
                                :fill="health >= 5 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 21.2 67.9A99 99 0 0 1 -21.2 67.9Z"
                                :fill="health >= 6 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 -21.2 67.9A99 99 0 0 1 -56.6 41.4Z"
                                :fill="health >= 7 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 -56.6 41.4A99 99 0 0 1 -70 0Z"
                                :fill="health >= 8 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 -70 0A99 99 0 0 1 -56.6 -41.4Z"
                                :fill="health >= 9 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 -56.6 -41.4A99 99 0 0 1 -21.2 -67.9Z"
                                :fill="health >= 10 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 -21.2 -67.9A99 99 0 0 1 21.2 -67.9Z"
                                :fill="health >= 1 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 21.2 -67.9A99 99 0 0 1 56.6 -41.4Z"
                                :fill="health >= 2 ? 'var(--green)' : 'var(--half-dark)'" />
                            <path d="M0 0 56.6 -41.4A99 99 0 0 1 70 0Z"
                                :fill="health >= 3 ? 'var(--green)' : 'var(--half-dark)'" />
                        </g>
                    </svg>
                </div>

                <div class="shield2 w-full">
                    <div>
                        <p class="w-max">{{ t('screens.cars.first.rating') }}</p>
                    </div>
                    <img src="../components/icons/rating-icon.png" alt="">
                </div>
                <div @click="router.push('/cars/characteristics')" class="shield2 w-full">
                    <div>
                        <p class="w-max">{{ t('screens.cars.main.characteristics') }}</p>
                    </div>
                    <img src="../components/icons/characteristics-icon.png" alt="">
                </div>

            </div>

            <button @click="startRace" v-if="carsStore.raceState.startTime == 0">{{ t('screens.cars.main.startBtn')
                }}</button>
            <button @click="endRace" v-else>{{ t('screens.cars.main.stopBtn') }}</button>
        </div>
    </section>
</template>


<style scoped>
@keyframes wiggle {
    0% {
        transform: translateY(0) rotate(0deg);
    }

    1% {
        transform: translateY(-5px) rotate(0deg);
    }

    2% {
        transform: translateY(0) rotate(0deg);
    }

    3% {
        transform: translateY(-5px) rotate(0deg);
    }

    4% {
        transform: translateY(0) rotate(0deg);
    }

    5% {
        transform: translateY(-5px) rotate(0deg);
    }

    6% {
        transform: translateY(0) rotate(0deg);
    }

    7% {
        transform: translateY(-5px) rotate(0deg);
    }

    8% {
        transform: translateY(0) rotate(0deg);
    }

    9% {
        transform: translateY(-5px) rotate(0deg);
    }

    10% {
        transform: translateY(-5px) rotate(-2deg);
    }

    11% {
        transform: translateY(0) rotate(0deg);
    }

    12% {
        transform: translateY(-5px) rotate(0deg);
    }

    13% {
        transform: translateY(0) rotate(0deg);
    }

    14% {
        transform: translateY(-5px) rotate(0deg);
    }

    15% {
        transform: translateY(0) rotate(0deg);
    }

    16% {
        transform: translateY(-5px) rotate(0deg);
    }

    17% {
        transform: translateY(0) rotate(0deg);
    }

    18% {
        transform: translateY(-5px) rotate(0deg);
    }

    19% {
        transform: translateY(0) rotate(0deg);
    }

    20% {
        transform: translateY(-5px) rotate(0deg);
    }

    21% {
        transform: translateY(0) rotate(0deg);
    }

    22% {
        transform: translateY(-5px) rotate(0deg);
    }

    23% {
        transform: translateY(0) rotate(0deg);
    }

    24% {
        transform: translateY(-5px) rotate(0deg);
    }

    25% {
        transform: translateY(0) rotate(0deg);
    }

    26% {
        transform: translateY(-5px) rotate(0deg);
    }

    27% {
        transform: translateY(0) rotate(0deg);
    }

    28% {
        transform: translateY(-10px) rotate(2deg);
    }

    29% {
        transform: translateY(0) rotate(0deg);
    }

    30% {
        transform: translateY(0) rotate(0deg);
    }

    31% {
        transform: translateY(-5px) rotate(0deg);
    }

    32% {
        transform: translateY(0) rotate(0deg);
    }

    33% {
        transform: translateY(-5px) rotate(0deg);
    }

    34% {
        transform: translateY(0) rotate(0deg);
    }

    35% {
        transform: translateY(-5px) rotate(0deg);
    }

    36% {
        transform: translateY(0) rotate(0deg);
    }

    37% {
        transform: translateY(-5px) rotate(0deg);
    }

    38% {
        transform: translateY(0) rotate(0deg);
    }

    39% {
        transform: translateY(-5px) rotate(0deg);
    }

    40% {
        transform: translateY(0) rotate(0deg);
    }

    41% {
        transform: translateY(-5px) rotate(0deg);
    }

    42% {
        transform: translateY(0) rotate(0deg);
    }

    43% {
        transform: translateY(-5px) rotate(0deg);
    }

    44% {
        transform: translateY(0) rotate(0deg);
    }

    45% {
        transform: translateY(-5px) rotate(0deg);
    }

    46% {
        transform: translateY(0) rotate(0deg);
    }

    47% {
        transform: translateY(-5px) rotate(0deg);
    }

    48% {
        transform: translateY(0) rotate(0deg);
    }

    49% {
        transform: translateY(-5px) rotate(0deg);
    }

    50% {
        transform: translateY(-15px) rotate(-5deg);
    }

    51% {
        transform: translateY(0) rotate(0deg);
    }

    52% {
        transform: translateY(-5px) rotate(0deg);
    }

    53% {
        transform: translateY(0) rotate(0deg);
    }

    54% {
        transform: translateY(-5px) rotate(0deg);
    }

    55% {
        transform: translateY(0) rotate(0deg);
    }

    56% {
        transform: translateY(-5px) rotate(0deg);
    }

    57% {
        transform: translateY(0) rotate(0deg);
    }

    58% {
        transform: translateY(-5px) rotate(0deg);
    }

    59% {
        transform: translateY(0) rotate(0deg);
    }

    60% {
        transform: translateY(-5px) rotate(0deg);
    }

    61% {
        transform: translateY(0) rotate(0deg);
    }

    62% {
        transform: translateY(-5px) rotate(0deg);
    }

    63% {
        transform: translateY(0) rotate(0deg);
    }

    64% {
        transform: translateY(-5px) rotate(0deg);
    }

    65% {
        transform: translateY(0) rotate(0deg);
    }

    66% {
        transform: translateY(-5px) rotate(0deg);
    }

    67% {
        transform: translateY(0) rotate(0deg);
    }

    68% {
        transform: translateY(-5px) rotate(0deg);
    }

    69% {
        transform: translateY(0) rotate(0deg);
    }

    70% {
        transform: translateY(-20px) rotate(5deg);
    }

    71% {
        transform: translateY(0) rotate(0deg);
    }

    72% {
        transform: translateY(-5px) rotate(0deg);
    }

    73% {
        transform: translateY(0) rotate(0deg);
    }

    74% {
        transform: translateY(-5px) rotate(0deg);
    }

    75% {
        transform: translateY(0) rotate(0deg);
    }

    76% {
        transform: translateY(-5px) rotate(0deg);
    }

    77% {
        transform: translateY(0) rotate(0deg);
    }

    78% {
        transform: translateY(-5px) rotate(0deg);
    }

    79% {
        transform: translateY(0) rotate(0deg);
    }

    80% {
        transform: translateY(-5px) rotate(0deg);
    }

    81% {
        transform: translateY(0) rotate(0deg);
    }

    82% {
        transform: translateY(-5px) rotate(0deg);
    }

    83% {
        transform: translateY(0) rotate(0deg);
    }

    84% {
        transform: translateY(-5px) rotate(0deg);
    }

    85% {
        transform: translateY(0) rotate(0deg);
    }

    86% {
        transform: translateY(-5px) rotate(0deg);
    }

    87% {
        transform: translateY(0) rotate(0deg);
    }

    88% {
        transform: translateY(-5px) rotate(0deg);
    }

    89% {
        transform: translateY(0) rotate(0deg);
    }

    90% {
        transform: translateY(-25px) rotate(-10deg);
    }

    91% {
        transform: translateY(0) rotate(0deg);
    }

    92% {
        transform: translateY(-5px) rotate(0deg);
    }

    93% {
        transform: translateY(0) rotate(0deg);
    }

    94% {
        transform: translateY(-5px) rotate(0deg);
    }

    95% {
        transform: translateY(0) rotate(0deg);
    }

    96% {
        transform: translateY(-5px) rotate(0deg);
    }

    97% {
        transform: translateY(0) rotate(0deg);
    }

    98% {
        transform: translateY(-5px) rotate(0deg);
    }

    99% {
        transform: translateY(0) rotate(0deg);
    }

    100% {
        transform: translateY(0) rotate(0deg);
    }
}

#car {
    animation: wiggle 10s infinite;
}

.shields {
    @apply flex gap-2
}

.shields>* {
    @apply sm-shadow
}

.shield {
    @apply bg-half_dark rounded-2xl p-2 px-4
}

.shield2 {
    @apply bg-half_dark text-sm rounded-2xl flex flex-col items-center
}

.shield2>div {
    @apply px-4 py-2;
}

.shield2>img {
    @apply h-6 object-cover object-top w-12
}

.main-shields>.shield:first-child {
    @apply rounded-l-none
}

.main-shields>.shield:last-child {
    @apply rounded-r-none
}
</style>