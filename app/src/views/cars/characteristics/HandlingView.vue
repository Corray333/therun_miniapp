<script lang="ts" setup>

import { useCarsStore } from '@/stores/cars'
import { useI18n } from 'vue-i18n'
import TuningCard from './TuningCard.vue'
import { computed } from 'vue'

const { t } = useI18n()
const carsStore = useCarsStore()

const characteristic = 'handling'

const constModule = computed(()=>{
    return carsStore.mainCar.modules?.find(module => module.isTemp === false && module.characteristic === characteristic)
})

const tempModule = computed(()=>{
    return carsStore.mainCar.modules?.find(module => module.isTemp === true && module.characteristic === characteristic)
})

</script>

<template>
    <section class=" pb-20 p-4">
        <div class="header">
            <img src="../../../assets/images/cars/characteristics/handling.png" alt="">
            <h1>{{ t('screens.cars.characteristics.handling.name') }}</h1>
            <p class="font-bold">{{ carsStore.mainCar.handling }}/100</p>
            <p>{{ t('screens.cars.characteristics.tuning.description') }}</p>
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

</style>