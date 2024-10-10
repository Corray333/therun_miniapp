<script lang="ts" setup>

import { ref, onBeforeMount, computed } from 'vue'
import { Warehouse, Resource } from '@/types/types'
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
	"img": "https://notably-great-coyote.ngrok-free.app/static/images/buildings/warehouse1.png",
	"type": "warehouse",
	"level": 0,
	"state": "build",
	"lastStateChange": 1728414635,
	"stateUntil": 1728415235,
	"resources": [
		{
			"name": "quartz",
			"type": "mineral",
			"amount": 0
		},
		{
			"name": "titan",
			"type": "ore",
			"amount": 0
		}
	],
	"currentLevel": null,
	"nextLevel": {
		"capacity": 1000,
		"cost": [
			{
				"currency": "point",
				"amount": -1000
			},
			{
				"currency": "blue_key",
				"amount": -1
			}
		],
		"requirements": null,
		"buildingDuration": 600
	},
	"moreLevel": {
		"capacity": 2000,
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

const loaded = ref<boolean>(false)


const resourceDictionary = computed(() => {
    return warehouse.value.resources.reduce((acc, resource) => {
        if (acc[resource.type]) {
            acc[resource.type].totalAmount += resource.amount;
            acc[resource.type].resources.push(resource);
        } else {
            acc[resource.type] = {
                totalAmount: resource.amount,
                resources: [resource],
            };
        }

        return acc;
    }, {} as { [key: string]: { totalAmount: number; resources: Resource[] } });
});

const capacityUsed = computed(() => {
    let used = 0
    for (const resource of warehouse.value.resources) {
        used += resource.amount
    }
    return used
})
const maxLevel = computed(() => warehouse.value.nextLevel === null)
const notBought = computed(() => warehouse.value.level === 0)
const balanceEnoughForUpgrade = computed(() => {
    if (warehouse.value.nextLevel === null) return false
    for (const cost of warehouse.value.nextLevel.cost) {
        if ((accStore.user as any)[`${cost.currency}Balance`] < -cost.amount) return false
    }
    return true
})

const getWarehouse = async () => {
    try {
        const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/city/warehouse`, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })
        warehouse.value = data
        loaded.value = true
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
    loading.value = true
    try {
        await axios.patch(`${import.meta.env.VITE_API_URL}/city/warehouse/upgrade`, {}, {
            withCredentials: true,
            headers: {
                Authorization: accStore.token
            }
        })

        getWarehouse()
        getUser(accStore.user.id)
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
    } finally {
        loading.value = false
        showUpgradeModal.value = false
    }
}

onBeforeMount(() => {
    getWarehouse()
})

const showUpgradeModal = ref<boolean>(false)
const loading = ref<boolean>(false)

const sprintf = (format: string, ...args: any[]): string => {
    for (let i = 0; i < args.length; i++) {
        format = format.replace(/%s/, args[i])
    }
    return format
}

const resourcesToDisplay = ref<Resource[]>([])



const remainingTime = ref<string>('00:00:00')
const remainingSeconds = ref<number>(0)

const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600).toString().padStart(2, '0');
    const minutes = Math.floor((seconds % 3600) / 60).toString().padStart(2, '0');
    const secondsRemaining = (seconds % 60).toString().padStart(2, '0');
    return `${hours}:${minutes}:${secondsRemaining}`;
};

const calculateRemainingTimeAndPoints = () => {
    const now = Date.now() / 1000
    let secondsLeft = warehouse.value.stateUntil - now
    remainingSeconds.value = secondsLeft

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

const showUpdateDetails = ref<boolean>(false)

</script>

<template>
    <section class=" pb-20">
        <Transition name="slide-down">
            <section v-if="showUpdateDetails"
                class="z-50 w-full h-screen bg-white fixed top-0 left-0">
                <section class=" overflow-y-scroll h-full relative p-4">
                    <i @click="showUpdateDetails = false"
                        class=" pi fixed pi-times bg-dark text-white aspect-square p-1 rounded-full top-4 right-4"></i>
                    <p class="text-center text-dark font-bold my-4">
                        {{ t(`screens.chibi.city.buildings.level`) }} {{ warehouse.level }}
                        <i class="pi pi-arrow-right"></i>
                        {{ t(`screens.chibi.city.buildings.level`) }} {{ warehouse.level+1 }}
                    </p>

                    <div class="flex gap-2 flex-col">
                        <div class=" resource-mini font-bold flex rounded-2xl bg-half_dark p-4 gap-2 items-center">
                            <img class=" w-16 object-contain aspect-square"
                                src="../../../assets/images/chibi/city/box.png" alt="">
                            <p class="w-full">{{ t(`screens.chibi.city.buildings.warehouse.capacity`) }}</p>
                            <div class=" flex flex-col">
                                <p class=" text-primary flex gap-1 items-center">{{ warehouse.nextLevel?.capacity }}<i class="pi pi-arrow-up"></i></p>
                                <p>{{ warehouse.currentLevel?.capacity }}</p>
                            </div>
                        </div>

                    </div>
                </section>
            </section>
        </Transition>

        <Transition name="slide-down">
            <section v-if="resourcesToDisplay.length > 0"
                class="z-50 w-full h-screen bg-white fixed top-0 left-0">
                <section class=" overflow-y-scroll h-full relative p-4">
                    <i @click="resourcesToDisplay = []"
                        class=" pi fixed pi-times bg-dark text-white aspect-square p-1 rounded-full top-4 right-4"></i>
                    <p class="text-center text-dark font-bold my-4">{{
                    t(`screens.chibi.city.buildings.${warehouse.type}.name`) }}</p>
                    <div class="flex gap-2 flex-col">
                        
                        <div v-for="(res, i) of resourcesToDisplay" :key="i"
                            class=" resource-mini font-bold flex rounded-2xl bg-half_dark p-4 gap-2 items-center">
                            <img class=" w-16 object-contain aspect-square"
                                :src="`${baseURL}/static/images/resources/${res.name}.png`" alt="">
                            <p class="w-full">{{ t(`resources.meterials.${res.name}`) }}</p>
                            <p>{{ res.amount }}</p>
                        </div>

                    </div>
                </section>
            </section>
        </Transition>

        <Transition name="delay">
            <section v-show="showUpgradeModal" @click.self="showUpgradeModal = false" class=" wrapper">
                <Transition name="scale">
                    <section v-if="showUpgradeModal" class=" modal">

                        <p class="font-bold text-center">
                            {{ sprintf(t('screens.chibi.city.buildings.buildApprove'),
                t(notBought ? 'screens.chibi.city.buildings.actions.build' :
                    'screens.chibi.city.buildings.actions.upgrade'),
                t(`screens.chibi.city.buildings.${warehouse.type}.name`))
                            }}
                        </p>

                        <div class="flex gap-2 w-full">

                            <button @click="showUpgradeModal = false" class=" py-2 text-primary bg-white">
                                {{ t('screens.chibi.city.buildings.buildApproveCancel') }}
                            </button>

                            <button @click="upgrade" class=" py-2">
                                <p v-if="loading"><i class="pi pi-spinner pi-spin"></i></p>
                                <p v-else>
                                    {{
                t(notBought ? 'screens.chibi.city.buildings.build' :
                    'screens.chibi.city.buildings.upgrade')
            }}
                                </p>
                            </button>

                        </div>

                    </section>
                </Transition>
            </section>
        </Transition>

        <p v-if="!loaded" class="text-center font-dark mt-4"><i class=" pi pi-spinner pi-spin" style="font-size: 1.5rem; color: var(--dark)"></i></p>

        <section class=" p-4 flex flex-col gap-4" v-else>
            <Balances />

            <img class=" mx-auto max-w-80"
                :src="`${baseURL}/static/images/buildings/${warehouse.type}${warehouse.level > 0 ? warehouse.level : 1}.png`"
                alt="">
            <h1>{{ t(`screens.chibi.city.buildings.${warehouse.type}.name`) }}</h1>
            <div class="card">
                <span>
                    <p class=" font-bold">{{ t(`screens.chibi.city.buildings.${warehouse.type}.header`) }}</p>
                    <p class=" text-dark">{{ t(`screens.chibi.city.buildings.${warehouse.type}.fullDescription`) }}</p>
                </span>
                <div v-if="warehouse.currentLevel" class="flex gap-4">
                    <div class="shield">
                        <div class="text-center">
                            <p class="level">{{ warehouse.level }}</p>
                            <p class="label">{{ t('screens.chibi.city.buildings.level') }}</p>
                        </div>
                        <!-- <i class=" text-dark pi pi-chevron-right"></i> -->
                    </div>
                    <div class="shield">
                        <div class="text-center">
                            <p class=" font-bold">{{ capacityUsed }}/{{ warehouse.currentLevel.capacity }}</p>
                            <p class="label">{{ t('screens.chibi.city.buildings.warehouse.filled') }}</p>
                        </div>
                    </div>
                </div>
                <div v-else-if="warehouse.nextLevel" class="flex gap-4">
                    <div class="shield">
                        <div class="text-center">
                            <p class="level">{{ warehouse.level+1 }}</p>
                            <p class="label">{{ t('screens.chibi.city.buildings.level') }}</p>
                        </div>
                        <!-- <i class=" text-dark pi pi-chevron-right"></i> -->
                    </div>
                    <div class="shield">
                        <div class="text-center">
                            <p class=" font-bold">{{ capacityUsed }}/{{ warehouse.nextLevel.capacity }}</p>
                            <p class="label">{{ t('screens.chibi.city.buildings.warehouse.filled') }}</p>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="!notBought && !maxLevel" class="card">
                <span>
                    <p class=" font-bold">{{ t(`screens.chibi.city.buildings.upgrading`) }}</p>
                    <p class=" text-dark">{{ t(`screens.chibi.city.buildings.${warehouse.type}.upgrade`) }}</p>
                </span>
                <div class="flex gap-4">
                    <div @click="showUpdateDetails = true" class="shield">
                        <div class="text-center">
                            <div class="flex items-center gap-2">
                                <p class="level">{{ warehouse.level }}</p>
                                <i v-if="!maxLevel" class=" pi pi-arrow-right" style="color:var(--blue)"></i>
                                <p v-if="!maxLevel" class="level">{{ warehouse.level + 1 }}</p>
                            </div>
                            <p class="label">{{ t('screens.chibi.city.buildings.level') }}</p>
                        </div>
                        <i class=" text-dark pi pi-chevron-right"></i>
                    </div>
                    <div class="shield">
                        <div class="text-center">
                            <div class="flex gap-1">
                                <span class="flex items-center gap-1" v-for="(cost, i) of warehouse.nextLevel?.cost" :key="i">
                                    <img class="h-6" :src="`${baseURL}/static/images/resources/${cost.currency}.png`"
                                        alt="">
                                    <p>{{ -cost.amount }}</p>
                                </span>
                            </div>
                            <p class="label">{{ t('screens.chibi.city.buildings.cost') }}</p>
                        </div>
                    </div>
                </div>

                <div v-if="!maxLevel && warehouse.nextLevel?.requirements?.length" class="requirements">
                    <p class=" font-bold">{{ t('screens.chibi.city.buildings.requirements') }}</p>
                    <BuildingCardTiny class="requirement" v-for="(building, i) of warehouse.nextLevel?.requirements"
                        :key="i" :building="building" />
                </div>
                <button :disabled="!balanceEnoughForUpgrade" @click="showUpgradeModal = true"
                    class="flex justify-center">
                    <p v-if="remainingSeconds > 0" class="flex gap-2 items-center h-6">{{ warehouse.level == 0 ? t('screens.chibi.city.buildings.building') : t('screens.chibi.city.buildings.upgrading') }} <p class=" font-mono">{{ remainingTime }}</p> </p>

                    <span v-else-if="balanceEnoughForUpgrade" class="flex gap-2">
                        {{ t('screens.chibi.city.buildings.upgrade') }}
                        <div class="flex gap-1">
                            <span class="flex items-center gap-1" v-for="(cost, i) of warehouse.nextLevel?.cost" :key="i">
                                <img class="h-6" :src="`${baseURL}/static/images/resources/${cost.currency}.png`" alt="">
                                <p>{{ -cost.amount }}</p>
                            </span>
                        </div>
                    </span>
                    <p v-else class="h-6">{{ t('screens.chibi.city.buildings.notAvailible') }}</p>
                </button>
            </div>
            
            <div v-if="notBought && warehouse.moreLevel" class="card opacity-50">
                <span>
                    <p class=" font-bold">{{ t(`screens.chibi.city.buildings.upgrading`) }}</p>
                    <p class=" text-dark">{{ t(`screens.chibi.city.buildings.${warehouse.type}.upgrade`) }}</p>
                </span>
                <div class="flex gap-4">
                    <div class="shield">
                        <div class="text-center">
                            <div class="flex items-center gap-2">
                                <p class="level">{{ warehouse.level+1 }}</p>
                                <i class=" pi pi-arrow-right" style="color:var(--blue)"></i>
                                <p class="level">{{ warehouse.level + 2 }}</p>
                            </div>
                            <p class="label">{{ t('screens.chibi.city.buildings.level') }}</p>
                        </div>
                        <i class=" text-dark pi pi-chevron-right"></i>
                    </div>
                    <div class="shield">
                        <div class="text-center">
                            <div class="flex gap-1">
                                <span class="flex items-center gap-1" v-for="(cost, i) of warehouse.moreLevel?.cost" :key="i">
                                    <img class="h-6" :src="`${baseURL}/static/images/resources/${cost.currency}.png`"
                                        alt="">
                                    <p>{{ -cost.amount }}</p>
                                </span>
                            </div>
                            <p class="label">{{ t('screens.chibi.city.buildings.cost') }}</p>
                        </div>
                    </div>
                </div>

                <button disabled
                    class="flex justify-center">
                    <span class="flex gap-2">
                        {{ t('screens.chibi.city.buildings.upgrade') }}
                        <div class="flex gap-1">
                            <span class="flex items-center gap-1" v-for="(cost, i) of warehouse.moreLevel?.cost" :key="i">
                                <img class="h-6" :src="`${baseURL}/static/images/resources/${cost.currency}.png`" alt="">
                                <p>{{ -cost.amount }}</p>
                            </span>
                        </div>
                    </span>
                </button>
            </div>

            <p v-else-if="maxLevel" class=" text-dark text-center font-bold">{{
                t('screens.chibi.city.buildings.maxLevel') }}</p>

            <button v-if="notBought" :disabled="!balanceEnoughForUpgrade" @click="showUpgradeModal = true"
                class="flex justify-center">
                <span v-if="balanceEnoughForUpgrade" class="flex gap-2">
                    {{ t('screens.chibi.city.buildings.build') }}

                    <div class="flex gap-1">
                        <span class="flex items-center gap-1" v-for="(cost, i) of warehouse.nextLevel?.cost" :key="i">
                            <img class="h-6" :src="`${baseURL}/static/images/resources/${cost.currency}.png`" alt="">
                            <p>{{ -cost.amount }}</p>
                        </span>
                    </div>
                </span>
                <p v-else class="h-6">{{ t('screens.chibi.city.buildings.notAvailible') }}</p>
            </button>

            <section v-if="!notBought" class="resources">
                <div class="resource-card" v-for="(v, k) of resourceDictionary" :key="k"
                    @click="resourcesToDisplay = v.resources">
                    <img class=" mx-auto max-w-80" :src="`${baseURL}/static/images/resources/${k}.png`" alt="">
                    <p>{{ t(`resources.${k}`) }}</p>
                    <p class=" resource-amount-label">{{ v.totalAmount }}</p>
                </div>
            </section>
        </section>

    </section>
</template>


<style scoped>
.slide-down-enter-active,
.slide-down-leave-active {
    transition: transform 0.5s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
    transform: translateY(100%);
}

.wrapper {
    @apply fixed z-50 w-full h-screen top-0 left-0 flex items-center justify-center p-4 drop-shadow-lg;
}

.modal {
    @apply max-w-80 gap-4 w-full rounded-2xl bg-white p-4 flex flex-col justify-center items-center shadow-lg;
}

.card {
    @apply p-4 rounded-2xl bg-half_dark flex flex-col gap-4;
}

.shield {
    @apply bg-white rounded-2xl flex items-center w-full p-2;
}

.shield>div {
    @apply w-full flex flex-col items-center;
}

.level {
    @apply px-2 aspect-square flex items-center justify-center bg-custom_blue text-white rounded-md;
}

.requirements{
    @apply flex flex-col gap-2
}

.resources {
    @apply grid grid-cols-2 gap-4;
}

.resource-card {
    @apply bg-half_dark relative rounded-2xl p-4 flex flex-col items-center;
}

.resource-card>img {
    @apply w-full p-4
}

.resource-amount-label {
    @apply absolute top-4 right-4 bg-dark text-white p-1 rounded-full;
}
</style>