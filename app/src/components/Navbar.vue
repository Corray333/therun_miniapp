<script lang="ts" setup>
import gift from '@/components/icons/gift-icon.vue'
import chibi from '@/components/icons/chibi-icon.vue'
import battles from '@/components/icons/battles-icon.vue'
import friens from '@/components/icons/friens-icon.vue'
import tasks from '@/components/icons/tasks-icon.vue'
import more from '@/components/icons/more-icon.vue'
import { useRoute } from 'vue-router'

const route = useRoute()

</script>

<template>
    <nav class=" w-full bg-white">
        <ul class=" flex justify-around py-4 w-full list-none text-sm">
            <li >
                <router-link to="/farm" class=" flex flex-col gap-1 text-dark items-center">
                    <gift :color="route.path == '/farm' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ $t('menu.bonuses') }}</p>
                </router-link>
            </li>
            <li >
                <router-link to="/chibi" class=" flex flex-col gap-1 text-dark items-center">
                    <chibi :color="route.path == '/chibi' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ $t('menu.chibi') }}</p>
                </router-link>
            </li>
            <li >
                <router-link to="/battles" class=" flex flex-col gap-1 text-dark items-center">
                    <battles :color="route.path == '/battles' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ $t('menu.battles') }}</p>
                </router-link>
            </li>
            <li >
                <router-link to="/friens" class=" flex flex-col gap-1 text-dark items-center">
                    <friens :color="route.path == '/friens' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ $t('menu.friens') }}</p>
                </router-link>
            </li>
            <li >
                <router-link to="/tasks" class=" flex flex-col gap-1 text-dark items-center">
                    <tasks :color="route.path == '/tasks' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ $t('menu.tasks') }}</p>
                </router-link>
            </li>
            <li >
                <router-link to="/more" class=" flex flex-col gap-1 text-dark items-center">
                    <more :color="route.path == '/more' ? 'var(--primary)' : 'var(--dark)'" />
                    <p>{{ $t('menu.more') }}</p>
                </router-link>
            </li>
        </ul>
    </nav>
</template>


<style scoped>

.router-link-active{
    color: var(--primary)
}

nav{
    box-shadow: 0 -0.25rem 1rem 0 rgba(0, 0, 0, 0.1);
}

</style>