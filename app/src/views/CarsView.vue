<script lang="ts" setup>
import { ref, onBeforeMount, computed } from 'vue'
import { useAccountStore } from '@/stores/account'
import { useI18n } from 'vue-i18n'
import axios, { isAxiosError } from 'axios'
import { auth, getUser } from '@/utils/helpers'
import { useComponentsStore } from '@/stores/components'
import {useRoute, useRouter} from 'vue-router'

import type { Car } from '@/types/types'

import roundCard from '@/components/RoundCard.vue'

const componentsStore = useComponentsStore()


const { t } = useI18n()

const accStore = useAccountStore()
const baseURL = import.meta.env.VITE_BASE_URL

const route = useRoute()
const router = useRouter()
const element = route.query.element

const mainCar = ref<Car>({
    "element": "desert",
    "img": "",
    "acceleration": 0,
    "hendling": 0,
    "brakes": 0,
    "strength": 0,
    "tank": 0,
    "fuel": 0,
    "health": 0,
    "modules": null
})

const getMainCar = async () : Promise<boolean> => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/cars/main`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        if (data) {
            mainCar.value = data
            return true
        } else {
            return false
        }
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                return await getMainCar()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
                return false
            }
        } else {
            return false
        }
    }
}

onBeforeMount(async ()=>{
    if (!await getMainCar()) {
        router.push('/cars/all')
    }
})

</script>

<template>
    <section class=" pb-20 p-4">
        <div class=" shields">
            <roundCard :round="componentsStore.round" />
            <div class="shield">
                <p>{{ t('screens.cars.pitStop.pitStop')}}</p>
                <img src="../components/icons/pit-stop-icon.png" alt="">
            </div>
        </div>
    </section>
</template>


<style scoped>

.shields{
    @apply flex gap-2
}
.shields>*{
    @apply sm-shadow
}

.shield{
    @apply bg-half_dark rounded-2xl p-2 px-4
}

</style>