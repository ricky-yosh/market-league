import { NgFor } from '@angular/common';
import { CommonModule } from '@angular/common';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { LeagueService } from '../services/league.service';
import { League } from '../../models/league.model';
import { Subscription } from 'rxjs';
import { VerifyUserService } from '../services/verify-user.service';
import { User } from '../../models/user.model';
import { LeaderboardUser } from '../../models/websocket-responses/league/leaderboard-user';

@Component({
  selector: 'app-league-completed-screen',
  standalone: true,
  imports: [NgFor, CommonModule],
  templateUrl: './league-completed-screen.component.html',
  styleUrl: './league-completed-screen.component.scss'
})
export class LeagueCompletedScreenComponent implements OnInit, OnDestroy {
  
  selectedLeague: League | null = null;
  
  // Leaderboard data
  leaderboard: LeaderboardUser[] = [];
  completeLeaderboard: {username: string, total_value: number, rank: number, isWinner: boolean}[] = [];
  currentUser: string = '';
  showConfetti: boolean = false;
  leaderboardLoaded = false;

  private subscriptions: Subscription = new Subscription();

  constructor(
    private leagueService: LeagueService,
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
    
    // Leaderboard data
    this.subscriptions.add(
      this.leagueService.leaderboard$.subscribe((leaderboardData) => {
        this.leaderboard = leaderboardData;
        this.leaderboardLoaded = true;
        
        // Process leaderboard data
        this.processLeaderboard();
      })
    );
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

    // Add rankings and identify winner
    this.completeLeaderboard = [];
    let currentRank = 1;
    let previousValue: number | null = null;
    
    sortedLeaderboard.forEach((member, index) => {
      // If not the first member and value is different from previous
      if (index > 0 && member.total_value !== previousValue) {
        currentRank = index + 1;
      }
      
      const isWinner = currentRank === 1;
      
      // Add to complete leaderboard
      this.completeLeaderboard.push({
        username: member.username,
        total_value: member.total_value,
        rank: currentRank,
        isWinner: isWinner
      });
      
      // If this is the current user and they won, show confetti
      if (member.username === this.currentUser && isWinner) {
        this.showConfetti = true;
      }
      
      previousValue = member.total_value;
    });
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