import { Injectable } from '@angular/core';
import { Stock } from '../../models/stock.model';
import { Observable, Subject } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { StockWithHistory } from '../../models/stock-with-history.model';
import { WebSocketMessageTypes } from './websocket-message-types';
import { devLog } from '../../../environments/development/devlog';
import { WebSocketService } from './websocket.service';

@Injectable({
  providedIn: 'root'
})
export class StockService {

  // * Observables

  // User Portfolio
  private selectedStockDetailsSubject = new Subject<StockWithHistory>();
  selectedStockDetails$ = this.selectedStockDetailsSubject.asObservable();

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

  // * Helper Functions to Websocket Responses

  handleSuccessfulGetStockInformationResponse(stockDetails: StockWithHistory): void {
    this.selectedStockDetailsSubject.next(stockDetails);
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

  getAllStocks(): void {
    const websocketMessage = {
      type: WebSocketMessageTypes.MessageType_Stock_GetAllStocks,
      data: {} // No additional data needed
    };
    this.webSocketService.sendMessage(websocketMessage);

    // Listen for the response from the WebSocket
    this.webSocketService.getMessages().subscribe((message: any) => {
        if (message.type === WebSocketMessageTypes.MessageType_Stock_GetAllStocks) {
            console.log("Received all stocks:", message.data); // Should log all stocks
            // You can now use message.data (the list of stocks) as needed
        }
    });
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
