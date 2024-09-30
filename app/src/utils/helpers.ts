import axios from 'axios'
import { useAccountStore } from '@/stores/account'
declare const Telegram: any



export const auth = async ():Promise<boolean> => {
    const store = useAccountStore()
    let initData = ""
    let refCode = ""
    if (typeof Telegram !== 'undefined' && Telegram.WebApp) {
        const tg = Telegram.WebApp;
        if (tg.initData) {
            initData = tg.initData;
        }
        if (tg.initDataUnsafe){
            refCode = tg.initDataUnsafe.start_param
        }
    } else {
        console.log("Telegram Web App SDK не доступен.");
    }
    if (initData == ""){
        return false
    }

    try {
        const {data} = await axios.post(`${import.meta.env.VITE_API_URL}/users/auth`, {
            initData: initData,
            refCode: refCode
        }, {
            withCredentials: true
        })
        store.token = data.accessToken
        return data.isNew
    } catch (error) {
        console.log(error)
        return false
    }
}

export const getUser = async (uid: number) => {
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
		return false
	}
}