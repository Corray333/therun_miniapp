<script setup lang="ts">
import { onBeforeMount, ref } from 'vue'
import { RouterView, useRouter, useRoute } from 'vue-router'
import { auth } from '@/utils/helpers'
import { useAccountStore } from './stores/account'
import axios from 'axios'
import { useI18n } from 'vue-i18n'
declare const Telegram: any

const i18n = useI18n()


const router = useRouter()
const temp = ref<string>("")

const getUser = async (uid: number) => {
	const store = useAccountStore()

	try {
		const { data } = await axios.get(`${import.meta.env.VITE_API_URL}/users/${uid}`, {
			withCredentials: true,
			headers: {
				Authorization: store.token
			}
		})
		store.user = data
		return true
	} catch (error) {
		console.log(error)
	}
}

onBeforeMount(async () => {

	const locale = localStorage.getItem('locale') || 'en'
	i18n.locale.value = locale

	
	const tg = Telegram.WebApp
	temp.value = tg.initData
	const uid = tg.initDataUnsafe.user.id
	await auth()
	const user = await getUser(uid)

	Telegram.WebApp.onEvent('backButtonClicked', function () {
		router.back()
	})
})
</script>

<template>
	<RouterView />
</template>

<style scoped></style>
