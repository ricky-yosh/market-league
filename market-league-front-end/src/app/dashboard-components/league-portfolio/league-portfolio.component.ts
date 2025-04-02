import { Component, HostListener, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs';
import { Portfolio } from '../../models/portfolio.model';
import { PortfolioService } from '../services/portfolio.service';
import { PortfolioPointsHistoryEntry } from '../../models/points-history-entry.model';
import { StockHistoryEntry } from '../../models/stock-history-entry.model';
import { CommonModule, NgFor, NgIf } from '@angular/common';
import { PortfolioPointsHistoryChartComponent } from './portfolio-points-history-chart/portfolio-points-history-chart.component';

@Component({
  selector: 'app-league-portfolio',
  standalone: true,
  imports: [NgFor, NgIf, CommonModule, PortfolioPointsHistoryChartComponent],
  templateUrl: './league-portfolio.component.html',
  styleUrl: './league-portfolio.component.scss'
})
export class LeaguePortfolioComponent {
  @ViewChild(PortfolioPointsHistoryChartComponent) chart!: PortfolioPointsHistoryChartComponent;
  
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

  // Track window resize events
  @HostListener('window:resize', ['$event'])
  onResize(event: any) {
    // Instead of triggering another resize, just update the chart directly
    this.updateChartSize();
  }
  
  ngAfterViewInit() {
    // Short delay to ensure DOM is ready
    setTimeout(() => {
      this.updateChartSize();
    }, 100);
  }
  
  private updateChartSize() {
    // Call a resize or redraw method on your chart component instead of dispatching a resize event
    if (this.chart) {
      // Assuming your chart component has a resize or redraw method
      this.chart.resize(); // or this.chart.redraw();
    }
  }
}