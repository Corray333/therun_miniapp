<script lang="ts" setup>

import { useI18n } from 'vue-i18n'

const { t } = useI18n()

</script>

<template>
    <section class=" pb-20">
        <section class="p-4 flex flex-col items-center gap-4">
            <h1>{{ t('screens.chibi.header') }}</h1>
            <router-link to="/chibi/cases">
                <div class="banner p-4 rounded-2xl flex gap-4 items-center bg-cover">
                    <img src="../assets/images/chibi/slider/box-3.png" alt="" class=" h-20">
                    <div class="flex flex-col w-full text-white">
                        <h2>{{ t('screens.chibi.banner.header') }}</h2>
                        <p class=" text-sm min-h-10">{{ t('screens.chibi.banner.description') }}</p>
                    </div>
                    <i class="pi pi-chevron-right text-white" style="font-size:1.25rem"></i>
                </div>
            </router-link>
            <div class="city card">
                <div class=" flex">
                    <div class="flex flex-col gap-2 w-full">
                        <h2>{{ t('screens.chibi.city.header') }}</h2>
                        <p>{{ t('screens.chibi.city.description') }}</p>
                        <button>{{ t('screens.chibi.city.btn') }}</button>
                    </div>
                    <img src="../assets/images/chibi/city.png" alt="" class="block w-2/5 object-contain">
                </div>
            </div>
            <div class="flex gap-4 w-full">
                <div class="card flex flex-col w-full">
                    <div class="flex justify-between w-full">
                        <img src="../assets/images/chibi/battles.png" class=" w-20" alt="">
                        <i class=" pi pi-arrow-right bg-calm h-fit p-2 rounded-full text-white"></i>
                    </div>
                    <h2 class=" mt-4">{{ t('screens.chibi.fights.header') }}</h2>
                    <p>{{ t('screens.chibi.fights.description') }}</p>
                </div>
                <div class="card flex flex-col w-full">
                    <div class="flex justify-between w-full">
                        <img src="../assets/images/chibi/team.png" class=" w-20" alt="">
                        <i class=" pi pi-arrow-right bg-calm h-fit p-2 rounded-full text-white"></i>
                    </div>
                    <h2 class=" mt-4">{{ t('screens.chibi.team.header') }}</h2>
                    <p>{{ t('screens.chibi.team.description') }}</p>
                </div>
            </div>
        </section>
    </section>
</template>


<style scoped>
.card {
    @apply rounded-2xl bg-half_dark p-4;

    p {
        @apply text-sm text-dark;
        line-height: 1rem;
    }

}

button {
    @apply w-fit px-6 py-2 bg-calm;
}

.banner {
    background-image: url(../assets/images/chibi/case-banner-bg.png);
}
</style>