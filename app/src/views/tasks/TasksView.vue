<script lang="ts" setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

</script>

<template>
    <section class=" flex flex-col gap-4 py-4">
        <div class="daily flex flex-col gap-4">
            <div class="daily-check daily p-4 rounded-2xl flex border-4 border-primary_light gap-4 items-center">
                <img src="../../assets/images/tasks/gift-icon.png" class=" w-16 h-16" alt="">
                <div class="flex flex-col">
                    <h2>{{ t('screens.tasks.tasks.dailyCheck.header') }}</h2>
                    <p class="label">{{ t('screens.tasks.tasks.dailyCheck.description') }}</p>
                </div>
                <i class="pi pi-chevron-right text-dark" style="font-size:1.25rem"></i>
            </div>
            <div class="daily-combo daily p-4 rounded-2xl flex border-4 border-primary_light gap-4 items-center">
                <img src="../../assets/images/tasks/calendar-icon.png" class=" w-16 h-16" alt="">
                <div class="flex flex-col">
                    <h2>{{ t('screens.tasks.tasks.dailyCombo.header') }}</h2>
                    <p class="label">{{ t('screens.tasks.tasks.dailyCombo.description') }}</p>
                </div>
                <i class="pi pi-chevron-right text-dark" style="font-size:1.25rem"></i>
            </div>
        </div>
        <h1>{{ t('screens.tasks.tasks.header') }}</h1>
    </section>
</template>


<style>
</style>