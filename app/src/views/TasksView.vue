<script lang="ts" setup>
import {RouterView} from 'vue-router'

</script>

<template>
    <section class=" pb-20 min-h-screen">
        <section class=" flex flex-col h-full w-full p-4">
            <ul class=" flex w-full justify-between">
                <li class=" w-full text-center">
                    <router-link to="/tasks/tasks" class=" w-full block">Tasks</router-link>
                </li>
                <li class=" w-full text-center">
                    <router-link to="/tasks/quests" class=" w-full block">Quests</router-link>
                </li>
                <li class=" w-full text-center">
                    <router-link to="/tasks/special" class=" w-full block">Special</router-link>
                </li>
            </ul>
            <RouterView />
        </section>
    </section>
</template>


<style scoped>

.router-link-active{
    border-bottom: solid 2px var(--primary);
}

</style>