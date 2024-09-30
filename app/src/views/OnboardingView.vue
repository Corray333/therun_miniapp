<script lang="ts" setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import Slider from '@/components/Slider.vue'

const { t } = useI18n()
const router = useRouter()

const currentSlide = ref<number>(0)

const sliderVal = ref<number>(0)

watch(sliderVal, (val) => {
    if (val === 85) {
        router.push('/')
    }
})

</script>

<template>
    <section>
        <div class="slider">
            <div class="slide-track" :style="`transform:translateX(${-100 * currentSlide}vw)`">
                <div class="slide">
                    <div class="pic-container flex-col p-8 gap-4 p-4">
                        <div class="pic" id="onboarding-pic-1"></div>
                        <div class="flex gap-4">
                            <a href="https://play.google.com/store/apps/details?id=com.therun.app" target="_blank"><img
                                    src="../assets/images/onboarding/google-play-btn.png" alt=""></a>
                            <a href="https://apps.apple.com/us/app/therun/id1634366310" target="_blank"><img
                                    src="../assets/images/onboarding/appstore-btn.png" alt=""></a>
                        </div>
                    </div>
                    <div class="info">
                        <h1>{{ t('screens.onboarding.slide1.title') }}</h1>
                        <p>{{ t('screens.onboarding.slide1.description') }}</p>
                        <button @click="currentSlide = 1">{{ t('screens.onboarding.nextButton') }}</button>
                    </div>
                </div>
                <div class="slide">
                    <div class="pic-container p-4">
                        <div class="pic" id="onboarding-pic-2"></div>
                    </div>
                    <div class="info">
                        <h1>{{ t('screens.onboarding.slide2.title') }}</h1>
                        <p>{{ t('screens.onboarding.slide2.description') }}</p>
                        <button @click="currentSlide = 2">{{ t('screens.onboarding.nextButton') }}</button>
                    </div>
                </div>
                <div class="slide">
                    <div class="pic-container p-4">
                        <div class="pic" id="onboarding-pic-3"></div>
                    </div>
                    <div class="info">
                        <h1>{{ t('screens.onboarding.slide3.title') }}</h1>
                        <p>{{ t('screens.onboarding.slide3.description') }}</p>
                        <button @click="currentSlide = 3">{{ t('screens.onboarding.nextButton') }}</button>
                    </div>
                </div>
                <div class="slide">
                    <div class="pic-container pt-4 px-4">
                        <div class="pic" id="onboarding-pic-4"></div>
                    </div>
                    <div class="info">
                        <h1>{{ t('screens.onboarding.slide4.title') }}</h1>
                        <p>{{ t('screens.onboarding.slide4.description') }}</p>
                        <button @click="currentSlide = 4">{{ t('screens.onboarding.nextButton') }}</button>
                    </div>
                </div>

                <div class="slide">
                    <div class="pic-container pt-4 px-4">
                        <div class="pic" id="onboarding-pic-5"></div>
                    </div>
                    <div class="info">
                        <h1>{{ t('screens.onboarding.slide5.title') }}</h1>
                        <p>{{ t('screens.onboarding.slide5.description') }}</p>
                        <button @click="currentSlide = 5">{{ t('screens.onboarding.nextButton') }}</button>
                    </div>
                </div>

                <div class="slide">
                    <div class="pic-container p-4">
                        <div class="pic" id="onboarding-pic-6"></div>
                    </div>
                    <div class="info">
                        <h1 v-html="t('screens.onboarding.slide6.title')"></h1>
                        <p>{{ t('screens.onboarding.slide6.description') }}</p>
                        <Slider v-model="sliderVal" :min="0" :max="100" :step="1" id="slider-finish" />
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
    background-position: center;
    background-repeat: no-repeat;
    background-size: contain;
}

.info {
    padding: 2rem 1rem;
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

#onboarding-pic-6 {
    background-image: url('../assets/images/onboarding/slide6.png');
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