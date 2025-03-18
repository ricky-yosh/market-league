import { Injectable } from '@angular/core';
import { BehaviorSubject, Subject } from 'rxjs';
import { League } from '../../models/league.model';
import { devLog } from '../../../environments/development/devlog';
import { WebSocketService } from './websocket.service';
import { WebSocketMessageTypes } from './websocket-message-types';
import { CreateLeagueResponse } from '../../models/websocket-responses/league/create-league-response.model';
import { RemoveLeagueResponse } from '../../models/websocket-responses/league/remove-league-response.model';
import { AddUserToLeagueResponse } from '../../models/websocket-responses/league/add-user-to-league-response.model';
import { LeaderboardUser } from '../../models/websocket-responses/league/leaderboard-user';
import { User } from '../../models/user.model';
import { VerifyUserService } from './verify-user.service';

@Injectable({
  providedIn: 'root'
})
export class LeagueService {
  
  // * Observables

  // Selected League
  private selectedLeagueSource = new BehaviorSubject<League | null>(this.getStoredLeague());
  selectedLeague$ = this.selectedLeagueSource.asObservable();
  // User Leagues List
  private userLeaguesSubject = new Subject<League[]>();
  userLeagues$ = this.userLeaguesSubject.asObservable();
  // League Members
  private leagueMembersSubject = new Subject<User[]>();
  leagueMembers$ = this.leagueMembersSubject.asObservable();
  // Leaderboard
  private leaderboardSubject = new Subject<LeaderboardUser[]>();
  leaderboard$ = this.leaderboardSubject.asObservable();

  // * Constructor

  // Routes Websocket Messages
  constructor(
    private webSocketService: WebSocketService,
    private verifyUserService: VerifyUserService
  ) {
    this.webSocketService.getMessages().subscribe((message) => {
      switch (message.type) {
        case WebSocketMessageTypes.MessageType_League_CreateLeague:
          devLog("Received CreateLeague Response: " + message.data);
          this.handleCreateLeagueResponse(message.data);
          break;
        case WebSocketMessageTypes.MessageType_League_RemoveLeague:
          devLog("Received RemoveLeague Response: " + message.data);
          this.handleRemoveLeagueResponse(message.data);
          break;
        case WebSocketMessageTypes.MessageType_League_AddUserToLeague:
          devLog("Received AddUserToLeague Response: " + message.data);
          this.handleAddUserToLeagueResponse(message.data);
          break;
        case WebSocketMessageTypes.MessageType_League_GetDetails:
          devLog("Received GetLeagueDetails Response: " + message.data);
          this.handleGetLeagueDetailsResponse(message.data);
          break;
        case WebSocketMessageTypes.MessageType_League_GetLeaderboard:
          devLog("Received GetLeaderboard Response: " + message.data);
          this.handleGetLeaderboardResponse(message.data);
          break;
        case WebSocketMessageTypes.MessageType_User_UserLeagues:
          devLog("Received GetUserLeagues Response: " + message.data);
          this.handleGetUserLeaguesResponse(message.data);
          break;
        default:
          // devLog("League Service unable to route Websocket Message properly! " + message.data);
      }
    });
  }

  // * Websocket Response Handler Functions

  handleCreateLeagueResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulCreateLeagueResponse(responseData);
  }

  handleRemoveLeagueResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulRemoveLeagueResponse(responseData);
  }

  handleAddUserToLeagueResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulAddUserToLeagueResponse(responseData);
  }

  handleGetLeagueDetailsResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulGetLeagueDetailsResponse(responseData);
  }
  
  handleGetLeaderboardResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulGetLeaderboardResponse(responseData);
  }

  handleGetUserLeaguesResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulGetUserLeaguesResponse(responseData);
  }

  // * Helper Functions to Websocket Responses

  handleSuccessfulCreateLeagueResponse(response: CreateLeagueResponse): void {
    devLog("Added League: " + response.league.id)
    this.getUserLeagues();
  }

  handleSuccessfulRemoveLeagueResponse(response: RemoveLeagueResponse): void {
    devLog("Removed League: " + response.message)
    this.getUserLeagues();
  }

  handleSuccessfulAddUserToLeagueResponse(response: AddUserToLeagueResponse): void {
    const users = response.league.users;
    if (users) {
      this.leagueMembersSubject.next(users);
    }
  }

  handleSuccessfulGetLeagueDetailsResponse(response: League): void {
    const users = response.users;
    if (users) {
      this.leagueMembersSubject.next(users);
    }
    this.selectedLeagueSource.next(response)
  }

  handleSuccessfulGetLeaderboardResponse(response: LeaderboardUser[]): void {
    const leaderboard = response;
    this.leaderboardSubject.next(leaderboard)
  }

  handleSuccessfulGetUserLeaguesResponse(response: League[]): void {
    const leagues = response;
    this.userLeaguesSubject.next(leagues);
  }

  // * Websocket Call Functions

  // Create League
  createLeague(leagueName: string, ownerUser: number, endDate: string): void {
    const data = {
      league_name: leagueName,
      owner_user: ownerUser,
      end_date: endDate
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_League_CreateLeague,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  // Remove League
  removeLeague(leagueId: number): void {
    const data = {
      league_id: leagueId
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_League_RemoveLeague,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  // Get league details using the league ID
  getLeagueDetails(): void {
    const selectedLeague = this.getSelectedLeagueValue();
    if (!selectedLeague) {
      return
    }
    const data = {
      league_id: selectedLeague.id
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_League_GetDetails,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  // Get members of a league using the league ID
  getLeagueMembers(): void {
    const selectedLeague = this.getSelectedLeagueValue();
    if (!selectedLeague) {
      return
    }
    const data = {
      league_id: selectedLeague.id
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_League_GetDetails,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  // Method to get members of a league using the league ID
  getLeagueLeaderboard(leagueId: number): void {
    const data = {
      league_id: leagueId
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_League_GetLeaderboard,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  // Get the user's leagues using their user ID
  getUserLeagues(): void {
    const currentUser = this.verifyUserService.getCurrentUserValue()
    if (!currentUser) {
      return
    }
    const data = {
      user_id: currentUser.id
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_User_UserLeagues,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  // * Setters

  // Method to set the selected league
  setSelectedLeague(league: League | null): void {
    devLog("Selected League: ", league);
    this.selectedLeagueSource.next(league); // Set the selected league as the full League object
    if (league) {
      // Store the entire league object as a JSON string in localStorage
      localStorage.setItem('selectedLeague', JSON.stringify(league)); 
    } else {
      localStorage.removeItem('selectedLeague');
    }
  }

  // * Getter Functions
  
  // Get current user value
  getSelectedLeagueValue(): League | null {
    return this.selectedLeagueSource.value;
  }

  // Retrieve the stored league from localStorage (if it exists)
  getStoredLeague(): League | null {
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