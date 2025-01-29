import { CommonModule, NgFor, NgIf } from '@angular/common';
import { Component } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { VerifyUserService } from '../services/verify-user.service';
import { LeagueService } from '../services/league.service';
import { User } from '../../models/user.model';
import { League } from '../../models/league.model';
import { Subscription } from 'rxjs';
import { devLog } from '../../../environments/development/devlog';
import { PortfolioService } from '../services/portfolio.service';
import { TradeService } from '../services/trade.service';
import { FormsModule } from '@angular/forms';
import { StockService } from '../services/stock.service';

@Component({
  selector: 'app-dashboard-main',
  standalone: true,
  imports: [RouterOutlet, CommonModule, NgFor, NgIf, FormsModule],
  templateUrl: './dashboard-main.component.html',
  styleUrl: './dashboard-main.component.scss'
})
export class DashboardMainComponent {

  leagues: League[] = [];
  selectedLeague: League | null = null;
  user: string = "User"
  showMenu = false

  searchQuery: string = ''; 
  searchResults: { name: string; symbol: string }[] = [];

  private subscription!: Subscription;

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
    
    this.subscription = this.leagueService.userLeagues$.subscribe((leagues) => {
      this.leagues = leagues;
    });

    // * Get Starting Values for Dashboard
    
    // Loads the user
    this.loadUser();
    // Gets all of the user's leagues
    this.loadUserLeagues();
  }

  ngOnDestroy(): void {
    // Unsubscribe to avoid memory leaks
    this.subscription.unsubscribe();
  }

  // * Template Methods

  // Method to handle league selection
  selectLeague(league: League) {
    this.leagueService.setSelectedLeague(league)
    this.redirectToDashboard();
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
    const stocks = this.stockService.getAllStocks();
    console.log("helo", stocks)

    if (this.searchQuery.trim()) {
      console.log('Searching for:', this.searchQuery);
      // Mock search results (Replace with an API call in the future)
      this.searchResults = [
        { name: 'Apple', symbol: 'AAPL' },
        { name: 'Tesla', symbol: 'TSLA' },
        { name: 'Microsoft', symbol: 'MSFT' },
      ].filter((stock) =>
        stock.name.toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    } else {
      this.searchResults = [];
    }
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