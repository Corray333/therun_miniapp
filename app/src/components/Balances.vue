<script lang="ts" setup>
import race from '@/components/icons/race-icon.vue'
import bcoin from '@/components/icons/bcoin-icon.vue'
import key from '@/components/icons/key-icon.vue'
import { useBalanceStore } from '@/stores/counter'

const store = useBalanceStore()

</script>

<template>
  <section class=" flex w-full gap-2">
    <div class=" bg-half_dark p-4 rounded-xl w-full flex flex-col justify-center items-center">
      <span class=" flex gap-1">
        <race/>
        <p class=" text-left font-bold">{{ store.race }}</p>
      </span>
      <p class=" label">{{ $t('wallet.balance') }}</p>
    </div>
    <div class=" bg-half_dark p-4 rounded-xl w-full flex flex-col justify-center items-center">
      <span class=" flex gap-1">
        <bcoin/>
        <p class=" text-left font-bold">{{ store.bonuses }}</p>
      </span>
      <p class=" label">{{ $t('wallet.balance') }}</p>
    </div>
    <div class=" bg-half_dark p-4 rounded-xl flex flex-col justify-center items-center">
      <span class=" flex gap-1">
        <key color="black" />
        <p class=" text-left font-bold">{{ store.keys }}</p>
      </span>
      <p class=" label">{{ $t('wallet.keys') }}</p>
    </div>
  </section>
</template>


<style></style>