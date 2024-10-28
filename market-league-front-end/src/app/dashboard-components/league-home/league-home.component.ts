import { NgFor } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { UserLeaguesService } from '../league-services/user-leagues/user-leagues.service';
import { User } from '../../models/user.model';
import { League } from '../../models/league.model';
import { Stock } from '../../models/stock.model';
import { Portfolio } from '../../models/portfolio.model';
import { VerifyUserService } from '../../user-verification/verify-user.service';

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
    this.leagueService.selectedLeague$.subscribe({
      next: (league) => {
        console.log('Selected league updated:', league);
        this.selectedLeague = league;
        console.log("Selected League: " + this.selectedLeague?.league_name)
        this.cd.detectChanges(); // Detect changes after updating selectedLeague

        // Load league members once the selected league has been updated
        if (this.selectedLeague?.id) {
          this.loadLeagueMembers(this.selectedLeague.id);
          this.loadUser();
          console.log("User ID: " + this.user?.id)
          if (this.user?.id) {
            this.loadUserPortfolio(this.user.id, this.selectedLeague.id)
          }
        }
      },
      error: (error) => {
        console.error('Failed to fetch selected league:', error);
      }
    });
  }

  // Method to load members of a selected league
  private loadLeagueMembers(leagueId: number): void {
    this.leagueService.getLeagueMembers(leagueId).subscribe({
      next: (members) => {
        console.log('League members fetched successfully:', members);
        this.leagueMembers = members.map((member: User) => member.username);
      },
      error: (error) => {
        console.error('Failed to fetch league members:', error);
      }
    });
  }

  // Load the user's portfolio for a specific league
  private loadUserPortfolio(userId: number, leagueId: number) {
    this.leagueService.getUserPortfolio(userId, leagueId).subscribe({
      next: (response: Portfolio) => {
        // Extract the stocks array from the response
        this.userPortfolio = response.stocks;
        console.log('User portfolio fetched successfully:', this.userPortfolio);
      },
      error: (error) => {
        console.error('Failed to fetch user portfolio:', error);
      }
    });
  }

  private loadUser(): void {
    this.userService.getUserFromToken().subscribe({
      next: (user: User) => {
        console.log('User fetched successfully:', user);
        this.user = user;
      },
      error: (error) => {
        console.error('Failed to fetch user from token:', error);
      }
    });
  }

}