<script lang="ts" setup>

import { ref, onBeforeMount, watch } from 'vue'
import { Case, Reward } from '@/types/types'
import { useI18n } from 'vue-i18n'
import axios, { isAxiosError } from 'axios'
import { auth, getUser } from '@/utils/helpers'
import { useAccountStore } from '@/stores/account'
import { useComponentsStore } from '@/stores/components'


import Carousel from 'primevue/carousel'
import Navbar from '@/components/Navbar.vue'
import Balances from '@/components/Balances.vue'
import key from '@/components/icons/key-icon.vue'



const { t } = useI18n()
const baseURL = import.meta.env.VITE_BASE_URL
const componentsStore = useComponentsStore()
const accStore = useAccountStore()

const cases = ref<Case[]>([])



const getCases = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/cases`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        for (let i = 0; i < data.length; i++) {
            cases.value.push(data[i])

        }
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            if (error.response?.status === 401) {
                await auth()
                try {
                    await getCases()
                } catch (error) {
                    if (isAxiosError(error)) {
                        componentsStore.addError(error.message)
                    }
                }
            } else {
                componentsStore.addError(error.message)
            }
        }
    }
}

const openCase = async () => {
    loading.value = true
    try {
        const { data } = await axios.post(`${import.meta.env.VITE_API_URL}/cases/open`,{
            caseType: casePick.value
        }, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        getUser(accStore.user.id)
        casePick.value = undefined
        rewardsGot.value = data
        setTimeout(() => {
            rewardsGot.value = undefined
        }, 2500)
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            if (error.response?.status === 401) {
                await auth()
                try {
                    await getCases()
                } catch (error) {
                    if (isAxiosError(error)) {
                        componentsStore.addError(error.message)
                    }
                }
            } else {
                componentsStore.addError(error.message)
            }
        }
    } finally{
        loading.value = false
    }
}

onBeforeMount(async ()=>{
    await getCases()
})



const balanceEnough = (caseToOpen: Case) => {
    for (let i = 0; i < caseToOpen.keys.length; i++) {
        let key = caseToOpen.keys[i].type + 'Balance';
        if (key in accStore.user && typeof (accStore.user as { [key: string]: any })[key] === 'number') {
            if ((accStore.user as { [key: string]: any })[key] < caseToOpen.keys[i].amount) {
                return false
            }
        } else {
            return false
        }
    }
    return true
}

const casePick = ref<string>()
const rewardsGot = ref<Reward>()
const loading = ref<boolean>(false)



</script>


<template>
    <section class=" pb-20">

        <Transition name="delay">
        <section v-show="casePick" @click.self="casePick = ''"
            class=" wrapper fixed z-50 w-full h-screen top-0 left-0 flex items-center justify-center p-4 drop-shadow-lg">
            <Transition name="scale">
                <section v-if="casePick"
                    class=" modal max-w-80 gap-4 w-full rounded-2xl bg-white p-4 flex flex-col justify-center items-center shadow-lg">
                    <p class="font-bold text-center">{{ t('screens.chibi.case.openApprove') }}{{ t(`screens.chibi.case.types.${casePick}`) }}?</p>
                    <div class="flex gap-2 w-full">
                        <button @click="casePick = ''" class=" py-2 text-primary bg-white">{{ t('screens.chibi.case.openApproveCancel') }}</button>
                        <button @click="openCase()" class=" py-2">
                            <p v-if="loading"><i class="pi pi-spinner pi-spin"></i></p>
                            <p v-else>
                                {{ t('screens.chibi.case.openApproveOk') }}
                            </p>
                        </button>
                    </div>
                </section>
            </Transition>
        </section>
    </Transition>

    <Transition name="delay">
        <section v-show="rewardsGot" @click.self="rewardsGot = undefined"
            class=" wrapper fixed z-50 w-full h-screen top-0 left-0 flex items-center justify-center p-4 drop-shadow-lg">
            <Transition name="scale">
                <section v-if="rewardsGot"
                    class=" modal max-w-80 gap-4 w-full rounded-2xl bg-white p-4 flex flex-col justify-center items-center shadow-lg">
                    <div class="hero w-full aspect-square relative flex justify-center items-center">
                        <img class=" absolute w-24 h-24 animate-wiggle" :src="`${baseURL}/static/images/resources/${rewardsGot.type}.png`" alt="">
                        <img style="animation-delay: -6.4s; rotate: -28deg;" class=" absolute mt-32 mr-32 blur-[0.1rem] w-16 h-16 animate-wiggle" :src="`${baseURL}/static/images/resources/${rewardsGot.type}.png`" alt="">
                        <img style="animation-delay: -3.7s; rotate: 24deg;" class=" absolute -mt-32 mr-36 blur-[0.15rem] w-12 h-12 animate-wiggle" :src="`${baseURL}/static/images/resources/${rewardsGot.type}.png`" alt="">
                        <img style="animation-delay: -8.2s; rotate: -24deg;" class=" absolute -mt-24 ml-44 blur-[0.15rem] w-12 h-12 animate-wiggle" :src="`${baseURL}/static/images/resources/${rewardsGot.type}.png`" alt="">
                        <img style="animation-delay: -10.8s; rotate: -36deg;" class=" absolute mt-32 ml-36 blur-[0.15rem] w-14 h-14 animate-wiggle" :src="`${baseURL}/static/images/resources/${rewardsGot.type}.png`" alt="">
                        <img style="animation-delay: -2.6s; rotate: 32deg;" class=" absolute mt-4 mr-36 blur-[0.15rem] w-6 h-6 animate-wiggle" :src="`${baseURL}/static/images/resources/${rewardsGot.type}.png`" alt="">
                        <img style="animation-delay: -7.3s; rotate: -36deg;" class=" absolute -mt-44 ml-8 blur-[0.15rem] w-8 h-8 animate-wiggle" :src="`${baseURL}/static/images/resources/${rewardsGot.type}.png`" alt="">
                        <img style="animation-delay: -9.1s; rotate: 18deg;" class=" absolute mt-44 ml-4 blur-[0.15rem] w-8 h-8 animate-wiggle" :src="`${baseURL}/static/images/resources/${rewardsGot.type}.png`" alt="">
                    </div>
                    <p class="font-bold flex justify-center w-full gap-1 text-2xl items-center">
                        <img class=" w-8 h-8" :src="`${baseURL}/static/images/resources/${rewardsGot.type}.png`" alt="">{{ rewardsGot.amount }}
                    </p>
                    <button @click="rewardsGot = undefined">{{ t('screens.chibi.case.openApproveCancel') }}</button>
                </section>
            </Transition>
        </section>
    </Transition>

        <section class=" flex flex-col p-4 h-full w-full">
            <Balances />
        </section>
        <section class=" flex flex-col gap-6 pb-12 items-center p-4">   

            <Carousel :value="cases" :numVisible="1" :circular="false" :numScroll="1" :showIndicators="false" class="w-full">
                <template #item="value" class="w-full">
                    <div class="flex flex-col w-full items-center gap-4">
                        <p class=" font-bold text-2xl">{{ t(`screens.chibi.case.types.${value.data.type}`) }}</p>
                        <p class="flex items-center gap-1"><p class=" text-dark">{{ t('screens.chibi.case.inside') }}</p> <img class=" w-6 h-6"
                                :src="`${baseURL}/static/images/resources/${value.data.rewardType}.png`" alt=""> {{
                value.data.min_rewards }}-{{ value.data.max_rewards }}</p>
                        <img id="case" :src="`${baseURL}/static/images/cases/${value.data.type}.png`" alt="">
                        <div class="flex gap-2">
                            <span v-for="key in value.data.keys" :key="key.type"
                                class="flex gap-2 w-full justify-center">
                                <span v-for="i of key.amount" :key="i" class="rounded-full bg-white p-4 sm-shadow">
                                    <key class=" opacity-50" :color="`var(--${key.type})`" />
                                </span>
                            </span>
                        </div>
                        <button @click="casePick = value.data.type" class="" :disabled="!balanceEnough(value.data)">{{ t('screens.chibi.case.openButton')}}</button>
                    </div>
                </template>
            </Carousel>
        </section>

        <section class=" info p-4 rounded-t-2xl bg-white ">
            <div class="game text-center flex flex-col gap-6">
                <span class=" flex flex-col gap-2">
                    <h1 class=" font-bold text-2xl">{{ t("screens.chibi.game.header") }}</h1>
                    <p>{{ t("screens.chibi.game.subheader") }}</p>
                </span>
                <img src="../../assets/images/chibi/phones.png" alt="">
                <div class="flex w-full gap-4">
                    <a href="https://play.google.com/store/apps/details?id=com.therun.app" target="_blank"><img
                            src="../../assets/images/chibi/google-play-btn.png" alt=""></a>
                    <a href="https://apps.apple.com/us/app/therun/id1634366310" target="_blank"><img
                            src="../../assets/images/chibi/appstore-btn.png" alt=""></a>
                </div>
                <p class=" p-4 bg-half_dark rounded-2xl">{{ t("screens.chibi.game.description") }}</p>
            </div>
        </section>
        <div class="slider">
            <div class="slide-track">
                <img src="../../assets/images/chibi/slider/box.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-1.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-2.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-3.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-4.png" alt="" class="slide">
                <!-- Дублируем изображения для бесконечного эффекта -->
                <img src="../../assets/images/chibi/slider/box.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-1.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-2.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-3.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-4.png" alt="" class="slide">
            </div>
        </div>
    </section>
</template>

<style>

.delay-enter-active,
.delay-leave-active {
    transition: opacity 0.5s ease;
}

.delay-enter-from,
.delay-leave-to {
    opacity: 0;
}

.scale-enter-active,
.scale-leave-active {
    transition: transform 0.5s ease;
}

.scale-enter-from,
.scale-leave-to {
    transform: scale(0);
}

.p-button {
    background-color: var(--primary) !important;
    min-width: 2.5rem !important;
    min-height: 2.5rem !important;
    max-width: 2.5rem !important;
    max-height: 2.5rem !important;
    color: white !important;
}
</style>

<style scoped>
.info {
    box-shadow: 0 -1rem 1rem 0 rgba(0, 0, 0, 0.1);
}



#case {
    animation: case-breath 2s infinite;
    max-width: 14rem;
}

.slider {
    overflow: hidden;
    width: 100%;
}

.slide-track {
    display: flex;
    gap: 1.4rem;
    animation: scroll 20s linear infinite;
    width: calc(7rem * 10 + 1.4rem * 10);
}

.slide {
    height: 7rem;
    width: 7rem;
    flex-shrink: 0;
}

@keyframes scroll {
    0% {
        transform: translateX(0);
    }

    100% {
        transform: translateX(calc(calc(7rem * 10 + 1.4rem * 10) / -2));
    }
}

@keyframes case-breath {
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
</style>