import { Component } from '@angular/core';
import { DraftService } from '../services/draft.service';
import { League } from '../../models/league.model';
import { Subscription } from 'rxjs';
import { LeagueService } from '../services/league.service';
import { DraftStatus } from '../../models/draft-status.enum';
import { LeagueState } from '../../models/league-state.model';
import { Router } from '@angular/router';

@Component({
  selector: 'app-league-draft-queue',
  standalone: true,
  imports: [],
  templateUrl: './league-draft-queue.component.html',
  styleUrl: './league-draft-queue.component.scss'
})
export class LeagueDraftQueueComponent {

  selectedLeague: League | null = null;

  private subscription!: Subscription;

  constructor(
    private router: Router,
    private draftService: DraftService,
    private leagueService: LeagueService,
  ) {}

  ngOnInit(): void {
    // Selected League
    this.subscription = this.leagueService.selectedLeague$.subscribe((league) => {
      this.selectedLeague = league;
      switch(league?.league_state) {
        case LeagueState.InDraft: {
          this.redirectToDraft()
          break;
        }
        case LeagueState.PostDraft: {
          this.redirectToDashboard()
          break;
        }
        default: {
          // Stay on draft
        }
      }
    });

    this.loadLeagueDetails()
  }

  ngOnDestroy(): void {
    // Unsubscribe to avoid memory leaks
    this.subscription.unsubscribe();
  }

  queuePlayer() {
    this.draftService.queuePlayerForDraft();
  }

  getPlayerCount(): number {
    if (!this.selectedLeague?.league_players) return 0;
    
    return this.selectedLeague.league_players.filter(
      player => player.draft_status === DraftStatus.DraftNotReady || 
                player.draft_status === DraftStatus.DraftReady
    ).length;
  }

  getQueuedPlayersCount(): number {
    if (!this.selectedLeague?.league_players) return 0;
    
    return this.selectedLeague.league_players.filter(
      player => player.draft_status === DraftStatus.DraftReady
    ).length;
  }

  // Load the league details for a specific league
  private loadLeagueDetails(): void {
    this.leagueService.getLeagueDetails();
  }

  // ** Routes **

  redirectToDraft() {
    this.router.navigate(['/dashboard/draft']);
  }

  redirectToDashboard() {
    this.router.navigate(['/dashboard']);
  }

}
