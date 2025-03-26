import { CommonModule } from '@angular/common';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { LeagueService } from '../services/league.service';
import { User } from '../../models/user.model';
import { VerifyUserService } from '../services/verify-user.service';
import { Stock } from '../../models/stock.model';
import { Subscription } from 'rxjs';
import { Trade } from '../../models/trade.model';
import { guard } from '../../utils/guard';
import { PortfolioService } from '../services/portfolio.service';
import { TradeService } from '../services/trade.service';
import { DraftService } from '../services/draft.service';
import { Portfolio } from '../../models/portfolio.model';

@Component({
  selector: 'app-league-trades',
  standalone: true,
  imports: [FormsModule, CommonModule],
  templateUrl: './league-trades.component.html',
  styleUrl: './league-trades.component.scss'
})
export class LeagueTradesComponent implements OnInit, OnDestroy {
  currentUsersStocks: Stock[] = [];
  leagueUsers: User[] = [];
  selectedUserStocks: Stock[] = [];
  currentUserTrades: Trade[] | null = null;
  pendingTrades: Trade[] = [];
  completedTrades: Trade[] = [];
  allPortfolios: Portfolio[] = [];

  formInput: { user2: User | null; stocks1: Stock[]; stocks2: Stock[] } = {
    user2: null,
    stocks1: [],
    stocks2: []
  };

  private subscriptions: Subscription[] = [];

  constructor(
    private portfolioService: PortfolioService,
    private tradeService: TradeService,
    private leagueService: LeagueService,
    private verifyUserService: VerifyUserService,
    private draftService: DraftService,
  ) {}

  ngOnInit() {
    // Load the user's portfolio and league users first
    this.loadLeagueUsers();
    this.portfolioService.getCurrentUserPortfolio();
    
    // Subscribe to the portfolios
    this.subscriptions.push(
      this.draftService.playerPortfoliosForLeague$.subscribe((portfolios) => {
        this.allPortfolios = portfolios;
      })
    );

    // User's Portfolio
    this.subscriptions.push(
      this.portfolioService.userPortfolio$.subscribe((portfolio) => {
        this.currentUsersStocks = portfolio ? portfolio.stocks : [];
      })
    );
    
    // Subscribe to trades and then explicitly request them
    this.subscriptions.push(
      this.tradeService.leagueTrades$.subscribe((leagueTrades) => {
        const currentUser = this.verifyUserService.getCurrentUserValue();
        
        if (currentUser && leagueTrades) {
          const userTrades = leagueTrades.filter(trade => 
            trade.user2 && trade.user2.id === currentUser.id
          );
          
          // Separate into pending and completed trades
          this.pendingTrades = userTrades.filter(trade => 
            !trade.user1_confirmed || !trade.user2_confirmed
          );
          
          this.completedTrades = userTrades.filter(trade => 
            trade.user1_confirmed && trade.user2_confirmed
          );
          
          this.currentUserTrades = userTrades; // Keep this for backward compatibility
        } else {
          this.pendingTrades = [];
          this.completedTrades = [];
          this.currentUserTrades = [];
          
          if (!currentUser) {
            console.warn('Current user is null, no trades to display');
          }
          if (!leagueTrades) {
            console.warn('League trades is null, no trades to display');
          }
        }
      })
    );
    
    // Get all portfolios in the league
    this.draftService.getAllPortfolios();
    
    // Now fetch trades after setting up subscriptions and other data
    this.tradeService.getTrades(true, false);
  }

  ngOnDestroy(): void {
    // Unsubscribe from all subscriptions
    this.subscriptions.forEach(sub => sub.unsubscribe());
  }

  // * Form Functions

  onSubmit() {
    const user2_id = this.formInput.user2?.id;
    const stocks1_ids = this.formInput.stocks1.map((stock: Stock) => stock.id);
    const stocks2_ids = this.formInput.stocks2.map((stock: Stock) => stock.id);

    if (user2_id &&
        stocks1_ids.length > 0 &&
        stocks2_ids.length > 0) {
      
      this.tradeService.createTrade(user2_id, stocks1_ids, stocks2_ids);
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
    // When user selection changes, find the selected user's portfolio
    if (this.formInput.user2) {
      const selectedUserPortfolio = this.allPortfolios.find(
        portfolio => portfolio.user_id === this.formInput.user2?.id
      );
      this.selectedUserStocks = selectedUserPortfolio ? selectedUserPortfolio.stocks : [];
    } else {
      this.selectedUserStocks = [];
    }
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
    
    const currentUser = this.verifyUserService.getCurrentUserValue();
    
    // Filter out the current user from the list if currentUser exists
    if (currentUser) {
      this.leagueUsers = leagueUsers.filter(user => user.id !== currentUser.id);
    } else {
      this.leagueUsers = leagueUsers;
      console.warn('Current user is null, showing all league users');
    }
  }
  
  // Manually reload trades - can be called by a refresh button if needed
  refreshTrades(): void {
    this.tradeService.getTrades(true, false);
  }
}