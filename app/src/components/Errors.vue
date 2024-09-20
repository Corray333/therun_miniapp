<script lang="ts" setup>

import {useComponentsStore} from '@/stores/components'

const componentsStore = useComponentsStore()

</script>

<template>
    <div class=" flex flex-col w-full p-4 z-60">
        <TransitionGroup name="slide-down">
            <p v-for="(error, i) of componentsStore.errors" :key="i" class=" w-full text-center p-4 bg-white border-primary border-2 rounded-2xl shadow-lg">{{ error }}</p>
        </TransitionGroup>
    </div>
</template>


<style scoped>

.slide-down-enter-active,
.slide-down-leave-active {
    transition: transform 0.25s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
    transform: translateY(-100%);
}

</style>