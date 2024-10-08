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
    "lastStateChange": 0,
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
    }
})

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
        const { data } = await axios.patch(`${import.meta.env.VITE_API_URL}/city/warehouse/upgrade`, {}, {
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

</script>

<template>
    <section class=" pb-20">
        <Transition name="delay">
            <section v-show="showUpgradeModal" @click.self="showUpgradeModal = false"
                class=" wrapper fixed z-50 w-full h-screen top-0 left-0 flex items-center justify-center p-4 drop-shadow-lg">
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

        <section class=" p-4 flex flex-col gap-4">
            <Balances />

            <img class=" mx-auto max-w-80"
                :src="`${baseURL}/static/images/buildings/${warehouse.type}${warehouse.level > 0 ? warehouse.level : 1}.png`"
                alt="">
            <h1>{{ t(`screens.chibi.city.buildings.${warehouse.type}.name`) }}</h1>
            <div class="card">
                <p>
                <p class=" font-bold">{{ t(`screens.chibi.city.buildings.${warehouse.type}.header`) }}</p>
                <p class=" text-dark">{{ t(`screens.chibi.city.buildings.${warehouse.type}.fullDescription`) }}</p>
                </p>
                <div v-if="warehouse.currentLevel" class="flex gap-4">
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

            <div v-if="!notBought && !maxLevel" class="card">
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
                                <p v-if="!maxLevel" class="level">{{ warehouse.level + 1 }}</p>
                            </div>
                            <p class="label">{{ t('screens.chibi.city.buildings.level') }}</p>
                        </div>
                        <i class=" text-dark pi pi-chevron-right"></i>
                    </div>
                    <div class="shield">
                        <div class="text-center">
                            <div class="flex gap-1">
                                <p class="flex items-center gap-1" v-for="(cost, i) of warehouse.nextLevel?.cost">
                                    <img class="h-6" :src="`${baseURL}/static/images/resources/${cost.currency}.png`"
                                        alt="">
                                <p>{{ -cost.amount }}</p>
                                </p>
                            </div>
                            <p class="label">{{ t('screens.chibi.city.buildings.cost') }}</p>
                        </div>
                    </div>
                </div>

                <div v-if="!maxLevel" class="requirements">
                    <p class=" font-bold">{{ t('screens.chibi.city.buildings.requirements') }}</p>
                    <BuildingCardTiny class="requirement" v-for="(building, i) of warehouse.nextLevel?.requirements"
                        :key="i" :building="building" />
                </div>

                <button :disabled="!balanceEnoughForUpgrade" @click="showUpgradeModal = true"
                    class="flex justify-center">
                    <p v-if="balanceEnoughForUpgrade" class="flex gap-2">
                        {{ t('screens.chibi.city.buildings.upgrade') }}

                    <div class="flex gap-1">
                        <p class="flex items-center gap-1" v-for="(cost, i) of warehouse.nextLevel?.cost">
                            <img class="h-6" :src="`${baseURL}/static/images/resources/${cost.currency}.png`" alt="">
                        <p>{{ -cost.amount }}</p>
                        </p>
                    </div>
                    </p>
                    <p v-else>{{ t('screens.chibi.city.buildings.notAvailible') }}</p>
                </button>
            </div>
            <p v-else-if="maxLevel" class=" text-dark text-center font-bold">{{
                t('screens.chibi.city.buildings.maxLevel') }}</p>

            <button v-if="notBought" :disabled="!balanceEnoughForUpgrade" @click="showUpgradeModal = true"
                class="flex justify-center">
                <p v-if="balanceEnoughForUpgrade" class="flex gap-2">
                    {{ t('screens.chibi.city.buildings.build') }}

                <div class="flex gap-1">
                    <p class="flex items-center gap-1" v-for="(cost, i) of warehouse.nextLevel?.cost">
                        <img class="h-6" :src="`${baseURL}/static/images/resources/${cost.currency}.png`" alt="">
                    <p>{{ -cost.amount }}</p>
                    </p>
                </div>
                </p>
                <p v-else>{{ t('screens.chibi.city.buildings.notAvailible') }}</p>
            </button>

            <section v-if="!notBought" class="resources">
                <div class="resource-card" v-for="(v,k) of resourceDictionary" :key="k">
                    <img class=" mx-auto max-w-80" :src="`${baseURL}/static/images/resources/${k}.png`" alt="">
                    <p>{{ t(`resources.${k}`) }}</p>
                    <p class=" resource-amount-label">{{ v.totalAmount }}</p>
                </div>
            </section>
        </section>

    </section>
</template>


<style scoped>
.modal {
    @apply max-w-80 gap-4 w-full rounded-2xl bg-white p-4 flex flex-col justify-center items-center shadow-lg
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

.resources{
    @apply grid grid-cols-2 gap-4;
}

.resource-card {
    @apply bg-half_dark relative rounded-2xl p-4 flex flex-col items-center;
}

.resource-card>img{
    @apply w-full p-4
}

.resource-amount-label {
    @apply absolute top-4 right-4 bg-dark text-white p-1 rounded-full;
}

</style>