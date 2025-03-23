import { Component, HostListener } from '@angular/core';
import { Subscription } from 'rxjs';
import { Portfolio } from '../../models/portfolio.model';
import { PortfolioService } from '../services/portfolio.service';
import { PortfolioPointsHistoryEntry } from '../../models/points-history-entry.model';
import { StockHistoryEntry } from '../../models/stock-history-entry.model';
import { CommonModule, NgFor, NgIf } from '@angular/common';
import { PortfolioPointsHistoryChartComponent } from './portfolio-points-history-chart/portfolio-points-history-chart.component';
import { devLog } from '../../../environments/development/devlog';


@Component({
  selector: 'app-league-portfolio',
  standalone: true,
  imports: [NgFor, NgIf, CommonModule, PortfolioPointsHistoryChartComponent],
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

  // Track window resize events
  @HostListener('window:resize', ['$event'])
  onResize(event: any) {
    // Trigger chart redraw when window is resized
    this.triggerChartResize();
  }
  
  ngAfterViewInit() {
    // Short delay to ensure DOM is ready
    setTimeout(() => {
      this.triggerChartResize();
    }, 100);
  }
  
  private triggerChartResize() {
    // Create and dispatch a resize event using the modern approach
    window.dispatchEvent(new Event('resize'));
  }

}
