import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ActivatedRoute } from '@angular/router';
import { Stock } from '../../models/stock.model';
import { StockService } from '../services/stock.service';
import { StockWithHistory } from '../../models/stock-with-history.model';
import { StockChartComponent } from './stock-chart/stock-chart.component';
import { Subscription, Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Component({
  selector: 'app-stock-details',
  standalone: true,
  imports: [StockChartComponent],
  templateUrl: './stock-details.component.html',
  styleUrl: './stock-details.component.scss'
})  
export class StockDetailsComponent implements OnInit{
  selectedStock: Stock | null = null;
  stockDetails: StockWithHistory | null = null

  private subscription!: Subscription;

  constructor(
    private router: Router,
    private stockService: StockService,
    private route: ActivatedRoute,  
  ) {}
  
  ngOnInit() {
    
    // * Subscribe to the observables to listen for changes
    
    this.subscription = this.stockService.selectedStockDetails$.subscribe((stockDetails) => {
      this.stockDetails = stockDetails;
    });

    // * Get Starting Values for Dashboard
    this.route.params.subscribe(params => {
      this.selectedStock = this.stockService.getStock()
      
      this.getStockWithSymbol(params['ticker_symbol']).subscribe(stock => {
        if (stock) {
          this.selectedStock = stock;  // Store the stock in a variable if needed
        }
      });
      this.loadStockDetails(this.selectedStock);
    });
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

  private getStockWithSymbol(ticker_symbol: string): Observable<Stock | undefined> {
    return this.stockService.getAllStocks().pipe(
      map(stocks => stocks.find(stock => stock.ticker_symbol === ticker_symbol))
    );
  }
}