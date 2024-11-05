import { CommonModule, NgFor, NgIf } from '@angular/common';
import { Component, EventEmitter, Output } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { VerifyUserService } from '../../user-verification/verify-user.service';
import { UserLeaguesService } from '../league-services/user-leagues/user-leagues.service';
import { firstValueFrom } from 'rxjs';
import { User } from '../../models/user.model';
import { League } from '../../models/league.model';
import { Leagues } from '../../models/leagues.model';

@Component({
  selector: 'app-dashboard-main',
  standalone: true,
  imports: [RouterOutlet, CommonModule, NgFor, NgIf],
  templateUrl: './dashboard-main.component.html',
  styleUrl: './dashboard-main.component.scss'
})
export class DashboardMainComponent {

    leagues: Leagues = { leagues: [] };
    selectedLeague: League | null = null;
    user: string = "User"
    showMenu = false

  constructor(
    private router: Router,
    private userService: VerifyUserService,
    private leagueService: UserLeaguesService
  ) {}

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
  }

  redirectToPortfolio() {
    this.router.navigate(['dashboard/portfolio']);
  }  

  redirectToTrades() {
    this.router.navigate(['dashboard/trades']);
  }

  redirectToCreateLeague() {
    this.router.navigate(['dashboard/create-league']);
  }

  redirectToRemoveLeague() {
    this.router.navigate(['dashboard/remove-league']);
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
            // Assuming 'response' has a 'leagues' property that is an array of 'League' objects
            this.leagues = response;
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
  selectLeague(league: League) {
    this.leagueService.setSelectedLeague(league)
    this.redirectToDashboard();
  }

  // Method to load the user data asynchronously
  private loadUser(): void {
    this.userService.getUserFromToken().subscribe({
      next: (user: User) => {
        console.log('User fetched successfully:', user);
        this.user = user.username;
      },
      error: (error) => {
        console.error('Failed to fetch user from token:', error);
      }
    });
  }

  toggleMenu() {
    this.showMenu = !this.showMenu;
  }

}
