<script lang="ts" setup>

import { defineProps, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const baseURL = import.meta.env.VITE_BASE_URL

defineProps({
    currentElement: String
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
                    <img class=" rounded-2xl" :src="`${baseURL}/static/images/cars/${currentElement}-banner.png`"
                        alt="">
                    <div class="header flex gap-4">
                        <img class="w-12" :src="`${baseURL}/static/images/round/${currentElement}-icon.png`" alt="">
                        <h3 class="font-bold text-lg">{{ t(`screens.cars.carTypes.header.${currentElement}`) }}</h3>
                    </div>
                </div>

                <span class="">
                    <p class="font-bold inline">{{ t('screens.cars.carTypes.note') }}</p>
                    <p class="inline">{{ t('screens.cars.carTypes.description') }}</p>
                </span>

            </section>
        </section>
    </Transition>

    <div class="shield" @click="showModal = true">
        <img :src="`${baseURL}/static/images/round/${currentElement}-icon.png`" alt="">
        <span>
            <p>{{ t('screens.cars.carType') }}</p>
            <i class="pi pi-chevron-right" style="color:var(--dark); font-size: 0.75rem;"></i>
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

.shield {
    @apply w-full bg-half_dark rounded-2xl p-2 px-4 flex flex-col gap-2 items-center;
}

.shield>span {
    @apply label items-center flex;
}

.shield>img {
    @apply w-8 h-full;
}
</style>