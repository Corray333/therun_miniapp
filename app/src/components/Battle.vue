<script lang="ts" setup>

import { Battle } from '@/types/types'
import bcoin from '@/components/icons/bcoin-icon.vue'
import key from '@/components/icons/key-icon.vue'
import miles from '@/components/icons/miles-icon.vue'
import { useI18n } from 'vue-i18n'
import { useAccountStore } from '@/stores/account'
import axios, { isAxiosError } from 'axios'
import { useComponentsStore } from '@/stores/components'
import { auth, getUser } from '@/utils/helpers'
import { ref } from 'vue'

const { t } = useI18n()
const accStore = useAccountStore()
const componentsStore = useComponentsStore()

const props = defineProps({
    battle: {
        type: Battle,
        required: true
    },
    showMiles: {
        type: Boolean,
        default: false
    }
})

const loading = ref<boolean>(false)
const loading2 = ref<boolean>(false)

const makeBet = async (pick: number) => {
    if (accStore.user.pointBalance<300){
        componentsStore.addError(t('screens.battles.battle.errNotEnoughPoints'))
        return
    }

    if (pick === 1) {
        loading.value = true
    } else {
        loading2.value = true
    }

    try {
        await axios.post(`${import.meta.env.VITE_API_URL}/battles/${props.battle.id}/bet`, {
            pick: pick
        }, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        props.battle.pick = pick
        getUser(accStore.user.id)
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
        } else if (isAxiosError(error)) {
            componentsStore.addError(error.message)
        }
    } finally {
        loading.value = false
        prePick.value = 0
    }
}


const prePick = ref<number>(0)


</script>

<template>
    <Transition name="delay">
        <section v-show="prePick != 0" @click.self="prePick = 0"
            class=" wrapper fixed z-50 w-full h-screen top-0 left-0 flex items-center justify-center p-4 backdrop-blur-md">
            <Transition name="slide-down">
                <section v-if="prePick != 0"
                    class=" modal gap-4 w-full rounded-2xl bg-white p-4 flex flex-col justify-center items-center shadow-lg">
                    <p class="font-bold">{{ t('screens.battles.battle.pickApprove') }} @{{ prePick ==1 ? battle.user.username : battle.opponent.username }}</p>
                    <div class="flex gap-2 w-full">
                        <button @click="prePick = 0"  class=" py-2 text-primary bg-white">{{ t('screens.battles.battle.pickApproveCancel') }}</button>
                        <button @click="makeBet(prePick)"  class=" py-2">
                            <p v-if="loading"><i class="pi pi-spinner pi-spin"></i></p>
                            <p v-else>
                                {{ t('screens.battles.battle.pickApproveOk') }}
                            </p>
                        </button>
                    </div>
                </section>
            </Transition>
        </section>
    </Transition>

    <div class=" bg-half_dark rounded-2xl">
        <div class="flex w-full">
            <div class="flex flex-col items-center justify-center w-full border-half_gray border-b-1 border-r-1 p-4">
                <p class=" font-bold flex gap-2 items-center">
                    <!-- TODO: get from server -->
                    <bcoin />300
                </p>
                <p class="label">{{ t('screens.battles.battle.participationFee') }}</p>
            </div>
            <div class="flex flex-col items-center justify-center w-full border-half_gray border-b-1 p-4">
                <p class=" font-bold flex gap-2 items-center">
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
                        <span v-show="showMiles" class=" w-full flex items-center justify-center gap-2 font-bold text-xl">
                            <miles/>{{ Math.floor(battle.userResult) }}
                        </span>
                        <span class="w-full relative rounded-2xl bg-dark text-white text-4xl font-bold aspect-square flex justify-center items-center">
                            <img v-if="battle.user.photo" :src="battle.user.photo" class=" w-full h-full absolute rounded-2xl" alt="">
                            {{ battle.user.username[0] }}
                        </span>
                        <p>@{{ battle.user.username }}</p>
                    </div>
                    <button v-show="!showMiles" @click="prePick = 1" v-if="battle.pick == 0" class=" py-1">
                        <p>{{ t('screens.battles.battle.choose') }}</p>
                    </button>
                    <button v-if="battle.pick == 1" class=" py-1">{{ t('screens.battles.battle.chosen') }}</button>
                </div>
                <img src="../assets/images/battles/vs-sign.png" class=" w-16 object-contain" alt="">
                <div class="player w-full">
                    <div class="user-info text-center py-2 px-4 gap-2 flex flex-col">
                        <span v-show="showMiles" class=" w-full flex items-center justify-center gap-2 font-bold text-xl">
                            <miles/>{{ Math.floor(battle.opponentResult) }}
                        </span>
                        <span  class="w-full relative rounded-2xl bg-dark text-white text-4xl font-bold aspect-square flex justify-center items-center">
                            <img v-if="battle.opponent.photo" :src="battle.opponent.photo" class=" w-full h-full absolute rounded-2xl" alt="">
                            {{ battle.opponent.username[0] }}
                        </span>
                        <p>@{{ battle.opponent.username }}</p>
                    </div>
                    <button v-show="!showMiles" @click="prePick = 2" v-if="battle.pick == 0" class=" py-1">
                        <p>{{ t('screens.battles.battle.choose') }}</p>
                    </button>
                    <button v-if="battle.pick == 2" class=" py-1" disabled>{{ t('screens.battles.battle.chosen')
                        }}</button>
                </div>
            </div>
        </div>
    </div>
</template>


<style scoped>

.delay-enter-active,
.delay-leave-active {
    transition: opacity 0.5s ease;
}

.delay-enter-from,
.delay-leave-to {
    opacity: 0;
}

.slide-down-enter-active,
.slide-down-leave-active {
    transition: transform 0.5s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
    transform: scale(0);
}

.user-picked {
    border: 4px solid var(--primary);
}
</style>