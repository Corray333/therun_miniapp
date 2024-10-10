<script lang="ts" setup>

import { Requirement } from '@/types/types'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const baseURL = import.meta.env.VITE_BASE_URL


defineProps({
    building: {
        type: Requirement
    }
})

</script>

<template>
    <router-link :to="`/chibi/city/${building?.type}`">
        <section class="building" v-if="building">
            <img class="w-16" :src="`${baseURL}/static/images/buildings/${building.type}${building.level>0?building.level:1}.png`" alt="">
            <span class="w-full flex gap-2 items-center">{{ t(`screens.chibi.city.buildings.${building.type}.name`) }} <p class="level">{{ building.level }}</p></span>
            <i class="pi pi-chevron-right text-dark"></i>
        </section>
    </router-link>
</template>


<style scoped>


.building{
    @apply p-2 bg-white rounded-2xl flex items-center gap-2;
}

.level {
    @apply px-2 w-fit aspect-square flex items-center justify-center bg-custom_blue text-white rounded-md;
}

</style>