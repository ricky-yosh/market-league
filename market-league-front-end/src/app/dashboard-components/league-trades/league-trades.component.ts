import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-league-trades',
  standalone: true,
  imports: [FormsModule, CommonModule],
  templateUrl: './league-trades.component.html',
  styleUrl: './league-trades.component.scss'
})
export class LeagueTradesComponent {
  availableStocks: string[] = ['AAPL', 'TSLA', 'MSFT', 'NVDA', 'GOOGL', 'AMZN', 'NFLX', 'META'];
  leagueUsers: string[] = ['User1', 'User2', 'User3', 'User4'];

  trade: { user2: string; stocks1: string[]; stocks2: string[] } = {
    user2: '',
    stocks1: [],
    stocks2: []
  };

  constructor(private http: HttpClient) {}

  onSubmit() {
    if (this.trade.user2 && this.trade.stocks1.length > 0 && this.trade.stocks2.length > 0) {
      // For simplicity, we're just logging the trade instead of making an actual HTTP request.
      console.log('Trade details:', this.trade);
      alert('Trade successfully created!');
      this.resetForm();
    } else {
      alert('Please complete the form before submitting.');
    }
  }

  resetForm() {
    this.trade = {
      user2: '',
      stocks1: [],
      stocks2: []
    };
  }

  toggleStockSelection(stockList: string[], stock: string) {
    const index = stockList.indexOf(stock);
    if (index === -1) {
      stockList.push(stock);
    } else {
      stockList.splice(index, 1);
    }
  }

}
