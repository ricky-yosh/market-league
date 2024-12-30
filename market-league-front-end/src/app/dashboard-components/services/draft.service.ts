import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { LeaguePortfolio } from '../../models/league-portfolio';
import { Observable } from 'rxjs';
import { devLog } from '../../../environments/development/devlog';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class DraftService {
  private baseUrl = environment.api_url;
  private getLeaguePortfolioInfoUrl = `${this.baseUrl}/api/league-portfolio/get-league-portfolio-info`;
  private draftStockUrl = `${this.baseUrl}/api/league-portfolio/draft-stock`;

  constructor(private http: HttpClient) {}

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
