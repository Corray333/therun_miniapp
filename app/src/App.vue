<script setup lang="ts">
import { onBeforeMount, ref, onMounted } from 'vue'
import { RouterView, useRouter, useRoute } from 'vue-router'
import { auth, getUser } from '@/utils/helpers'
import { useI18n } from 'vue-i18n'
import Navbar from './components/Navbar.vue'
import Errors from './components/Errors.vue'
declare const Telegram: any

const i18n = useI18n()


const router = useRouter()
const route = useRoute()

const temp = ref<string>("")

const excludedRoutes = ['onboarding'];

const showStart = ref<boolean>(true)
const loggingIn = ref<boolean>(false)



onMounted(() => {
	setTimeout(() => {
		showStart.value = false
	}, 2000)
})

onBeforeMount(async () => {

	const locale = localStorage.getItem('locale') || 'en'
	i18n.locale.value = locale


	const tg = Telegram.WebApp
	try {
		tg.expand()
		tg.disableVerticalSwipes()
	} catch (error) {
		alert(error)
	}
	temp.value = tg.initData
	const uid = tg.initDataUnsafe.user.id
	const isNew = await auth()
	const logged = await getUser(uid)
	if (logged) loggingIn.value = false
	if (isNew) {
		router.push('/onboarding')
	}

	Telegram.WebApp.onEvent('backButtonClicked', function () {
		router.back()
	})
})
</script>

<template>
	<Transition>
		<img v-show="showStart || loggingIn" class=" fixed w-screen h-screen object-cover z-60" src="./assets/images/start.png" alt="">
	</Transition>
	<section>
		<Errors class=" fixed top-0 left-0" />
		<RouterView />
		<Navbar v-show="!excludedRoutes.includes(String(route.name))" class=" fixed bottom-0 left-0 z-20" />
	</section>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
	transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
	opacity: 0;
}
</style>
