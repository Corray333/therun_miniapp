<script lang="ts" setup>
import gift from '@/components/icons/gift-icon.vue'
import chibi from '@/components/icons/chibi-icon.vue'
import battles from '@/components/icons/battles-icon.vue'
import friens from '@/components/icons/friens-icon.vue'
import tasks from '@/components/icons/tasks-icon.vue'
import more from '@/components/icons/more-icon.vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { ref } from 'vue'

const { t } = useI18n()

const route = useRoute()

const showMore = ref<boolean>(false)

</script>

<template>
    <nav class=" w-full bg-white">
        <ul class=" flex justify-around py-4 relative w-full list-none text-sm">
            <Transition name="delay">
                <div class=" wrapper fixed z-50 w-full h-screen top-0 left-0 flex items-end" v-show="showMore" @click="showMore = false">
                    <Transition name="resize">
                        <aside v-if="showMore"
                            class=" text-dark absolute mb-24 mr-4 right-0 p-4 rounded-2xl bg-white sm-shadow">
                            <ul class=" flex flex-col more-list">
                                <li class=" flex gap-2 p-2">
                                    <router-link to="/settings" class=" flex items-center gap-2">
                                        <i class="pi pi-cog"></i> Settings
                                    </router-link>
                                </li>
                                <li class=" flex gap-2 p-2">
                                    <router-link to="/onboarding" class=" flex items-center gap-2">
                                        <i class="pi pi-info-circle"></i> Info
                                    </router-link>
                                </li>
                            </ul>
                        </aside>
                    </Transition>
                </div>
            </Transition>
            <li>
                <router-link to="/" class=" flex flex-col gap-1 text-dark items-center router-link">
                    <gift :color="route.path == '/' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ t('menu.bonuses') }}</p>
                </router-link>
            </li>
            <li>
                <router-link to="/chibi" class=" flex flex-col gap-1 text-dark items-center router-link">
                    <chibi :color="route.path == '/chibi' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ t('menu.chibi') }}</p>
                </router-link>
            </li>
            <li>
                <router-link to="/battles" class=" flex flex-col gap-1 text-dark items-center router-link">
                    <battles :color="route.path == '/battles' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ t('menu.battles') }}</p>
                </router-link>
            </li>
            <li>
                <router-link to="/friens" class=" flex flex-col gap-1 text-dark items-center router-link">
                    <friens :color="route.path == '/friens' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ t('menu.friens') }}</p>
                </router-link>
            </li>
            <li>
                <router-link to="/tasks" class=" flex flex-col gap-1 text-dark items-center router-link">
                    <tasks :color="route.path == '/tasks' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ t('menu.tasks') }}</p>
                </router-link>
            </li>
            <li>
                <p @click="showMore = true" class=" flex flex-col gap-1 text-dark items-center">
                    <more :color="route.path == '/more' ? 'var(--primary)' : 'var(--dark)'" />
                <p>{{ t('menu.more') }}</p>
                </p>
            </li>
        </ul>
    </nav>
</template>


<style scoped>
.router-link-active {
    color: var(--primary);
}

.more-list i {
    font-size: 1.5rem;
}

.router-link svg {
    transition: all 0.25s ease;
}

.router-link-active svg {
    transform: scale(1.5);
}

.router-link-active svg#friens {
    transform: scale(1.75);
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

.delay-enter-active,
.delay-leave-active {
    transition: opacity 0.5s ease;
}

.delay-enter-from,
.delay-leave-to {
    opacity: 1;
}
</style>