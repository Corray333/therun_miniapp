<script lang="ts" setup>
import { ref, onBeforeMount } from 'vue'
import { useAccountStore } from '@/stores/account'
import { Referal } from '@/types/types'
import bcoinXL from '@/components/icons/bcoin-icon-xl.vue'
import CopyIcon from '@/components/icons/copy-icon.vue'
import { useI18n } from 'vue-i18n'
import axios, { isAxiosError } from 'axios'
import { auth } from '@/utils/helpers'
import SlideUpDown from 'vue-slide-up-down'
import { useComponentsStore } from '@/stores/components'

const componentsStore = useComponentsStore()


const { t } = useI18n()

const accStore = useAccountStore()

const infoLoaded = ref<boolean>(false)

const appUrl = import.meta.env.VITE_APP_URL

const friends = ref<Referal[]>([
])

const friendsInfo = ref({
    rewardsFrozen: 3000,
    rewardsAvailible: 1000,
    refsActivated: 16,
    refsFrozen: 10,
})

const copyRefUrl = () => {
    navigator.clipboard.writeText(`${appUrl}?startapp=${accStore.user.refCode}`)
}


const getFriends = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/users/0/referals`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        friends.value.push(...data)
        return true
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await getFriends()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        }
    }
}

const getFriendsInfo = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/users/0/referals/info`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        friendsInfo.value = data
        return true
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await getFriendsInfo()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        }
    }
}

onBeforeMount(async () => {
    getFriendsInfo()
    await getFriends()
    infoLoaded.value = true
})

const shareOnTelegram = () => {
    const url = encodeURIComponent(`${appUrl}?startapp=${accStore.user.refCode}`);
    const text = encodeURIComponent('Check out this cool app!');
    const telegramUrl = `https://telegram.me/share/url?text=${text}&url=${url}`;
    window.open(telegramUrl, '_blank');
}

const showInfo = ref<boolean>(false)

</script>

<template>
    <section class=" pb-20">
        <Transition>
            <section v-if="!infoLoaded"
                class=" fixed z-20 top-0 left-0 w-full h-screen bg-white flex justify-center items-center">
                <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
            </section>
        </Transition>
        <section class=" flex flex-col p-4 gap-4 h-full w-full">
            <h1>{{ t("screens.friens.call") }}</h1>
            <div class="frozen bg-half_dark p-4 rounded-2xl flex flex-col items-center gap-2">
                <div class="flex  items-center gap-2 text-2xl font-bold">
                    <bcoinXL />{{ friendsInfo.rewardsFrozen }}
                </div>
                <p class="label">{{ t('screens.friens.frozen') }}</p>
            </div>
            <div class=" bg-half_dark rounded-2xl">
                <button @click="showInfo = !showInfo"
                    class=" bg-half_dark text-black flex items-center justify-start gap-2"><i :style="{transform: showInfo ? 'rotate(180deg)' : 'rotate(0deg)'}"
                        class=" text-dark pi pi-angle-down duration-300"></i>{{ t('screens.friens.info.title') }}</button>
                <SlideUpDown :active="showInfo">
                    <div class="info p-4 pt-0 flex flex-col gap-2">
                        <div class="flex gap-4">
                            <span class=" w-2 h-2 bg-primary rounded-full mt-2"></span>
                            <div class="flex flex-col">
                                <p>{{ t('screens.friens.info.shareTitle') }}</p>
                                <p class="label">{{ t('screens.friens.info.shareDescription') }}</p>
                            </div>
                        </div>
                        <div class="flex gap-4">
                            <span class=" w-2 h-2 bg-primary rounded-full mt-2"></span>
                            <div class="flex flex-col">
                                <p>{{ t('screens.friens.info.activateTitle') }}</p>
                                <p class="label">{{ t('screens.friens.info.activateDescription') }}</p>
                            </div>
                        </div>
                    </div>
                </SlideUpDown>
            </div>
            <div class="frozen bg-half_dark p-4 rounded-2xl flex flex-col items-center gap-2">
                <div class="flex  items-center gap-2 text-2xl font-bold">
                    <bcoinXL />{{ friendsInfo.rewardsAvailible }}
                </div>
                <p class="label">{{ t('screens.friens.availibleForClaim') }}</p>
                <button v-show="friendsInfo.rewardsAvailible>0" class=" btn-type-2">{{ t('screens.friens.claimBtn') }}</button>
            </div>
            <div class="friends flex flex-col gap-2">
                <p class=" w-full flex justify-between px-2">
                    <h1>{{ t('screens.friens.listHeader') }}</h1>
                    <p class=" text-dark font-bold">{{ friendsInfo.refsActivated+friendsInfo.refsFrozen }}</p>
                </p>
                <button class=" bg-half_dark text-black flex justify-between font-medium items-center">
                    <p>{{ t('screens.friens.activated') }}</p>
                    <p class="flex items-center">
                        <p class=" text-green-400 font-bold">{{ friendsInfo.refsActivated }}</p>
                        <i class="pi pi-angle-right text-dark" style="font-size: 1.25rem;"></i>
                    </p>
                </button>
                <button class=" bg-half_dark text-black flex justify-between font-medium items-center">
                    <p>{{ t('screens.friens.waitingForActivation') }}</p>
                    <p class="flex items-center">
                        <p class=" text-orange-400 font-bold">{{ friendsInfo.refsFrozen }}</p>
                        <i class="pi pi-angle-right text-dark" style="font-size: 1.25rem;"></i>
                    </p>
                </button>
            </div>
            <div class="flex gap-2">
                <button @click="shareOnTelegram">{{ t('screens.friens.inviteBtn') }}</button>
                <button @click="copyRefUrl" class=" w-fit aspect-square"><CopyIcon color="white" /></button>
            </div>
        </section>
    </section>
</template>


<style>
.friend-list>div:first-child {
    border-top-left-radius: 1rem;
    border-top-right-radius: 1rem;
}

.friend-list>div:last-child {
    border-bottom-left-radius: 1rem;
    border-bottom-right-radius: 1rem;
}

.v-enter-active,
.v-leave-active {
    transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
    opacity: 0;
}

.slide-up-enter-active,
.slide-up-leave-active {
    transition: all .5s ease;
}

.slide-up-enter,
.slide-up-leave-to {
    height: 0;
    overflow: hidden;
}
</style>