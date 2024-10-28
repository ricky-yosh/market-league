import { CommonModule, NgFor, NgIf } from '@angular/common';
import { Component, EventEmitter, Output } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { VerifyUserService } from '../../user-verification/verify-user.service';
import { UserLeaguesService } from '../league-services/user-leagues.service';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'app-dashboard-main',
  standalone: true,
  imports: [RouterOutlet, CommonModule, NgFor, NgIf],
  templateUrl: './dashboard-main.component.html',
  styleUrl: './dashboard-main.component.scss'
})
export class DashboardMainComponent {

  constructor(
    private router: Router,
    private userService: VerifyUserService,
    private leagueService: UserLeaguesService
  ) {}

  leagues: string[] = [];
  selectedLeague: string | null = null;
  user: string = "User"

  ngOnInit(): void {
    this.loadUserLeagues();
    this.loadUser();
  }

  // Routing
  redirectToHome() {
    this.router.navigate(['/home']);
  }

  redirectToDashboard() {
    this.router.navigate(['/dashboard']);
  }
  
  redirectToLeaderboard() {
    this.router.navigate(['dashboard/leaderboard']);
    // current_league = id_of_league
  }

  redirectToPortfolio() {
    this.router.navigate(['dashboard/portfolio']);
  }  

  redirectToTrades() {
    this.router.navigate(['dashboard/trades']);
  }

  // Method to load the leagues for the user
  private loadUserLeagues(): void {
    // Step 1: Get the user from the token
    this.userService.getUserFromToken().subscribe({
      next: (user) => {
        const userId = user.id;

        // Step 2: Fetch leagues based on the user's ID
        this.leagueService.getUserLeagues(userId).subscribe({
          next: (response) => {
            this.leagues = response.leagues.map((league: any) => league.league_name);
          },
          error: (error) => {
            console.error('Failed to fetch user leagues:', error);
          }
        });
      },
      error: (error) => {
        console.error('Failed to fetch user from token:', error);
      }
    });
  }

  // Method to handle league selection
  selectLeague(league: string) {
    this.leagueService.setSelectedLeague(league)
  }

  // Method to fetch the user using async/await
  private async getUser(): Promise<any> {
    return firstValueFrom(this.userService.getUserFromToken());
  }

  // Method to load the user data asynchronously
  private loadUser(): void {
    this.userService.getUserFromToken().subscribe({
      next: (user) => {
        console.log('User fetched successfully:', user);
        this.user = user.username;
      },
      error: (error) => {
        console.error('Failed to fetch user from token:', error);
      }
    });
  }

}
