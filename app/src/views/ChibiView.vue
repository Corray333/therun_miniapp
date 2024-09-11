<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import Navbar from '@/components/Navbar.vue'

const { t } = useI18n()

</script>


<template>
    <section class=" pb-24">
        <section class=" flex flex-col p-5 h-full w-full">
            <div class="game text-center flex flex-col gap-4">
                <h1 class=" font-bold text-2xl">{{ t("screens.chibi.game.header") }}</h1>
                <p>{{ t("screens.chibi.game.subheader") }}</p>
                <img src="../assets/images/chibi/phones.png" alt="">
                <p>{{ t("screens.chibi.game.description") }}</p>
            </div>
        </section>
        <div class="slider">
            <div class="slide-track">
                <img src="../assets/images/chibi/slider/box.png" alt="" class="slide">
                <img src="../assets/images/chibi/slider/box-1.png" alt="" class="slide">
                <img src="../assets/images/chibi/slider/box-2.png" alt="" class="slide">
                <img src="../assets/images/chibi/slider/box-3.png" alt="" class="slide">
                <img src="../assets/images/chibi/slider/box-4.png" alt="" class="slide">
                <!-- Дублируем изображения для бесконечного эффекта -->
                <img src="../assets/images/chibi/slider/box.png" alt="" class="slide">
                <img src="../assets/images/chibi/slider/box-1.png" alt="" class="slide">
                <img src="../assets/images/chibi/slider/box-2.png" alt="" class="slide">
                <img src="../assets/images/chibi/slider/box-3.png" alt="" class="slide">
                <img src="../assets/images/chibi/slider/box-4.png" alt="" class="slide">
            </div>
        </div>
        <Navbar class=" fixed bottom-0" />
    </section>
</template>


<style scoped>
.slider {
    overflow: hidden;
    width: 100%;
}

.slide-track {
    display: flex;
    gap: 1.4rem;
    animation: scroll 20s linear infinite;
    width: calc(7rem * 10 + 1.4rem * 9);
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
        transform: translateX(calc((7rem * 10 + 1.4rem * 9) / -2));
    }
}
</style>