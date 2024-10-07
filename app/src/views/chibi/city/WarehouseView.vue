<script lang="ts" setup>

import { ref, onBeforeMount } from 'vue'
import { Warehouse } from '@/types/types'
import { useI18n } from 'vue-i18n'
import axios, { isAxiosError } from 'axios'
import { auth, getUser } from '@/utils/helpers'
import { useAccountStore } from '@/stores/account'
import { useComponentsStore } from '@/stores/components'


import Carousel from 'primevue/carousel'
import Balances from '@/components/Balances.vue'
import BuildingCardTiny from '@/components/chibi/BuildingCardTiny.vue'
import KeyIcon from '@/components/icons/key-icon.vue'

const accStore = useAccountStore()
const componentsStore = useComponentsStore()
const { t } = useI18n()

const baseURL = import.meta.env.VITE_BASE_URL


const warehouse = ref<Warehouse>({
	"img": "https://notably-great-coyote.ngrok-free.app/static/images/buildings/warehouse2.png",
	"type": "warehouse",
	"level": 2,
	"state": "idle",
	"lastStateChange": 0,
	"resources": [
		{
			"name": "",
			"type": "quartz",
			"amount": 10
		},
		{
			"name": "",
			"type": "titan",
			"amount": 0
		}
	],
	"currentLevel": {
		"capacity": 2000,
		"resources": [
			"titan",
			"quartz"
		],
		"cost": [
			{
				"currency": "point",
				"amount": -2000
			},
			{
				"currency": "blue_key",
				"amount": -2
			}
		],
		"requirements": [
			{
				"type": "warehouse",
				"level": 1
			}
		],
		"buildingDuration": 7200
	},
	"nextLevel": {
		"capacity": 2000,
		"resources": [
			"titan",
			"quartz"
		],
		"cost": [
			{
				"currency": "point",
				"amount": -2000
			},
			{
				"currency": "blue_key",
				"amount": -2
			}
		],
		"requirements": [
			{
				"type": "warehouse",
				"level": 1
			}
		],
		"buildingDuration": 7200
	}
})

const capacityUsed = ref(0)
const maxLevel = ref<boolean>(false)

const getWarehouse = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/city/warehouse`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        warehouse.value = data
        capacityUsed.value = 0
        for (const resource of data.resources) {
            capacityUsed.value += resource.amount
        }

        if (data.nextLevel == null){
            maxLevel.value = true
        }
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            if (error.response?.status === 401) {
                await auth()
                try {
                    await getWarehouse()
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

const upgrade = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/city/warehouse/upgrade`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        
        getWarehouse()
    } catch (error) {
        if (isAxiosError(error) && error.response?.status === 401) {
            if (error.response?.status === 401) {
                await auth()
                try {
                    await upgrade()
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
    getWarehouse()
})

</script>

<template>
    <section class=" pb-20">
        <section class=" p-4 flex flex-col gap-4">
            <Balances />

            <img class=" mx-auto max-w-80" :src="`${baseURL}/static/images/buildings/${warehouse.type}${warehouse.level>0?warehouse.level:1}.png`" alt="">
            <h1>{{ t(`screens.chibi.city.buildings.${warehouse.type}.name`) }}</h1>
            <div class="card">
                <p>
                    <p class=" font-bold">{{ t(`screens.chibi.city.buildings.${warehouse.type}.header`) }}</p>
                    <p class=" text-dark">{{ t(`screens.chibi.city.buildings.${warehouse.type}.fullDescription`) }}</p>
                </p>
                <div class="flex gap-4">
                    <div class="shield">
                        <div class="text-center">
                            <p class="level">{{ warehouse.level }}</p>
                            <p class="label">{{ t('screens.chibi.city.buildings.level') }}</p>
                        </div>
                        <i class=" text-dark pi pi-chevron-right"></i>
                    </div>
                    <div class="shield">
                        <div class="text-center">
                            <p class=" font-bold">{{ capacityUsed }}/{{ warehouse.currentLevel.capacity }}</p>
                            <p class="label">{{ t('screens.chibi.city.buildings.warehouse.filled') }}</p>
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="card">
                <p>
                    <p class=" font-bold">{{ t(`screens.chibi.city.buildings.upgrading`) }}</p>
                    <p class=" text-dark">{{ t(`screens.chibi.city.buildings.${warehouse.type}.upgrade`) }}</p>
                </p>
                <div class="flex gap-4">
                    <div class="shield">
                        <div class="text-center">
                            <div class="flex items-center gap-2">
                                <p class="level">{{ warehouse.level }}</p>
                                <i v-if="!maxLevel" class=" pi pi-arrow-right" style="color:var(--blue)"></i>
                                <p v-if="!maxLevel" class="level">{{ warehouse.level+1 }}</p>
                            </div>
                            <p class="label">{{ t('screens.chibi.city.buildings.level') }}</p>
                        </div>
                        <i class=" text-dark pi pi-chevron-right"></i>
                    </div>
                    <div class="shield">
                        <div class="text-center">
                            <div class="flex gap-1">
                                <p class="flex items-center gap-1" v-for="(cost, i) of warehouse.nextLevel?.cost">
                                    <img class="h-6" :src="`${baseURL}/static/images/resources/${cost.currency}.png`" alt="">
                                    <p>{{ -cost.amount }}</p>
                                </p>
                            </div>
                            <p class="label">{{ t('screens.chibi.city.buildings.cost') }}</p>
                        </div>
                    </div>
                </div>

                <div v-if="!maxLevel" class="requirements">
                    <p class=" font-bold mb-2">{{ t('screens.chibi.city.buildings.requirements') }}</p>
                    <BuildingCardTiny class="requirement" v-for="(building, i) of warehouse.nextLevel?.requirements" :key="i" :building="building" />
                </div>

                <button class=" py-2">{{ t('screens.chibi.city.buildings.upgrade') }}</button>
            </div>
        </section>

    </section>
</template>


<style scoped>

.card{
    @apply p-4 rounded-2xl bg-half_dark flex flex-col gap-4;
}

.shield{
    @apply bg-white rounded-2xl flex items-center w-full p-2;
}
.shield>div{
    @apply w-full flex flex-col items-center;
}

.level{
    @apply  px-2 aspect-square flex items-center justify-center bg-custom_blue text-white rounded-md;
}

</style>