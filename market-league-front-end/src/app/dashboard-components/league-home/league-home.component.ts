import { NgFor } from '@angular/common';
import { CommonModule } from '@angular/common';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { LeagueService } from '../services/league.service';
import { League } from '../../models/league.model';
import { Stock } from '../../models/stock.model';
import { Subscription, combineLatest } from 'rxjs';
import { Trade } from '../../models/trade.model';
import { PortfolioService } from '../services/portfolio.service';
import { TradeService } from '../services/trade.service';
import { VerifyUserService } from '../services/verify-user.service';
import { User } from '../../models/user.model';
import { LeaderboardUser } from '../../models/websocket-responses/league/leaderboard-user';

@Component({
  selector: 'app-league-home',
  standalone: true,
  imports: [NgFor, CommonModule],
  templateUrl: './league-home.component.html',
  styleUrl: './league-home.component.scss'
})
export class LeagueHomeComponent implements OnInit, OnDestroy {
  
  selectedLeague: League | null = null;
  userPortfolio: Stock[] | null = null;
  leagueTrades: Trade[] | null = null;
  leagueMembers: string[] | null = null;
  
  // Leaderboard data
  leaderboard: LeaderboardUser[] = [];
  leaderboardWithRank: {username: string, total_value: number, rank: number}[] = [];
  currentUser: string = '';
  leaderboardLoaded = false;

  private subscriptions: Subscription = new Subscription();

  constructor(
    private leagueService: LeagueService,
    private portfolioService: PortfolioService,
    private tradeService: TradeService,
    private userService: VerifyUserService
  ) {}

  ngOnInit(): void {
    // Get current user
    this.loadUser();

    // Selected League
    this.subscriptions.add(
      this.leagueService.selectedLeague$.subscribe((league) => {
        this.selectedLeague = league;
        if (league && league.id) {
          // Load leaderboard when league changes
          this.leagueService.getLeagueLeaderboard(league.id);
        }
      })
    );
    
    // User Portfolio
    this.subscriptions.add(
      this.portfolioService.userPortfolio$.subscribe((portfolio) => {
        this.userPortfolio = portfolio ? portfolio.stocks : null;
      })
    );
    
    // League Trades
    this.subscriptions.add(
      this.tradeService.leagueTrades$.subscribe((trades) => {
        this.leagueTrades = trades;
      })
    );
    
    // League Members
    this.subscriptions.add(
      this.leagueService.leagueMembers$.subscribe((members) => {
        const memberNames = members.map(user => user.username);
        this.leagueMembers = memberNames;
      })
    );
    
    // Leaderboard data
    this.subscriptions.add(
      this.leagueService.leaderboard$.subscribe((leaderboardData) => {
        this.leaderboard = leaderboardData;
        this.leaderboardLoaded = true;
        
        // Process leaderboard data
        this.processLeaderboard();
      })
    );

    // Load initial data
    this.loadLeagueMembers();
    this.loadUserPortfolio();
    this.loadTrades();
  }
  
  ngOnDestroy(): void {
    // Unsubscribe to avoid memory leaks
    this.subscriptions.unsubscribe();
  }

  // Process leaderboard data - sort and add rankings
  private processLeaderboard(): void {
    if (!this.leaderboard || this.leaderboard.length === 0) {
      return;
    }

    // Sort by total_value in descending order
    const sortedLeaderboard = [...this.leaderboard].sort((a, b) => {
      const valueA = typeof a.total_value === 'number' ? a.total_value : 0;
      const valueB = typeof b.total_value === 'number' ? b.total_value : 0;
      return valueB - valueA;
    });

    // Add rankings
    this.leaderboardWithRank = [];
    let currentRank = 1;
    let previousValue: number | null = null;
    
    sortedLeaderboard.forEach((member, index) => {
      // If not the first member and value is different from previous
      if (index > 0 && member.total_value !== previousValue) {
        currentRank = index + 1;
      }
      
      this.leaderboardWithRank.push({
        username: member.username,
        total_value: member.total_value,
        rank: currentRank
      });
      
      previousValue = member.total_value;
    });
    
    // Limit to top 5 for the abbreviated view
    this.leaderboardWithRank = this.leaderboardWithRank.slice(0, 5);
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
  
  // Load the current user
  private loadUser(): void {
    this.userService.getUserFromToken().subscribe({
      next: (user: User) => {
        this.currentUser = user.username;
      },
      error: (error) => {
        console.error('Failed to fetch user from token:', error);
      }
    });
  }
}