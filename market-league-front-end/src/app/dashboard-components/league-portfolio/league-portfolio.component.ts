import { Component } from '@angular/core';
import { Subscription } from 'rxjs';
import { Portfolio } from '../../models/portfolio.model';
import { PortfolioService } from '../services/portfolio.service';
import { PortfolioPointsHistoryEntry } from '../../models/points-history-entry.model';
import { StockHistoryEntry } from '../../models/stock-history-entry.model';
import { NgFor, NgIf } from '@angular/common';

@Component({
  selector: 'app-league-portfolio',
  standalone: true,
  imports: [NgFor, NgIf],
  templateUrl: './league-portfolio.component.html',
  styleUrl: './league-portfolio.component.scss'
})
export class LeaguePortfolioComponent {

  portfolio: Portfolio | null = null
  stockHistoryList: StockHistoryEntry[] = []
  portfolioPointsHistoryList: PortfolioPointsHistoryEntry[] = []

  private subscription!: Subscription;

  constructor(
    private portfolioService: PortfolioService
  ) {}

  ngOnInit(): void {
    
    // * Subscribe to the observables to listen for changes
    
    this.subscription = this.portfolioService.userPortfolio$.subscribe((portfolio) => {
      this.portfolio = portfolio;
    });
    this.subscription = this.portfolioService.stockHistoryList$.subscribe((stockHistoryList) => {
      this.stockHistoryList = stockHistoryList;
    });
    this.subscription = this.portfolioService.portfolioPointsHistoryList$.subscribe((portfolioPointsHistoryList) => {
      this.portfolioPointsHistoryList = portfolioPointsHistoryList;
    });

    // * Get Starting Values for Dashboard
    this.portfolioService.getCurrentUserPortfolio();
    // Gets the points history of a portfolio
    this.portfolioService.getPortfolioPointsHistory();
    // Gets the change in stock value for a portfolio
    this.portfolioService.getStocksValueChange();

  }

  ngOnDestroy(): void {
    // Unsubscribe to avoid memory leaks
    this.subscription.unsubscribe();
  }

}
