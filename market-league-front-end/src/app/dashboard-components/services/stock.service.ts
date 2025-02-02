import { Injectable } from '@angular/core';
import { Stock } from '../../models/stock.model';
import { Observable, Subject } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { StockWithHistory } from '../../models/stock-with-history.model';
import { WebSocketMessageTypes } from './websocket-message-types';
import { devLog } from '../../../environments/development/devlog';
import { WebSocketService } from './websocket.service';
import { filter, map } from 'rxjs/operators';



@Injectable({
  providedIn: 'root'
})
export class StockService {

  // * Observables

  // User Portfolio
  private selectedStockDetailsSubject = new Subject<StockWithHistory>();
  selectedStockDetails$ = this.selectedStockDetailsSubject.asObservable();

  private allStockSubject = new Subject<Stock[]>();
  allStock$ = this.allStockSubject.asObservable();

  // * Constructor
  
  constructor(
    private webSocketService: WebSocketService,
  ) {
    this.webSocketService.getMessages().subscribe((message) => {
      switch (message.type) {
        case WebSocketMessageTypes.MessageType_Stock_GetStockInformation:
          devLog("Received GetStockInformation Response: " + message.data);
          this.handleGetStockInformationResponse(message.data);
          break;
        case WebSocketMessageTypes.MessageType_Stock_GetAllStocks:
          devLog("Received GetAllStocks Response: " + message.data);
          this.handleGetAllStocks(message.data);
          break;
        default:
          // devLog("Stock Service unable to route Websocket Message properly! " + message.data);
      }
    });
  }

  private readonly SELECTED_STOCK = "selected_stock";

  // * Websocket Response Handler Functions

  handleGetStockInformationResponse(responseData: any): void {
    // Check for error message
    const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
    if (didErrorOccur) {
      devLog("Error occurred: " + responseData.message)
      return
    }
    this.handleSuccessfulGetStockInformationResponse(responseData);
  }
  
    handleGetAllStocks(responseData: any): void {
      // Check for error message
      const didErrorOccur = this.webSocketService.didErrorOccur(responseData);
      if (didErrorOccur) {
        devLog("Error occurred: " + responseData.message)
        return
      }
      this.handleSuccessfulGetAllStocks(responseData);
    }

  // * Helper Functions to Websocket Responses

  handleSuccessfulGetStockInformationResponse(stockDetails: StockWithHistory): void {
    this.selectedStockDetailsSubject.next(stockDetails);
  }

  handleSuccessfulGetAllStocks(allStocks: Stock[]): void{
    this.allStockSubject.next(allStocks)
  }

  // * Websocket Call Functions

  getStockDetails(stockId: number): void {
    const data = {
      stock_id: stockId
    };
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_Stock_GetStockInformation,
      data: data
    };
    this.webSocketService.sendMessage(websocketMessage);
  }

  getAllStocks(): Observable<Stock[]> {
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_Stock_GetAllStocks,
      data: {} // No additional data needed
    };
    this.webSocketService.sendMessage(websocketMessage);

    return this.webSocketService.getMessages().pipe(
      filter((message: any) => message.type === WebSocketMessageTypes.MessageType_Stock_GetAllStocks),
      map((message: any) => message.data as Stock[]) // Ensure correct type
  );
  }  

  // * Helper Functions

  setStock(stock: Stock): void {
    localStorage.setItem(this.SELECTED_STOCK, JSON.stringify(stock));
  }

  getStock(): Stock {
    const stockData = localStorage.getItem(this.SELECTED_STOCK);
    return stockData ? JSON.parse(stockData) : null;
  }

}
