import { ref } from 'vue'
import { defineStore } from 'pinia'
import { Car, RaceState } from '@/types/types'


export const useCarsStore = defineStore('cars', () => {
    const cars = ref<Car[]>([])
    const mainCar = ref<Car>({
        "id": 1,
        "isMain": true,
        "element": "desert",
        "img": "https://notably-great-coyote.ngrok-free.app/static/images/cars/desert-car.png",
        "acceleration": 58,
        "handling": 61,
        "brakes": 94,
        "strength": 69,
        "tank": 69,
        "fuel": 9.995014245014245,
        "health": 9.997507122507123,
        "modules": [
            {
                "characteristic": "acceleration",
                "boost": 110,
                "name": "Thing",
                "isTemp": false,
                "roundId": null,
                "carId": null,
                "userModuleId": 2,
                "img": "https://notably-great-coyote.ngrok-free.app/static/images/cars/modules/1.png"
            }
        ],
        "speed": 28.78,
        "fuelWasting": 129.07829999999998,
        "healthWasting": 258.15659999999997
    })
    const raceState = ref<RaceState>({
        "currentMiles": 0,
        "startTime": 0,
        "place": 0
    })

    return { cars, mainCar, raceState }
})
