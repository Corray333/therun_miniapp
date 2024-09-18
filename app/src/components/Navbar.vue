<script lang="ts" setup>
import gift from '@/components/icons/gift-icon.vue'
import chibi from '@/components/icons/chibi-icon.vue'
import battles from '@/components/icons/battles-icon.vue'
import friens from '@/components/icons/friens-icon.vue'
import tasks from '@/components/icons/tasks-icon.vue'
import more from '@/components/icons/more-icon.vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import {ref} from 'vue'

const { t } = useI18n()

const route = useRoute()

const showMore = ref<boolean>(false)

</script>

<template>
    <nav class=" w-full bg-white">
        <ul class=" flex justify-around py-4 relative w-full list-none text-sm">
            <Transition name="resize">
                <aside v-if="showMore" class=" text-dark absolute bottom-full mb-4 mr-4 right-0 p-4 rounded-2xl bg-white sm-shadow">
                    <ul class=" flex flex-col">
                        <li class=" flex gap-2 p-2">
                            <router-link to="/settings" class=" flex items-center gap-2">
                                <i class="pi pi-cog text-2xl"></i> Settings
                            </router-link>
                        </li>
                        <li class=" flex gap-2 p-2">
                            <router-link to="/onboarding" class=" flex items-center gap-2">
                                <i class="pi pi-info-circle text-2xl"></i> Info
                            </router-link>
                        </li>
                    </ul>
                </aside>
            </Transition>
            <li>
                <router-link to="/" class=" flex flex-col gap-1 text-dark items-center">
                    <gift :color="route.path == '/' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ t('menu.bonuses') }}</p>
                </router-link>
            </li>
            <li>
                <router-link to="/chibi" class=" flex flex-col gap-1 text-dark items-center">
                    <chibi :color="route.path == '/chibi' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ t('menu.chibi') }}</p>
                </router-link>
            </li>
            <li>
                <router-link to="/battles" class=" flex flex-col gap-1 text-dark items-center">
                    <battles :color="route.path == '/battles' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ t('menu.battles') }}</p>
                </router-link>
            </li>
            <li>
                <router-link to="/friens" class=" flex flex-col gap-1 text-dark items-center">
                    <friens :color="route.path == '/friens' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ t('menu.friens') }}</p>
                </router-link>
            </li>
            <li>
                <router-link to="/tasks" class=" flex flex-col gap-1 text-dark items-center">
                    <tasks :color="route.path == '/tasks' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ t('menu.tasks') }}</p>
                </router-link>
            </li>
            <li>
                <p @click="showMore = !showMore" class=" flex flex-col gap-1 text-dark items-center">
                    <more :color="route.path == '/more' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ t('menu.more') }}</p>
                </p>
            </li>
        </ul>
    </nav>
</template>


<style>
.router-link-active {
    color: var(--primary);
}

.router-link-active svg{
    transform: scale(1.5);
    transition: transform 0.25s ease;
}


nav {
    box-shadow: 0 -0.25rem 1rem 0 rgba(0, 0, 0, 0.1);
}

.resize-enter-active,
.resize-leave-active {
    transition: transform 0.25s ease;
    transform-origin: right bottom;
}

.resize-enter-from,
.resize-leave-to {
    transform: scale(0);
}
</style>