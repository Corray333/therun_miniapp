<script lang="ts" setup>
import race from '@/components/icons/race-icon.vue'
import bcoin from '@/components/icons/bcoin-icon.vue'
import key from '@/components/icons/key-icon.vue'
import {Task} from '@/types/types'

defineProps({
    task: Task
})

</script>

<template>
    <div class=" p-4 flex items-center justify-between gap-4 rounded-2xl bg-half_dark w-full">
        <img :src="task?.icon" class=" w-16 h-16 rounded-full" alt="icon">
        <div class="flex flex-col gap-2">
            <p class=" line-clamp-2">{{ task?.description }}</p>
            <span class="flex gap-4">
                <p v-if="task?.raceReward != 0" class=" flex gap-2"><race/>{{ task?.raceReward }}</p>
                <p v-if="task?.keysReward != 0" class=" flex gap-2"><key color="var(--primary)"/>{{ task?.keysReward }}</p>
                <p v-if="task?.pointsReward != 0" class=" flex gap-2"><bcoin/>{{ task?.pointsReward }}</p>
            </span>
        </div>
        <i v-if="!task?.claimed" class="pi pi-chevron-right text-dark" style="font-size:1.25rem"></i>
        <i v-else  class="pi pi-check-circle text-dark text-green-600" style="font-size:1.25rem"></i>
    </div>
</template>


<style></style>