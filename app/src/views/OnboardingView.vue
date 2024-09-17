<script lang="ts" setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {useRouter} from 'vue-router'

const { t } = useI18n()
const router = useRouter()

const currentSlide = ref<number>(0)


</script>

<template>
    <section>
        <div class="slider">
            <div class="slide-track" :style="`transform:translateX(${-100*currentSlide}vw)`">
                <div class="slide">
                    <div class="pic-container">
                        <div class="pic" id="onboarding-pic-1"></div>
                    </div>
                    <div class="info">
                        <h1>{{ t('screens.onboarding.slide1.title') }}</h1>
                        <p>{{ t('screens.onboarding.slide1.description') }}</p>
                        <button @click="currentSlide = 1">{{ t('screens.onboarding.nextButton') }}</button>
                    </div>
                </div>
                <div class="slide">
                    <div class="pic-container">
                        <div class="pic" id="onboarding-pic-3"></div>
                    </div>
                    <div class="info">
                        <h1>{{ t('screens.onboarding.slide3.title') }}</h1>
                        <p>{{ t('screens.onboarding.slide3.description') }}</p>
                        <button @click="currentSlide = 2">{{ t('screens.onboarding.nextButton') }}</button>
                    </div>
                </div>
                <div class="slide">
                    <div class="pic-container">
                        <div class="pic bg-bottom" id="onboarding-pic-4"></div>
                    </div>
                    <div class="info">
                        <h1>{{ t('screens.onboarding.slide4.title') }}</h1>
                        <p>{{ t('screens.onboarding.slide4.description') }}</p>
                        <button @click="currentSlide = 3">{{ t('screens.onboarding.nextButton') }}</button>
                    </div>
                </div>
                <div class="slide">
                    <div class="pic-container">
                        <div class="pic" id="onboarding-pic-5"></div>
                    </div>
                    <div class="info">
                        <h1>{{ t('screens.onboarding.slide5.title') }}</h1>
                        <p>{{ t('screens.onboarding.slide5.description') }}</p>
                        <button @click="router.push('/farm')">{{ t('screens.onboarding.startButton') }}</button>
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>


<style>
.slider {
    overflow: hidden;
}

.slide-track {
    display: flex;
    transition: transform 0.5s; 
}

.slide {
    width: 100vw;
    height: 100vh;
    display: flex;
    flex-direction: column;
    flex: 0 0 auto;
}

.pic-container {
    height: 100%;
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
}

.pic {
    height: 100%;
    width: 100%;
    max-width: 80%;
    background-position: center;
    background-repeat: no-repeat;
    background-size: contain;
}

.info {
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    text-align: center;
    border-radius: 1rem 1rem 0 0;
    box-shadow: 0 -0.25rem 1rem 0 rgba(0, 0, 0, 0.1);
}

#onboarding-pic-1 {
    background-image: url('../assets/images/onboarding/slide1.png');
}

#onboarding-pic-2 {
    background-image: url('../assets/images/onboarding/slide2.png');
}

#onboarding-pic-3 {
    background-image: url('../assets/images/onboarding/slide3.png');
}

#onboarding-pic-4 {
    background-image: url('../assets/images/onboarding/slide4.png');
}

#onboarding-pic-5 {
    background-image: url('../assets/images/onboarding/slide5.png');
}


.slide-enter-active,
.slide-leave-active {
    transition: transform 0.5s;
}

.slide-enter,
.slide-leave-to {
    transform: translateX(100%);
}
</style>