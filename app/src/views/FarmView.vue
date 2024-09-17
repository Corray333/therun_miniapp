<script lang="ts" setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
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
    const now = Date.now() / 1000
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

    remainingTime.value = formatTime(Math.floor(secondsLeft))
    const elapsedTime = totalDuration - secondsLeft;
    currentPoints.value = Math.round((elapsedTime / totalDuration) * accStore.user.maxPoints * 100) / 100;
}

const startAnimations = () => {
    coinsGainInterval = setInterval(createSmallCoin, 500);
};

const stopAnimations = () => {
    if (coinsGainInterval) clearInterval(coinsGainInterval);
};

const animateCoin = ref<boolean>(false)
const componentsStore = useComponentsStore()

watch(() => accStore.user.farmingFrom, () => {
    if (accStore.user.farmingFrom > accStore.user.lastClaim) {
        startAnimations()
    } else {
        stopAnimations()
    }
})

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
        startAnimations()
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

let scale = 1
let style: HTMLStyleElement
let breathInterval: number | null = null

onMounted(() => {
    style = document.createElement('style');
    document.head.appendChild(style);

    startBreath()

    calculateRemainingTimeAndPoints()
    if (accStore.user.farmingFrom > accStore.user.lastClaim) {
        startAnimations()
    }
    setInterval(calculateRemainingTimeAndPoints, 500)
})

const startBreath = ()=>{
    breathInterval = setInterval(() => {
        scale = scale === 1 ? 0.9 : 1
        if (tapCoin.value) tapCoin.value.style.transform = `scale(${scale})`
    }, 2500)
}


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


    setTimeout(() => {
        animateCoin.value = false
        componentsStore.animateBonuses = false
    }, 510 * 4)
}

const tapCoin = ref<HTMLElement>()

const tap = () => {
    if (tapCoin.value) {
        tapCoin.value.style.transition = 'all 0.2s'
        tapCoin.value.style.transform = 'scale(0.85)'
        setTimeout(() => {
            if (tapCoin.value){
                tapCoin.value.style.transition = 'all 2.5s'
                tapCoin.value.style.transform = 'scale(1)'
                scale = 1
                if (breathInterval)clearInterval(breathInterval)
                startBreath()
            }
        }, 200)
    }
}

onUnmounted(() => {
    stopAnimations()
});
</script>


<template>
    <section class=" h-screen flex overflow-hidden flex-col">
        <section class=" h-full flex flex-col gap-4 p-4">
            <Balances />
            <!-- <div class=" banner flex items-center bg-cover rounded-2xl h-28 w-full">
                <h2 class=" ml-4 absolute italic font-bold text-[#523810]">{{ t('screens.farming.banner') }}</h2>
            </div> -->
            <span class=" flex justify-center">
                <div ref="coinsFarmedEl" class=" flex gap-2 p-2 bg-white rounded-full items-center">
                    <bcoinXL />
                    <bcoinXL class="absolute anim-coin duration-500" id="anim-coin-1"
                        :class="{ 'animate-coin': animateCoin }" />
                    <p class=" text-left text-4xl font-bold w-24">{{ currentPoints }}</p>
                </div>
            </span>
            <section class=" coins-container flex h-full items-center relative justify-center" ref="coinsContainer">
                <img ref="tapCoin" @click="tap" id="coin" src="../components/coin-tap.png" class=" h-full absolute" alt="">
            </section>
            <div class="more-points grid grid-cols-3 gap-2">
                <div class=" bg-half_dark rounded-2xl flex flex-col items-center">
                    <p class=" mt-2">{{ t('screens.farming.earning.getMore') }}</p>
                    <bcoin id="more-btn-coin" />
                </div>
                <div class=" bg-half_dark rounded-2xl flex flex-col items-center">
                    <p class=" mt-2">{{ t('screens.farming.earning.upgrade') }}</p>
                    <img class=" h-8 object-contain object-bottom" src="../assets/images/farming/robot.png" alt="">
                </div>
                <div class=" bg-full_dark rounded-2xl flex flex-col text-white text-bold items-center">
                    <p class=" mt-2 font-bold">{{ t('screens.farming.earning.ninja') }}</p>
                    <img class=" h-8 object-contain object-bottom" src="../assets/images/farming/spy.png" alt="">
                </div>
            </div>
            <section class=" flex flex-col gap-4">
                <button
                    v-if="accStore.user.farmingFrom > accStore.user.lastClaim && accStore.user.farmingFrom + accStore.user.farmingTime - Math.floor(Date.now() / 1000) > 1"
                    class="flex items-center justify-between" @click="claim" disabled>
                    <p class=" flex gap-2">
                        {{ t('screens.farming.button.farming') }}
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
                    {{ t('screens.farming.button.claim') }}
                    <p class="flex items-center gap-1">
                        <bcoin />
                        {{ accStore.user.maxPoints }}
                    </p>
                </button>
                <button v-else class="flex items-center justify-center gap-2" @click="start"
                    :disabled="accStore.user.farmingFrom + accStore.user.farmingTime - Math.floor(Date.now() / 1000) > 1">
                    {{ t('screens.farming.button.start') }}
                </button>
            </section>
        </section>
        <Navbar class="" />
    </section>
</template>

<style>

.banner{
    background-image:url(../assets/images/farming/upgrade-banner.png);
}

#more-btn-coin{
    animation: rotate-coin 2s infinite;
}

#coin {
    /* animation: coin-breath 5s infinite; */
    transition: all 2.5s;
    max-height: 17rem;
}

.animate-coin {
    animation: animate-coins 0.5s 4;
}


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

@keyframes shaking {
    0% {
        transform: translateX(0);
    }

    25% {
        transform: translateX(-0.5rem);
    }

    50% {
        transform: translateX(0);
    }

    75% {
        transform: translateX(0.5rem);
    }

    100% {
        transform: translateX(0);
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

@keyframes rotate-coin {
    0% {
        transform: rotateY(0deg);
    }
    100% {
        transform: rotateY(360deg);
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