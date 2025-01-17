import { LeaguePortfolio } from "../../league-portfolio.model";
import { League } from "../../league.model";
import { Portfolio } from "../../portfolio.model";

export interface CreateLeagueResponse {
    league: League;
    league_portfolio: LeaguePortfolio;
    user_portfolio: Portfolio;
}