<script lang="ts" setup>
import { ref, onBeforeMount } from 'vue'
import { useAccountStore } from '@/stores/account'
import bcoin from '@/components/icons/bcoin-icon.vue'
import key from '@/components/icons/key-icon.vue'
import bcoinXL from '@/components/icons/bcoin-icon-xl.vue'
import CopyIcon from '@/components/icons/copy-icon.vue'
import { useI18n } from 'vue-i18n'
import axios, { isAxiosError } from 'axios'
import { auth, getUser } from '@/utils/helpers'
import SlideUpDown from 'vue-slide-up-down'
import { useComponentsStore } from '@/stores/components'
import { onMounted } from 'vue'
import { Vue3Lottie } from 'vue3-lottie'
import coinBlast from '@/assets/animations/coin-blast.json'

const componentsStore = useComponentsStore()


const { t } = useI18n()

const accStore = useAccountStore()
const baseURL = import.meta.env.VITE_BASE_URL

const infoLoaded = ref<boolean>(false)

const appUrl = import.meta.env.VITE_APP_URL

class Reward {
    currency!: string;
    amount!: number;
}

class FriendsInfo {
    refsActivated!: number;
    refsFrozen!: number;
    rewardsFrozen!: Reward[];
    rewardsAvailible!: Reward[];
}

const friendsInfo = ref<FriendsInfo>({
	"refsActivated": 2,
	"refsFrozen": 3,
	"rewardsFrozen": [
		{
			"currency": "point",
			"amount": 5000
		},
		{
			"currency": "blue_key",
			"amount": 4
		},
		{
			"currency": "red_key",
			"amount": 3
		}
	],
	"rewardsAvailible": []
})

const copyRefUrl = () => {
    navigator.clipboard.writeText(`${appUrl}?startapp=${accStore.user.refCode}`)
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
    await getFriendsInfo()
    infoLoaded.value = true
})

let style: HTMLStyleElement

onMounted(()=>{
    style = document.createElement('style')
    document.head.appendChild(style)
})

const shareOnTelegram = () => {
    const url = encodeURIComponent(`${appUrl}?startapp=${accStore.user.refCode}`);
    const text = encodeURIComponent('Check out this cool app!');
    const telegramUrl = `https://telegram.me/share/url?text=${text}&url=${url}`;
    window.open(telegramUrl, '_blank');
}

const showInfo = ref<boolean>(false)


const claimRewards = async () => {
    try {
        await axios.post(`${import.meta.env.VITE_API_URL}/users/0/referals/claim`,{}, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        friendsInfo.value.rewardsAvailible = []
        coinBlastAnimation.value?.goToAndStop(0, true)
        coinBlastAnimation.value?.play()
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
    } finally{
        getFriendsInfo()
        getUser(accStore.user.id)
    }
}

const coinBlastAnimation = ref<typeof Vue3Lottie>()

</script>

<template>
    <section class=" pb-20">
        <Vue3Lottie ref="coinBlastAnimation" :animation-data="coinBlast" class="fixed bottom-0 left-0 w-full pointer-events-none" :auto-play="false" :loop="false"/>
        <Transition>
            <section v-if="!infoLoaded"
                class=" fixed z-20 top-0 left-0 w-full h-screen bg-white flex justify-center items-center">
                <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
            </section>
        </Transition>
        
        <section class=" flex flex-col p-4 gap-4 h-full w-full">
            <h1 class="text-center">{{ t("screens.friens.call") }}</h1>
            <div class="frozen bg-half_dark p-4 rounded-2xl flex flex-col items-center gap-2">
                <div v-if="friendsInfo.rewardsFrozen.length" class="flex items-center gap-2 text-2xl font-bold">
                    <p v-for="(reward, i) of friendsInfo.rewardsFrozen" :key="i" class="flex gap-1 items-center">
                        <img class="h-6" :src="`${baseURL}/static/images/resources/${reward.currency}.png`" alt="">
                        {{ reward.amount }}
                    </p>
                </div>
                <p v-else class="text-4xl font-bold">0</p>
                
                <p class="label">{{ t('screens.friens.frozen') }}</p>
            </div>

            <div class="flex flex-col w-full p-4 bg-half_dark rounded-2xl gap-4">

                <div class="flex gap-4 items-center">
                    <img class="h-16" src="../assets/images/friends/tg-logo.png" alt="">
                    <div class=" flex flex-col gap-2">
                        <p class="font-bold">{{ t('screens.friens.invite.commonInviteHeader') }}</p>
                        <span class="flex flex-col gap-1">
                            <p class="text-dark">
                                {{ t('screens.friens.invite.commonInviteDescription') }}
                            </p>
                            <p class="flex gap-1"><bcoin />1000 <key color="var(--blue_key)"/>1</p>
                        </span>
                    </div>
                </div>
                
                <div class="flex gap-4 items-center">
                    <img class="h-16" src="../assets/images/friends/premium-star.gif" alt="">
                    <div class=" flex flex-col gap-2">
                        <p class="font-bold">{{ t('screens.friens.invite.premiumInviteHeader') }}</p>
                        <span class="flex flex-col gap-1">
                            <p class="text-dark">
                                {{ t('screens.friens.invite.premiumInviteDescription') }}
                            </p>
                            <p class="flex gap-1"><bcoin />3000 <key color="var(--red_key)"/>3 <key color="var(--blue_key)"/>2</p>
                        </span>
                    </div>
                </div>

                <div class="flex gap-2">
                    <button @click="shareOnTelegram">{{ t('screens.friens.inviteBtn') }}</button>
                    <button @click="copyRefUrl" class=" w-fit aspect-square">
                        <CopyIcon color="white" />
                    </button>
                </div>

            </div>

            <div class=" bg-half_dark rounded-2xl">
                <button @click="showInfo = !showInfo"
                    class=" bg-half_dark text-black flex items-center justify-start gap-2"><i
                        :style="{ transform: showInfo ? 'rotate(180deg)' : 'rotate(0deg)' }"
                        class=" text-dark pi pi-angle-down duration-300"></i>{{ t('screens.friens.info.title')
                    }}</button>
                <SlideUpDown :active="showInfo">
                    <div class="info p-4 pt-0 flex flex-col gap-2">
                        <div class="flex gap-4">
                            <span class=" w-2 h-2 bg-primary rounded-full min-w-2 mt-2"></span>
                            <div class="flex flex-col">
                                <p>{{ t('screens.friens.info.shareTitle') }}</p>
                                <p class="label">{{ t('screens.friens.info.shareDescription') }}</p>
                            </div>
                        </div>
                        <div class="flex gap-4">
                            <span class=" w-2 h-2 bg-primary rounded-full min-w-2 mt-2"></span>
                            <div class="flex flex-col">
                                <p>{{ t('screens.friens.info.tellYourFriendTitle') }}</p>
                                <p class="label">{{ t('screens.friens.info.tellYourFriendDescription') }}</p>
                            </div>
                        </div>
                        <div class="flex gap-4">
                            <span class=" w-2 h-2 bg-primary rounded-full min-w-2 mt-2"></span>
                            <div class="flex flex-col">
                                <p class="font-bold">{{ t('screens.friens.info.activate1Title') }}</p>
                                <span class="label">
                                    {{ t('screens.friens.info.activate1Description1') }}
                                    <bcoin class=" scale-[80%] inline" />
                                    {{ t('screens.friens.info.activate1Description2') }}
                                </span>
                            </div>
                        </div>
                        <div class="flex gap-4">
                            <span class=" w-2 h-2 bg-primary rounded-full min-w-2 mt-2"></span>
                            <div class="flex flex-col">
                                <p class="font-bold">{{ t('screens.friens.info.activate2Title') }}</p>
                                <span class="label">
                                    {{ t('screens.friens.info.activate2Description1') }}
                                    <key class=" scale-[80%] inline" color="var(--primary)" /> 1
                                    {{ t('screens.friens.info.activate2Description2') }}
                                </span>
                            </div>
                        </div>
                    </div>
                </SlideUpDown>
            </div>

            <div class="frozen bg-half_dark p-4 rounded-2xl flex flex-col items-center gap-2">
                
                <div v-if="friendsInfo.rewardsAvailible.length" class="flex items-center gap-2 text-2xl font-bold">
                    <p v-for="(reward, i) of friendsInfo.rewardsAvailible" :key="i" class="flex gap-1 items-center">
                        <img class="h-6" :src="`${baseURL}/static/images/resources/${reward.currency}.png`" alt="">
                        {{ reward.amount }}
                    </p>
                </div>
                <p v-else class="text-4xl font-bold">0</p>

                <p class="label">{{ t('screens.friens.availibleForClaim') }}</p>
                <button @click="claimRewards" v-show="friendsInfo.rewardsAvailible.length" class=" btn-type-2">{{
                t('screens.friens.claimBtn') }}</button>
            </div>

            <div class="friends flex flex-col gap-2">
                <span class=" w-full flex justify-between px-2">
                    <h1>{{ t('screens.friens.listHeader') }}</h1>
                    <p class=" text-dark font-bold">{{ friendsInfo.refsActivated + friendsInfo.refsFrozen }}</p>
                </span>
                <router-link to="/friens/activated">
                    <button class=" bg-half_dark text-black flex justify-between font-medium items-center">
                        <p>{{ t('screens.friens.activated') }}</p>
                        <span class="flex items-center">
                            <p class=" text-green-400 font-bold">{{ friendsInfo.refsActivated }}</p>
                            <i class="pi pi-angle-right text-dark" style="font-size: 1.25rem;"></i>
                        </span>
                    </button>
                </router-link>

                <router-link to="/friens/not-activated">
                    <button class=" bg-half_dark text-black flex justify-between font-medium items-center">
                        <p>{{ t('screens.friens.waitingForActivation') }}</p>
                        <span class="flex items-center">
                            <p class=" text-orange-400 font-bold">{{ friendsInfo.refsFrozen }}</p>
                            <i class="pi pi-angle-right text-dark" style="font-size: 1.25rem;"></i>
                        </span>
                    </button>
                </router-link>
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