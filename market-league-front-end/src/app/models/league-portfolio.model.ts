import { League } from "./league.model";
import { Stock } from "./stock.model";

// Define an interface for a League
export interface LeaguePortfolio {
    id: number;
    league_id: number;
    league: League;
    name: string;
    stocks: Stock[] | null;
}