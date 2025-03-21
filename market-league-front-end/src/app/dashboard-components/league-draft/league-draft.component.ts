import { CommonModule } from '@angular/common';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { Stock } from '../../models/stock.model';
import { Portfolio } from '../../models/portfolio.model';
import { Router } from '@angular/router';
import { StockService } from '../services/stock.service';
import { PortfolioService } from '../services/portfolio.service';
import { DraftService } from '../services/draft.service';
import { Subscription } from 'rxjs';
import { LeagueState } from '../../models/league-state.model';
import { LeagueService } from '../services/league.service';
import { WebSocketService } from '../services/websocket.service';
import { WebSocketMessageTypes } from '../services/websocket-message-types';
import { DraftUpdateResponse } from '../../models/websocket-responses/draft/draft-update-response.model';
import { DraftPickResponse } from '../../models/websocket-responses/draft/draft-pick-response.model';
import { VerifyUserService } from '../services/verify-user.service';
import { User } from '../../models/user.model';

interface DraftPick {
  player_id: number;
  stock_id: number;
  timestamp: Date;
}

@Component({
  selector: 'app-league-draft',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './league-draft.component.html',
  styleUrl: './league-draft.component.scss'
})
export class LeagueDraftComponent implements OnInit, OnDestroy {
  // Data for UI
  leagueStocks: Stock[] = [];
  userPortfolio: Portfolio | null = null;
  leaguePortfolios: Portfolio[] = [];
  
  // Draft state
  currentPlayerID: number = 0;
  currentUserID: number = 0;
  remainingTime: number = 0;
  maxDraftTime: number = 30; // 30 seconds per turn
  draftPicks: DraftPick[] = [];
  
  // User data
  players: Map<number, User> = new Map();
  
  // Stock data for quick lookup
  stocksMap: Map<number, Stock> = new Map();
  
  // Subscriptions
  private subscriptions: Subscription[] = [];
  
  // Timer interval
  private timerInterval: any;

  constructor(
    private router: Router,
    private draftService: DraftService,
    private portfolioService: PortfolioService,
    private stockService: StockService,
    private leagueService: LeagueService,
    private verifyUserService: VerifyUserService
  ) {}

  ngOnInit(): void {
    // Get current user ID
    const currentUser = this.verifyUserService.getCurrentUserValue();
    if (currentUser) {
      this.currentUserID = currentUser.id;
    }

    // Subscribe to selected league changes
    this.subscriptions.push(
      this.leagueService.selectedLeague$.subscribe((league) => {
        if (!league) return;
        
        switch(league.league_state) {
          case LeagueState.PreDraft:
            this.redirectToDraftQueue();
            break;
          case LeagueState.PostDraft:
            this.redirectToDashboard();
            break;
        }
        
        // Load league players
        if (league.users) {
          league.users.forEach(user => {
            this.players.set(user.id, user);
          });
        }
      })
    );

    // Subscribe to league's stocks
    this.subscriptions.push(
      this.draftService.leaguePortfolio$.subscribe((leaguePortfolio) => {
        if (!leaguePortfolio || !leaguePortfolio.stocks) return;
        
        this.leagueStocks = leaguePortfolio.stocks;
        
        // Create map for quick lookup
        this.leagueStocks.forEach(stock => {
          this.stocksMap.set(stock.id, stock);
        });
      })
    );
    
    // Subscribe to user's portfolio
    this.subscriptions.push(
      this.portfolioService.userPortfolio$.subscribe((portfolio) => {
        this.userPortfolio = portfolio;
        // Make sure stocks array exists, even if empty
        if (this.userPortfolio && !this.userPortfolio.stocks) {
          this.userPortfolio.stocks = [];
        }
      })
    );

    this.subscriptions.push(
      this.portfolioService.userPortfolio$.subscribe((portfolio) => {
        this.userPortfolio = portfolio;
        // Make sure stocks array exists, even if empty
        if (this.userPortfolio && !this.userPortfolio.stocks) {
          this.userPortfolio.stocks = [];
        }
      })
    );
    
    // Subscribe to all portfolios
    this.subscriptions.push(
      this.draftService.playerPortfoliosForLeague$.subscribe((portfolios) => {
        if (!portfolios) return;
        this.leaguePortfolios = portfolios;
      })
    );
    
    // Subscribe to current draft player updates
    this.subscriptions.push(
      this.draftService.currentDraftPlayer$.subscribe((draftUpdate: DraftUpdateResponse) => {
        if (!draftUpdate) return;
        
        // Update current player
        this.currentPlayerID = draftUpdate.playerID;
        
        // Reset and start timer
        this.remainingTime = this.maxDraftTime;
        this.startTimer();
      })
    );
    
    // Subscribe to draft pick updates
    this.subscriptions.push(
      this.draftService.draftPick$.subscribe((draftPick: DraftPickResponse) => {
        if (!draftPick) return;
        
        // Add to draft picks history
        this.draftPicks.push({
          player_id: draftPick.player_id,
          stock_id: draftPick.stock_id,
          timestamp: new Date()
        });
        
        // Refresh portfolios after a pick
        this.getUserPortfolio();
        this.getAllPortfolios();
      })
    );

    // Get initial data
    this.getLeaguePortfolio();
    this.getUserPortfolio();
    this.getAllPortfolios();
    this.leagueService.getLeagueDetails();
  }
  
  ngOnDestroy(): void {
    // Clear timer
    if (this.timerInterval) {
      clearInterval(this.timerInterval);
    }
    
    // Unsubscribe from all subscriptions
    this.subscriptions.forEach(sub => sub.unsubscribe());
  }
  
  // Start countdown timer
  private startTimer(): void {
    // Clear existing timer if any
    if (this.timerInterval) {
      clearInterval(this.timerInterval);
    }
    
    this.timerInterval = setInterval(() => {
      if (this.remainingTime > 0) {
        this.remainingTime--;
      } else {
        // Time's up logic - could notify the user or auto-skip
        clearInterval(this.timerInterval);
      }
    }, 1000);
  }
  
  // Check if it's the current user's turn
  isUsersTurn(): boolean {
    return this.currentPlayerID === this.currentUserID;
  }
  
  // Draft a stock
  draftStock(stock: Stock): void {
    if (!this.isUsersTurn()) {
      return; // Not your turn
    }
    this.draftService.draftStock(stock.id);
  }

  // View stock details
  stockDetails(stock: Stock): void {
    this.stockService.setStock(stock);
    this.router.navigate(['dashboard/stock-details']);
  }
  
  // Get player name by ID
  getPlayerName(playerID: number): string {
    const player = this.players.get(playerID);
    return player ? player.username : `Player ${playerID}`;
  }
  
  // Get stock ticker by ID
  getStockTicker(stockID: number): string {
    const stock = this.stocksMap.get(stockID);
    return stock ? stock.ticker_symbol : `Stock ${stockID}`;
  }
  
  // Get display text for current player turn
  getCurrentPlayerText(): string {
    if (this.isUsersTurn()) {
      return 'Your turn! Pick a stock.';
    } else {
      return `Waiting for ${this.getPlayerName(this.currentPlayerID)} to pick...`;
    }
  }
  
  // Get empty slots for portfolio visualization
  getEmptySlots(stockCount: number): number[] {
    // Assuming each portfolio can have up to 5 stocks
    const maxStocks = 5;
    const emptyCount = Math.max(0, maxStocks - stockCount);
    return Array(emptyCount).fill(0);
  }

  // Helper Methods
  private getLeaguePortfolio(): void {
    this.draftService.getLeaguePortfolioInfo();
  }

  private getUserPortfolio(): void {
    this.portfolioService.getCurrentUserPortfolio();
  }
  
  private getAllPortfolios(): void {
    this.draftService.getAllPortfolios();
  }

  // Routes
  redirectToDraftQueue(): void {
    this.router.navigate(['/dashboard/draft-queue']);
  }

  redirectToDashboard(): void {
    this.router.navigate(['/dashboard']);
  }
}