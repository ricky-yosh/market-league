import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class UserPortfolioService {
  private baseUrl = environment.api_url;
  private findUserLeagues = `${this.baseUrl}/api/users/user-leagues`;
  private findLeagueMembers = `${this.baseUrl}/api/leagues/details`;

  // BehaviorSubject for managing the selected league
  private selectedLeagueSource = new BehaviorSubject<League | null>(this.getStoredLeague());
  selectedLeague$ = this.selectedLeagueSource.asObservable();

  constructor(private http: HttpClient) {}
}