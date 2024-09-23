<script lang="ts" setup>

import { useComponentsStore } from '@/stores/components'

const componentsStore = useComponentsStore()

</script>

<template>
    <Transition name="delay">
        <div v-show="componentsStore.errors.length > 0" class=" flex flex-col w-full p-4 z-60">
            <div id="errors-container" class="relative">
                <TransitionGroup name="slide-down">
                    <p v-for="(error, i) of componentsStore.errors" :key="i"
                        class=" w-full text-center p-4 bg-white border-primary border-2 rounded-2xl shadow-lg absolute">
                        Error: {{ error }}
                    </p>
                </TransitionGroup>
            </div>
        </div>
    </Transition>
</template>


<style scoped>
.delay-enter-active,
.delay-leave-active {
    transition: opacity 0.25s ease;
}

.delay-enter-from,
.delay-leave-to {
    opacity: 0;
}

#errors-container>p {
    transition: all 0.25s;
    transform: translateY(30px) scale(0.9);
}

#errors-container>:nth-last-child(2) {
    transform: translateY(15px) scale(0.95);
}

#errors-container>:nth-last-child(1) {
    transform: translateY(0px) scale(1);
}

.slide-down-enter-active,
.slide-down-leave-active {
    transition: transform 0.25s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
    transform: translateY(-100%);
}
</style>