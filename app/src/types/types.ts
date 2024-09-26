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
    clicked: boolean = false;
}

export class Social {
    name: string = '';
    url: string = '';
}

export class MinBet {
    currencyName: string = '';
    currencyId: number = 0;
    minBet: number = 0;
}

export class League {
    id: number = 0;
    name: string = '';
    bottomThreshold: number = 0;
    topThreshold: number = 0;
}

export class ExternalUser {
    id: number = 0;
    username: string = '';
    email: string = '';
    phone: string = '';
    sex_id: string = '';
    birth_date: string = '';
    firstname: string = '';
    middlename: string = '';
    lastname: string = '';
    photo: string = '';
    cover: string = '';
    bio: string = '';
    isFollower: boolean = false;
    followersCount: number = 0;
    followingCount: number = 0;
    favoritesCount: number = 0;
    friendsCount: number = 0;
    activeTournamentsCount: number = 0;
    isFriend: boolean = false;
    isFavorite: boolean = false;
    city: string = '';
    exp: number = 0;
    expYear: number = 0;
    expOverall: number = 0;
    country: string = '';
    expRatingPlace: number = 0;
    transport: number[] = [];
    fPoints: number = 0;
    league: League = new League();
    refCode: string = '';
    inviteCode: string = '';
    appliedRefCode: string = '';
    codeGenerationsLeft: number = 0;
    social: Social[] = [];
    minBet: MinBet[] = [];
}

export class Battle {
    id: number = 0;
    user: ExternalUser = new ExternalUser();
    opponent: ExternalUser = new ExternalUser();
    status: string = '';
    userResult: number = 0;
    opponentResult: number = 0;
    pick: number = 0;
}

export class Round {
    endTime: number = 0;
    id: number = 0;
    battles: Battle[] = [];
}