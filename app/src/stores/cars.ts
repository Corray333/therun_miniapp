import { ref } from 'vue'
import { defineStore } from 'pinia'
import { Car, RaceState } from '@/types/types'


export const useCarsStore = defineStore('cars', () => {
    const cars = ref<Car[]>([])
    const mainCar = ref<Car>({
        "id": 4,
        "isMain": true,
        "element": "track",
        "img": "https://notably-great-coyote.ngrok-free.app/static/images/cars/track-car.png",
        "acceleration": 78,
        "hendling": 85,
        "brakes": 93,
        "strength": 83,
        "tank": 69,
        "fuel": 9.1296605970519,
        "health": 9.638232416846876,
        "modules": null,
        "speed": 84.75,
        "fuelWasting": 152.04149999999998,
        "healthWasting": 365.781
    })
    const raceState = ref<RaceState>({
        "currentMiles": 0,
        "startTime": 0,
        "place": 0
    })

    return { cars, mainCar, raceState }
})
