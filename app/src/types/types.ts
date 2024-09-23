export class User {
    id: number=0;
    avatar: string="";
    username?: string;
    pointBalance: number=0;
    raceBalance: number=0;
    keyBalance: number=0;
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
}