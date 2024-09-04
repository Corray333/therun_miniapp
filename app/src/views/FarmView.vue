<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import Balances from '@/components/Balances.vue';
import clock from '@/components/icons/clock-icon.vue'
import bcoin from '@/components/icons/bcoin-icon.vue'
import Navbar from '@/components/Navbar.vue'
import { Vue3Lottie } from 'vue3-lottie'
import RunningJSON from '@/assets/animations/running.json'
import { useBalanceStore } from '@/stores/counter'
import { storeToRefs } from 'pinia';

const balanceStore = useBalanceStore()


class Farming {
    maxPoints: number = 200;
    endTime: number = 1725091363; // Unix timestamp in seconds
}

const farming = ref<Farming>({
    maxPoints: 200,
    endTime: 1725091363,
});

const remainingTime = ref<string>('');
const currentPoints = ref<number>(0);

const totalDuration = 7200

// Helper function to format time
const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600).toString().padStart(2, '0');
    const minutes = Math.floor((seconds % 3600) / 60).toString().padStart(2, '0');
    const secondsRemaining = (seconds % 60).toString().padStart(2, '0');
    return `${hours}:${minutes}:${secondsRemaining}`;
};

// Function to calculate remaining time and current points
const calculateRemainingTimeAndPoints = () => {
    const now = Math.floor(Date.now() / 1000); // Current time in seconds (Unix timestamp)
    const secondsLeft = farming.value.endTime - now;

    if (secondsLeft <= 0) {
        remainingTime.value = '00:00:00';
        currentPoints.value = farming.value.maxPoints; // All points earned
        clearInterval(coinsGainInterval.value)
        return;
    }

    remainingTime.value = formatTime(secondsLeft);

    const elapsedTime = totalDuration - secondsLeft;
    currentPoints.value = Math.round((elapsedTime / totalDuration) * farming.value.maxPoints * 100) / 100;
};

const coinsGainInterval = ref(0)

onMounted(() => {
    calculateRemainingTimeAndPoints(); // Initial calculation
    const interval = setInterval(calculateRemainingTimeAndPoints, 1000); // Update every second

    onUnmounted(() => {
        clearInterval(interval); // Clear the interval when the component is unmounted
    })

    coinsGainInterval.value = setInterval(createSmallCoin, 500);

})

const claim = () => {
    balanceStore.bonuses += farming.value.maxPoints
    farming.value.endTime = Math.floor(Date.now() / 1000) + totalDuration
    coinsGainInterval.value = setInterval(createSmallCoin, 500);
}

const coinsContainer = ref<HTMLDivElement>()

function createSmallCoin() {
    const smallCoin = document.createElement('div');
    smallCoin.classList.add('small-coin');

    // Генерация случайного места появления монеты
    const startX = Math.random() * 200 - 50 + 'vw';
    const startY = Math.random() * 200 - 50 + 'vh';

    smallCoin.style.setProperty('--start-x', startX);
    smallCoin.style.setProperty('--start-y', startY);

    smallCoin.style.animation = 'move-to-center 2s ease-out infinite'

    if (coinsContainer.value == null) return
    coinsContainer.value.appendChild(smallCoin);

    // Удаляем монету после завершения анимации
    setTimeout(() => {
        smallCoin.remove();
    }, 2000);
}


</script>

<template>
    <section class=" h-screen flex overflow-hidden flex-col">
        <section class=" h-full flex flex-col justify-between p-4">
            <Balances />
            <span class=" flex gap-2 justify-center items-center">
                <bcoin />
                <p class=" text-left text-4xl font-bold">{{ currentPoints }}</p>
            </span>
            <span class=" relative mx-10">
                <Vue3Lottie class=" absolute  duration-500 -top-12 left-0"
                    :style="`margin-left: calc(${currentPoints / (farming.maxPoints / 100)}% - 32px)`"
                    :animationData="RunningJSON" :height="48" :width="48" />
                <span
                    class=" relative overflow-hidden flex w-full justify-between py-1 px-4 font-bold bg-half_dark rounded-full">
                    <span class=" absolute duration-500 bg-green-500 h-full left-0 top-0 rounded-full"
                        :style="`width: ${currentPoints / (farming.maxPoints / 100)}%`"></span>
                    <p class=" relative z-10">0</p>
                    <p class=" relative z-10">{{ farming.maxPoints }}</p>
                </span>
            </span>
            <section class=" flex relative items-center justify-center" ref="coinsContainer">
                <div id="glow" class=" absolute rounded-full"></div>
                <img id="coin" src="../components/coin-tap.png" class=" relative z-10" alt="">
            </section>
            <div class=" w-full flex justify-end">
                <div class=" flex flex-col items-center w-fit bg-half_dark p-4 rounded-xl">
                    <span class="flex gap-2">
                        <clock color="var(--primary)" />
                        <span class=" font-bold">{{ remainingTime }}</span>
                    </span>
                    <p class=" text-dark text-sm">{{ $t('screens.bonuses.time') }}</p>
                </div>
            </div>
            <button @click="claim" :disabled="farming.endTime - Math.floor(Date.now() / 1000) > 1">Claim {{
                    farming.maxPoints }}</button>
        </section>
        <Navbar class="" />
    </section>
</template>

<style>
#coin {
    animation: coin-breath 5s infinite;
}

@keyframes coin-breath {
    0% {
        transform: scale(1);
    }

    50% {
        transform: scale(1.1);
    }

    100% {
        transform: scale(1);
    }
}

.small-coin {
    position: absolute;
    width: 2rem;
    height: 2rem;
    background-image: url(../components/icons/bcoin.svg);
    background-size: contain;
    border-radius: 50%;
    z-index: -10;

}

#glow {
    left: 50%;
    top: 50%;
    width: 250px;
    height: 250px;
    background-color: rgba(255, 215, 0, 0.5);
    /* Золотистый цвет с прозрачностью */
    z-index: 1;
    /* Чтобы свечение было позади монеты */
    filter: blur(20px);
    /* Размытие для эффекта свечения */
    animation: glow-animation 2s infinite;
}

@keyframes glow-animation {
    0% {
        transform: translate(-50%, -50%) scale(1);
        opacity: 1;
    }

    50% {
        transform: translate(-50%, -50%) scale(1.2);
        opacity: 0.8;
    }

    100% {
        transform: translate(-50%, -50%) scale(1);
        opacity: 1;
    }
}

@keyframes move-to-center {
    0% {
        transform: translate(var(--start-x), var(--start-y));
    }

    100% {
        transform: translate(0, 0);
    }
}
</style>
