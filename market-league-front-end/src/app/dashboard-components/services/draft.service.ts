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
import { DraftPick } from '../../models/websocket-responses/draft/draft-pick.model';

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

  // Track the last received draft update to handle re-connections
  private lastDraftUpdate: DraftUpdateResponse | null = null;

  // * Constructor
  
  // Routes Websocket Messages
  constructor(
    private portfolioService: PortfolioService,
    private leagueService: LeagueService,
    private verifyUserService: VerifyUserService,
    private webSocketService: WebSocketService,
  ) {
    // Monitor connection status for reconnection handling
    this.webSocketService.connectionStatus.subscribe((isConnected) => {
      if (isConnected && this.lastDraftUpdate) {
        // If we have a stored draft update, rebroadcast it on reconnection
        setTimeout(() => {
          devLog("Rebroadcasting last draft update after reconnection");
          // Add null check to prevent TypeScript error
          if (this.lastDraftUpdate) {
            this.currentDraftPlayerSubject.next(this.lastDraftUpdate);
          }
        }, 1000); // Small delay to ensure other services are ready
      }
    });

    this.webSocketService.getMessages().subscribe((message) => {
      switch (message.type) {

        case WebSocketMessageTypes.MessageType_LeaguePortfolio_GetLeaguePortfolioInfo:
          devLog("Received GetLeaguePortfolioInfo Response: " + JSON.stringify(message.data));
          this.handleGetLeaguePortfolioInfoResponse(message.data);
          break;

        case WebSocketMessageTypes.MessageType_LeaguePortfolio_DraftStock:
          devLog("Received DraftStock Response: " + JSON.stringify(message.data));
          this.handleDraftStockResponse(message.data);
          break;  

        case WebSocketMessageTypes.MessageType_League_QueueUp:
          devLog("Received QueueUp Response: " + JSON.stringify(message.data));
          this.handleQueueUpResponse(message.data);
          break;

        case WebSocketMessageTypes.MessageType_League_Portfolios:
          devLog("Received PlayerPortfolios Response: " + JSON.stringify(message.data));
          this.handleGetLeaguePortfoliosResponse(message.data);
          break;

        case WebSocketMessageTypes.MessageType_League_DraftUpdate:
          devLog("Received DraftUpdate Response: " + JSON.stringify(message.data));
          this.handleDraftUpdateResponse(message.data);
          break;

        case WebSocketMessageTypes.MessageType_League_DraftPick:
          devLog("Received DraftPick Response: " + JSON.stringify(message.data));
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
    devLog(`Got league portfolio with ${leaguePortfolio.stocks?.length || 0} stocks`);
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
    devLog(`Received ${playerPortfolios.length} player portfolios`);
    this.playerPortfoliosForLeagueSubject.next(playerPortfolios);
  }

  handleSuccessfulDraftUpdateResponseResponse(currentDraftPlayer: DraftUpdateResponse): void {
    devLog(`Current draft player updated: ${currentDraftPlayer.playerID}, Remaining time: ${currentDraftPlayer.remainingTime}s`);
    // Store the last draft update for reconnection handling
    this.lastDraftUpdate = currentDraftPlayer;
    this.currentDraftPlayerSubject.next(currentDraftPlayer);
  }

  handleSuccessfulDraftPickResponseResponse(draftPick: DraftPickResponse): void {
    devLog(`Draft pick processed: Player ${draftPick.player_id} picked stock ${draftPick.stock_id}`);
    
    // Get current league ID
    const selectedLeague = this.leagueService.getSelectedLeagueValue();
    if (!selectedLeague) {
      devLog("Cannot process draft pick: No league selected");
      return;
    }
    
    // Add pick to draft history
    this.addDraftPick(selectedLeague.id, {
      player_id: draftPick.player_id,
      stock_id: draftPick.stock_id,
      timestamp: new Date()
    });
    
    // Broadcast the draft pick
    this.draftPickSubject.next(draftPick);
    
    // Refresh data
    this.getLeaguePortfolioInfo();
    this.getAllPortfolios();
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
    devLog(`Sent request for league portfolio info for league ${selectedLeague.id}`);
  }

  draftStock(stockId: number): void {
    const selectedLeague = this.leagueService.getSelectedLeagueValue();
    guard(selectedLeague != null, "The selected League is null!");
    const currentUser = this.verifyUserService.getCurrentUserValue();
    guard(currentUser != null, "Current User is null!");

    devLog(`Sending draft stock request: Stock ID ${stockId}, User ID ${currentUser.id}, League ID ${selectedLeague.id}`);
    
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
    devLog(`Sent queue up request for player ${currentUser.id} in league ${selectedLeague.id}`);
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
    devLog(`Sent request for all portfolios in league ${selectedLeague.id}`);
  }

  // ** Storage **
  // Add a class property to track draft picks by league
  private draftPicksByLeague: Map<number, DraftPick[]> = new Map();

  // Modified saveDraftState to work with multiple leagues
  saveDraftState(leagueId: number, draftPicks: DraftPick[]): void {
    if (!leagueId) {
      devLog("Cannot save draft state: Invalid league ID");
      return;
    }
    
    // Update in-memory cache
    this.draftPicksByLeague.set(leagueId, [...draftPicks]);
    
    // Update localStorage
    const storageKey = `draftPicks_${leagueId}`;
    localStorage.setItem(storageKey, JSON.stringify(draftPicks));
    devLog(`Saved ${draftPicks.length} draft picks for league ${leagueId}`);
  }

  // Modified loadDraftState to work with multiple leagues
  loadDraftState(leagueId: number): DraftPick[] {
    if (!leagueId) {
      devLog("Cannot load draft state: Invalid league ID");
      return [];
    }
    
    // Check in-memory cache first
    if (this.draftPicksByLeague.has(leagueId)) {
      return [...this.draftPicksByLeague.get(leagueId)!]; // Return a copy
    }
    
    // Fall back to localStorage
    const storageKey = `draftPicks_${leagueId}`;
    const storedPicks = localStorage.getItem(storageKey);
    
    if (!storedPicks) {
      return [];
    }
    
    try {
      const picks: DraftPick[] = JSON.parse(storedPicks);
      // Update the in-memory cache
      this.draftPicksByLeague.set(leagueId, [...picks]);
      devLog(`Loaded ${picks.length} draft picks for league ${leagueId}`);
      return picks;
    } catch (e) {
      devLog(`Error parsing draft picks for league ${leagueId}: ${e}`);
      return [];
    }
  }

  // Add a method to add a single draft pick
  addDraftPick(leagueId: number, pick: DraftPick): void {
    if (!leagueId) {
      devLog("Cannot add draft pick: Invalid league ID");
      return;
    }
    
    // Get current picks
    const currentPicks = this.loadDraftState(leagueId);
    
    // Add new pick
    currentPicks.push({
      ...pick,
      timestamp: new Date() // Ensure timestamp is set
    });
    
    // Save updated picks
    this.saveDraftState(leagueId, currentPicks);
  }

  // Modified clearDraftState to work with multiple leagues
  clearDraftState(leagueId: number): void {
    if (!leagueId) {
      devLog("Cannot clear draft state: Invalid league ID");
      return;
    }
    
    // Clear from in-memory cache
    this.draftPicksByLeague.delete(leagueId);
    
    // Clear from localStorage
    const storageKey = `draftPicks_${leagueId}`;
    localStorage.removeItem(storageKey);
    
    devLog(`Cleared draft picks for league ${leagueId}`);
  }

  // Add method to reset all draft states (useful for testing/debugging)
  resetAllDraftStates(): void {
    // Clear in-memory cache
    this.draftPicksByLeague.clear();
    
    // Clear all draft-related items from localStorage
    const keysToRemove: string[] = [];
    
    // Find all draft-related keys
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key && key.startsWith('draftPicks_')) {
        keysToRemove.push(key);
      }
    }
    
    // Remove the keys
    keysToRemove.forEach(key => localStorage.removeItem(key));
    
    devLog(`Cleared all draft states (${keysToRemove.length} leagues)`);
  }
}