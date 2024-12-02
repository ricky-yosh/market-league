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

  setStock(stock: Stock): void {
    this.selectedStock = stock;
  }

  getStock(): Stock | null {
    return this.selectedStock;
  }

  getStockDetails(stockId: number): Observable<StockWithHistory> {
    const payload = {
      stock_id: stockId
    };
    return this.http.post<StockWithHistory>(this.getStockDetailsURL, payload);
  }

}
