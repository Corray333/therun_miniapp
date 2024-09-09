<script lang="ts" setup>
import { ref, onBeforeMount } from 'vue'
import { useAccountStore } from '@/stores/account'
import { Referal } from '@/types/types'
import Navbar from '@/components/Navbar.vue'
import CopyIcon from '@/components/icons/copy-icon.vue'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

const { t } = useI18n()

const accStore = useAccountStore()

const appUrl = import.meta.env.VITE_APP_URL

const friends = ref<Referal[]>([
])

const friendsInfo = ref({
    count: 0,
    level: 1,
    nextLevelCount: 3,
    previousLevelCount: 0
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
        alert(error)
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
        alert(error)
    }
}

onBeforeMount(() => {
    getFriends()
    getFriendsInfo()
})

const shareOnTelegram = () => {
    const url = encodeURIComponent(`${appUrl}?startapp=${accStore.user.refCode}`);
    const text = encodeURIComponent('Check out this cool app!');
    const telegramUrl = `https://telegram.me/share/url?text=${text}&url=${url}`;
    window.open(telegramUrl, '_blank');
}

</script>

<template>
    <section>
        <section class=" flex flex-col p-4 gap-4 h-full w-full">
            <div class="info w-full bg-half_dark p-4 rounded-2xl flex flex-col gap-2">
                <h1>{{ t("screens.friens.call") }}</h1>
                <div class="level flex flex-col gap-2">
                    <div class="flex justify-between">
                        <p>{{ t("screens.friens.levelLabel") }} {{ friendsInfo.level }}</p>
                        <p>{{ t("screens.friens.friendsLabel") }} {{ friendsInfo.count }}/{{ friendsInfo.nextLevelCount
                            }}</p>
                    </div>
                    <div class="progress relative h-6 rounded-full w-full bg-white overflow-hidden">
                        <div :style="`width: ${Math.round((friendsInfo.count - friendsInfo.previousLevelCount) / (friendsInfo.nextLevelCount - friendsInfo.previousLevelCount) * 100)}%;`"
                            class="progress-bar relative z-10 bg-green h-full"></div>
                        <div class=" absolute top-0 left-0 w-full h-full flex justify-around">
                            <span class=" h-full w-0.5 bg-half_dark z-20"></span>
                            <span class=" h-full w-0.5 bg-half_dark z-20"></span>
                            <span class=" h-full w-0.5 bg-half_dark z-20"></span>
                            <span class=" h-full w-0.5 bg-half_dark z-20"></span>
                        </div>
                    </div>
                    <p class=" bg-white rounded-xl p-2 relative flex items-center">
                    <p>{{ appUrl }}?startapp={{ accStore.user.refCode }}</p>
                    <CopyIcon @click="copyRefUrl" class=" absolute right-0 mr-2 bg-white" />
                    </p>
                    <button @click="shareOnTelegram">{{ t("screens.friens.shareLink") }}</button>
                </div>
            </div>
            <div class="friends w-full flex flex-col gap-2">
                <h1>{{ t("screens.friens.header") }}</h1>
                <div class="friend-list flex flex-col gap-1">
                    <div v-for="(friend, i) of friends" :key="i" class=" flex p-2 bg-half_dark gap-2 items-center">
                        <img :src="friend.avatar" alt="avatar" class=" w-12 h-12 rounded-full object-cover">
                        <p>{{ friend.username }}</p>
                    </div>
                </div>
            </div>
        </section>
        <Navbar class=" fixed bottom-0" />
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
</style>