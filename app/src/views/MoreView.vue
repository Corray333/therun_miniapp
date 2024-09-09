<script lang="ts" setup>
import { ref, watch } from 'vue'

import Navbar from '@/components/Navbar.vue'
import Select from 'primevue/select'
import { useI18n } from 'vue-i18n'

const i18n = useI18n()

const languages = [
    { name: 'English', code: 'en' },
    { name: 'Russian', code: 'ru' },
    // { name: 'Spanish', code: 'es' },
    // { name: 'German', code: 'de'}
]

const language = ref(languages.find((lang) => lang.code === i18n.locale.value) || languages[0])

watch(language, (newLanguage) => {
    i18n.locale.value = newLanguage.code
    localStorage.setItem('locale', newLanguage.code)
})

</script>

<template>
    <section>
        <section class=" flex flex-col min-h-screen p-5 h-full w-full">
            <Select v-model="language" :options="languages" optionLabel="name" placeholder="Select a City" class="w-full bg-primary text-black" />
        </section>
        <Navbar class=" fixed bottom-0" />
    </section>
</template>


<style>
.p-select-option-selected {
    background-color: var(--primary) !important;
    color: white !important;
}

.p-select:not(.p-disabled).p-focus {
    border-color: inherit !important;
}
.p-select-option-label {
    color: inherit !important;
}

</style>