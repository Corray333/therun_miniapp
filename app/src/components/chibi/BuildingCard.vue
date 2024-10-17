<script lang="ts" setup>

import { Building } from '@/types/types'
import { useI18n } from 'vue-i18n'
import { ref, onBeforeMount } from 'vue'

const { t } = useI18n()
const baseURL = import.meta.env.VITE_BASE_URL


const props = defineProps({
    building: {
        type: Building
    }
})


const remainingTime = ref<string>('00:00:00')
const remainingSeconds = ref<number>(0)

const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600).toString().padStart(2, '0');
    const minutes = Math.floor((seconds % 3600) / 60).toString().padStart(2, '0');
    const secondsRemaining = (seconds % 60).toString().padStart(2, '0');
    return `${hours}:${minutes}:${secondsRemaining}`;
};

const calculateRemainingTimeAndPoints = () => {
    if (!props.building){
        return
    }
    const now = Date.now() / 1000
    let secondsLeft = props.building.stateUntil - now
    remainingSeconds.value = secondsLeft

    if (secondsLeft <= 0) {
        remainingTime.value = '00:00:00'
        return
    }

    remainingTime.value = formatTime(Math.floor(secondsLeft))
}

onBeforeMount(() => {
    calculateRemainingTimeAndPoints()
    setInterval(calculateRemainingTimeAndPoints, 1000)
})

</script>

<template>
    <section class="building" v-if="building">
        <div>
            <div>
                <p class=" font-bold">{{ t(`screens.chibi.city.buildings.${building.type}.name`) }}</p>
                <p class=" label">{{ t(`screens.chibi.city.buildings.${building.type}.description`) }}</p>
            </div>
            <button v-if="building.level==0" class=" btn-type-4 bg-custom_blue">
                <p>{{ t(`screens.chibi.city.buildBtn`) }}</p>
            </button>
            <div v-else class="flex gap-2">
                <p class="level">{{ building.level  }}</p>
                <button v-if="building.upgradeCost != null || remainingSeconds > 0" class=" btn-type-4 state-2">
                    <p v-if="remainingSeconds < 0" class="flex gap-1 justify-center items-center">{{ t(`screens.chibi.city.upgradeBtn`) }}<i class=" pi pi-arrow-up"></i></p>
                    <p v-else class="flex items-center gap-2 font-mono">{{ remainingTime }}<i class="pi pi-clock"></i></p>
                </button>
            </div>
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
    @apply h-32 object-contain;
}

.level {
    @apply p-2 max-h-full font-bold aspect-square flex items-center justify-center bg-custom_blue text-white rounded-2xl;
}

</style>