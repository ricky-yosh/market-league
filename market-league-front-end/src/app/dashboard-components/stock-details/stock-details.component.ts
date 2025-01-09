import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { Stock } from '../../models/stock.model';
import { StockService } from '../services/stock.service';
import { StockWithHistory } from '../../models/stock-with-history.model';
import { StockChartComponent } from './stock-chart/stock-chart.component';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-stock-details',
  standalone: true,
  imports: [StockChartComponent],
  templateUrl: './stock-details.component.html',
  styleUrl: './stock-details.component.scss'
})
export class StockDetailsComponent {
  selectedStock: Stock | null = null;
  stockDetails: StockWithHistory | null = null

  private subscription!: Subscription;

  constructor(
    private router: Router,
    private stockService: StockService,
  ) {}
  
  ngOnInit() {
    
    // * Subscribe to the observables to listen for changes
    
    this.subscription = this.stockService.selectedStockDetails$.subscribe((stockDetails) => {
      this.stockDetails = stockDetails;
    });

    // * Get Starting Values for Dashboard

    this.selectedStock = this.stockService.getStock();
    this.loadStockDetails(this.selectedStock);
  }

  ngOnDestroy(): void {
    // Unsubscribe to avoid memory leaks
    this.subscription.unsubscribe();
  }

  returnToDraft() {
    this.router.navigate(['dashboard/draft']);
  }

  private loadStockDetails(stock: Stock | null): void {
    if (stock == null) {
      this.stockDetails = null;
      return
    }
    const stockId = stock.id
    this.stockService.getStockDetails(stockId)
  }

}