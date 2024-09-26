<script lang="ts" setup>

import { Battle } from '@/types/types'
import bcoin from '@/components/icons/bcoin-icon.vue'
import key from '@/components/icons/key-icon.vue'
import miles from '@/components/icons/miles-icon.vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps({
    battle: {
        type: Battle,
        required: true
    }
})

</script>

<template>
    <div class=" bg-half_dark rounded-2xl">
        <div class="flex w-full">
            <div class="flex flex-col items-center justify-center w-full border-half_gray border-b-1 border-r-1 p-4">
                <p class=" font-bold flex gap-2">
                    <bcoin />2000
                </p>
                <p class="label">{{ t('screens.battles.battle.participationFee') }}</p>
            </div>
            <div class="flex flex-col items-center justify-center w-full border-half_gray border-b-1 p-4">
                <p class=" font-bold flex gap-2">
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
                            <miles />{{ battle.userResult }}
                        </span>
                        <img :src="battle.user.photo" :class="battle.pick == 1 ? 'user-picked':''" class=" w-full rounded-2xl" alt="">
                        <p>@{{ battle.user.username }}</p>
                    </div>
                    <button v-if="battle.pick == 0" class=" py-1">{{ t('screens.battles.battle.choose') }}</button>
                    <button v-if="battle.pick == 1" class=" py-1">{{ t('screens.battles.battle.chosen') }}</button>
                </div>
                <img src="../assets/images/battles/vs-sign.png" class=" w-16 object-contain" alt="">
                <div class="player w-full">
                    <div class="user-info text-center py-2 px-4 gap-2 flex flex-col">
                        <span class=" w-full flex items-center justify-center gap-2 font-bold text-xl">
                            <miles />{{ battle.opponentResult }}
                        </span>
                        <img :src="battle.opponent.photo" :class="battle.pick == 2 ? 'user-picked':''" class=" w-full rounded-2xl" alt="">
                        <p>@{{ battle.opponent.username }}</p>
                    </div>
                    <button v-if="battle.pick == 0" class=" py-1">{{ t('screens.battles.battle.choose') }}</button>
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