import { Injectable } from '@angular/core';
import { devLog } from '../../../environments/development/devlog';
import { Portfolio } from '../../models/portfolio.model';
import { WebSocketService } from './websocket.service';
import { WebSocketMessageTypes } from './websocket-message-types';
import { BehaviorSubject, Subject } from 'rxjs';
import { LeagueService } from './league.service';
import { VerifyUserService } from './verify-user.service';
import { StockHistoryEntry } from '../../models/stock-history-entry.model';
import { PortfolioPointsHistoryEntry } from '../../models/points-history-entry.model';

@Injectable({
  providedIn: 'root'
})
export class PortfolioService {

  // * Observables

  // User Portfolio
  private userPortfolioSubject = new BehaviorSubject<Portfolio | null>(this.getStoredPortfolio());
  userPortfolio$ = this.userPortfolioSubject.asObservable();
  // Stock History
  private stockHistoryListSubject = new Subject<StockHistoryEntry[]>();
  stockHistoryList$ = this.stockHistoryListSubject.asObservable();
  // Stock History
  private portfolioPointsHistoryListSubject = new Subject<PortfolioPointsHistoryEntry[]>();
  portfolioPointsHistoryList$ = this.portfolioPointsHistoryListSubject.asObservable();
  
  // * Constructor

  constructor(
    private webSocketService: WebSocketService,
    private leagueService: LeagueService,
    private verifyUserService: VerifyUserService,
  ) {
      this.webSocketService.getMessages().subscribe((message) => {
        switch (message.type) {
          case WebSocketMessageTypes.MessageType_Portfolio_LeaguePortfolio:
            devLog("Received GetUserPortfolio Response: " + message.data);
            this.handleGetUserPortfolioResponse(message.data);
            break;
          case WebSocketMessageTypes.MessageType_Portfolio_GetStocksValueChange:
            devLog("Received GetStocksValueChange Response: " + message.data);
            this.handleGetStocksValueChange(message.data);
            break;
          case WebSocketMessageTypes.MessageType_Portfolio_GetPortfolioPointsHistory:
            devLog("Received GetPortfolioPointsHistory Response: " + message.data);
            this.handleGetPortfolioPointsHistory(message.data);
            break;
          default:
      }
    });
  }

  // * Websocket Response Handler Functions
  
  handleGetUserPortfolioResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulGetUserPortfolioResponse(responseData);
  }

  handleGetStocksValueChange(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleGetStocksValueChangeResponse(responseData);
  }

  handleGetPortfolioPointsHistory(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleGetPortfolioPointsHistoryResponse(responseData);
  }

  // * Helper Functions to Websocket Responses

  handleSuccessfulGetUserPortfolioResponse(portfolio: Portfolio): void {
    this.userPortfolioSubject.next(portfolio);
  }

  handleGetStocksValueChangeResponse(stockHistoryList: StockHistoryEntry[]): void {
    this.stockHistoryListSubject.next(stockHistoryList)
  }

  handleGetPortfolioPointsHistoryResponse(portfolioPointsHistoryList: PortfolioPointsHistoryEntry[]): void {
    this.portfolioPointsHistoryListSubject.next(portfolioPointsHistoryList)
  }

  // * Websocket Call Functions

  // Method to fetch the user's portfolio for the selected league
  getCurrentUserPortfolio(): void {
    const selectedLeague = this.leagueService.getSelectedLeagueValue();
    const currentUser = this.verifyUserService.getCurrentUserValue();
    if (!selectedLeague || !currentUser) {
      return
    }
    const data = {
      user_id: currentUser.id,
      league_id: selectedLeague.id
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_Portfolio_LeaguePortfolio,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  // Method to get the change in the stock values in a portfolio
  getStocksValueChange(): void {
    this.userPortfolio$.subscribe((portfolio) => {
      if (!portfolio) {
        console.error("Unable to fetch portfolio_id!");
        return 
      }
      const data = {
        portfolio_id: portfolio.id
      };
      const websocketMessage = {
        type: WebSocketMessageTypes.MessageType_Portfolio_GetStocksValueChange,
        data: data
      };
      this.webSocketService.sendMessage(websocketMessage);
    });
  }

  // Method to get the points history of a portfolio
  getPortfolioPointsHistory(): void {
    this.userPortfolio$.subscribe((portfolio) => {
      if (!portfolio) {
        console.error("Unable to fetch portfolio_id!");
        return 
      }
      const data = {
        portfolio_id: portfolio.id
      };
      const websocketMessage = {
        type: WebSocketMessageTypes.MessageType_Portfolio_GetPortfolioPointsHistory,
        data: data
      };
      this.webSocketService.sendMessage(websocketMessage);
    });
  }

  // * Setters

  // Method to set the selected league
  setSelectedPortfolio(portfolio: Portfolio | null): void {
    devLog("Selected Portfolio: ", portfolio);
    this.userPortfolioSubject.next(portfolio); // Set the selected league as the full League object
    if (portfolio) {
      // Store the entire league object as a JSON string in localStorage
      localStorage.setItem('selectedPortfolio', JSON.stringify(portfolio)); 
    } else {
      localStorage.removeItem('selectedPortfolio');
    }
  }

  // * Getter Functions

  // Retrieve the stored league from localStorage (if it exists)
  getStoredPortfolio(): Portfolio | null {
    const storedPortfolio = localStorage.getItem('selectedPortfolio');
    
    // Check if storedLeague is a valid JSON
    if (storedPortfolio) {
      try {
        return JSON.parse(storedPortfolio) as Portfolio;
      } catch (e) {
        console.error("Error parsing stored league JSON:", e);
        localStorage.removeItem('selectedPortfolio'); // Clean up invalid entry
        return null;
      }
    }

    return null;
  }

}

