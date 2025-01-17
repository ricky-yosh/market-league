import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { LeagueService } from '../services/league.service';
import { User } from '../../models/user.model';
import { VerifyUserService } from '../services/verify-user.service';
import { League } from '../../models/league.model';
import { devLog } from '../../../environments/development/devlog';
import { Stock } from '../../models/stock.model';
import { Observable, Subscription, catchError, map, of, tap } from 'rxjs';
import { Trade } from '../../models/trade.model';
import { guard } from '../../utils/guard';
import { PortfolioService } from '../services/portfolio.service';
import { TradeService } from '../services/trade.service';

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
  currentUserTrades: Trade[] | null = null;

  formInput: { user2: User | null; stocks1: Stock[]; stocks2: Stock[] } = {
    user2: null,
    stocks1: [],
    stocks2: []
  };

  private subscription!: Subscription;

  constructor(
    private portfolioService: PortfolioService,
    private tradeService: TradeService,
    private leagueService: LeagueService,
    private verifyUserService: VerifyUserService,
  ) {}

  ngOnInit() {

    // * Subscribe to the observables to listen for changes

    // User's Portfolio
    this.subscription = this.portfolioService.userPortfolio$.subscribe((portfolio) => {
      this.currentUsersStocks = portfolio.stocks;
    });
    this.subscription = this.tradeService.leagueTrades$.subscribe((leagueTrades) => {
      const filteredCurrentUserFromTrades = leagueTrades.filter(trade => trade.user2 == this.verifyUserService.getCurrentUserValue());
      this.currentUserTrades = filteredCurrentUserFromTrades;
    });
    
    // * Get Starting Values for Dashboard
    
    // Get all the users in the league
    this.loadLeagueUsers();
  }

  ngOnDestroy(): void {
    // Unsubscribe to avoid memory leaks
    this.subscription.unsubscribe();
  }

  // * Form Functions

  onSubmit() {
    const user2_id = this.formInput.user2?.id;
    const stocks1_ids = this.formInput.stocks1.map((stock: Stock) => stock.id);
    const stocks2_ids = this.formInput.stocks2.map((stock: Stock) => stock.id);

    if (user2_id &&
        stocks1_ids.length > 0 &&
        stocks2_ids.length > 0) {
      
      this.tradeService.createTrade(user2_id, stocks1_ids, stocks2_ids)
      this.resetForm();
    } else {
      alert('Please complete the form before submitting.');
    }
  }

  toggleStockSelection(stockList: Stock[], stock: Stock) {
    const index = stockList.indexOf(stock);
    if (index === -1) {
      stockList.push(stock);
    } else {
      stockList.splice(index, 1);
    }
  }

  onUserSelectionChange(): void {
    // Fetch user's portfolio for the selected league
    this.portfolioService.getCurrentUserPortfolio();
  }

  confirmTrade(tradeId: number): void {
    guard(tradeId != null, "tradeId is null");

    this.tradeService.confirmTradeForUser(tradeId);
    alert('Trade successfully confirmed!');
  }
  
  // * Helper Functions

  private resetForm() {
    this.formInput = {
      user2: null,
      stocks1: [],
      stocks2: []
    };
    this.selectedUserStocks = [];
  }

  // Load the user's trades for a specific league
  private loadLeagueUsers(): void {
    const selectedLeague = this.leagueService.getSelectedLeagueValue();
    const leagueUsers = selectedLeague?.users;
    guard(leagueUsers != null, "Selected League Users list is null!");
    
    this.leagueUsers = leagueUsers;
  }

}