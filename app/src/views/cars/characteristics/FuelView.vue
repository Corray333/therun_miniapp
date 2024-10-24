<script lang="ts" setup>

import { useCarsStore } from '@/stores/cars'
import { useI18n } from 'vue-i18n'
import TuningCard from './TuningCard.vue'
import { computed, onBeforeMount, ref } from 'vue'
import MilesinIcon from '@/components/icons/milesin-icon.vue';

const { t } = useI18n()
const carsStore = useCarsStore()

const characteristic = 'fuel'

const constModule = computed(()=>{
    return carsStore.mainCar.modules?.find(module => module.isTemp === false && module.characteristic === characteristic)
})

const tempModule = computed(()=>{
    return carsStore.mainCar.modules?.find(module => module.isTemp === true && module.characteristic === characteristic)
})


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

</script>

<template>
    <section class=" pb-20 p-4">
        <div class="header">
            <img src="../../../assets/images/cars/characteristics/fuel.png" alt="">
            <h1>{{ t('screens.cars.characteristics.fuel.name') }}</h1>
            <p class="font-bold">{{ carsStore.mainCar.tank }}/100</p>
            <p>{{ t('screens.cars.characteristics.fuel.description') }}</p>
        </div>
        <div class="flex gap-2">
            <div class="shield">
                <div class="span">
                    <MilesinIcon class="text-xl mr-1"/>
                    <p>{{ carsStore.mainCar.fuelWasting.toFixed(1) }} / {{  t('screens.cars.characteristics.strength.slot')  }}</p>
                </div>
                <p class="label">{{  t('screens.cars.characteristics.strength.wasting')  }}</p>
            </div>

            <div class="shield">
                <div class="span">
                    <MilesinIcon class="text-xl mr-1"/>
                    <p>{{ (fuel*carsStore.mainCar.fuelWasting).toFixed(1) }}</p>
                </div>
                <p class="label">{{  t('screens.cars.characteristics.strength.remaining')  }}</p>
            </div>
        </div>
        <h2 class="mt-4 mb-2">{{ t('screens.cars.characteristics.tuning.header') }}</h2>
        <div class="tuning">
            <TuningCard v-if="constModule" :module="constModule"/>
            <div v-else class="no-module tuning-card">
                <i class="pi pi-plus text-custom_blue" style="font-size: 2rem;"></i>
                <p class="label">{{ t(`screens.cars.characteristics.tuning.constModule`) }}</p>
            </div>
            <TuningCard v-if="tempModule" :module="tempModule"/>
            <div v-else class="no-module tuning-card">
                <i class="pi pi-plus text-secondary" style="font-size: 2rem;"></i>
                <p class="label">{{ t(`screens.cars.characteristics.tuning.tempModule`) }}</p>
            </div>
        </div>
    </section>
</template>


<style scoped>

.header{
    @apply flex flex-col items-center gap-2;
}

.header>p{
    @apply w-4/5 text-center;
}

.header>img{
    @apply w-2/3;
}

.tuning{
    @apply flex gap-2;
}

.no-module{
    @apply flex flex-col items-center justify-center gap-2 p-4 h-full w-full text-center;
}

.tuning-card{
    @apply flex flex-col w-full rounded-2xl bg-half_dark aspect-square;
}

.shield{
    @apply bg-half_dark p-2 w-full rounded-2xl flex flex-col items-center justify-center
}
.shield p,svg{
    @apply inline
}

</style>