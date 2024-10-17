<template>
    <div class="slider-container" @mousedown="startDrag" @mousemove="drag" @mouseup="endDrag" @mouseleave="endDrag"
        @touchstart="startDrag" @touchmove="drag" @touchend="endDrag">
        <div class="slider-track">
            <div class="slider-fill" :style="{ width: `calc(${thumbPosition}% + 3.5rem / 2)` }">
            </div>
            <div class="slider-thumb" :style="{ left: `${thumbPosition}%` }" @mousedown.prevent @touchstart.prevent>
                <div class="bg-white w-full h-full rounded-full flex justify-center items-center">
                    <logo />
                </div>
            </div>
            <p class="font-bold text-white">{{ t('screens.onboarding.startButton') }}</p>
        </div>
    </div>
</template>

<script>
import logo from "@/components/icons/logo-icon.vue";
import { useI18n } from "vue-i18n";

export default {
    name: "CustomSlider",
    components: {
        logo
    },
    props: {
        modelValue: {
            type: Number,
            required: true,
            validator(value) {
                return value >= 0 && value <= 100;
            }
        },
        min: {
            type: Number,
            default: 0
        },
        max: {
            type: Number,
            default: 100
        },
        step: {
            type: Number,
            default: 1
        },
    },
    data() {
        return {
            dragging: false,
            thumbPosition: this.calculateThumbPosition(this.modelValue),
            t: useI18n().t
        };
    },
    methods: {
        calculateThumbPosition(value) {
            // Вычисляем процентное значение позиции ползунка
            const percentage = ((value - this.min) / (this.max - this.min)) * 100;
            return Math.min(Math.max(percentage, 0), 100); // Ограничиваем до диапазона [0, 100]
        },
        startDrag(event) {
            this.$emit('drag-started');
            this.dragging = true;
            this.updateValue(event);
        },
        drag(event) {
            if (this.dragging) {
                this.updateValue(event);
            }
        },
        endDrag() {
            this.$emit('drag-ended');
            this.dragging = false;
        },
        updateValue(event) {
            const clientX = event.touches ? event.touches[0].clientX : event.clientX;
            const rect = this.$el.getBoundingClientRect();
            const offsetX = clientX - rect.left;
            // Ограничиваем позицию в пределах слайдера
            const clampedOffsetX = Math.max(0, Math.min(offsetX, rect.width));
            const newValue = Math.min(((clampedOffsetX / rect.width) * (this.max - this.min)) + this.min, 85);
            const steppedValue = Math.round(newValue / this.step) * this.step;
            const clampedValue = Math.min(Math.max(steppedValue, this.min), this.max);
            this.$emit("update:modelValue", clampedValue); // Обновляем значение через v-model
            this.thumbPosition = this.calculateThumbPosition(clampedValue); // Обновляем позицию ползунка
        }
    },
    watch: {
        modelValue(newValue) {
            this.thumbPosition = this.calculateThumbPosition(newValue); // Синхронизация слайдера с внешним значением
        }
    }
};
</script>

<style scoped>
.slider-container {
    position: relative;
    width: 100%;
    height: 3.5rem;
    cursor: pointer;
    user-select: none;
}

.slider-track {
    border-radius: 99px;
    position: relative;
    width: 100%;
    height: 100%;
    background-color: var(--half-gray);

    display: flex;
    align-items: center;
    justify-content: center;
}

.slider-fill {
    @apply rounded-l-full;
    position: absolute;
    height: 100%;
    top: 0;
    left: 0;
    background-color: var(--primary);
}

.slider-thumb {
    position: absolute;
    padding: 0.25rem;
    width: 3.5rem;
    height: 3.5rem;
    border-radius: 99px;
    cursor: pointer;
    transition: background-color 0.3s;
    background-color: var(--primary);
    touch-action: none;
}
</style>
