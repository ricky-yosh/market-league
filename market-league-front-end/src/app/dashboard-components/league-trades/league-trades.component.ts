import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { LeagueService } from '../services/league.service';
import { User } from '../../models/user.model';
import { VerifyUserService } from '../../user-verification/verify-user.service';
import { League } from '../../models/league.model';
import { devLog } from '../../../environments/development/devlog';

@Component({
  selector: 'app-league-trades',
  standalone: true,
  imports: [FormsModule, CommonModule],
  templateUrl: './league-trades.component.html',
  styleUrl: './league-trades.component.scss'
})
export class LeagueTradesComponent {
  availableStocks: string[] = ['AAPL', 'TSLA', 'MSFT', 'NVDA', 'GOOGL', 'AMZN', 'NFLX', 'META'];
  leagueUsers: User[] = [];
  selectedUserStocks: string[] = [];
  currentUser: User | null = null
  currentLeague: League | null = null

  trade: { user2: User | null; stocks1: string[]; stocks2: string[] } = {
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
    if (this.trade.user2 && this.trade.stocks1.length > 0 && this.trade.stocks2.length > 0) {
      // For simplicity, we're just logging the trade instead of making an actual HTTP request.
      devLog('Trade details:', this.trade);
      alert('Trade successfully created!');
      this.resetForm();
    } else {
      alert('Please complete the form before submitting.');
    }
  }

  resetForm() {
    this.trade = {
      user2: null,
      stocks1: [],
      stocks2: []
    };
    this.selectedUserStocks = [];
  }

  toggleStockSelection(stockList: string[], stock: string) {
    const index = stockList.indexOf(stock);
    if (index === -1) {
      stockList.push(stock);
    } else {
      stockList.splice(index, 1);
    }
  }

  onUserSelectionChange(selectedUser: User | null) {
    this.selectedUserStocks = []; // empty list to clear earlier portfolio
  
    if (!selectedUser || !this.currentLeague) {
      return; // if user or league is null, return early
    }
  
    const selectedUserId = selectedUser.id;
    const selectedLeagueId = this.currentLeague.id;
  
    devLog("selectedUserId & selectedLeagueId: ", selectedUserId, selectedLeagueId);
    
    // Fetch user's portfolio for the selected league
    this.leagueService.getUserPortfolio(selectedUserId, selectedLeagueId).subscribe(portfolio => {
      devLog("selectedUserId's Portfolio: ", portfolio);
      this.selectedUserStocks = portfolio.stocks.map(stock => stock.ticker_symbol);
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
      },
      error: (error) => {
        devLog('Failed to fetch user from token:', error);
      }
    });
  }

  private getCurrentLeague(): void {
    this.currentLeague = this.leagueService.getStoredLeague();
  }

}