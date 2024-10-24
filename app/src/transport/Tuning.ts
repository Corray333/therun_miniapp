import axios, {isAxiosError} from "axios"
import { Module } from "@/types/types"
import {auth} from "@/utils/helpers"
import { useAccountStore } from "@/stores/account"
import { useComponentsStore } from "@/stores/components"

const api = axios.create({ baseURL: import.meta.env.VITE_API_URL })
const accStore = useAccountStore()
const componentsStore = useComponentsStore()

class TuningTransport {
    getTuningModules = async (characteristic: string) : Promise<Module[]> => {
        try {
            let url = ""
            if (characteristic == ""){
                url = `/modules`
            } else {
                url = `/modules?characteristic=${characteristic}`
            }
            const { data } = await api.get(url, {
                withCredentials: true,
                headers: {
                    Authorization: accStore.token
                }
            })
            return data
        } catch (error) {
            if (isAxiosError(error) && error.response?.status === 401) {
                await auth()
                try {
                    return await this.getTuningModules(characteristic)
                } catch (error) {
                    if (isAxiosError(error)) {
                        componentsStore.addError(error.message)
                    }
                }
            } else {
                componentsStore.addError(`${error}`)
            }
            return []
        }
    }
}

export {TuningTransport}