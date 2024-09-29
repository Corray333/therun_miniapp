<script lang="ts" setup>

import { Battle } from '@/types/types'
import bcoin from '@/components/icons/bcoin-icon.vue'
import key from '@/components/icons/key-icon.vue'
import miles from '@/components/icons/miles-icon.vue'
import { useI18n } from 'vue-i18n'
import { useAccountStore } from '@/stores/account'
import axios, {isAxiosError} from 'axios'
import { useComponentsStore } from '@/stores/components'
import { auth } from '@/utils/helpers'
import {ref} from 'vue'

const { t } = useI18n()
const accStore = useAccountStore()
const componentsStore = useComponentsStore()

const props = defineProps({
    battle: {
        type: Battle,
        required: true
    }
})

const loading1 = ref<boolean>(false)
const loading2 = ref<boolean>(false)

const makeBet = async (pick: number) => {
    if (pick === 1) {
        loading1.value = true
    } else {
        loading2.value = true
    }

    try {
        await axios.post(`${import.meta.env.VITE_API_URL}/battles/${props.battle.id}/bet`,{
            pick: pick
        }, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        props.battle.pick = pick
        return true
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await makeBet(pick)
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        } else if (isAxiosError(error) && error.response?.status === 400) {
            componentsStore.addError(t('screens.battles.battle.errBetTooLate'))
        } else if(isAxiosError(error)) {
            componentsStore.addError(error.message)
        }
    } finally{
        loading1.value = false
        loading2.value = false
    }
}

</script>

<template>
    <div class=" bg-half_dark rounded-2xl">
        <div class="flex w-full">
            <div class="flex flex-col items-center justify-center w-full border-half_gray border-b-1 border-r-1 p-4">
                <p class=" font-bold flex gap-2">
                    <!-- TODO: get from server -->
                    <bcoin />300
                </p>
                <p class="label">{{ t('screens.battles.battle.participationFee') }}</p>
            </div>
            <div class="flex flex-col items-center justify-center w-full border-half_gray border-b-1 p-4">
                <p class=" font-bold flex gap-2">
                    <!-- TODO: get from server -->
                    <key color="var(--primary)" />1
                </p>
                <p class="label">{{ t('screens.battles.battle.prize') }}</p>
            </div>
        </div>

        <div class="p-4">
            <div class="players flex gap-2">
                <div class="player w-full">
                    <div class="user-info text-center py-2 px-4 gap-2 flex flex-col">
                        <span class=" w-full flex items-center justify-center gap-2 font-bold text-xl">
                            <miles />{{ Math.floor(battle.userResult) }}
                        </span>
                        <img v-if="battle.user.photo" :src="battle.user.photo" :class="battle.pick == 1 ? 'user-picked':''" class=" w-full rounded-2xl" alt="">
                        <span v-else class="w-full rounded-2xl bg-dark text-white text-4xl font-bold aspect-square flex justify-center items-center">{{ battle.user.username[0] }}</span>
                        <p>@{{ battle.user.username }}</p>
                    </div>
                    <button @click="makeBet(1)" v-if="battle.pick == 0" class=" py-1">
                        <p v-if="loading1"><i class="pi pi-spinner pi-spin"></i></p>
                        <p v-else>{{ t('screens.battles.battle.choose') }}</p>
                    </button>
                    <button v-if="battle.pick == 1" class=" py-1">{{ t('screens.battles.battle.chosen') }}</button>
                </div>
                <img src="../assets/images/battles/vs-sign.png" class=" w-16 object-contain" alt="">
                <div class="player w-full">
                    <div class="user-info text-center py-2 px-4 gap-2 flex flex-col">
                        <span class=" w-full flex items-center justify-center gap-2 font-bold text-xl">
                            <miles />{{ Math.floor(battle.opponentResult) }}
                        </span>
                        <img v-if="battle.opponent.photo" :src="battle.opponent.photo" :class="battle.pick == 2 ? 'user-picked':''" class=" w-full rounded-2xl" alt="">
                        <span v-else class="w-full rounded-2xl bg-dark text-white text-4xl font-bold aspect-square flex justify-center items-center">{{ battle.opponent.username[0] }}</span>
                        <p>@{{ battle.opponent.username }}</p>
                    </div>
                    <button @click="makeBet(2)" v-if="battle.pick == 0" class=" py-1">
                        <p v-if="loading2"><i class="pi pi-spinner pi-spin"></i></p>
                        <p v-else>{{ t('screens.battles.battle.choose') }}</p>
                    </button>
                    <button v-if="battle.pick == 2" class=" py-1" disabled>{{ t('screens.battles.battle.chosen') }}</button>
                </div>
            </div>
        </div>
    </div>
</template>


<style scoped>

.user-picked{
    border: 4px solid var(--primary);
}

</style>