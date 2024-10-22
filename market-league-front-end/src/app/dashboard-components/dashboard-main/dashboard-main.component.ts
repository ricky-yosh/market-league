import { CommonModule, NgFor, NgIf } from '@angular/common';
import { Component, EventEmitter, Output } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-dashboard-main',
  standalone: true,
  imports: [RouterOutlet, CommonModule, NgFor, NgIf],
  templateUrl: './dashboard-main.component.html',
  styleUrl: './dashboard-main.component.scss'
})
export class DashboardMainComponent {
  constructor(private router: Router) {}

  leagues = ['Bull Market Champs', 'Wall Street Warriors', 'Market Mavericks', 'Stock Titans League', 'Bearish Bulls League', 'Dividend Dynamos', 'The Capital Gains Crew', 'Blue Chip Brawlers', 'The Stock Savvy Syndicate', 'Trading Titans League', 'Wall Street Whiz Kids', 'Penny Stock Prospectors', 'Portfolio Powerhouse League', 'The Equity Experts', 'Bullish Investors Club', 'Futures and Fortunes', 'The Risk Takers League', 'The Value Investors', 'Wall Street Wizards', 'The Growth Gurus League'];
  user = "Ricky"
  // Routing
  redirectToHome() {
    this.router.navigate(['/home']);
  }

  redirectToDashboard() {
    this.router.navigate(['/dashboard']);
  }
  
  redirectToLeaderboard() {
    this.router.navigate(['dashboard/leaderboard']);
    // current_league = id_of_league
  }

  redirectToPortfolio() {
    this.router.navigate(['dashboard/portfolio']);
  }  

  redirectToTrades() {
    this.router.navigate(['dashboard/trades']);
  }
}
