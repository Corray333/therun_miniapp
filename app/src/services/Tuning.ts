import { TuningTransport } from "@/transport/Tuning"
import { Module } from "@/types/types"

const tuningTransport = new TuningTransport()

const getTuningModules = async (characteristic: string) : Promise<Module[]> => {
    return await tuningTransport.getTuningModules(characteristic)
}

export { getTuningModules }