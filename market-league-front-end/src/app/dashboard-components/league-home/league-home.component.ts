import { NgFor } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { UserLeaguesService } from '../league-services/user-leagues/user-leagues.service';
import { User } from '../../models/user.model';
import { League } from '../../models/league.model';
import { Stock } from '../../models/stock.model';
import { Portfolio } from '../../models/portfolio.model';
import { VerifyUserService } from '../../user-verification/verify-user.service';
import { EMPTY, Observable, catchError, map, of, switchMap, tap } from 'rxjs';

@Component({
  selector: 'app-league-home',
  standalone: true,
  imports: [NgFor],
  templateUrl: './league-home.component.html',
  styleUrl: './league-home.component.scss'
})
export class LeagueHomeComponent implements OnInit {
  
  userPortfolio: Stock[] = [];
  leagueMembers: string[] = [];
  selectedLeague: League | null = null;
  user: User | null = null;

  constructor(
    private leagueService: UserLeaguesService,
    private userService: VerifyUserService,
    private cd: ChangeDetectorRef
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
              // Now that we have the user, load the portfolio
              if (user?.id && this.selectedLeague?.id) {
                return this.loadUserPortfolio(user.id, this.selectedLeague.id);
              }
              return of([]); // Return an empty array if user or league ID is not available
            })
          );
        }
        return EMPTY; // If no league selected, return an empty observable
      })
    ).subscribe({
      next: (portfolio) => {
        this.userPortfolio = portfolio;
        console.log('Portfolio loaded successfully');
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

}