<script lang="ts" setup>

import { useI18n } from 'vue-i18n'
import { ref } from 'vue'
import MilesinIcon from '../icons/milesin-icon.vue'

const { t } = useI18n()

defineProps({
    speed: {
        type: Number,
        required: true
    }
})

const showModal = ref<boolean>(false)

</script>

<template>
    <Transition name="slide-down">
        <section v-show="showModal" class="z-50 w-full h-screen bg-white fixed top-0 left-0">
            <section class=" modal relative p-4 flex flex-col gap-4">
                <i @click="showModal = false"
                        class=" pi fixed pi-times bg-dark text-white aspect-square p-1 rounded-full top-4 right-4"></i>
                
                <div class="flex flex-col gap-4 mt-10">
                    <img class=" rounded-2xl" src="../../assets/images/cars/characteristics/speed-header.png" alt="">
                    <div class="header flex gap-4">
                        <div class="shield">
                            <span>
                                <MilesinIcon class="text-xl mr-1"/>
                                <p>{{ speed.toFixed(1) }}</p>
                                <p> / {{ t('screens.cars.characteristics.speed.shortHours') }}</p>
                            </span>
                            <p>{{ t('screens.cars.characteristics.speed.name') }}</p>
                        </div>

                        <div class="shield">
                            <span>
                                <MilesinIcon class="text-xl mr-1"/>
                                <p>{{ (speed*26).toFixed(1) }}</p>
                            </span>
                            <p>{{ t('screens.cars.characteristics.speed.maxPerRound') }}</p>
                        </div>
                    </div>
                </div>

                <div class="info">
                    <h3>{{ t('screens.cars.characteristics.speed.name') }}</h3>
                    <span>
                        <p>{{ t('screens.cars.characteristics.speed.description1') }}</p>
                        <MilesinIcon class="text-xl"/>
                        <p>{{ t('screens.cars.characteristics.speed.description2') }}</p>
                    </span>
                </div>
            </section>
        </section>
    </Transition>

    <div @click="showModal = true" class="speed-shield">
        <span class="flex items-center font-bold">
            <MilesinIcon class="w-6 h-6 mr-1 text-xl" />
            <p>{{ speed.toFixed(1) }}</p>
            <p> / {{ t('screens.cars.characteristics.speed.shortHours') }}</p>
        </span>
        <span class=" label gap-2">
            <p>{{ t('screens.cars.characteristics.speed.name') }}</p>
            <i class="pi pi-chevron-right" style="font-size: 0.75rem;"></i>
        </span>
    </div>
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

.speed-shield {
    @apply flex flex-col items-center justify-between bg-blue-50 rounded-2xl p-2 px-4;
}

.shield *{
    @apply inline;
}

.label{
    @apply flex items-center gap-1;
}

.header>.shield{
    @apply bg-half_dark rounded-2xl p-2 px-4 flex flex-col items-center justify-between w-full;
}
.header>.shield>span{
    @apply font-bold;
}
.header>.shield>p{
    @apply label;
}

.info{
    @apply mt-4;
}
.info>h3{
    @apply text-xl font-bold;
}
.info>span *{
    @apply inline;
}

</style>