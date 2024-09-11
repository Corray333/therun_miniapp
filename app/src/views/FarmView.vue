<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import Balances from '@/components/Balances.vue';
import clock from '@/components/icons/clock-icon.vue'
import bcoin from '@/components/icons/bcoin-icon.vue'
import Navbar from '@/components/Navbar.vue'
import { Vue3Lottie } from 'vue3-lottie'
import RunningJSON from '@/assets/animations/running.json'
import { useAccountStore } from '@/stores/account'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

const { t } = useI18n()


const accStore = useAccountStore()



const remainingTime = ref<string>('');
const currentPoints = ref<number>(0);

const totalDuration = 7200

const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600).toString().padStart(2, '0');
    const minutes = Math.floor((seconds % 3600) / 60).toString().padStart(2, '0');
    const secondsRemaining = (seconds % 60).toString().padStart(2, '0');
    return `${hours}:${minutes}:${secondsRemaining}`;
};

// Function to calculate remaining time and current points
const calculateRemainingTimeAndPoints = () => {
    const now = Math.floor(Date.now() / 1000); // Current time in seconds (Unix timestamp)
    const secondsLeft = accStore.user.lastClaim + accStore.user.farmingTime - now;

    if (secondsLeft <= 0) {
        remainingTime.value = '00:00:00';
        currentPoints.value = accStore.user.maxPoints; // All points earned
        clearInterval(coinsGainInterval.value)
        return;
    }

    remainingTime.value = formatTime(secondsLeft);

    const elapsedTime = totalDuration - secondsLeft;
    currentPoints.value = Math.round((elapsedTime / totalDuration) * accStore.user.maxPoints * 100) / 100;
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

const claim = async () => {
    try {
        const { data } = await axios.post(`${import.meta.env.VITE_API_URL}/farming/claim`, {}, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        accStore.user.farmingTime = data.farmingTime
        accStore.user.lastClaim = data.lastClaim
        accStore.user.pointBalance = data.pointBalance
        accStore.user.maxPoints
        coinsGainInterval.value = setInterval(createSmallCoin, 500)
        return true
    } catch (error) {
        console.error(error)
    }
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
            <span class=" flex justify-center">
                <div class="flex gap-2 p-2 bg-white rounded-full items-center">
                    <bcoin />
                <p class=" text-left text-2xl font-bold w-16">{{ currentPoints }}</p>
                </div>
            </span>
            <span class=" relative mx-10">
                <Vue3Lottie class=" absolute  duration-500 -top-8 left-0"
                    :style="`margin-left: calc(${currentPoints / (accStore.user.maxPoints / 100)}% - 18px)`"
                    :animationData="RunningJSON" :height="36" :width="36" />
                <span
                    class=" relative overflow-hidden flex w-full justify-between text-white px-4 font-bold bg-half_dark rounded-full">
                    <span class=" absolute duration-500 bg-green-400 h-full left-0 top-0"
                        :style="`width: ${currentPoints / (accStore.user.maxPoints / 100)}%`"></span>
                    <p class=" text-green-800 relative z-10">0</p>
                    <p class=" text-green-800 relative z-10">{{ accStore.user.maxPoints }}</p>
                </span>
            </span>
            <section class=" flex relative items-center justify-center" ref="coinsContainer">
                <img id="pulsing" src="../assets/images/farming/pulsing.png" class=" absolute z-10" alt="">
                <img id="coin" src="../components/coin-tap.png" class=" relative z-20" alt="">
            </section>
            <div class=" w-full flex justify-end">
                <div class=" flex flex-col items-center w-32 bg-half_dark p-4 rounded-2xl">
                    <span class="flex gap-2">
                        <clock color="var(--primary)" />
                        <span class=" font-bold">{{ remainingTime }}</span>
                    </span>
                    <p class=" text-dark text-sm text-center">{{ t('screens.bonuses.time') }}</p>
                </div>
            </div>
            <button class="flex items-center justify-center gap-2" @click="claim"
                :disabled="accStore.user.lastClaim + accStore.user.farmingTime - Math.floor(Date.now() / 1000) > 1">
                Claim
                <p class="flex items-center gap-1">
                    <bcoin />
                    {{ accStore.user.maxPoints }}
                </p>
            </button>
        </section>
        <Navbar class="" />
    </section>
</template>

<style>
#coin {
    animation: coin-breath 5s infinite;
    height: 13rem;
    width: 13rem;
}

#pulsing {
    animation: coin-breath 3s infinite;
    opacity: 100%;
    border: solid 2px red;
    height: 17rem;
    width: 17rem;
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
    background-repeat: no-repeat;
    border-radius: 50%;
    z-index: -10;

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
@/stores/balance