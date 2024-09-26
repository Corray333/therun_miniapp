<script lang="ts" setup>

import Balances from '@/components/Balances.vue'

import { ref, onBeforeMount } from 'vue'
import { Battle as BattleType } from '@/types/types'
import Battle from '@/components/Battle.vue' 

const battles = ref<BattleType[]>([
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
        "opponentResult": 10.5
    }
])


</script>

<template>
    <section class="pb-20">
        <section class=" p-4 flex flex-col">
            <Battle v-for="(battle, i) of battles" :key="i" :battle="battle" />
        </section>
    </section>
</template>


<style></style>