import { Injectable } from '@angular/core';
import { devLog } from '../../../environments/development/devlog';
import { Trade } from '../../models/trade.model';
import { Subject } from 'rxjs';
import { WebSocketMessageTypes } from './websocket-message-types';
import { WebSocketService } from './websocket.service';
import { ConfirmTradeResponse } from '../../models/websocket-responses/trade/confirm-trade-response.model';
import { VerifyUserService } from './verify-user.service';
import { LeagueService } from './league.service';
import { guard } from '../../utils/guard';

@Injectable({
  providedIn: 'root'
})
export class TradeService {

  // * Observables

  // League Trades
  private leagueTradesSubject = new Subject<Trade[]>();
  leagueTrades$ = this.leagueTradesSubject.asObservable();

  // * Constructor

  constructor(
    private webSocketService: WebSocketService,
    private verifyUserService: VerifyUserService,
    private leagueService: LeagueService,
  ) {
      this.webSocketService.getMessages().subscribe((message) => {
        switch (message.type) {
          case WebSocketMessageTypes.MessageType_Trade_GetTrades:
            devLog("Received GetTrades Response: " + message.data);
            this.handleGetTradesResponse(message.data);
            break;
          case WebSocketMessageTypes.MessageType_Trade_CreateTrade:
            devLog("Received CreateTrade Response: " + message.data);
            this.handleCreateTradeResponse(message.data);
            break;
          case WebSocketMessageTypes.MessageType_Trade_ConfirmTrade:
            devLog("Received ConfirmTradeForUser Response: " + message.data);
            this.handleConfirmTradeForUserResponse(message.data);
            break;
          default:
            // devLog("Trade Service unable to route Websocket Message properly! " + message.data);
      }
    });
  }

  // * Websocket Response Handler Functions
  
  handleGetTradesResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulGetTradesResponse(responseData);
  }

  handleCreateTradeResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulCreateTradeResponse(responseData);
  }

  handleConfirmTradeForUserResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulConfirmTradeForUserResponse(responseData);
  }

  // * Helper Functions to Websocket Responses

  handleSuccessfulGetTradesResponse(trades: Trade[]): void {
    this.leagueTradesSubject.next(trades);
  }

  handleSuccessfulCreateTradeResponse(trade: Trade): void {
    devLog("Created Trade: " + trade);
    this.confirmTradeForUser(trade.id);
    this.getTrades();
  }

  handleSuccessfulConfirmTradeForUserResponse(response: ConfirmTradeResponse): void {
    devLog("Confirmed Trade for User: " + response.message);
    this.getTrades();
  }

  // * Websocket Call Functions

  // Fetch user trades based on userId and leagueId
  getTrades(receiving_trade: boolean = false, sending_trade: boolean = false): void {
    const currentUser = this.verifyUserService.getCurrentUserValue();
    const selectedLeague = this.leagueService.getSelectedLeagueValue();
    guard(currentUser != null, "Current user is null!");
    guard(selectedLeague != null, "Selected League is null!");

    const data = {
      user_id: currentUser.id,
      league_id: selectedLeague.id,
      receiving_trade: receiving_trade,
      sending_trade: sending_trade
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_Trade_GetTrades,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  createTrade(user2Id: number, stocks1Id: number[], stocks2Id: number[]): void {
    const currentUser = this.verifyUserService.getCurrentUserValue();
    const selectedLeague = this.leagueService.getSelectedLeagueValue();
    guard(currentUser != null, "Current user is null!");
    guard(selectedLeague != null, "Selected League is null!");
    const data = {
      league_id: selectedLeague.id,
      user1_id: currentUser.id,
      user2_id: user2Id,
      stocks1_ids: stocks1Id,
      stocks2_ids: stocks2Id
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_Trade_CreateTrade,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  confirmTradeForUser(tradeId: number): void {
    const currentUser = this.verifyUserService.getCurrentUserValue();
    guard(currentUser != null, "Current user is null!");
    const data = {
      trade_id: tradeId,
      user_id: currentUser.id
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_Trade_ConfirmTrade,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

}
