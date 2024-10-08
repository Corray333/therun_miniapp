<script lang="ts" setup>

import { ref, onBeforeMount } from 'vue'
import { Buildings } from '@/types/types'
import { useI18n } from 'vue-i18n'
import axios, { isAxiosError } from 'axios'
import { auth, getUser } from '@/utils/helpers'
import { useAccountStore } from '@/stores/account'
import { useComponentsStore } from '@/stores/components'


import Carousel from 'primevue/carousel'
import Balances from '@/components/Balances.vue'
import BuildingCard from '@/components/chibi/BuildingCard.vue'
import KeyIcon from '@/components/icons/key-icon.vue'

const accStore = useAccountStore()
const componentsStore = useComponentsStore()
const { t } = useI18n()

const buildings = ref<Buildings>({
    "fabric": {
        state: "build",
        "lastStateChange": 1728409781,
        "stateUntil": 1728417099,
        "img": "https://notably-great-coyote.ngrok-free.app/static/images/buildings/fabric0.png",
        "type": "fabric",
        "level": 0,
        "upgradeCost": null
    },
    "mine": {
        state: "build",
        "lastStateChange": 1728409781,
        "stateUntil": 1728417099,
        "img": "https://notably-great-coyote.ngrok-free.app/static/images/buildings/mine0.png",
        "type": "mine",
        "level": 0,
        "upgradeCost": null
    },
    "warehouse": {
        state: "build",
        "lastStateChange": 1728409781,
        "stateUntil": 1728431136,
        "img": "https://notably-great-coyote.ngrok-free.app/static/images/buildings/warehouse0.png",
        "type": "warehouse",
        "level": 1,
        "upgradeCost": [
            {
                "currency": "point",
                "amount": 1000
            },
            {
                "currency": "blue_key",
                "amount": 1
            }
        ]
    }
})

const getCity = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/city`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        buildings.value = data
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            if (error.response?.status === 401) {
                await auth()
                try {
                    await getCity()
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

onBeforeMount(()=>{
    getCity()
})

</script>

<template>
    <section class=" pb-20">
        <section class=" p-4 flex flex-col gap-4">
            <Balances />

            <h1 class="text-center">{{ t('screens.chibi.city.yourCityHeader') }}</h1>

            <div class="buildings">
                <router-link to="/chibi/city/warehouse">
                    <BuildingCard :building="buildings?.warehouse" />
                </router-link>
                <BuildingCard :building="buildings?.mine" />
            </div>
        </section>

    </section>
</template>


<style scoped>

.buildings{
    @apply flex flex-col gap-4
}

</style>