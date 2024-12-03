import { Injectable } from '@angular/core';
import { Stock } from '../../models/stock.model';
import { environment } from '../../../environments/environment';
import { Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { StockWithHistory } from '../../models/stock-with-history';

@Injectable({
  providedIn: 'root'
})
export class StockService {
  private baseUrl = environment.api_url;
  private getStockDetailsURL = `${this.baseUrl}/api/stocks/get-stock-information`;

  private selectedStock: Stock | null = null;

  constructor(private http: HttpClient) {}

  private readonly selected_stock = '';

  setStock(stock: Stock): void {
    localStorage.setItem(this.selected_stock, JSON.stringify(stock));
  }

  getStock(): Stock {
    const stockData = localStorage.getItem(this.selected_stock);
    return stockData ? JSON.parse(stockData) : null;
  }

  getStockDetails(stockId: number): Observable<StockWithHistory> {
    const payload = {
      stock_id: stockId
    };
    return this.http.post<StockWithHistory>(this.getStockDetailsURL, payload);
  }

}
