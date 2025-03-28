import { Injectable } from '@angular/core';
import { LeaguePortfolio } from '../../models/league-portfolio.model';
import { Subject } from 'rxjs';
import { devLog } from '../../../environments/development/devlog';
import { WebSocketService } from './websocket.service';
import { WebSocketMessageTypes } from './websocket-message-types';
import { DraftStockResponse } from '../../models/websocket-responses/draft/draft-stock-response.model';
import { PortfolioService } from './portfolio.service';
import { guard } from '../../utils/guard';
import { LeagueService } from './league.service';
import { VerifyUserService } from './verify-user.service';
import { Portfolio } from '../../models/portfolio.model';
import { DraftUpdateResponse } from '../../models/websocket-responses/draft/draft-update-response.model';
import { DraftPickResponse } from '../../models/websocket-responses/draft/draft-pick-response.model';

@Injectable({
  providedIn: 'root'
})
export class DraftService {

  // * Observables

  // League Portfolio List
  private leaguePortfolioSubject = new Subject<LeaguePortfolio>();
  leaguePortfolio$ = this.leaguePortfolioSubject.asObservable();

  private playerPortfoliosForLeagueSubject = new Subject<Portfolio[]>();
  playerPortfoliosForLeague$ = this.playerPortfoliosForLeagueSubject.asObservable();

  private currentDraftPlayerSubject = new Subject<DraftUpdateResponse>();
  currentDraftPlayer$ = this.currentDraftPlayerSubject.asObservable();
  
  private draftPickSubject = new Subject<DraftPickResponse>();
  draftPick$ = this.draftPickSubject.asObservable();

  // * Constructor
  
  // Routes Websocket Messages
  constructor(
    private portfolioService: PortfolioService,
    private leagueService: LeagueService,
    private verifyUserService: VerifyUserService,
    private webSocketService: WebSocketService,
  ) {
    this.webSocketService.getMessages().subscribe((message) => {
      switch (message.type) {

        case WebSocketMessageTypes.MessageType_LeaguePortfolio_GetLeaguePortfolioInfo:
          devLog("Received GetLeaguePortfolioInfo Response: " + message.data);
          this.handleGetLeaguePortfolioInfoResponse(message.data);
          break;

        case WebSocketMessageTypes.MessageType_LeaguePortfolio_DraftStock:
          devLog("Received DraftStock Response: " + message.data);
          this.handleDraftStockResponse(message.data);
          break;  

        case WebSocketMessageTypes.MessageType_League_QueueUp:
          devLog("Received QueueUp Response: " + message.data);
          this.handleQueueUpResponse(message.data);
          break;

        case WebSocketMessageTypes.MessageType_League_Portfolios:
          devLog("Received PlayerPortfolios Response: " + message.data);
          this.handleGetLeaguePortfoliosResponse(message.data);
          break;

        case WebSocketMessageTypes.MessageType_League_DraftUpdate:
          devLog("Received DraftUpdate Response: " + message.data);
          this.handleDraftUpdateResponse(message.data);
          break;

        case WebSocketMessageTypes.MessageType_League_DraftPick:
          devLog("Received DraftPick Response: " + message.data);
          this.handleDraftPickResponse(message.data);
          break;

        default:
          // devLog("League Service unable to route Websocket Message properly! " + message.data);
      }
    });
  }

  // * Websocket Response Handler Functions

  handleGetLeaguePortfolioInfoResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulGetLeaguePortfolioInfoResponse(responseData);
  }

  handleDraftStockResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulDraftStockResponse(responseData);
  }

  handleQueueUpResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulQueueUpResponse(responseData);
  }

  handleGetLeaguePortfoliosResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulGetLeaguePortfoliosResponse(responseData);
  }

  handleDraftUpdateResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulDraftUpdateResponseResponse(responseData);
  }

  handleDraftPickResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulDraftPickResponseResponse(responseData);
  }

  // * Helper Functions to Websocket Responses

  handleSuccessfulGetLeaguePortfolioInfoResponse(leaguePortfolio: LeaguePortfolio): void {
    this.leaguePortfolioSubject.next(leaguePortfolio);
  }

  handleSuccessfulDraftStockResponse(response: DraftStockResponse): void {
    devLog("Drafted Stock: " + response.message)
    this.portfolioService.getCurrentUserPortfolio();
  }

  handleSuccessfulQueueUpResponse(response: any): void {
    devLog("Player queued up: " + response.message)
  }

  handleSuccessfulGetLeaguePortfoliosResponse(playerPortfolios: Portfolio[]): void {
    this.playerPortfoliosForLeagueSubject.next(playerPortfolios);
  }

  handleSuccessfulDraftUpdateResponseResponse(currentDraftPlayer: DraftUpdateResponse): void {
    this.currentDraftPlayerSubject.next(currentDraftPlayer);
  }

  handleSuccessfulDraftPickResponseResponse(draftPick: DraftPickResponse): void {
    this.draftPickSubject.next(draftPick);
  }

  // * Websocket Call Functions

  getLeaguePortfolioInfo(): void {
    const selectedLeague = this.leagueService.getSelectedLeagueValue();
    guard(selectedLeague, "The selected League is null!");

    const data = {
      league_id: selectedLeague.id
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_LeaguePortfolio_GetLeaguePortfolioInfo,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  draftStock(stockId: number): void {
    const selectedLeague = this.leagueService.getSelectedLeagueValue();
    guard(selectedLeague != null, "The selected League is null!");
    const currentUser = this.verifyUserService.getCurrentUserValue();
    guard(currentUser != null, "Current User is null!");

    const data = {
      league_id: selectedLeague.id,
      user_id: currentUser.id,
      stock_id: stockId
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_LeaguePortfolio_DraftStock,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  queuePlayerForDraft(): void {
    const selectedLeague = this.leagueService.getSelectedLeagueValue();
    guard(selectedLeague != null, "The selected League is null!");
    const currentUser = this.verifyUserService.getCurrentUserValue();
    guard(currentUser != null, "Current User is null!");

    const data = {
      league_id: selectedLeague.id,
      player_id: currentUser.id
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_League_QueueUp,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  getAllPortfolios(): void {
    const selectedLeague = this.leagueService.getSelectedLeagueValue();
    if (!selectedLeague) {
      devLog("No league selected!");
      return;
    }

    const data = {
      league_id: selectedLeague.id
    };

    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_League_Portfolios,
      data: data
    };

    this.webSocketService.sendMessage(websocketMessage);
  }

}
