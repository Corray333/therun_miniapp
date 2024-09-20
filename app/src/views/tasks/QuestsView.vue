<script lang="ts" setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

</script>

<template>
    <section class="w-full min-h-screen h-full flex justify-center items-center">
        <h1 class=" text-dark">{{ t('comingSoon') }}</h1>
    </section>
</template>


<style></style>