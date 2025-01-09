import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Stock } from '../../models/stock.model';
import { LeagueService } from '../services/league.service';
import { LeaguePortfolio } from '../../models/league-portfolio.model';
import { League } from '../../models/league.model';
import { guard } from '../../utils/guard';
import { User } from '../../models/user.model';
import { VerifyUserService } from '../services/verify-user.service';
import { Portfolio } from '../../models/portfolio.model';
import { Router } from '@angular/router';
import { StockService } from '../services/stock.service';
import { PortfolioService } from '../services/portfolio.service';
import { DraftService } from '../services/draft.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-league-draft',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './league-draft.component.html',
  styleUrl: './league-draft.component.scss'
})
export class LeagueDraftComponent {
  leagueStocks: Stock[] = [];
  userPortfolio: Portfolio | null = null

  private subscription!: Subscription;

  constructor(
    private router: Router,
    private draftService: DraftService,
    private portfolioService: PortfolioService,
    private stockService: StockService,
  ) {}

  ngOnInit(): void {

    // * Subscribe to the observables to listen for changes
    
    // League's Portfolio
    this.subscription = this.draftService.leaguePortfolio$.subscribe((leaguePortfolio) => {
      const leaguePortfolioStocks = leaguePortfolio.stocks;
      guard(leaguePortfolioStocks != null, "LeaguePortfolio is nil!");
      this.leagueStocks = leaguePortfolioStocks;
    });
    // User's Portfolio
    this.subscription = this.portfolioService.userPortfolio$.subscribe((portfolio) => {
      this.userPortfolio = portfolio;
    });

    // * Get Starting Values for Dashboard

    // Get League's Portfolio
    this.getLeaguePortfolio();
    // Get User's Portfolio
    this.getUserPortfolio();

  }
  
  ngOnDestroy(): void {
    // Unsubscribe to avoid memory leaks
    this.subscription.unsubscribe();
  }
  
  draftStock(stock: Stock) {
    this.draftService.draftStock(stock.id);
  }

  stockDetails(stock: Stock) {
    this.stockService.setStock(stock);
    this.router.navigate(['dashboard/stock-details']);
  }

  // * Helper Functions

  private getLeaguePortfolio() {
    this.draftService.getLeaguePortfolioInfo();
  }

  // Load the user's portfolio for a specific league
  private getUserPortfolio(): void {
    this.portfolioService.getCurrentUserPortfolio();
  }

}
