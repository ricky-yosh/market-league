import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Stock } from '../../models/stock.model';
import { LeagueService } from '../services/league.service';
import { LeaguePortfolio } from '../../models/league-portfolio';
import { League } from '../../models/league.model';
import { guard } from '../../utils/guard';
import { catchError, EMPTY, map, Observable, of, tap } from 'rxjs';
import { User } from '../../models/user.model';
import { devLog } from '../../../environments/development/devlog';
import { VerifyUserService } from '../../user-verification/verify-user.service';
import { Portfolio } from '../../models/portfolio.model';
import { Router } from '@angular/router';
import { StockService } from '../services/stock.service';

@Component({
  selector: 'app-league-draft',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './league-draft.component.html',
  styleUrl: './league-draft.component.scss'
})
export class LeagueDraftComponent {
  stocks: Stock[] = [];
  currentLeague: League | null = null
  user: User | null = null
  userPortfolio: Portfolio | null = null
  userPortfolioStocks: Stock[] = []

  constructor(
    private router: Router,
    private leagueService: LeagueService,
    private userService: VerifyUserService,
    private stockService: StockService,
  ) {}

  ngOnInit(): void {
    this.getCurrentLeague();
    this.loadUser();
    this.getLeaguePortfolioInfo(this.currentLeague);
  }

  private getCurrentLeague(): void {
    this.currentLeague = this.leagueService.getStoredLeague();
  }

  private getLeaguePortfolioInfo(league: League | null) {
    guard(league != null, "LeagueId cannot be null!")

    const leagueId = league.id
    this.leagueService.getLeaguePortfolioInfo(leagueId).subscribe({
      next: (data: LeaguePortfolio) => {
        this.stocks = data.stocks || []; // Assuming 'stocks' is a property in the response
      },
      error: (error) => {
        console.error('Error fetching league portfolio info:', error);
      }
    });
  }

  private loadUser(): void {
    this.userService.getUserFromToken().subscribe({
      next: (user: User) => {
        console.log('User fetched successfully:', user);
        this.user = user;
        guard(this.currentLeague != null, "Current League cannot be null!");
        this.loadUserPortfolio(user.id, this.currentLeague.id);
      },
      error: (error) => {
        console.error('Failed to fetch user from token:', error);
      }
    });
  }

  // Load the user's portfolio for a specific league
  private loadUserPortfolio(userId: number, leagueId: number): void {
    this.leagueService.getUserPortfolio(userId, leagueId).pipe(
      map((response: Portfolio) => {
        devLog('User portfolio fetched successfully:', response);
        this.userPortfolio = response;
        this.userPortfolioStocks = this.userPortfolio.stocks
      }),
      catchError((error) => {
        console.error('Failed to fetch user portfolio:', error);
        return EMPTY; // Return an empty observable to handle the error gracefully
      })
    ).subscribe();
  }
  
  draftStock(stock: Stock) {
    guard(this.user != null, "User cannot be null!");
    guard(this.currentLeague != null, "User cannot be null!");

    this.leagueService.draftStock(this.currentLeague.id, this.user.id, stock.id).subscribe({
      next: (data: LeaguePortfolio) => {
        this.getLeaguePortfolioInfo(this.currentLeague)
        this.userPortfolioStocks.push(stock)
      },
      error: (error) => {
        console.error('Error fetching league portfolio info:', error);
      }
    });
  }

  stockDetails(stock: Stock) {
    this.stockService.setStock(stock);
    this.router.navigate(['dashboard/stock-details']);
  }

}
