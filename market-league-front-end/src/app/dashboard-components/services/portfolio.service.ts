import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { devLog } from '../../../environments/development/devlog';
import { Observable } from 'rxjs';
import { Portfolio } from '../../models/portfolio.model';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class PortfolioService {
  private baseUrl = environment.api_url;
  private findUserPortfolioUrl = `${this.baseUrl}/api/portfolio/league-portfolio`;

  constructor(private http: HttpClient) {}

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

}
