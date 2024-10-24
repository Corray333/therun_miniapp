<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import { computed, ref, onBeforeMount } from 'vue'
import raceIcon from '@/components/icons/race-icon.vue';

import { useAccountStore } from '@/stores/account'
import { useComponentsStore } from '@/stores/components'
import { useCarsStore } from '@/stores/cars'
import { useRoute, useRouter } from 'vue-router'
import { Car, RaceState } from '@/types/types'

import SpeedShield from '@/components/cars/SpeedShield.vue'
import CarTypeShield from '@/components/cars/CarTypeShield.vue'

const { t } = useI18n()
const accStore = useAccountStore()
const componentsStore = useComponentsStore()
const carsStore = useCarsStore()

const router = useRouter()

const totalCarScore = computed(() => {
    return carsStore.mainCar.acceleration + carsStore.mainCar.strength + carsStore.mainCar.handling + carsStore.mainCar.tank + carsStore.mainCar.brakes
})

</script>

<template>
    <section class="pb-20 p-4 relative text-center flex flex-col gap-4">
        <div class="top-row">
            <SpeedShield :speed="carsStore.mainCar.speed" />
            <div class="shield">
                <p class="font-bold">{{ totalCarScore }}</p>
                <p class="label">{{ t('screens.cars.characteristics.totalScore') }}</p>
            </div>
            <CarTypeShield :current-element="carsStore.mainCar.element" />
        </div>

        <div class="banner tuning-shop-banner">
            <span>
                <p>{{ t('screens.cars.characteristics.tuningShop') }}</p>
                <i class="pi pi-chevron-right"></i>
            </span>
        </div>
        
        <div @click="router.push('/cars/tuning')" class="banner inventory-banner">
            <span>
                <p>{{ t('screens.cars.characteristics.myInventory') }}</p>
                <i class="pi pi-chevron-right"></i>
            </span>
        </div>

        <div class="characteristics">
            <div @click="router.push('/cars/characteristics/acceleration')" class="characteristic-card acceleration">
                <div class="left">
                    <img src="../../assets/images/cars/characteristics/acceleration.png" alt="">
                    <p>{{ t('screens.cars.characteristics.acceleration.name') }}</p>
                </div>
                <div class="right">
                    <p>{{ carsStore.mainCar.acceleration }}/100</p>
                    <i class="pi pi-chevron-right text-dark"></i>
                </div>
            </div>

            <div @click="router.push('/cars/characteristics/handling')" class="characteristic-card handling">
                <div class="left">
                    <img src="../../assets/images/cars/characteristics/handling.png" alt="">
                    <p>{{ t('screens.cars.characteristics.handling.name') }}</p>
                </div>
                <div class="right">
                    <p>{{ carsStore.mainCar.handling }}/100</p>
                    <i class="pi pi-chevron-right text-dark"></i>
                </div>
            </div>

            <div @click="router.push('/cars/characteristics/brakes')" class="characteristic-card brakes">
                <div class="left">
                    <img src="../../assets/images/cars/characteristics/brakes.png" alt="">
                    <p>{{ t('screens.cars.characteristics.brakes.name') }}</p>
                </div>
                <div class="right">
                    <p>{{ carsStore.mainCar.brakes }}/100</p>
                    <i class="pi pi-chevron-right text-dark"></i>
                </div>
            </div>

            <div @click="router.push('/cars/characteristics/strength')" class="characteristic-card strength">
                <div class="left">
                    <img src="../../assets/images/cars/characteristics/strength.png" alt="">
                    <p>{{ t('screens.cars.characteristics.strength.name') }}</p>
                </div>
                <div class="right">
                    <p>{{ carsStore.mainCar.strength }}/100</p>
                    <i class="pi pi-chevron-right text-dark"></i>
                </div>
            </div>

            <div @click="router.push('/cars/characteristics/fuel')" class="characteristic-card fuel">
                <div class="left">
                    <img src="../../assets/images/cars/characteristics/fuel.png" alt="">
                    <p>{{ t('screens.cars.characteristics.fuel.name') }}</p>
                </div>
                <div class="right">
                    <p>{{ carsStore.mainCar.tank }}/100</p>
                    <i class="pi pi-chevron-right text-dark"></i>
                </div>
            </div>
        </div>

    </section>
</template>


<style scoped>

.top-row {
    @apply flex gap-2 w-full
}

.top-row>* {
    @apply w-full;
}

.shield {
    @apply bg-half_dark rounded-2xl p-2 px-4 flex flex-col items-center justify-between;
}

.banner{
    @apply p-4 rounded-2xl text-white font-bold bg-cover bg-top;
}
.banner span{
    @apply flex items-center text-lg mt-12;
}
.tuning-shop-banner{
    background-image: url(../../assets/images/cars/characteristics/tuning-shop-banner.png);
}
.inventory-banner{
    background-image: url(../../assets/images/cars/characteristics/inventory-banner.png);
}


.characteristics{
    @apply flex flex-col gap-2
}
.characteristic-card{
    @apply flex justify-between items-center w-full bg-half_dark rounded-2xl p-4;
}

.characteristic-card>.left>img{
    @apply h-12;
}
.characteristic-card>div{
    @apply flex gap-2 items-center font-bold
}


</style>