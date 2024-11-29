import { NgFor } from '@angular/common';
import { CommonModule } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { LeagueService } from '../services/league.service';
import { User } from '../../models/user.model';
import { League } from '../../models/league.model';
import { Stock } from '../../models/stock.model';
import { Portfolio } from '../../models/portfolio.model';
import { VerifyUserService } from '../../user-verification/verify-user.service';
import { EMPTY, Observable, catchError, map, of, switchMap, tap } from 'rxjs';
import { Trade } from '../../models/trade.model';
import { devLog } from '../../../environments/development/devlog';
import { Router } from '@angular/router';

@Component({
  selector: 'app-league-home',
  standalone: true,
  imports: [NgFor, CommonModule],
  templateUrl: './league-home.component.html',
  styleUrl: './league-home.component.scss'
})
export class LeagueHomeComponent implements OnInit {
  
  selectedLeague: League | null = null;
  user: User | null = null;
  userPortfolio: Stock[] | null = null
  userTrades: Trade[] | null = null
  leagueMembers: string[] | null = null;

  constructor(
    private leagueService: LeagueService,
    private userService: VerifyUserService,
    private cd: ChangeDetectorRef,
    private router: Router
  ) {}

  ngOnInit(): void {
    // Subscribe to the league changes
    this.leagueService.selectedLeague$.pipe(
      switchMap((league) => {
        this.selectedLeague = league;
        this.cd.detectChanges(); // Detect changes after updating selectedLeague
  
        if (this.selectedLeague?.id) {
          // Load league members and wait for the user to be loaded
          return this.loadLeagueMembers(this.selectedLeague.id).pipe(
            switchMap((members) => {
              this.leagueMembers = members;
              return this.loadUser(); // Load user after league members are loaded
            }),
            switchMap((user) => {
              // Store the user information
              this.user = user;
  
              // Check if both user and selectedLeague are defined
              if (user?.id && this.selectedLeague?.id) {
                return this.loadUserPortfolio(user.id, this.selectedLeague.id).pipe(
                  switchMap((portfolio) => {
                    this.userPortfolio = portfolio;
                    console.log('Portfolio loaded successfully');
  
                    // Check again if both user and selectedLeague are defined before loading trades
                    if (user.id && this.selectedLeague?.id) {
                      return this.loadUserTrades(user.id, this.selectedLeague.id);
                    }
                    return of([]); // Return an empty array if conditions are not met
                  })
                );
              }
              return of([]); // Return an empty array if user or league ID is not available
            })
          );
        }
        return EMPTY; // If no league selected, return an empty observable
      })
    ).subscribe({
      next: (trades) => {
        this.userTrades = trades;
        console.log('Trades loaded successfully:', this.userTrades);
      },
      error: (error) => console.error('Failed to load data:', error)
    });
  }  

  // Method to load members of a selected league
  private loadLeagueMembers(leagueId: number): Observable<string[]> {
    return this.leagueService.getLeagueMembers(leagueId).pipe(
      map((members) => {
        console.log('League members fetched successfully:', members);
        return members.map((member: User) => member.username);
      }),
      catchError((error) => {
        console.error('Failed to fetch league members:', error);
        return of([]); // Return an empty array on error
      })
    );
  }

  // Load the user's portfolio for a specific league
  private loadUserPortfolio(userId: number, leagueId: number): Observable<Stock[]> {
    return this.leagueService.getUserPortfolio(userId, leagueId).pipe(
      map((response: Portfolio) => {
        console.log('User portfolio fetched successfully:', response.stocks);
        return response.stocks;
      }),
      catchError((error) => {
        console.error('Failed to fetch user portfolio:', error);
        return of([]); // Return an empty array on error
      })
    );
  }

  // Method to load the user from a token
  private loadUser(): Observable<User> {
    return this.userService.getUserFromToken().pipe(
      tap((user: User) => {
        console.log('User fetched successfully:', user);
        this.user = user;
      }),
      catchError((error) => {
        console.error('Failed to fetch user from token:', error);
        return EMPTY; // Return an empty observable on error
      })
    );
  }

  // Load the user's trades for a specific league
  private loadUserTrades(userId: number, leagueId: number): Observable<Trade[]> {
    return this.leagueService.getTrades(userId, leagueId).pipe(
      tap((response) => {
        devLog('User trades fetched successfully:', response);
      }),
      catchError((error) => {
        console.error('Failed to fetch user trades:', error);
        return of([]); // Return an empty array on error
      })
    );
  }

  redirectToDraft() {
    this.router.navigate(['dashboard/draft']);
  }

}