import { User } from './user.model';
import { League } from './league.model';
import { Stock } from './stock.model';

export interface Portfolio {
    id: number;
    user_id: number;
    user: User;
    league_id: number;
    league: League;
    stocks: Stock[];
    points: number;
    created_at: string; // ISO date string
}