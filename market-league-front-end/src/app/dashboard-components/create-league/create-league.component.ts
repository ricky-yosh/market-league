import { Component } from '@angular/core';
import { VerifyUserService } from '../services/verify-user.service';
import { LeagueService } from '../services/league.service';
import { devLog } from '../../../environments/development/devlog';
import { FormsModule } from '@angular/forms';
import { guard, guardRFC3339 } from '../../utils/guard';
import { combineLatest, Subscription } from 'rxjs';
import { League } from '../../models/league.model';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-create-league',
  standalone: true,
  imports: [FormsModule, CommonModule],
  templateUrl: './create-league.component.html',
  styleUrl: './create-league.component.scss'
})
export class CreateLeagueComponent {
  leagueName: string = '';
  endDate: Date | null = null;
  allLeagues: League[] = [];
  userLeagues: League[] = [];
  availableLeagues: League[] = [];

  private subscription: Subscription = new Subscription();

  constructor(
    private userService: VerifyUserService,
    private leagueService: LeagueService
  ) {}

  ngOnInit(): void {
    // Subscribe to all leagues updates
    this.subscription.add(
      combineLatest([
        this.leagueService.allLeagues$,
        this.leagueService.userLeagues$
      ]).subscribe(([allLeagues, userLeagues]) => {
        this.allLeagues = allLeagues;
        this.userLeagues = userLeagues;
        
        // Filter out leagues the user is already in
        this.availableLeagues = this.filterAvailableLeagues(allLeagues, userLeagues);
        devLog('Available leagues:', this.availableLeagues);
      })
    );
    
    // Request initial data
    this.leagueService.getAllLeagues();
    this.leagueService.getUserLeagues();
  }

  ngOnDestroy(): void {
    // Unsubscribe to avoid memory leaks
    this.subscription.unsubscribe();
  }

  onSubmit() {
    devLog('Create League Request:', {
      leagueName: this.leagueName,
      endDate: this.endDate,
    });

    guard(this.endDate != null, "End Date is null!");
    const endDateObj = new Date(this.endDate);
    guard(!isNaN(endDateObj.getTime()), "End Date is Invalid!");

    const formattedEndDate = endDateObj.toISOString();

    this.userService.getUserFromToken().subscribe({
      next: (user) => {
        this.createLeague(this.leagueName, user.id, formattedEndDate);
      },
      error: (error) => {
        devLog('Failed to fetch user from token:', error);
      }
    });
  }

  // Join a league when the user clicks the join button
  joinLeague(leagueId: number): void {
    this.leagueService.addUserToLeague(leagueId);
  }

  private createLeague(leagueName: string, user_id: number, endDate: string | null): void {    
    guard(leagueName != '', "League name is empty!");
    guardRFC3339(endDate, "End date is required and must be in RFC3339 format");

    devLog('Create League Inputs\n---\nLeague Name: ', leagueName, "\nUser ID: ", user_id, "\nEnd Date: ", endDate)
    this.leagueService.createLeague(leagueName, user_id, endDate)
  }

  private filterAvailableLeagues(allLeagues: League[], userLeagues: League[]): League[] {
    // Get all user league IDs for quick lookup
    const userLeagueIds = new Set(userLeagues.map(league => league.id));
    
    // Return leagues not in the user's leagues
    return allLeagues.filter(league => !userLeagueIds.has(league.id));
  }
  
  
}
