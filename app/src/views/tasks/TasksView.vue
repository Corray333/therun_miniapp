<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import { ref, onBeforeMount } from 'vue'
import { Task } from '@/types/types'
import { useAccountStore } from '@/stores/account'
import { auth } from '@/utils/helpers'
import axios, { isAxiosError } from 'axios'
import { useComponentsStore } from '@/stores/components'

import TaskCard from '@/components/Task.vue'
import race from '@/components/icons/race-icon.vue'
import bcoin from '@/components/icons/bcoin-icon.vue'
import key from '@/components/icons/key-icon.vue'

const { t } = useI18n()
const accStore = useAccountStore()
const componentsStore = useComponentsStore()

const tasks = ref<Task[]>([
    {
        "id": 1,
        "description": "Subscribe to Jun so Soon channel",
        "type": "tg",
        "link": "https://t.me/jun_so_soon",
        "expireAt": 1727596313,
        "pointsReward": 200,
        "keysReward": 0,
        "raceReward": 0,
        "data": {},
        "icon": "https://store-images.s-microsoft.com/image/apps.55245.13537716651231321.3067a421-6c2f-48a9-b77c-1e38e19146e6.10e2aa49-52ca-4e79-9a61-b6422978afb9?h=210"
    }
])

const getTasks = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/tasks`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        tasks.value = data
        return true
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await getTasks()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        }
    }
}

onBeforeMount(async () => {
    await getTasks()
})

const pickedTask = ref<Task | null>(null)

</script>

<template>
    <section class=" flex flex-col gap-4 py-4">
        <Transition name="delay">
            <section v-show="pickedTask" @click.self="pickedTask = null"
                class=" wrapper fixed z-50 w-full h-screen top-0 left-0 flex items-end">
                <Transition name="slide-down">
                    <section v-if="pickedTask"
                        class=" modal w-full rounded-t-2xl bg-white p-4 py-8 flex flex-col justify-center items-center shadow-lg gap-4">
                        <img :src="pickedTask?.icon" class=" w-16 h-16 rounded-full" alt="icon">
                        <p class=" font-bold text-center">{{ pickedTask?.description }}</p>
                        <span class="flex gap-4">
                            <p v-if="pickedTask?.raceReward != 0" class=" flex gap-2 font-bold items-center"><i class="pi pi-plus"></i><race/>{{ pickedTask?.raceReward }}</p>
                            <p v-if="pickedTask?.keysReward != 0" class=" flex gap-2 font-bold items-center"><i class="pi pi-plus"></i><key color="var(--primary)"/>{{ pickedTask?.keysReward }}</p>
                            <p v-if="pickedTask?.pointsReward != 0" class=" flex gap-2 font-bold items-center"><i class="pi pi-plus"></i><bcoin/>{{ pickedTask?.pointsReward }}</p>
                        </span>
                        <a target="_blank":href="pickedTask?.link" class="w-full"><button>{{ t('screens.tasks.startBtn') }}</button></a>
                        <button class=" btn-type-2">{{ t('screens.tasks.checkBtn') }}</button>
                    </section>
                </Transition>
            </section>
        </Transition>
        <div class="daily flex flex-col gap-4">
            <div class="daily-check daily p-4 rounded-2xl flex border-2 border-primary   gap-4 items-center">
                <img src="../../assets/images/tasks/gift-icon.png" class=" w-16 h-16" alt="">
                <div class="flex flex-col">
                    <h2>{{ t('screens.tasks.tasks.dailyCheck.header') }}</h2>
                    <p class="label">{{ t('screens.tasks.tasks.dailyCheck.description') }}</p>
                </div>
                <i class="pi pi-chevron-right text-dark" style="font-size:1.25rem"></i>
            </div>
            <div class="daily-combo daily p-4 rounded-2xl flex border-2 border-primary   gap-4 items-center">
                <img src="../../assets/images/tasks/calendar-icon.png" class=" w-16 h-16" alt="">
                <div class="flex flex-col">
                    <h2>{{ t('screens.tasks.tasks.dailyCombo.header') }}</h2>
                    <p class="label">{{ t('screens.tasks.tasks.dailyCombo.description') }}</p>
                </div>
                <i class="pi pi-chevron-right text-dark" style="font-size:1.25rem"></i>
            </div>
        </div>
        <h1 class=" mt-4">{{ t('screens.tasks.tasks.header') }}</h1>
        <TaskCard v-for="task in tasks" :task="task" :key="task.id" @click="pickedTask = task" />
    </section>
</template>


<style scoped>

.delay-enter-active,
.delay-leave-active {
    transition: opacity 0.5s ease;
}

.delay-enter-from,
.delay-leave-to {
    opacity: 1;
}

.slide-down-enter-active,
.slide-down-leave-active {
    transition: transform 0.5s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
    transform: translateY(100%);
}

.modal{
    box-shadow: 0 -0.25rem 1rem 0 rgba(0, 0, 0, 0.1);
}

</style>