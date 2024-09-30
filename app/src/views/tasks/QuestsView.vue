<script lang="ts" setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

</script>

<template>
    <section class="w-full flex justify-center min-h-screen items-center py-4">
        <h1 class=" text-dark">{{ t('screens.tasks.noNewTasks') }}</h1>
    </section>
</template>


<style></style>