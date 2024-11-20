import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { LeagueService } from '../services/league.service';
import { User } from '../../models/user.model';
import { VerifyUserService } from '../../user-verification/verify-user.service';
import { League } from '../../models/league.model';
import { devLog } from '../../../environments/development/devlog';
import { Stock } from '../../models/stock.model';
import { EMPTY, Observable, catchError, of, tap } from 'rxjs';
import { Trade } from '../../models/trade.model';
import { guard } from '../../utils/guard';

@Component({
  selector: 'app-league-trades',
  standalone: true,
  imports: [FormsModule, CommonModule],
  templateUrl: './league-trades.component.html',
  styleUrl: './league-trades.component.scss'
})
export class LeagueTradesComponent {
  currentUsersStocks: Stock[] = [];
  leagueUsers: User[] = [];
  selectedUserStocks: Stock[] = [];
  currentUser: User | null = null;
  currentLeague: League | null = null;
  currentUsersTrades: Trade[] | null = null;

  formInput: { user2: User | null; stocks1: Stock[]; stocks2: Stock[] } = {
    user2: null,
    stocks1: [],
    stocks2: []
  };

  constructor(
    private leagueService: LeagueService,
    private userService: VerifyUserService
  ) {}

  ngOnInit() {
    this.populateLeagueUsers();
    this.loadUser();
    this.getCurrentLeague();
  }

  onSubmit() {
    const league_id = this.currentLeague?.id;
    const user1_id = this.currentUser?.id;
    const user2_id = this.formInput.user2?.id;
    const stocks1_ids = this.formInput.stocks1.map((stock: Stock) => stock.id);
    const stocks2_ids = this.formInput.stocks2.map((stock: Stock) => stock.id);

    if (league_id &&
      user1_id &&
      user2_id &&
      stocks1_ids.length > 0 &&
      stocks2_ids.length > 0) {
      
      this.leagueService.createTrade(league_id, user1_id, user2_id, stocks1_ids, stocks2_ids).subscribe(response => {
        devLog('Trade successfully created:', response);
        alert('Trade successfully created!');
        this.resetForm();
      });
    } else {
      alert('Please complete the form before submitting.');
    }
  }

  resetForm() {
    this.formInput = {
      user2: null,
      stocks1: [],
      stocks2: []
    };
    this.selectedUserStocks = [];
  }

  toggleStockSelection(stockList: Stock[], stock: Stock) {
    const index = stockList.indexOf(stock);
    if (index === -1) {
      stockList.push(stock);
    } else {
      stockList.splice(index, 1);
    }
  }

  onUserSelectionChange(selectedUser: User | null) {
  
    if (!selectedUser || !this.currentLeague) {
      return; // if user or league is null, return early
    }
  
    const selectedUserId = selectedUser.id;
    const selectedLeagueId = this.currentLeague.id;
  
    devLog("selectedUserId & selectedLeagueId: ", selectedUserId, selectedLeagueId);
    
    // Fetch user's portfolio for the selected league
    this.leagueService.getUserPortfolio(selectedUserId, selectedLeagueId).subscribe(portfolio => {
      devLog("selectedUserId's Portfolio: ", portfolio);
      this.selectedUserStocks = portfolio.stocks;
    });
  }

  populateLeagueUsers() {
    // Fetching the selected league from the service.
    const selectedLeague = this.leagueService.getStoredLeague();
    if (selectedLeague) {
      const leagueId = selectedLeague.id;
      this.leagueService.getLeagueMembers(leagueId).subscribe(users => {
        this.leagueUsers = users;
      });
    } else {
      console.warn('No league selected');
    }
  }

  private loadUser(): void {
    this.userService.getUserFromToken().subscribe({
      next: (user: User) => {
        devLog('User fetched successfully:', user);
        this.currentUser = user;
        this.getCurrentUsersPortfolio(); // load portfolio for logged in user
        this.loadUserTrades(user.id, this.currentLeague?.id ? this.currentLeague?.id : null).subscribe(userTrades => {
          this.currentUsersTrades = userTrades;
        })
      },
      error: (error) => {
        devLog('Failed to fetch user from token:', error);
      }
    });
  }

  private getCurrentLeague(): void {
    this.currentLeague = this.leagueService.getStoredLeague();
  }

  private getCurrentUsersPortfolio() {

    devLog("Current User: ", this.currentUser)
    if (!this.currentUser || !this.currentLeague) {
      return; // if user or league is null, return early
    }
  
    const currentUserId = this.currentUser.id;
    const selectedLeagueId = this.currentLeague.id;
  
    devLog("currentUserId & selectedLeagueId: ", currentUserId, selectedLeagueId);
    
    // Fetch user's portfolio for the selected league
    this.leagueService.getUserPortfolio(currentUserId, selectedLeagueId).subscribe(portfolio => {
      devLog("currentUserId's Portfolio: ", portfolio);
      this.currentUsersStocks = portfolio.stocks;
      this.removeCurrentPlayerFromTradeList();
    });
  }

  private removeCurrentPlayerFromTradeList() {
    const currentUserId = this.currentUser?.id
    if (currentUserId != null) {
      this.leagueUsers = this.leagueUsers.filter(user => user.id !== currentUserId);
    }
  }

  // Load the user's trades for a specific league
  private loadUserTrades(userId: number, leagueId: number | null): Observable<Trade[]> {
    guard(leagueId != null, "LeagueId is Null");

    const receivingTrade: boolean = true;
    const sendingTrade: boolean = false;
    return this.leagueService.getTrades(userId, leagueId, receivingTrade, sendingTrade).pipe(
      tap((response) => {
        devLog('User trades fetched successfully:', response);
      }),
      catchError((error) => {
        console.error('Failed to fetch user trades:', error);
        return of([]);
      })
    );
  }

  confirmTrade(tradeId: number): void {
    let currentUserId = this.currentUser?.id
    guard(tradeId != null, "tradeId is null");
    guard(currentUserId != null, "UserID is null!")

    this.leagueService.confirmTradeForUser(tradeId, currentUserId).subscribe(response => {
      devLog('Trade successfully confirmed:', response);
      alert('Trade successfully confirmed!');
    });
  }

}