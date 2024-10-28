import { Injectable } from '@angular/core';
import { environment } from '../../../../environments/environment';
import { BehaviorSubject, Observable, map } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { League } from '../../../models/league.model'; // Assuming you have this model defined
import { User } from '../../../models/user.model';
import { Leagues } from '../../../models/leagues.model';
import { Stock } from '../../../models/stock.model';
import { Portfolio } from '../../../models/portfolio.model';

@Injectable({
  providedIn: 'root'
})
export class UserLeaguesService {
  private baseUrl = environment.api_url;
  private findUserLeagues = `${this.baseUrl}/api/users/user-leagues`;
  private findLeagueMembers = `${this.baseUrl}/api/leagues/details`;
  private findUserPortfolio = `${this.baseUrl}/api/portfolio/league-portfolio`;

  // BehaviorSubject for managing the selected league
  private selectedLeagueSource = new BehaviorSubject<League | null>(this.getStoredLeague());
  selectedLeague$ = this.selectedLeagueSource.asObservable();

  constructor(private http: HttpClient) {}

  // Method to get the user's leagues using their user ID
  getUserLeagues(userId: number): Observable<Leagues> {
    return this.http.post<Leagues>(this.findUserLeagues, { user_id: userId });
  }

  // Method to get members of a league using the league ID
  getLeagueMembers(leagueId: number): Observable<User[]> {
    return this.http.post<League>(this.findLeagueMembers, { league_id: leagueId }).pipe(
      map((league: League) => league.users || []) // Ensure 'users' is not null or undefined, return empty array if it is
    );
  }

  // Method to fetch the user's portfolio for a specific league
  getUserPortfolio(userId: number, leagueId: number): Observable<Portfolio> {
    // Send a POST request with the user ID and league ID as the request payload
    return this.http.post<Portfolio>(this.findUserPortfolio, { user_id: userId, league_id: leagueId });
  }

  // Method to set the selected league
setSelectedLeague(league: League | null): void {
  console.log(league);
  this.selectedLeagueSource.next(league); // Set the selected league as the full League object
  if (league) {
    // Store the entire league object as a JSON string in localStorage
    localStorage.setItem('selectedLeague', JSON.stringify(league)); 
  } else {
    localStorage.removeItem('selectedLeague');
  }
}

  // Retrieve the stored league from localStorage (if it exists)
// Retrieve the stored league from localStorage (if it exists)
private getStoredLeague(): League | null {
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


}
