import { LeaguePlayer } from "./league-player.model";
import { LeagueState } from "./league-state.model";
import { User } from "./user.model";

// Define an interface for a League
export interface League {
    id: number;
    league_name: string;
    start_date: string;
    end_date: string;
    league_state: LeagueState;
    users: User[] | null;
    max_players: number;
    league_players: LeaguePlayer[] | null;
}