<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from 'vue'
import Balances from '@/components/Balances.vue';
import clock from '@/components/icons/clock-icon.vue'
import bcoin from '@/components/icons/bcoin-icon.vue'
import Navbar from '@/components/Navbar.vue'
import { Vue3Lottie } from 'vue3-lottie'
import { type AnimationItem } from 'lottie-web'
import RunningJSON from '@/assets/animations/running.json'
import { useAccountStore } from '@/stores/account'
import { useComponentsStore } from '@/stores/components'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

const { t } = useI18n()
const accStore = useAccountStore()

const runningAnimation = ref<AnimationItem | null>(null)
const remainingTime = ref<string>('00:00:00');
const currentPoints = ref<number>(0);

const totalDuration = 7200
let coinsGainInterval: number | null = null;

const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600).toString().padStart(2, '0');
    const minutes = Math.floor((seconds % 3600) / 60).toString().padStart(2, '0');
    const secondsRemaining = (seconds % 60).toString().padStart(2, '0');
    return `${hours}:${minutes}:${secondsRemaining}`;
};

const calculateRemainingTimeAndPoints = () => {
    const now = Math.floor(Date.now() / 1000);
    let secondsLeft = accStore.user.farmingFrom + accStore.user.farmingTime - now;

    if (accStore.user.farmingFrom <= accStore.user.lastClaim) {
        stopAnimations();
        remainingTime.value = '00:00:00';
        currentPoints.value = 0;
        return;
    }

    if (secondsLeft <= 0) {
        stopAnimations();
        remainingTime.value = '00:00:00';
        currentPoints.value = accStore.user.maxPoints;
        return;
    }

    remainingTime.value = formatTime(secondsLeft);
    const elapsedTime = totalDuration - secondsLeft;
    currentPoints.value = Math.round((elapsedTime / totalDuration) * accStore.user.maxPoints * 100) / 100;
};

const startAnimations = () => {
    console.log("startAnimations");
    console.log("startAnimations");
    console.log("startAnimations");
    if (runningAnimation.value != null) runningAnimation.value.play();
    coinsGainInterval = setInterval(createSmallCoin, 500);
};

const stopAnimations = () => {
    if (runningAnimation.value != null) runningAnimation.value.goToAndStop(4, true);
    if (coinsGainInterval) clearInterval(coinsGainInterval);
};

const animateCoin = ref<boolean>(false)
const componentsStore = useComponentsStore()

const claim = async () => {
    try {
        const { data } = await axios.post(`${import.meta.env.VITE_API_URL}/farming/claim`, {}, {
            withCredentials: true,
            headers: { Authorization: accStore.token }
        });

        const diff = data.pointBalance - accStore.user.pointBalance
        const piece = Math.floor(diff / 4)
        alert(JSON.stringify([data.pointBalance, accStore.user.pointBalance, diff, piece]))
        claimAnimate()
        accStore.user.pointBalance += piece
        setTimeout(() => {
            accStore.user.pointBalance += piece
            setTimeout(() => {
                accStore.user.pointBalance += piece
                setTimeout(() => {
                    accStore.user.pointBalance = data.pointBalance
                }, 500);
            }, 500);
        }, 500)
        accStore.user.farmingTime = data.farmingTime
        accStore.user.lastClaim = data.lastClaim
    } catch (error) {
        alert(error);
    }
};

const start = async () => {
    try {
        const { data } = await axios.post(`${import.meta.env.VITE_API_URL}/farming/start`, {}, {
            withCredentials: true,
            headers: { Authorization: accStore.token }
        });

        accStore.user.farmingFrom = data.farmingFrom;
        startAnimations();
    } catch (error) {
        alert(error);
    }
};

const coinsContainer = ref<HTMLDivElement | null>(null);

const createSmallCoin = () => {
    if (!coinsContainer.value) return;

    const smallCoin = document.createElement('div');
    smallCoin.classList.add('small-coin');

    const startX = Math.random() * 200 - 50 + 'vw';
    const startY = Math.random() * 200 - 50 + 'vh';

    smallCoin.style.setProperty('--start-x', startX);
    smallCoin.style.setProperty('--start-y', startY);
    smallCoin.style.animation = 'move-to-center 2s ease-out';

    coinsContainer.value.appendChild(smallCoin);

    setTimeout(() => smallCoin.remove(), 2000);
};

onMounted(() => {
    calculateRemainingTimeAndPoints()
    if (runningAnimation.value != null && accStore.user.farmingFrom > accStore.user.lastClaim) {
        startAnimations()
    }
    setInterval(calculateRemainingTimeAndPoints, 1000)
});

const claimAnimate = async ()=>{
    animateCoin.value = true
    componentsStore.animateBonuses = true
    setTimeout(()=>{
        animateCoin.value = false
        componentsStore.animateBonuses = false
    }, 510*4)
}

onUnmounted(() => {
    stopAnimations()
});
</script>


<template>
    <section class=" h-screen flex overflow-hidden flex-col">
        <section class=" h-full flex flex-col p-4">
            <Balances />
            <section class=" h-full flex flex-col pt-4 justify-between">
                <section class=" flex flex-col gap-4">
                    <span class=" flex justify-center">
                        <div class=" relative flex gap-2 p-2 bg-white rounded-full items-center">
                            <bcoin />
                            <bcoin class="absolute anim-coin" id="anim-coin-1" :class="{'animate-coin':animateCoin}"/>
                        <p class=" text-left text-2xl font-bold w-16">{{ currentPoints }}</p>
                        </div>
                    </span>
                    <span class=" relative mx-10">
                        <Vue3Lottie class=" absolute  duration-500 -top-8 left-0" ref="runningAnimation"
                            :style="`margin-left: calc(${currentPoints / (accStore.user.maxPoints / 100)}% - 18px)`"
                            :animationData="RunningJSON" :height="36" :width="36"/>
                        <span
                            class=" relative overflow-hidden flex w-full justify-between text-white px-4 font-bold bg-half_dark rounded-full">
                            <span class=" absolute duration-500 bg-green-400 h-full left-0 top-0"
                                :style="`width: ${currentPoints / (accStore.user.maxPoints / 100)}%`"></span>
                            <p class=" text-green-800 relative z-10">0</p>
                            <p class=" text-green-800 relative z-10">{{ accStore.user.maxPoints }}</p>
                        </span>
                    </span>
                </section>
                <section class=" coins=container flex relative items-center justify-center" ref="coinsContainer">
                    <img id="pulsing" src="../assets/images/farming/pulsing.png" class=" absolute z-10" alt="">
                    <img id="coin" src="../components/coin-tap.png" class=" relative z-20" alt="">
                </section>
                <section class=" flex flex-col gap-4">
                    <div class=" w-full flex justify-end">
                        <div class=" flex flex-col justify-center items-center w-36 bg-half_dark p-4 rounded-2xl">
                            <span class="flex w-full gap-2">
                                <clock color="var(--primary)" />
                                <span class=" w-full text-left font-bold">{{ remainingTime }}</span>
                            </span>
                            <p class=" text-dark text-sm text-center">{{ t('screens.bonuses.time') }}</p>
                        </div>
                    </div>
                    <button v-if="accStore.user.farmingFrom > accStore.user.lastClaim" class="flex items-center justify-center gap-2" @click="claim"
                        :disabled="accStore.user.farmingFrom + accStore.user.farmingTime - Math.floor(Date.now() / 1000) > 1">
                        Claim
                        <p class="flex items-center gap-1">
                            <bcoin />
                            {{ accStore.user.maxPoints }}
                        </p>
                    </button>
                    <button v-else class="flex items-center justify-center gap-2" @click="start"
                        :disabled="accStore.user.farmingFrom + accStore.user.farmingTime - Math.floor(Date.now() / 1000) > 1">
                        Start
                    </button>
                </section>
            </section>
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



.animate-coin{
    animation: animate-coins 0.5s 4;
}

#pulsing {
    animation: coin-breath 3s infinite;
    opacity: 75%;
    height: 17rem;
    width: 17rem;
}

@keyframes animate-coins {
    0% {
        transform: translate(0, 0);
    }
    100% {
        transform: translate(100px, -75px);
    }
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