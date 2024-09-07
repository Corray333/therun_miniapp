export class User {
    id: number=0;
    avatar: string="";
    username?: string;
    pointBalance: number=0;
    raceBalance: number=0;
    keyBalance: number=0;
    lastClaim: number=0;
    maxPoints: number=0;
    farmingTime: number=0;
    refCode: string="";
    referer?: number;
}

export class Referal {
    avatar: string="";
    username: string="";
}