<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from 'vue'
import Balances from '@/components/Balances.vue';
import bcoinXL from '@/components/icons/bcoin-icon-xl.vue'
import bcoin from '@/components/icons/bcoin-icon.vue'
import Navbar from '@/components/Navbar.vue'
import { useAccountStore } from '@/stores/account'
import { useComponentsStore } from '@/stores/components'
import { useI18n } from 'vue-i18n'
import { auth } from '@/utils/helpers'
import axios, { isAxiosError } from 'axios'

const { t } = useI18n()
const accStore = useAccountStore()

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
    // if (runningAnimation.value != null) runningAnimation.value.play();
    coinsGainInterval = setInterval(createSmallCoin, 500);
};

const stopAnimations = () => {
    // if (runningAnimation.value != null) runningAnimation.value.goToAndStop(4, true);
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
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await claim()
            } catch (error) {
                alert(error)
            }
        }
        console.log(error);
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
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await start
            } catch (error) {
                alert(error)
            }
        }
        console.log(error);
    }
};

const coinsContainer = ref<HTMLDivElement>();
const coinsFarmedEl = ref<HTMLElement>();

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

let style: HTMLStyleElement
onMounted(() => {
    style = document.createElement('style');
    document.head.appendChild(style);

    calculateRemainingTimeAndPoints()
    if (accStore.user.farmingFrom > accStore.user.lastClaim) {
        startAnimations()
    }
    setInterval(calculateRemainingTimeAndPoints, 1000)
});




const claimAnimate = async () => {
    componentsStore.animateBonuses = true
    animateCoin.value = true

    let x = 0
    let y = 0
    if (coinsFarmedEl.value) {
        y = componentsStore.bonusesLabelPos[0] - coinsFarmedEl.value.offsetTop
        x = componentsStore.bonusesLabelPos[1] - coinsFarmedEl.value.offsetLeft
    }
    let keyframes = `
        @keyframes animate-coins {
        0% {transform: translate(0px, 0px);}
        100% { transform: translate(${x}px, ${y}px); }
    }`;
    style.innerHTML = keyframes;

    // Append the style element to the head of the document

    setTimeout(() => {
        animateCoin.value = false
        componentsStore.animateBonuses = false
    }, 510 * 4)
}

onUnmounted(() => {
    stopAnimations()
});
</script>


<template>
    <section class=" h-screen flex overflow-hidden flex-col">
        <section class=" h-full flex flex-col gap-4 p-4">
            <Balances />
            <div class="flex items-center relative">
                <h2 class=" ml-4 absolute italic font-bold text-[#523810]">Upgrade your Droid</h2>
                <img class="w-full" src="../assets/images/farming/upgrade-banner.png" alt="">
            </div>
            <span class=" flex justify-center">
                <div ref="coinsFarmedEl" class=" flex gap-2 p-2 bg-white rounded-full items-center">
                    <bcoinXL />
                    <bcoinXL class="absolute anim-coin duration-500" id="anim-coin-1"
                        :class="{ 'animate-coin': animateCoin }" />
                    <p class=" text-left text-4xl italic font-bold w-24">{{ currentPoints }}</p>
                </div>
            </span>
            <section class=" coins-container flex h-full items-center relative justify-center" ref="coinsContainer">
                <img id="coin" src="../components/coin-tap.png" class=" h-full absolute" alt="">
            </section>
            <section class=" flex flex-col gap-4">
                <button
                    v-if="accStore.user.farmingFrom > accStore.user.lastClaim && accStore.user.farmingFrom + accStore.user.farmingTime - Math.floor(Date.now() / 1000) > 1"
                    class="flex items-center justify-between" @click="claim" disabled>
                    <p class=" flex gap-2">
                        Farming
                    <p class="flex items-center gap-1">
                        <bcoin />
                        {{ Math.floor(currentPoints) }}/
                        {{ accStore.user.maxPoints }}
                    </p>
                    </p>
                    <p class=" text-left w-20">
                        {{ remainingTime }}
                    </p>
                </button>
                <button v-else-if="accStore.user.farmingFrom > accStore.user.lastClaim"
                    class="flex items-center justify-center gap-2" @click="claim">
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
        <Navbar class="" />
    </section>
</template>

<style>
#coin {
    animation: coin-breath 5s infinite;
    max-height: 17rem;
}

.animate-coin {
    animation: animate-coins 0.5s 4;
}

/* .anim-coin {
    transition: all 0.5s;
} */



/* @keyframes animate-coins {
    0% {
        transform: translate(0, 0);
    }

    100% {
        transform: translate(80px, -175px);
    }
} */

@keyframes coin-breath {
    0% {
        transform: scale(1);
    }

    50% {
        transform: scale(0.9);
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