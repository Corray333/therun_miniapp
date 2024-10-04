<script lang="ts" setup>

import { Building } from '@/types/types'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const baseURL = import.meta.env.VITE_BASE_URL


defineProps({
    building: {
        type: Building
    }
})

</script>

<template>
    <section class="building" v-if="building">
        <div>
            <p>
                <p class=" font-bold">{{ t(`screens.chibi.city.buildings.${building.type}.name`) }}</p>
                <p class=" label">{{ t(`screens.chibi.city.buildings.${building.type}.description`) }}</p>
            </p>
            <button class=" btn-type-4 state-2">
                <p v-if="building.level==0">{{ t(`screens.chibi.city.buildBtn`) }}</p>
                <p v-else>{{ t(`screens.chibi.city.upgradeBtn`) }}</p>
            </button>
        </div>
        <img :src="`${baseURL}/static/images/buildings/${building.type}${building.level>0?building.level:1}.png`" alt="">
    </section>
</template>


<style scoped>

.building{
    @apply flex gap-10 bg-half_dark rounded-2xl p-4
}
.building>div{
    @apply w-full flex flex-col justify-between;
}
.building>img{
    @apply h-24 object-contain scale-150 origin-bottom-right -mr-2;
}

</style>