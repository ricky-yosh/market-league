import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { Stock } from '../../models/stock.model';
import { StockService } from '../services/stock.service';
import { StockWithHistory } from '../../models/stock-with-history';
import { devLog } from '../../../environments/development/devlog';
import { StockChartComponent } from './stock-chart/stock-chart.component';

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

  constructor(
    private router: Router,
    private stockService: StockService,
  ) {}
  
  ngOnInit() {
    this.selectedStock = this.stockService.getStock();
    this.loadStockDetails(this.selectedStock);
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
    this.stockService.getStockDetails(stockId).subscribe({
      next: (stockDetails) => {
        this.stockDetails = stockDetails;
      },
      error: (error) => {
        devLog('Failed to fetch stock details from stockId:', error);
      }
    });
  }

}