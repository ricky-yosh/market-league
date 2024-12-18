import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { BehaviorSubject, Observable, map } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { League } from '../../models/league.model'; // Assuming you have this model defined
import { User } from '../../models/user.model';
import { Leagues } from '../../models/leagues.model';
import { Stock } from '../../models/stock.model';
import { Portfolio } from '../../models/portfolio.model';
import { Trade } from '../../models/trade.model';
import { devLog } from '../../../environments/development/devlog';
import { LeaguePortfolio } from '../../models/league-portfolio';

@Injectable({
  providedIn: 'root'
})
export class LeagueService {
  private baseUrl = environment.api_url;
  private findUserLeaguesUrl = `${this.baseUrl}/api/users/user-leagues`;
  private findLeagueMembersUrl = `${this.baseUrl}/api/leagues/details`;
  private findUserPortfolioUrl = `${this.baseUrl}/api/portfolio/league-portfolio`;
  private findTradesUrl = `${this.baseUrl}/api/trades/get-trades`;
  private createLeagueUrl = `${this.baseUrl}/api/leagues/create-league`;
  private removeLeagueUrl = `${this.baseUrl}/api/leagues/remove-league`;
  private createTradeUrl = `${this.baseUrl}/api/trades/create-trade`;
  private confirmTradeUrl = `${this.baseUrl}/api/trades/confirm-trade`;
  private getLeaguePortfolioInfoUrl = `${this.baseUrl}/api/league-portfolio/get-league-portfolio-info`;
  private draftStockUrl = `${this.baseUrl}/api/league-portfolio/draft-stock`;
  
  // BehaviorSubject for managing the selected league
  private selectedLeagueSource = new BehaviorSubject<League | null>(this.getStoredLeague());
  selectedLeague$ = this.selectedLeagueSource.asObservable();

  constructor(private http: HttpClient) {}

  // Method to get the user's leagues using their user ID
  getUserLeagues(userId: number): Observable<Leagues> {
    return this.http.post<Leagues>(this.findUserLeaguesUrl, { user_id: userId });
  }

  // Method to get members of a league using the league ID
  getLeagueMembers(leagueId: number): Observable<User[]> {
    return this.http.post<League>(this.findLeagueMembersUrl, { league_id: leagueId }).pipe(
      map((league: League) => league.users || []) // Ensure 'users' is not null or undefined, return empty array if it is
    );
  }
  
  // Method to fetch the user's portfolio for a specific league
  getUserPortfolio(userId: number, leagueId: number): Observable<Portfolio> {
    const payload = {
      user_id: userId,
      league_id: leagueId
    }
    devLog("Payload: ", payload)
    // Send a POST request with the user ID and league ID as the request payload
    return this.http.post<Portfolio>(this.findUserPortfolioUrl, payload);
  }

  // Method to set the selected league
  setSelectedLeague(league: League | null): void {
    devLog("Selected League: ", league);
    this.selectedLeagueSource.next(league); // Set the selected league as the full League object
    if (league) {
      // Store the entire league object as a JSON string in localStorage
      localStorage.setItem('selectedLeague', JSON.stringify(league)); 
    } else {
      localStorage.removeItem('selectedLeague');
    }
  }

  // Retrieve the stored league from localStorage (if it exists)
  getStoredLeague(): League | null {
    const storedLeague = localStorage.getItem('selectedLeague');
    
    // Check if storedLeague is a valid JSON
    if (storedLeague) {
      try {
        return JSON.parse(storedLeague) as League;
      } catch (e) {
        console.error("Error parsing stored league JSON:", e);
        localStorage.removeItem('selectedLeague'); // Clean up invalid entry
        return null;
      }
    }

    return null;
  }

  // Fetch user trades based on userId and leagueId
  getTrades(userId: number, leagueId: number, receiving_trade: boolean = false, sending_trade: boolean = false): Observable<Trade[]> {
    const payload = {
      user_id: userId,
      league_id: leagueId,
      receiving_trade: receiving_trade,
      sending_trade: sending_trade
    }
    return this.http.post<any>(this.findTradesUrl, payload);
  }

  // Create League
  createLeague(leagueName: string, ownerUser: number, endDate: string): Observable<any> {
    const payload = {
      league_name: leagueName,
      owner_user: ownerUser,
      end_date: endDate
    };
    return this.http.post<any>(this.createLeagueUrl, payload); // Send POST request
  }

  // Remove League
  removeLeague(leagueId: number): Observable<any> {
    const payload = { league_id: leagueId }; // Payload with league_id
    return this.http.post<any>(this.removeLeagueUrl, payload); // Send POST request
  }

  createTrade(leagueId: number, user1Id: number, user2Id: number, stocks1Id: number[], stocks2Id: number[]): Observable<Trade> {
    const payload = {
      league_id: leagueId,
      user1_id: user1Id,
      user2_id: user2Id,
      stocks1_ids: stocks1Id,
      stocks2_ids: stocks2Id
    }

    devLog("Payload: ", payload)
    return this.http.post<Trade>(this.createTradeUrl, payload); // Send POST request to create a trade
  }

  confirmTradeForUser(tradeId: number, userId: number): Observable<any> {
    const payload = {
      trade_id: tradeId,
      user_id: userId
    }

    devLog("Payload: ", payload)
    return this.http.post<any>(this.confirmTradeUrl, payload); // Send POST request to create a trade
  }

  getLeaguePortfolioInfo(leagueId: number): Observable<LeaguePortfolio> {
    const payload = {
      league_id: leagueId
    }

    devLog("Payload: ", payload)
    return this.http.post<LeaguePortfolio>(this.getLeaguePortfolioInfoUrl, payload);
  }

  draftStock(leagueId: number, userId: number, stockId: number): Observable<any> {
    const payload = {
      league_id: leagueId,
      user_id: userId,
      stock_id: stockId
    }

    devLog("Payload: ", payload)
    return this.http.post<any>(this.draftStockUrl, payload);
  }

}
