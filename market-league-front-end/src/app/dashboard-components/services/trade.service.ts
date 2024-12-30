import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { devLog } from '../../../environments/development/devlog';
import { Trade } from '../../models/trade.model';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TradeService {
  private baseUrl = environment.api_url;
  private findTradesUrl = `${this.baseUrl}/api/trades/get-trades`;
  private createTradeUrl = `${this.baseUrl}/api/trades/create-trade`;
  private confirmTradeUrl = `${this.baseUrl}/api/trades/confirm-trade`;

  constructor(private http: HttpClient) {}

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

}
