import { Injectable } from '@angular/core';
import { Subject, Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { WebSocketTransmission } from '../../models/websocket-transmission.model';
import { WebSocketMessageTypes } from './websocket-message-types';
import { devLog } from '../../../environments/development/devlog';

@Injectable({
  providedIn: 'root',
})
export class WebSocketService {
  private socket!: WebSocket;
  private messageSubject = new Subject<any>();
  private websocket_URL = environment.websocket_url;
  private messageQueue: string[] = [];
  private isConnected = false;

  constructor() {}

  connect(): void {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      const websocketUrlWithEndpoint = `${this.websocket_URL}/ws`;
      devLog("Websocket url: ", websocketUrlWithEndpoint);
      this.socket = new WebSocket(websocketUrlWithEndpoint);

      this.socket.onopen = () => {
          console.log('WebSocket connected');
          this.isConnected = true;

          // Flush the message queue
          while (this.messageQueue.length > 0) {
            const message = this.messageQueue.shift(); // Remove the first message
            if (message) this.socket.send(message); // Send it through WebSocket
          }
        };
      this.socket.onmessage = (event) => {
        try {
          const parsedData = JSON.parse(event.data);
          devLog(parsedData)
          this.messageSubject.next(parsedData);
        } catch (error) {
          console.error("Failed to parse WebSocket message:", error, event.data);
        }
      };
      this.socket.onclose = () => {
        console.log('WebSocket disconnected');
        this.isConnected = false;
      };
      this.socket.onerror = (error) => console.error('WebSocket error:', error);
    }
  }

  closeSocketConnection(): void {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.close(); // Close only if the socket is open
      console.log('WebSocket connection closed');
    } else {
      console.warn('WebSocket is already closed or not initialized.');
    }
  }

  // Listen for messages
  getMessages(): Observable<any> {
    return this.messageSubject.asObservable();
  }

  // Send messages to the backend
  sendMessage(message: WebSocketTransmission): void {
    const messageString = JSON.stringify(message);

    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(message));
    } else {
      console.warn('WebSocket is not connected! Queuing message.');
      this.messageQueue.push(messageString); // Queue the message if WebSocket is not ready
    }
  }
  
  // Check for errors in the 
  didErrorOccur(data: any): boolean {
    return data?.type == WebSocketMessageTypes.MessageType_Error
  }

}
