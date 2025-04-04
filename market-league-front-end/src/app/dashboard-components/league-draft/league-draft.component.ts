import { CommonModule } from '@angular/common';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { Stock } from '../../models/stock.model';
import { Portfolio } from '../../models/portfolio.model';
import { NavigationExtras, Router } from '@angular/router';
import { StockService } from '../services/stock.service';
import { PortfolioService } from '../services/portfolio.service';
import { DraftService } from '../services/draft.service';
import { Subscription } from 'rxjs';
import { LeagueState } from '../../models/league-state.model';
import { LeagueService } from '../services/league.service';
import { WebSocketService } from '../services/websocket.service';
import { DraftUpdateResponse } from '../../models/websocket-responses/draft/draft-update-response.model';
import { DraftPickResponse } from '../../models/websocket-responses/draft/draft-pick-response.model';
import { VerifyUserService } from '../services/verify-user.service';
import { User } from '../../models/user.model';
import { League } from '../../models/league.model';
import { DraftPick } from '../../models/websocket-responses/draft/draft-pick.model';

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
  selectedLeague: League | null = null;
  
  // Draft state
  currentPlayerID: number = 0;
  currentUserID: number = 0;
  remainingTime: number = 0;
  maxDraftTime: number = 30; // 30 seconds per turn
  draftPicks: DraftPick[] = [];
  private currentLeagueId: number = 0;
  
  // User data
  players: Map<number, User> = new Map();
  
  // Stock data for quick lookup
  stocksMap: Map<number, Stock> = new Map();
  
  // Subscriptions
  private subscriptions: Subscription[] = [];
  
  // Timer interval
  private timerInterval: any;
  
  // Reconnection flag
  private isReconnecting: boolean = false;

  constructor(
    private router: Router,
    private draftService: DraftService,
    private portfolioService: PortfolioService,
    private stockService: StockService,
    private leagueService: LeagueService,
    private verifyUserService: VerifyUserService,
    private websocketService: WebSocketService,
  ) {}

  ngOnInit(): void {
    // Subscribe to WebSocket connection status
    const connectionSub = this.websocketService.connectionStatus.subscribe(isConnected => {
      if (isConnected) {
        // When connection is established/re-established
        this.isReconnecting = true;
        console.log('WebSocket connected, initializing draft data...');

        // Get the stored league from localStorage first
        const storedLeague = this.leagueService.getStoredLeague();
        if (storedLeague) {
          // Explicitly set the selected league first
          this.leagueService.setSelectedLeague(storedLeague);

          const cacheKey = this.getStocksMapCacheKey();
          const cachedStocksStr = localStorage.getItem(cacheKey);
          console.log(`Looking for cached stocks with key: ${cacheKey}`);

          if (cachedStocksStr) {
            try {
              const cachedStocks = JSON.parse(cachedStocksStr);
              cachedStocks.forEach((stock: Stock) => {
                this.stocksMap.set(stock.id, stock);
              });
              console.log(`Loaded ${this.stocksMap.size} cached stocks from localStorage`);
            } catch (e) {
              console.error('Error loading cached stocks:', e);
            }
          }
          
          // Then subscribe to it
          this.leagueService.subscribeToLeague();
          
          // Also reload draft-specific state
          this.getLeaguePortfolio();
          this.getUserPortfolio();
          this.getAllPortfolios();
          
          // Get detailed league info including current draft status
          this.leagueService.getLeagueDetails();
          
          // Set isReconnecting to false after a short delay
          setTimeout(() => {
            this.isReconnecting = false;
            console.log('Draft data initialization completed');
          }, 2000);
        }
      } else {
        console.log('WebSocket disconnected');
      }
    });
    this.subscriptions.push(connectionSub);

    // Initial connection
    this.websocketService.connect();
    
    // Initial data load
    this.leagueService.getUserLeagues();
    this.leagueService.subscribeToLeague();

    // Get current user ID
    const currentUser = this.verifyUserService.getCurrentUserValue();
    if (currentUser) {
      this.currentUserID = currentUser.id;
    }

    // Subscribe to selected league changes
    this.subscriptions.push(
      this.leagueService.selectedLeague$.subscribe((league) => {
        this.selectedLeague = league;
        if (!league) return;
        // Store the current league ID
        this.currentLeagueId = league.id;

        // Load draft picks for this specific league
        this.loadDraftPicksForCurrentLeague();

        switch(league.league_state) {
          case LeagueState.PreDraft:
            this.redirectToDraftQueue();
            break;
          case LeagueState.PostDraft:
            this.redirectToDashboard();
            break;
          case LeagueState.Completed:
            this.redirectToCompletedLeague();
            break;
          default:
            // stay on in draft
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
        
        // Update stocksMap with any new stocks, but keep existing ones
        // This ensures we don't lose stock data when they're drafted
        leaguePortfolio.stocks.forEach(stock => {
          if (!this.stocksMap.has(stock.id)) {
            this.stocksMap.set(stock.id, stock);
          }
        });

        // Save updated stocksMap to localStorage
        this.saveStocksMapToLocalStorage();

        console.log(`Received league portfolio with ${this.leagueStocks.length} stocks`);
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
        
        console.log(`Received user portfolio with ${this.userPortfolio?.stocks?.length || 0} stocks`);
      })
    );
    
    // Subscribe to all portfolios
    this.subscriptions.push(
      this.draftService.playerPortfoliosForLeague$.subscribe((portfolios) => {
        if (!portfolios) return;
        this.leaguePortfolios = portfolios;
        console.log(`Received ${this.leaguePortfolios.length} player portfolios`);
      })
    );
    
    // Subscribe to current draft player updates
    this.subscriptions.push(
      this.draftService.currentDraftPlayer$.subscribe((draftUpdate: DraftUpdateResponse) => {
        if (!draftUpdate) return;
        
        console.log(`Draft update received: Current player ID = ${draftUpdate.playerID}, Remaining time = ${draftUpdate.remainingTime}s`);
        
        // Update current player
        this.currentPlayerID = draftUpdate.playerID;
        
        // Use the remaining time from the server instead of resetting to max
        if (draftUpdate.remainingTime !== undefined) {
          this.remainingTime = draftUpdate.remainingTime;
        } else {
          // Fallback to max time if for some reason the server doesn't send it
          this.remainingTime = this.maxDraftTime;
        }
        
        // Only start a new timer if one isn't already running
        if (!this.timerInterval) {
          this.startTimer();
        }
      })
    );
    
    // Subscribe to draft pick updates
    this.subscriptions.push(
      this.draftService.draftPick$.subscribe((draftPick: DraftPickResponse) => {
        if (!draftPick) return;
        
        console.log(`Draft pick received: Player ${draftPick.player_id} picked stock ${draftPick.stock_id}`);
        
        // Refresh the draft picks from the service
        // (the service already stored the pick in memory and localStorage)
        if (this.currentLeagueId) {
          this.draftPicks = this.draftService.loadDraftState(this.currentLeagueId);
        }
        
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
    this.leagueService.unsubscribeFromLeague();
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
        this.timerInterval = null; // Reset the timer variable
      }
    }, 1000);
  }

  // Add a new method to load draft picks for the current league
  private loadDraftPicksForCurrentLeague(): void {
    if (!this.currentLeagueId) {
      console.warn('Cannot load draft picks: No league selected');
      return;
    }
    
    // Load draft picks for this league
    this.draftPicks = this.draftService.loadDraftState(this.currentLeagueId);
    console.log(`Loaded ${this.draftPicks.length} draft picks for league ${this.currentLeagueId}`);
    
    // Check for stale draft history
    if (this.shouldResetDraftHistory()) {
      console.log('Detected stale draft history, resetting...');
      this.draftPicks = [];
      this.draftService.clearDraftState(this.currentLeagueId);
    }
  }

  // Update the shouldResetDraftHistory method
  private shouldResetDraftHistory(): boolean {
    // Get last draft pick timestamp if available
    if (!this.draftPicks.length || !this.currentLeagueId) {
      return false; // No history to reset or no league selected
    }

    const lastPick = this.draftPicks[this.draftPicks.length - 1];
    if (!lastPick.timestamp) {
      return true; // No timestamp, should reset to be safe
    }
    
    // Convert string date to Date object if needed
    const lastPickTime = typeof lastPick.timestamp === 'string' 
      ? new Date(lastPick.timestamp) 
      : lastPick.timestamp;
    
    // Check if the last pick is older than 24 hours
    const oneDayInMs = 24 * 60 * 60 * 1000;
    const isStale = (new Date().getTime() - lastPickTime.getTime()) > oneDayInMs;
    
    if (isStale) {
      console.log(`Draft history for league ${this.currentLeagueId} is stale. Last pick was at ${lastPickTime}`);
    }
    
    return isStale;
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
    console.log(`Drafting stock: ${stock.ticker_symbol} (ID: ${stock.id})`);
    this.draftService.draftStock(stock.id);
  }

  // View stock details
  stockDetails(stock: Stock): void {
    this.stockService.setStock(stock);
    // Add navigation extras to pass state during navigation
    const navigationExtras: NavigationExtras = {
      queryParams: { 'fromStockSelection': 'true' }
    };
    this.router.navigate(['dashboard/stock-details', stock.ticker_symbol], navigationExtras);
  }
  
  // Get player name by ID
  getPlayerName(playerID: number): string {
    const player = this.players.get(playerID);
    return player ? player.username : `Player ${playerID}`;
  }
  
  // Update the getStockTicker method as well
  getStockTicker(stockID: number): string {
    // First check in stocksMap
    const stock = this.stocksMap.get(stockID);
    if (stock) {
      return stock.ticker_symbol;
    }
    
    // If not found, check in leagueStocks array
    const leagueStock = this.leagueStocks.find(s => s.id === stockID);
    if (leagueStock) {
      // Add to map for future lookups and return
      this.stocksMap.set(stockID, leagueStock);
      return leagueStock.ticker_symbol;
    }
    
    // If we still don't have it, check if we have cached stock data
    const cacheKey = this.getStocksMapCacheKey();
    const cachedStocksStr = localStorage.getItem(cacheKey);
    if (cachedStocksStr) {
      try {
        const cachedStocks = JSON.parse(cachedStocksStr);
        const cachedStock = cachedStocks.find((s: Stock) => s.id === stockID);
        if (cachedStock) {
          // Add to map for future lookups and return
          this.stocksMap.set(stockID, cachedStock);
          return cachedStock.ticker_symbol;
        }
      } catch (e) {
        console.error(`Error parsing cached stocks for league ${this.currentLeagueId}:`, e);
      }
    }
    
    // If all else fails, just return the ID with a label
    return `Stock ${stockID}`;
  }

  // Add a method to save stocksMap to localStorage
  private saveStocksMapToLocalStorage(): void {
    if (!this.currentLeagueId) {
      console.warn('Cannot save stocks map: No league selected');
      return;
    }
    
    // Convert map to array for storage
    const stocksArray = Array.from(this.stocksMap.values());
    const cacheKey = this.getStocksMapCacheKey();
    localStorage.setItem(cacheKey, JSON.stringify(stocksArray));
    console.log(`Saved ${stocksArray.length} stocks to cache for league ${this.currentLeagueId}`);
  }
  
  // Get display text for current player turn
  getCurrentPlayerText(): string {
    if (this.currentPlayerID === 0) {
      return 'Waiting for draft to start...';
    }
    
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
    console.log('Fetching league portfolio');
    this.draftService.getLeaguePortfolioInfo();
  }

  private getUserPortfolio(): void {
    console.log('Fetching user portfolio');
    this.portfolioService.getCurrentUserPortfolio();
  }
  
  private getAllPortfolios(): void {
    console.log('Fetching all league portfolios');
    this.draftService.getAllPortfolios();
  }

  // Update stocksMap cache methods to use league ID
  private getStocksMapCacheKey(): string {
    if (!this.currentLeagueId) {
      return 'stocksMapCache';
    }
    return `stocksMapCache_${this.currentLeagueId}`;
  }

  // Routes
  redirectToDraftQueue(): void {
    // Don't redirect during reconnection
    if (this.isReconnecting) {
      console.log('Reconnection in progress, not redirecting to draft queue');
      return;
    }
    console.log('Redirecting to draft queue');
    this.router.navigate(['/dashboard/draft-queue']);
  }

  redirectToDashboard(): void {
    // Don't redirect during reconnection
    if (this.isReconnecting) {
      console.log('Reconnection in progress, not redirecting to dashboard');
      return;
    }
    console.log('Redirecting to dashboard');
    this.router.navigate(['/dashboard']);
  }

  redirectToCompletedLeague(): void {
    // Don't redirect during reconnection
    if (this.isReconnecting) {
      console.log('Reconnection in progress, not redirecting to completed league');
      return;
    }
    console.log('Redirecting to completed league');
    this.router.navigate(['/dashboard/league-completed']);
  }

  // Add a method to clear the stocks cache for a specific league
  clearStocksCache(): void {
    if (!this.currentLeagueId) {
      return;
    }
    
    const cacheKey = this.getStocksMapCacheKey();
    localStorage.removeItem(cacheKey);
    this.stocksMap.clear();
    console.log(`Cleared stocks cache for league ${this.currentLeagueId}`);
  }
}