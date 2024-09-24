<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import { ref, onBeforeMount } from 'vue'
import { Task } from '@/types/types'
import { useAccountStore } from '@/stores/account'
import { auth } from '@/utils/helpers'
import axios, { isAxiosError } from 'axios'
import { useComponentsStore } from '@/stores/components'
import SlideUpDown from 'vue-slide-up-down'

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
        "icon": "https://store-images.s-microsoft.com/image/apps.55245.13537716651231321.3067a421-6c2f-48a9-b77c-1e38e19146e6.10e2aa49-52ca-4e79-9a61-b6422978afb9?h=210",
        "done": false,
        "claimed": false,
        "clicked": false,
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
        for (let i = 0; i < tasks.value.length; i++) {
            tasks.value[i].claimed = false
            tasks.value[i].done = false
        }
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
        } else {
            if (isAxiosError(error)) {
                componentsStore.addError(error.message)
            }
        }
    }
}

const checkTask = async () => {
    if (!pickedTask.value) return

    pickedTaskLoading.value = true

    try {
        const { data } = await axios.post(`${import.meta.env.VITE_API_URL}/tasks/${pickedTask.value?.id}/check`, {}, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        pickedTask.value.done = data.done
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
        } else {
            if (isAxiosError(error)) {
                componentsStore.addError(error.message)
            }
        }
    } finally {
        pickedTaskLoading.value = false
    }
}

const claimTask = async () => {
    if (!pickedTask.value) return

    pickedTaskLoading.value = true

    try {
        const { data } = await axios.post(`${import.meta.env.VITE_API_URL}/tasks/${pickedTask.value?.id}/claim`, {}, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        accStore.user.pointBalance = data.pointBalance
        accStore.user.keyBalance = data.keyBalance
        accStore.user.raceBalance = data.raceBalance
        pickedTask.value.done = true
        pickedTask.value.claimed = true
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
        } else {
            if (isAxiosError(error)) {
                componentsStore.addError(error.message)
            }
        }
    } finally {
        setTimeout(() => {
            pickedTask.value = undefined
        }, 2000)
        pickedTaskLoading.value = false
    }
}

onBeforeMount(async () => {
    await getTasks()
})

const pickedTask = ref<Task>()
const pickedTaskLoading = ref<boolean>(false)

</script>

<template>
    <section class=" flex flex-col gap-4 py-4">
        <Transition name="delay">
            <section v-show="pickedTask" @click.self="pickedTask = undefined"
                class=" wrapper fixed z-50 w-full h-screen top-0 left-0 flex items-end">
                <Transition name="slide-down">
                    <section v-if="pickedTask"
                        class=" modal w-full rounded-t-2xl bg-white p-4 py-8 flex flex-col justify-center items-center shadow-lg gap-4">
                        <img :src="pickedTask?.icon" class=" w-16 h-16 rounded-full" alt="icon">
                        <p class=" font-bold text-center">{{ pickedTask?.description }}</p>
                        <span class="flex gap-4 text-2xl">
                            <p v-if="pickedTask?.raceReward != 0" class=" flex gap-2 font-bold items-center">
                                <race />+{{ pickedTask?.raceReward }}
                            </p>
                            <p v-if="pickedTask?.keysReward != 0" class=" flex gap-2 font-bold items-center">
                                <key color="var(--primary)" />+{{ pickedTask?.keysReward }}
                            </p>
                            <p v-if="pickedTask?.pointsReward != 0" class=" flex gap-2 font-bold items-center">
                                <bcoin />+{{ pickedTask?.pointsReward }}
                            </p>
                        </span>

                        <SlideUpDown :active="!pickedTask.done" class="w-full">
                            <div class="flex flex-col gap-4 w-full">
                                <a @click="pickedTask.clicked = true" target="_blank" :href="pickedTask?.link" class="w-full"><button>{{ t('screens.tasks.startBtn') }}</button></a>
                                <button :disabled="pickedTask.clicked" @click="checkTask" class=" btn-type-2">
                                    <p v-if="!pickedTaskLoading">{{ t('screens.tasks.checkBtn') }}</p>
                                    <i v-else class="pi pi-spin pi-spinner"></i>
                                </button>
                            </div>
                        </SlideUpDown>
                        <SlideUpDown :active="pickedTask.done && !pickedTask?.claimed" class="w-full">
                            <button @click="claimTask" class=" btn-type-2">
                                <p v-if="!pickedTaskLoading">{{ t('screens.tasks.claimBtn') }}</p>
                                <i v-else class="pi pi-spin pi-spinner"></i>
                            </button>
                        </SlideUpDown>
                        <SlideUpDown :active="pickedTask?.claimed" class="w-full">
                            <button class=" btn-type-3 gap-1">{{ t('screens.tasks.done') }}<i
                                    class="pi pi-check"></i></button>
                        </SlideUpDown>
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
        <p v-if="tasks.length == 0">{{ t('screens.tasks.noNewTasks') }}</p>
        <TaskCard v-for="task in tasks" :task="task" :key="task.id" @click="if (!task.claimed) pickedTask = task;" />
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

.modal {
    box-shadow: 0 -0.25rem 1rem 0 rgba(0, 0, 0, 0.1);
}
</style>