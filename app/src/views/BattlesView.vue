<script lang="ts" setup>

import Balances from '@/components/Balances.vue'

import { ref, onBeforeMount } from 'vue'
import { useI18n } from 'vue-i18n'
import { Round } from '@/types/types'
import Battle from '@/components/Battle.vue'

const { t } = useI18n()

const remainingTime = ref<string>('00:00:00');

const round = ref<Round>({
    endTime: 1727421132,
    id: 1,
    battles: [
        {
            "id": 1,
            "user": {
                "id": 1,
                "username": "incetro",
                "email": "email@example.com",
                "phone": "9990123456",
                "sex_id": "1",
                "birth_date": "05.05.1995",
                "firstname": "Иван",
                "middlename": "Иванович",
                "lastname": "Иванов",
                "photo": "https://store-images.s-microsoft.com/image/apps.55245.13537716651231321.3067a421-6c2f-48a9-b77c-1e38e19146e6.10e2aa49-52ca-4e79-9a61-b6422978afb9?h=210",
                "cover": "http://example.link/cover.png",
                "bio": "mybio",
                "isFollower": true,
                "followersCount": 100000,
                "followingCount": 100000,
                "favoritesCount": 100000,
                "friendsCount": 100000,
                "activeTournamentsCount": 100,
                "isFriend": false,
                "isFavorite": false,
                "city": "Сочи",
                "exp": 12.2,
                "expYear": 12.2,
                "expOverall": 12.2,
                "country": "Сочи",
                "expRatingPlace": 1,
                "transport": [
                    0
                ],
                "fPoints": 1,
                "league": {
                    "id": 0,
                    "name": "1",
                    "bottomThreshold": 100,
                    "topThreshold": 300
                },
                "refCode": "S4DGXC",
                "inviteCode": "S4DGXC",
                "appliedRefCode": "YNAUVG",
                "codeGenerationsLeft": 1,
                "social": [
                    {
                        "name": "VK",
                        "url": "https://vk.com/user"
                    }
                ],
                "minBet": [
                    {
                        "currencyName": "R",
                        "currencyId": 1,
                        "minBet": 50
                    }
                ]
            },
            "opponent": {
                "id": 1,
                "username": "incetro",
                "email": "email@example.com",
                "phone": "9990123456",
                "sex_id": "1",
                "birth_date": "05.05.1995",
                "firstname": "Иван",
                "middlename": "Иванович",
                "lastname": "Иванов",
                "photo": "https://store-images.s-microsoft.com/image/apps.55245.13537716651231321.3067a421-6c2f-48a9-b77c-1e38e19146e6.10e2aa49-52ca-4e79-9a61-b6422978afb9?h=210",
                "cover": "http://example.link/cover.png",
                "bio": "mybio",
                "isFollower": true,
                "followersCount": 100000,
                "followingCount": 100000,
                "favoritesCount": 100000,
                "friendsCount": 100000,
                "activeTournamentsCount": 100,
                "isFriend": false,
                "isFavorite": false,
                "city": "Сочи",
                "exp": 12.2,
                "expYear": 12.2,
                "expOverall": 12.2,
                "country": "Сочи",
                "expRatingPlace": 1,
                "transport": [
                    0
                ],
                "fPoints": 1,
                "league": {
                    "id": 0,
                    "name": "1",
                    "bottomThreshold": 100,
                    "topThreshold": 300
                },
                "refCode": "S4DGXC",
                "inviteCode": "S4DGXC",
                "appliedRefCode": "YNAUVG",
                "codeGenerationsLeft": 1,
                "social": [
                    {
                        "name": "VK",
                        "url": "https://vk.com/user"
                    }
                ],
                "minBet": [
                    {
                        "currencyName": "R",
                        "currencyId": 1,
                        "minBet": 50
                    }
                ]
            },
            "status": "active",
            "userResult": 10.5,
            "opponentResult": 10.5,
            "pick": 2
        }
    ]
})

const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600).toString().padStart(2, '0');
    const minutes = Math.floor((seconds % 3600) / 60).toString().padStart(2, '0');
    const secondsRemaining = (seconds % 60).toString().padStart(2, '0');
    return `${hours}:${minutes}:${secondsRemaining}`;
};

const calculateRemainingTimeAndPoints = () => {
    const now = Date.now() / 1000
    let secondsLeft = round.value.endTime - now;

    if (secondsLeft <= 0) {
        remainingTime.value = '00:00:00'
        return
    }

    remainingTime.value = formatTime(Math.floor(secondsLeft))
}

onBeforeMount(() => {
    calculateRemainingTimeAndPoints()
    setInterval(calculateRemainingTimeAndPoints, 1000)
})


</script>

<template>
    <section class="pb-20">
        <section class=" p-4 flex flex-col gap-4">
            <Balances/>

            <div class="timer flex flex-col items-center justify-center p-4 rounded-2xl bg-secondary">
                <p class=" text-4xl font-bold">{{ remainingTime }}</p>
                <p class="font-bold">{{ t('screens.battles.worldRound') }} #{{ round.id }}</p>
            </div>

            <div class="banner rounded-2xl w-full p-4 bg-cover text-white ">
                <p class="font-bold text-xl">{{ t('screens.battles.banner.header') }}</p>
                <p class="">{{ t('screens.battles.banner.description') }}</p>
            </div>

            <section>
                <Battle v-for="(battle, i) of round.battles" :key="i" :battle="battle" />
            </section>

            <span class="flex gap-2 text-dark">
                <i class=" pi pi-info-circle mt-1" style="font-size: 1.25rem;"></i>
                <p>{{ t('screens.battles.info') }}</p>
            </span>
        </section>
    </section>
</template>


<style scoped>

.banner{
    background-image: url(../assets/images/battles/banner-bg.png);
}

</style>