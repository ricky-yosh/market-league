import { CommonModule, NgFor, NgIf } from '@angular/common';
import { Component } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { VerifyUserService } from '../services/verify-user.service';
import { LeagueService } from '../services/league.service';
import { User } from '../../models/user.model';
import { League } from '../../models/league.model';
import { Stock } from '../../models/stock.model';
import { Subscription} from 'rxjs';
import { devLog } from '../../../environments/development/devlog';
import { PortfolioService } from '../services/portfolio.service';
import { TradeService } from '../services/trade.service';
import { FormsModule } from '@angular/forms';
import { StockService } from '../services/stock.service';
import { LeagueState } from '../../models/league-state.model';

@Component({
  selector: 'app-dashboard-main',
  standalone: true,
  imports: [RouterOutlet, CommonModule, NgFor, NgIf, FormsModule],
  templateUrl: './dashboard-main.component.html',
  styleUrl: './dashboard-main.component.scss'
})
export class DashboardMainComponent {
  leagues: League[] = [];
  stocks: Stock[] = [];
  selectedLeague: League | null = null;
  user: string = "User"
  showMenu = false

  searchQuery: string = ''; 
  searchResults: { company_name: string; ticker_symbol: string }[] = [];
  searchSuggestions: any[] = [];
  showDropdown: boolean = false;
  activeIndex: number = -1;

  private subscriptions: Subscription[] = [];

  constructor(
    private router: Router,
    private userService: VerifyUserService,
    private leagueService: LeagueService,
    private portfolioService: PortfolioService,
    private tradeService: TradeService,
    private stockService: StockService
  ) {}

  ngOnInit(): void {
    
    // * Subscribe to the observables to listen for changes
    
    // Subscribe to leagues
    const leaguesSub = this.leagueService.userLeagues$.subscribe((leagues) => {
      this.leagues = leagues;
    });
    this.subscriptions.push(leaguesSub);
    
    // Subscribe to stocks
    const stocksSub = this.stockService.allStock$.subscribe((stocks) => {
      this.stocks = stocks;
    });
    this.subscriptions.push(stocksSub);

    // * Get Starting Values for Dashboard
    
    // Loads the user
    this.loadUser();
    // Gets all of the user's leagues
    this.loadUserLeagues();
    // Gets all of the stocks
    this.stockService.getAllStocks();
  }

  ngOnDestroy(): void {
    // Unsubscribe to avoid memory leaks
    this.subscriptions.forEach(sub => sub.unsubscribe());
  }

  // * Template Methods

  // Method to handle league selection
  selectLeague(league: League) {
    this.leagueService.setSelectedLeague(league)
    switch (league.league_state) {

      case LeagueState.PreDraft:
        this.redirectToDraftQueue();
        break;

      case LeagueState.InDraft:
        this.redirectToDraft();
        break;

      case LeagueState.PostDraft:
        this.redirectToDashboard();
        break;

      default:
        devLog("An error has occurred! LeagueState is not a valid value: " + league.league_state)
        this.redirectToDashboard();
    }
  }

  // Method to hide/show side menu
  toggleMenu() {
    this.showMenu = !this.showMenu;
  }

  // * Helper Methods

  // Method to load the leagues for the user
  private loadUserLeagues(): void {
    // Fetch leagues
    this.leagueService.getUserLeagues();
  }

  // Method to load the user data asynchronously
  private loadUser(): void {
    this.userService.getUserFromToken().subscribe({
      next: (user: User) => {
        devLog('User fetched successfully:', user);
        this.user = user.username;
      },
      error: (error) => {
        console.error('Failed to fetch user from token:', error);
      }
    });
  }

  // * Routing

  redirectToHome() {
    this.router.navigate(['/home']);
  }

  redirectToDraftQueue() {
    this.router.navigate(['/dashboard/draft-queue']);
  }

  redirectToDraft() {
    this.router.navigate(['dashboard/draft']);
  }

  redirectToDashboard() {
    this.loadLeagueMembers();
    this.loadUserPortfolio();
    this.loadTrades();
    this.router.navigate(['/dashboard']);
  }
  
  redirectToLeaderboard() {
    this.router.navigate(['dashboard/leaderboard']);
  }

  redirectToPortfolio() {
    this.router.navigate(['dashboard/portfolio']);
  }  

  redirectToTrades() {
    this.router.navigate(['dashboard/trades']);
  }

  redirectToCreateLeague() {
    this.router.navigate(['dashboard/create-league']);
  }

  redirectToRemoveLeague() {
    this.router.navigate(['dashboard/remove-league']);
  }

  redirectToSettings() {
    this.router.navigate(['dashboard/settings']);
  }

  onSearch(): void {    
    const query = this.searchQuery.trim().toLowerCase()
    if (query) {
      console.log('Searching for:', query);
      // Mock search results (Replace with an API call in the future)
      this.searchResults = this.stocks.filter((stock) =>
        ((stock.company_name?.toLowerCase() || "").includes(query) ||
      (stock.ticker_symbol?.toLowerCase() || "").includes(query))
      );
      this.searchSuggestions = this.stocks
        .map(stock => stock.company_name) // Extract company names
        .filter(name => name.toLowerCase().startsWith(query)) // Suggest names starting with query
        .slice(0, 5); // Limit to 5 suggestions
      this.showDropdown = this.searchResults.length > 0 || this.searchSuggestions.length > 0;
    } else {
      this.searchResults = [];
      this.searchSuggestions = [];
      this.showDropdown = false;
    }
  }

  applySuggestion(suggestion: string): void {
    this.searchQuery = suggestion;  // Fill input with suggestion
    this.onSearch();  // Trigger search
  }
  
  selectStock(stock: any): void {
    this.searchQuery = stock.company_name;  // Fill input with stock name
    this.searchResults = [];  // Clear results
    this.showDropdown = false; // Hide dropdown
    this.activeIndex = -1;
    this.stockService.setStock(stock);
    this.router.navigate(['dashboard/stock-details', stock.ticker_symbol]);
  }

  navigateSuggestions(direction: number): void {
    if (this.searchResults.length > 0) {
      this.activeIndex = (this.activeIndex + direction + this.searchResults.length) % this.searchResults.length;
    }
  }

  showSuggestions(): void {
    this.showDropdown = this.searchResults.length > 0 || this.searchSuggestions.length > 0;
  }

  hideSuggestions(): void {
    setTimeout(() => this.showDropdown = false, 200);  // Small delay for selection clicks
  }

  // * Refresh Information
  // Method to load members of a selected league
  private loadLeagueMembers(): void {
    this.leagueService.getLeagueMembers();
  }

  // Load the user's portfolio for a specific league
  private loadUserPortfolio(): void {
    this.portfolioService.getCurrentUserPortfolio();
  }

  // Load the user's trades for a specific league
  private loadTrades(): void {
    this.tradeService.getTrades();
  }

}