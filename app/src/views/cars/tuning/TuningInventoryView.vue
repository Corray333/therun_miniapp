<script lang="ts" setup>

import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const { t } = useI18n()
const router = useRouter()

</script>

<template>
    <section class="pb-20 min-h-screen tuning">
        <h1>{{ t('screens.cars.tuning.header') }}</h1>
        <div class="inventory">
            <div @click="router.push('/cars/tuning/acceleration')" class="characteristic-card">
                <img src="../../../assets/images/cars/characteristics/acceleration.png" alt="">
                <p>{{ t('screens.cars.characteristics.acceleration.name') }}</p>
            </div>
    
            <div @click="router.push('/cars/tuning/handling')" class="characteristic-card">
                <img src="../../../assets/images/cars/characteristics/handling.png" alt="">
                <p>{{ t('screens.cars.characteristics.handling.name') }}</p>
            </div>
    
            <div @click="router.push('/cars/tuning/brakes')" class="characteristic-card">
                <img src="../../../assets/images/cars/characteristics/brakes.png" alt="">
                <p>{{ t('screens.cars.characteristics.brakes.name') }}</p>
            </div>
    
            <div @click="router.push('/cars/tuning/strength')" class="characteristic-card">
                <img src="../../../assets/images/cars/characteristics/strength.png" alt="">
                <p>{{ t('screens.cars.characteristics.strength.name') }}</p>
            </div>
    
            <div @click="router.push('/cars/tuning/fuel')" class="characteristic-card">
                <img src="../../../assets/images/cars/characteristics/fuel.png" alt="">
                <p>{{ t('screens.cars.characteristics.fuel.name') }}</p>
            </div>
        </div>
    </section>
</template>


<style scoped>

.tuning{
    @apply flex flex-col items-center gap-4 p-4;
}

.inventory{
    @apply grid grid-cols-2 gap-4;
}

.characteristic-card{
    @apply flex flex-col items-center gap-2 p-4 rounded-2xl bg-half_dark;
}

</style>