import { NgFor } from '@angular/common';
import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { LeagueService } from '../services/league.service';
import { League } from '../../models/league.model';
import { Stock } from '../../models/stock.model';
import { Subscription } from 'rxjs';
import { Trade } from '../../models/trade.model';
import { PortfolioService } from '../services/portfolio.service';
import { TradeService } from '../services/trade.service';

@Component({
  selector: 'app-league-home',
  standalone: true,
  imports: [NgFor, CommonModule],
  templateUrl: './league-home.component.html',
  styleUrl: './league-home.component.scss'
})
export class LeagueHomeComponent implements OnInit {
  
  selectedLeague: League | null = null;
  userPortfolio: Stock[] | null = null
  leagueTrades: Trade[] | null = null
  leagueMembers: string[] | null = null;

  private subscription!: Subscription;

  constructor(
    private leagueService: LeagueService,
    private portfolioService: PortfolioService,
    private tradeService: TradeService,
  ) {}

  ngOnInit(): void {

    // * Subscribe to the observables to listen for changes
    
    // Selected League
    this.subscription = this.leagueService.selectedLeague$.subscribe((league) => {
      this.selectedLeague = league;
    });
    // User Portfolio
    this.subscription = this.portfolioService.userPortfolio$.subscribe((portfolio) => {
      this.userPortfolio = portfolio ? portfolio.stocks : null;
    });
    // League Trades
    this.subscription = this.tradeService.leagueTrades$.subscribe((trades) => {
      this.leagueTrades = trades;
    });
    // League Members
    this.subscription = this.leagueService.leagueMembers$.subscribe((members) => {
      const memberNames = members.map(user => user.username);
      this.leagueMembers = memberNames;
    });

    // * Get Starting Values for Dashboard

    // Load League Members
    this.loadLeagueMembers();
    // Load User's Portfolio
    this.loadUserPortfolio();
    // Load League Trades
    this.loadTrades();

  }
  
  ngOnDestroy(): void {
    // Unsubscribe to avoid memory leaks
    this.subscription.unsubscribe();
  }

  // Method to load members of a selected league
  private loadLeagueMembers(): void {
    this.leagueService.getLeagueMembers();
  }

  // Load the user's portfolio for a specific league
  private loadUserPortfolio(): void {
    this.portfolioService.getCurrentUserPortfolio();
  }

  // Load the user's trades for a specific league
  private loadTrades(): void {
    this.tradeService.getTrades();
  }

}