<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import race from '@/components/icons/race-icon.vue'
import bcoin from '@/components/icons/bcoin-icon.vue'
import key from '@/components/icons/key-icon.vue'
import { useI18n } from 'vue-i18n'
import { useAccountStore } from '@/stores/account'
import { useComponentsStore } from '@/stores/components'

const showModal = ref<boolean>(false)
const pick = ref<string>("race")

const { t } = useI18n()
const store = useAccountStore()
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
  <section class=" flex w-full gap-2 z-10">
    <Transition name="delay">
      <section v-show="showModal" @click.self="showModal = false"
        class=" wrapper fixed z-50 w-full h-screen top-0 left-0 flex items-end">
        <Transition name="slide">
          <section v-if="showModal"
            class=" modal w-full rounded-t-2xl bg-white p-4 py-8 flex flex-col justify-center items-center shadow-lg">
            <div v-if="pick == 'race'" class=" w-full text-center flex flex-col items-center gap-2">
              <h2 class=" header flex items-center gap-2">
                <race />Race
              </h2>
              <p>{{ t('wallet.raceDescription') }}</p>
              <button disabled>{{ t('wallet.withdraw') }}</button>
            </div>
            <div v-else-if="pick == 'bonuses'" class=" w-full text-center flex flex-col items-center gap-2">
              <h2 class=" header flex items-center gap-2">
                <bcoin />{{ t('wallet.bonuses') }}
              </h2>
              <p>{{ t('wallet.bonusesDescription') }}</p>
              <button disabled>{{ t('wallet.withdraw') }}</button>
            </div>
            <div v-else class=" w-full text-center flex flex-col items-center gap-2">
              <h2 class=" header flex items-center gap-2"><key color="var(--primary)"/>{{ t('wallet.keys') }}</h2>
              <p>{{ t('wallet.keysDescription') }}</p>
              <router-link to="/chibi">
                <button class=" w-full">{{ t('wallet.use') }}</button>
              </router-link>
            </div>
          </section>
        </Transition>
      </section>
    </Transition>

    <div @click="pick = 'race'; showModal = true"
      class=" bg-half_dark p-4 w-full rounded-2xl flex flex-col justify-center items-center">
      <span class=" flex gap-1">
        <race color="var(--primary)" />
        <p class=" text-left font-bold">{{ store.user.raceBalance }}</p>
      </span>
      <p class=" label flex items-center gap-1">{{ t('wallet.balance') }}<i class="pi pi-chevron-right text-dark"
          style="font-size: 0.75rem;"></i></p>
    </div>

    <div @click="pick = 'bonuses'; showModal = true"
      class=" bg-half_dark p-4 w-full rounded-2xl flex flex-col justify-center items-center">
      <span class=" flex gap-1" ref="bonusesLabel">
        <bcoin color="var(--primary)" />
        <p class=" text-left font-bold">{{ store.user.pointBalance }}</p>
      </span>
      <p class=" label flex items-center gap-1">{{ t('wallet.balance') }}<i class="pi pi-chevron-right text-dark"
          style="font-size: 0.75rem;"></i></p>
    </div>

    <div @click="pick = 'key'; showModal = true"
      class=" bg-half_dark p-4 w-full rounded-2xl flex flex-col justify-center items-center">
      <span class=" flex gap-1" >
        <key color="var(--primary)" />
        <p class=" text-left font-bold">{{ store.user.keyBalance }}</p>
      </span>
      <p class=" label flex items-center gap-1">{{ t('wallet.keys') }}<i class="pi pi-chevron-right text-dark"
          style="font-size: 0.75rem;"></i></p>
    </div>
  </section>
</template>


<style>
.modal {
  box-shadow: 0 -0.25rem 1rem 0 rgba(0, 0, 0, 0.1);
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.5s ease;
}

.slide-enter-from,
.slide-leave-to {
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