<script lang="ts" setup>

import { useI18n } from 'vue-i18n'
import { getTuningModules } from '@/services/Tuning'
import { Module } from '@/types/types'
import { ref, onMounted } from 'vue'
import { Characteristic } from '@/services/Characteristics'
import { useRoute } from 'vue-router'

import TuningCard from '../characteristics/TuningCard.vue'

const { t } = useI18n()


const constModules = ref<Module[]>([])
const tempModules = ref<Module[]>([])

const route = useRoute()

const characteristic = route.params.characteristic as string

onMounted(async () => {
    const tuningModules = await getTuningModules(characteristic)
    constModules.value = tuningModules.filter(module => module.isTemp === false)
    tempModules.value = tuningModules.filter(module => module.isTemp === true)
})

const showConstModules = ref<boolean>(true)


</script>

<template>
    <section class="tuning-inventory">
        <h1 class="text-center">{{ t(`screens.cars.characteristics.${characteristic}.name`) }}</h1>
        <div class="buttons" >
            <button @click="showConstModules = true" :class="showConstModules ? 'choosed':''">{{ t('screens.cars.tuning.const') }}</button>
            <button @click="showConstModules = false" :class="!showConstModules ? 'choosed':''">{{ t('screens.cars.tuning.temp') }}</button>
        </div>

        <div v-if="showConstModules" >
            <div v-if="constModules.length > 0" class="tuning const">
                <TuningCard v-for="(module, i) in constModules" :module="module" :key="i" />
            </div>
            <p v-else class="text-center p-4 font-bold text-dark">{{ t('screens.cars.tuning.noTuning') }}</p>
        </div>
        <div v-else>
            <div v-if="tempModules.length > 0" class="tuning temp">
                <TuningCard v-for="(module, i) in tempModules" :module="module" :key="i" />
            </div>
            <p v-else class="text-center p-4 font-bold text-dark">{{ t('screens.cars.tuning.noTuning') }}</p>
        </div>

    </section>
</template>

<style scoped>

.tuning-inventory{
    @apply flex flex-col gap-2 p-4
}

.buttons{
    @apply flex 
}

.choosed{
    @apply border-b-2 border-primary
}

button{
    @apply bg-transparent text-black font-normal rounded-none p-2
}

.tuning{
    @apply grid grid-cols-2 gap-4
}

</style>