import { Injectable } from '@angular/core';
import { devLog } from '../../../environments/development/devlog';
import { Portfolio } from '../../models/portfolio.model';
import { WebSocketService } from './websocket.service';
import { WebSocketMessageTypes } from './websocket-message-types';
import { Subject } from 'rxjs';
import { LeagueService } from './league.service';
import { VerifyUserService } from './verify-user.service';

@Injectable({
  providedIn: 'root'
})
export class PortfolioService {

  // * Observables

  // User Portfolio
  private userPortfolioSubject = new Subject<Portfolio>();
  userPortfolio$ = this.userPortfolioSubject.asObservable();

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
          default:
            // devLog("Portfolio Service unable to route Websocket Message properly! " + message.data);
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

  // * Helper Functions to Websocket Responses

  handleSuccessfulGetUserPortfolioResponse(portfolio: Portfolio): void {
    this.userPortfolioSubject.next(portfolio);
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

}
