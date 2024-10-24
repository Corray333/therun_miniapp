<script lang="ts" setup>

import { Module } from '@/types/types'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps({
    module: {
        type: Object as () => Module,
        required: true
    }
})

</script>

<template>
    <div class="tuning-card">
        <div v-if="module != undefined" class="module">
            <img :src="module.img" alt="">
            <div class="info">
                <p>{{ module.name }}</p>
                <p>+{{ module.boost-100 }}% {{ t(`screens.cars.characteristics.${module.characteristic}.name`) }}</p>
            </div>
        </div>
    </div>
</template>


<style scoped>

.tuning-card{
    @apply flex flex-col w-full rounded-2xl bg-half_dark aspect-square;
}

.no-module{
    @apply flex flex-col items-center justify-center gap-2 p-4 h-full w-full text-center;
}

.module>img{
    @apply w-full rounded-2xl;
}

.info{
    @apply flex flex-col justify-center p-4;
}

</style>