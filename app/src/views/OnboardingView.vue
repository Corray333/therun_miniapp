<script lang="ts" setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const currentSlide = ref<number>(1)


</script>

<template>
    <section>
        <Transition name="slide">
            <div v-if="currentSlide == 1" class="slide">
                <div class="pic-container">
                    <div class="pic" id="onboarding-pic-1"></div>
                </div>
                <div class="info">
                    <h1>{{ t('screens.onboarding.slide1.title') }}</h1>
                    <p>{{ t('screens.onboarding.slide1.description') }}</p>
                    <button @click="currentSlide+=2">Next</button>
                </div>
            </div>
            <div v-if="currentSlide == 3" class="slide">
                <div class="pic-container">
                    <div class="pic" id="onboarding-pic-3"></div>
                </div>
                <div class="info">
                    <h1>{{ t('screens.onboarding.slide3.title') }}</h1>
                    <p>{{ t('screens.onboarding.slide3.description') }}</p>
                    <button @click="currentSlide++">Next</button>
                </div>
            </div>
        </Transition>
    </section>
</template>


<style>
.slide {
    width: 100vw;
    height: 100vh;
    display: flex;
    flex-direction: column;
}
.pic-container{
    height: 100%;
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
}
.pic{
    height: 100%;
    width: 100%;
    max-width: 80%;
    background-position: center;
    background-repeat: no-repeat;
    background-size: contain;
}
.info{
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    text-align: center;
    border-radius: 1rem 1rem 0 0;
    box-shadow: 0 -0.25rem 1rem 0 rgba(0, 0, 0, 0.1);
}

#onboarding-pic-1{
    background-image: url('../assets/images/onboarding/slide1.png');
}
#onboarding-pic-2{
    background-image: url('../assets/images/onboarding/slide2.png');
}
#onboarding-pic-3{
    background-image: url('../assets/images/onboarding/slide3.png');
}
#onboarding-pic-4{
    background-image: url('../assets/images/onboarding/slide4.png');
}
#onboarding-pic-5{
    background-image: url('../assets/images/onboarding/slide5.png');
}


.slide-left-enter-active,
.slide-right-leave-active {
    transition: all 0.3s ease;
}

.slide-left-enter,
.slide-left-leave-to,
.slide-right-enter,
.slide-right-leave-to {
    transform: translateX(100%);
}
</style>