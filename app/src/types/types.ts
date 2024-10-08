

export class User {
    id: number=0;
    avatar: string="";
    username?: string;
    pointBalance: number=0;
    raceBalance: number=0;

    red_keyBalance: number=0;
    blue_keyBalance: number=0;
    green_keyBalance: number=0;
    
    lastClaim: number=0;
    farmingFrom: number=0;
    maxPoints: number=0;
    farmingTime: number=0;
    refCode: string="";
    referer?: number;
}

export class Referal {
    avatar: string="";
    username: string="";
    isPremium: boolean = false;
}

export class Task {
    id: number = 0;
    description: string = '';
    type: string = '';
    link: string = '';
    expireAt: number = 0;
    pointsReward: number = 0;
    keysReward: number = 0;
    raceReward: number = 0;
    data: any = null;
    icon: string = '';
    done: boolean = false;
    claimed: boolean = false;
    clicked: boolean = false;
}

export class ExternalUser {
    id: number = 0;
    username: string = "";
    photo: string = "";
}

export class Battle {
    id: number = 0;
    roundID: number = 0;
    user: ExternalUser = new ExternalUser();
    opponent: ExternalUser = new ExternalUser();
    pick: number = 0;
    userResult: number = 0;
    opponentResult: number = 0;
}

export interface Round {
    id: number;
    endTime: number;
    battles: Battle[];
}

export class Key {
    type!: string;
    amount!: number;
}

export class Case {
    type!: string;
    keys!: Key[];
    rewardType!: string;
    min_rewards!: number;
    max_rewards!: number;
}

export class Reward{
    type!: string;
    amount!: number;
}

export class UpgradeCost {
    currency: string;
    amount: number;

    constructor(currency: string, amount: number) {
        this.currency = currency;
        this.amount = amount;
    }
}

export class Building {
    img: string;
    type: string;
    level: number;
    upgradeCost: UpgradeCost[] | null;

    state: string = '';
    lastStateChange: number = 0;
    stateUntil: number = 0;

    constructor(img: string, type: string, level: number, upgradeCost: UpgradeCost[] | null) {
        this.img = img;
        this.type = type;
        this.level = level;
        this.upgradeCost = upgradeCost;
    }
}

export class Buildings {
    fabric!: Building;
    mine!: Building;
    warehouse!: Building;

    constructor(fabric: Building, mine: Building, warehouse: Building) {
        this.fabric = fabric;
        this.mine = mine;
        this.warehouse = warehouse;
    }
}



// City

export class Resource {
    name: string = '';
    type: string = '';
    amount: number = 0;
}

class Cost {
    currency: string = '';
    amount: number = 0;
}

export class Requirement {
    type: string = '';
    level: number = 0;
}

class LevelDetail {
    capacity: number = 0;
    cost: Cost[] = [];
    requirements: Requirement[] | null = null;
    buildingDuration: number = 0;
}

export class Warehouse {
    img: string = '';
    type: string = '';
    level: number = 0;
    state: string = '';
    lastStateChange: number = 0;
    stateUntil: number = 0;
    resources: Resource[] = [];
    currentLevel: LevelDetail | null = null;
    nextLevel: LevelDetail | null = null;
}