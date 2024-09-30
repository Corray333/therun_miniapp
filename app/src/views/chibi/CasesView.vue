<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import Navbar from '@/components/Navbar.vue'
import Balances from '@/components/Balances.vue';
import key from '@/components/icons/key-icon.vue'

const { t } = useI18n()

</script>


<template>
    <section class=" pb-20">
        <section class=" flex flex-col p-4 h-full w-full">
            <Balances />
        </section>
        <section class=" flex flex-col gap-6 pb-12 items-center p-4">
            <span class="flex flex-col gap-1 text-center">
                <h1>{{ t('screens.chibi.case.header') }}</h1>
                <p>{{ t('screens.chibi.case.description') }}</p>
            </span>
            <img id="case" src="../../assets/images/chibi/slider/box-3.png" alt="">
            <div class="flex gap-2">
                <span class="rounded-full p-4 sm-shadow">
                    <key color="var(--dark)"/>
                </span>
                <span class="rounded-full p-4 sm-shadow">
                    <key color="var(--dark)"/>
                </span>
                <span class="rounded-full p-4 sm-shadow">
                    <key color="var(--dark)"/>
                </span>
            </div>
            <button class=" bg-dark" disabled>{{ t('comingSoon') }}</button>
        </section>
        <section class=" info p-4 rounded-t-2xl bg-white ">
            <div class="game text-center flex flex-col gap-6">
                <span class=" flex flex-col gap-2">
                    <h1 class=" font-bold text-2xl">{{ t("screens.chibi.game.header") }}</h1>
                    <p>{{ t("screens.chibi.game.subheader") }}</p>
                </span>
                <img src="../../assets/images/chibi/phones.png" alt="">
                <div class="flex w-full gap-4">
                    <a href="https://play.google.com/store/apps/details?id=com.therun.app" target="_blank"><img src="../../assets/images/chibi/google-play-btn.png" alt=""></a>
                    <a href="https://apps.apple.com/us/app/therun/id1634366310" target="_blank"><img src="../../assets/images/chibi/appstore-btn.png" alt=""></a>
                </div>
                <p class=" p-4 bg-half_dark rounded-2xl">{{ t("screens.chibi.game.description") }}</p>
            </div>
        </section>
        <div class="slider">
            <div class="slide-track">
                <img src="../../assets/images/chibi/slider/box.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-1.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-2.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-3.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-4.png" alt="" class="slide">
                <!-- Дублируем изображения для бесконечного эффекта -->
                <img src="../../assets/images/chibi/slider/box.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-1.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-2.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-3.png" alt="" class="slide">
                <img src="../../assets/images/chibi/slider/box-4.png" alt="" class="slide">
            </div>
        </div>
    </section>
</template>


<style scoped>



.info{
    box-shadow: 0 -1rem 1rem 0 rgba(0, 0, 0, 0.1);
}

#case{
    animation: case-breath 2s infinite;
    max-width: 14rem;
}

.slider {
    overflow: hidden;
    width: 100%;
}

.slide-track {
    display: flex;
    gap: 1.4rem;
    animation: scroll 20s linear infinite;
    width: calc(7rem * 10 + 1.4rem * 10);
}

.slide {
    height: 7rem;
    width: 7rem;
    flex-shrink: 0;
}

@keyframes scroll {
    0% {
        transform: translateX(0);
    }

    100% {
        transform: translateX(calc(calc(7rem * 10 + 1.4rem * 10) / -2));
    }
}

@keyframes case-breath {
    0% {
        transform: scale(1);
    }

    50% {
        transform: scale(0.9);
    }

    100% {
        transform: scale(1);
    }
}
</style>