import { NgFor } from '@angular/common';
import { Component } from '@angular/core';
import { Leagues } from '../../models/leagues.model';
import { League } from '../../models/league.model';
import { VerifyUserService } from '../../user-verification/verify-user.service';
import { UserLeaguesService } from '../league-services/user-leagues/user-leagues.service';
import { User } from '../../models/user.model';

@Component({
  selector: 'app-remove-league',
  standalone: true,
  imports: [NgFor],
  templateUrl: './remove-league.component.html',
  styleUrl: './remove-league.component.scss'
})
export class RemoveLeagueComponent {
  leagues: Leagues = { leagues: [] };
  selectedLeague: League | null = null;
  user: string = "User"

  constructor(
    private userService: VerifyUserService,
    private leagueService: UserLeaguesService
  ) {}

  ngOnInit(): void {
    this.loadUserLeagues();
    this.loadUser();
  }

  selectLeague(selectedLeague: League) {
    this.selectedLeague = selectedLeague
  }

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

  removeLeague(leagueToRemove: League) {

  }

}
