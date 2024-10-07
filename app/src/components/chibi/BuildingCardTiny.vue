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
    <section class="building" v-if="building">
        <img class="w-16" :src="`${baseURL}/static/images/buildings/${building.type}${building.level>0?building.level:1}.png`" alt="">
        <p class="w-full">{{ t(`screens.chibi.city.buildings.${building.type}.name`) }} lvl.{{ building.level }}</p>
        <i class="pi pi-chevron-right text-dark"></i>
    </section>
</template>


<style scoped>


.building{
    @apply p-2 bg-white rounded-2xl flex items-center gap-2;
}

</style>