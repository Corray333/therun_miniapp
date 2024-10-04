<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import race from '@/components/icons/race-icon.vue'
import bcoin from '@/components/icons/bcoin-icon.vue'
import keys from '@/components/icons/keys-icon.vue'
import KeyIcon from '@/components/icons/key-icon.vue'
import { useI18n } from 'vue-i18n'
import { useAccountStore } from '@/stores/account'
import { useComponentsStore } from '@/stores/components'

const showModal = ref<boolean>(false)
const pick = ref<string>("race")

const { t } = useI18n()
const accStore = useAccountStore()
const componentsStore = useComponentsStore()
const bonusesLabel = ref<HTMLBaseElement>()

onMounted(() => {
  if (bonusesLabel.value) {
    componentsStore.bonusesLabelPos[0] = bonusesLabel.value.offsetTop
    componentsStore.bonusesLabelPos[1] = bonusesLabel.value.offsetLeft
  }
})

</script>

<template>
  <section class=" flex w-full gap-2 z-30">
    <Transition name="delay">
      <section v-show="showModal" @click.self="showModal = false"
        class=" wrapper fixed z-50 w-full h-screen top-0 left-0 flex items-end">
        <Transition name="slide-down">
          <section v-if="showModal"
            class=" modal w-full rounded-t-2xl bg-white p-4 py-8 flex flex-col justify-center items-center shadow-lg">
            <div v-if="pick == 'race'" class=" w-full text-center flex flex-col items-center gap-2">
              <h2 class=" header flex items-center gap-2">
                <race />Race
              </h2>
              <p class="mb-14">{{ t('wallet.raceDescription') }}</p>
              <button disabled>{{ t('wallet.withdraw') }}</button>
            </div>
            <div v-else-if="pick == 'bonuses'" class=" w-full text-center flex flex-col items-center gap-2">
              <h2 class=" header flex items-center gap-2">
                <bcoin />{{ t('wallet.bonuses') }}
              </h2>
              <p class="mb-14">{{ t('wallet.bonusesDescription') }}</p>
              <button disabled>{{ t('wallet.withdraw') }}</button>
            </div>
            <div v-else class=" w-full text-center flex flex-col items-center gap-4">
              <div class="keys-container">
                <div class="balance-card">
                  <div>
                    <KeyIcon color="var(--red_key)" />
                    <p>{{ accStore.user.red_keyBalance }}</p>
                  </div>
                  <p>{{ t('wallet.redKeys') }}</p>
                </div>
                <div class="balance-card">
                  <div>
                    <KeyIcon color="var(--blue_key)" />
                    <p>{{ accStore.user.blue_keyBalance }}</p>
                  </div>
                  <p>{{ t('wallet.blueKeys') }}</p>
                </div>
                <div class="balance-card">
                  <div>
                    <KeyIcon color="var(--green_key)" />
                    <p>{{ accStore.user.green_keyBalance }}</p>
                  </div>
                  <p>{{ t('wallet.greenKeys') }}</p>
                </div>
              </div>
              <p class="mb-14">{{ t('wallet.keysDescription') }}</p>
              <button>{{ t('wallet.ok') }}</button>
            </div>
          </section>
        </Transition>
      </section>
    </Transition>

    <div @click="pick = 'race'; showModal = true"
      class=" bg-half_dark p-2 px-4 w-full rounded-2xl flex flex-col justify-center items-center">
      <span class=" flex gap-1">
        <race color="var(--primary)" />
        <p class=" text-left font-bold flex items-center">{{ accStore.user.raceBalance }}</p>
      </span>
      <p class=" label flex items-center gap-1">{{ t('wallet.balance') }}<i class="pi pi-chevron-right text-dark"
          style="font-size: 0.75rem;"></i></p>
    </div>

    <div @click="pick = 'bonuses'; showModal = true"
      class=" bg-half_dark p-2 px-4 w-full rounded-2xl flex flex-col justify-center items-center">
      <span class=" flex gap-1" ref="bonusesLabel" :class="componentsStore.animateBonuses ? 'animate-bonuses' : ''">
        <bcoin color="var(--primary)" />
        <p class=" text-left font-bold flex items-center">{{ accStore.user.pointBalance }}</p>
      </span>
      <p class=" label flex items-center gap-1">{{ t('wallet.balance') }}<i class="pi pi-chevron-right text-dark"
          style="font-size: 0.75rem;"></i></p>
    </div>

    <div @click="pick = 'key'; showModal = true"
      class=" bg-half_dark p-2 px-4 rounded-2xl flex flex-col justify-center items-center">
      <span class=" flex gap-1">
        <keys color="var(--primary)" />
      </span>
      <p class=" label flex items-center gap-1">{{ t('wallet.keys') }}<i class="pi pi-chevron-right text-dark"
          style="font-size: 0.75rem;"></i></p>
    </div>
  </section>
</template>


<style scoped>

.keys-container{
    @apply flex gap-2 w-full;
}
.balance-card {
    @apply w-full p-2 rounded-2xl bg-half_dark flex flex-col gap-1 text-center;
}

.balance-card>div {
    @apply flex gap-1 items-center justify-center font-bold;
}
.balance-card>p{
    @apply text-sm text-dark;
}

.modal {
  box-shadow: 0 -0.25rem 1rem 0 rgba(0, 0, 0, 0.1);
}

.slide-down-enter-active,
.slide-down-leave-active {
  transition: transform 0.5s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
  transform: translateY(100%);
}

.delay-enter-active,
.delay-leave-active {
  transition: opacity 0.5s ease;
}

.delay-enter-from,
.delay-leave-to {
  opacity: 1;
}

.animate-bonuses {
  animation: bonuses-animation 0.5s 4;
}

@keyframes bonuses-animation {
  0% {
    transform: scale(1)
  }

  50% {
    transform: scale(1.4)
  }

  0% {
    transform: scale(1)
  }
}
</style>