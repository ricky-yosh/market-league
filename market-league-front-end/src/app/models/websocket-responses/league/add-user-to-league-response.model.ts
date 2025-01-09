import { League } from "../../league.model";
import { Stock } from "../../stock.model";
import { User } from "../../user.model";

export interface AddUserToLeagueResponse {
    id: number;
    user_id: number;
    user: User;
    league_id: number;
    league: League;
    stocks: Stock[] | null;
    created_at: string;
}