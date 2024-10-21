<script lang="ts" setup>
import { onBeforeMount, ref, computed } from 'vue'
import { useAccountStore } from '@/stores/account'
import { useComponentsStore } from '@/stores/components'
import { auth } from '@/utils/helpers'
import { useI18n } from 'vue-i18n'
import axios, { isAxiosError } from 'axios'
import Carousel from 'primevue/carousel'
import carTypeShield from '@/components/cars/CarTypeShield.vue'
import { useRouter } from 'vue-router' 

import RoundCard from '@/components/RoundCard.vue'

const { t } = useI18n()
const router = useRouter()

import { Car } from '@/types/types'

const accStore = useAccountStore()
const componentsStore = useComponentsStore()
const baseURL = import.meta.env.VITE_BASE_URL

const cars = ref<Car[]>([
    {
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
    },
    {
        "element": "city",
        "img": "",
        "acceleration": 0,
        "hendling": 0,
        "brakes": 0,
        "strength": 0,
        "tank": 0,
        "fuel": 0,
        "health": 0,
        "modules": null
    },
    {
        "element": "track",
        "img": "",
        "acceleration": 0,
        "hendling": 0,
        "brakes": 0,
        "strength": 0,
        "tank": 0,
        "fuel": 0,
        "health": 0,
        "modules": null
    }
])

const getCars = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/cars/all`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })

        cars.value = data
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await getCars()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        }
    }
}

const getRound = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/round`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        componentsStore.round = data
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await getCars()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        }
    }
}

const chooseCar = async () => {
    try {
        await axios.post(`${import.meta.env.VITE_API_URL}/buy-car?element=${currentElement.value}`, {}, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        router.push('/cars')
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            await auth()
            try {
                await getCars()
            } catch (error) {
                if (isAxiosError(error)) {
                    componentsStore.addError(error.message)
                }
            }
        }
    }
}

const currentPage = ref<number>(0)
const currentElement = computed(() => cars.value[currentPage.value].element)

const pageChange = (event:  number) =>{
    currentPage.value = event
}

onBeforeMount(() => {
    getRound()
    getCars()
})

</script>

<template>
    <section class=" p-4 pb-20 flex flex-col justify-between h-screen">
        <RoundCard v-if="componentsStore.round" :round="componentsStore.round" />

        <div class="cars-info">
            <h3>{{ t('screens.cars.first.choose.header') }}</h3>
            <p>{{ t('screens.cars.first.choose.description') }}</p>
        </div>

        <Carousel :value="cars" :numVisible="1" @update:page="pageChange" :circular="true" :autoplay="false" :showIndicators="true">
            <template #item="value">
                <div class="car-card">
                    <img class=" w-full"
                        :src="`${baseURL}/static/images/cars/${value.data.element}-car.png`" alt="">
                </div>
            </template>
        </Carousel>

        <div class="flex flex-col gap-4">
            <div class="flex gap-2">
                <carTypeShield :currentElement="currentElement"/>
                <div class="shield rating">
                    <img class=" object-cover h-8 w-12" src="../../components/icons/rating-icon.png" alt="">
                    <span >
                        <p>{{ t('screens.cars.first.rating') }}</p>
                        <i class="pi pi-chevron-right" style="color:var(--dark); font-size: 0.75rem;"></i>
                    </span>
                </div>
            </div>
    
            <button @click="chooseCar">{{ t('screens.cars.first.choose.btn') }}</button>
        </div>

    </section>
</template>


<style>

.p-carousel-next-button, .p-carousel-prev-button{
    display: none !important;
}

.p-carousel-indicator-button{
    border-radius: 999px !important;
    width: 0.5rem !important;
    height: 0.5rem !important;
    background: var(--half-dark) !important;
}
.p-carousel-indicator-active{
    border-radius: 999px !important;
}
.p-carousel-indicator-active>.p-carousel-indicator-button{
    background: var(--primary) !important;
    border-radius: 999px !important;
}

</style>

<style scoped>

.car-card{
    @apply p-8 bg-half_dark rounded-2xl;
}

.cars-info{
    @apply text-center py-4;
}
.cars-info>h3{
    @apply font-bold text-xl;
}

.shield{
    @apply w-full bg-half_dark rounded-2xl p-2 px-4 flex flex-col gap-2 items-center;
}
.shield>span{
    @apply label items-center flex;
}
.shield>img{
    @apply w-8
}

.rating>img{
    @apply w-12 object-top
}

</style>