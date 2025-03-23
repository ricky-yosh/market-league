import { CommonModule, NgFor, NgIf } from '@angular/common';
import { Component } from '@angular/core';
import { VerifyUserService } from '../services/verify-user.service';
import { LeagueService } from '../services/league.service';
import { User } from '../../models/user.model';
import { League } from '../../models/league.model';
import { devLog } from '../../../environments/development/devlog';

@Component({
  selector: 'app-league-leaderboard',
  standalone: true,
  imports: [CommonModule, NgFor, NgIf],
  templateUrl: './league-leaderboard.component.html',
  styleUrl: './league-leaderboard.component.scss'
})
export class LeagueLeaderboardComponent {
  leagues: League[] = [];
  selectedLeague: League | null = null;
  user: string = "User"
  members: any[] = [];
  membersWithRank: {username: string, total_value: number, rank: number}[] = [];
  membersLoaded = false;


  constructor(
    private userService: VerifyUserService,
    private leagueService: LeagueService
  ) {
    this.leagueService.leaderboard$.subscribe((members: any[]) => {
      this.members = members;
      this.membersLoaded = true;
      
      // Only process members after they're loaded
      this.orderizeMembers();
      this.addRanking();
    });
  }

  async ngOnInit() {
    // Loads the user
    this.loadUser();
    // Gets all of the user's leagues
    this.loadUserLeagues();
    // Not sure if getStoredLeague better
    // Check if possibility of no league selected, or no league at all
    this.selectedLeague = this.leagueService.getSelectedLeagueValue();

    let leagueId = (typeof this.selectedLeague?.id === 'number') ? this.selectedLeague.id : 1;
    this.leagueService.getLeagueLeaderboard(leagueId);
  }

  // * Helper Methods
  
  // Method to load the leagues for the user
  private loadUserLeagues(): void {
    // Fetch leagues
    this.leagueService.getUserLeagues();
  }

  // Method to load the user data asynchronously
  private loadUser(): void {
    this.userService.getUserFromToken().subscribe({
      next: (user: User) => {
        devLog('User fetched successfully:', user);
        this.user = user.username;
      },
      error: (error) => {
        console.error('Failed to fetch user from token:', error);
      }
    });
  }

  
  private orderizeMembers(): void {
    // Make sure we have members to sort
    if (!this.members || this.members.length === 0) {
      console.warn("No members to order");
      return;
    }

    this.members = this.members.sort((member1, member2) => {
      const value1 = typeof member1.total_value === 'number' ? member1.total_value : 0;
      const value2 = typeof member2.total_value === 'number' ? member2.total_value : 0;
      return value2 - value1; // Descending order
    });
  }

  // getting ranking if same total value then same ranking
  private addRanking(): void {// Clear previous rankings
    this.membersWithRank = [];

    // if no member skipping
    if (!this.members || this.members.length === 0) 
    {
      console.warn("No members found, skipping ranking.");
      return;
    }

    for (let memberIndex = 0; memberIndex < this.members.length; memberIndex++)
    {
      // if there is a prior element
      if (this.members[memberIndex-1])
      {
        // if they hace the same value as the predecessor
        if (this.members[memberIndex-1].total_value == this.members[memberIndex].total_value)
        {
          // then push same ranking as prior
          this.membersWithRank.push({
              username: this.members[memberIndex].username,
              total_value: this.members[memberIndex].total_value,
              rank: this.membersWithRank[memberIndex-1].rank
            });
        }

        // If not equal them have their own ranking
        else 
        {
          this.membersWithRank.push({
              username: this.members[memberIndex].username,
              total_value: this.members[memberIndex].total_value,
              rank: memberIndex+1
            });
        }

      // else if first element, rank 1
      } else {
        this.membersWithRank.push({
            username: this.members[memberIndex].username,
            total_value: this.members[memberIndex].total_value,
            rank: 1
          });
      }
    }
  }
}
