import { Component, OnInit } from '@angular/core';
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
export class CreateLeagueComponent implements OnInit {
  leagueName: string = '';
  endDate: string = ''; // Change to string to match the input format
  allLeagues: League[] = [];
  userLeagues: League[] = [];
  availableLeagues: League[] = [];
  minEndDate: string = '';

  private subscription: Subscription = new Subscription();

  constructor(
    private userService: VerifyUserService,
    private leagueService: LeagueService
  ) {}

  ngOnInit(): void {
    // Set minimum end date to one week from today
    this.setMinEndDate();
    
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

  // Set the minimum end date to one week from today
  private setMinEndDate(): void {
    const today = new Date();
    const nextWeek = new Date(today);
    nextWeek.setDate(today.getDate() + 7);
    
    // Format the date as YYYY-MM-DD for the input element
    this.minEndDate = nextWeek.toISOString().split('T')[0];
    
    // Set the default end date to the minimum (one week from today)
    this.endDate = this.minEndDate;
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

    // Validate the date is provided and valid
    guard(this.endDate != null && this.endDate !== '', "End Date is required!");
    const endDateObj = new Date(this.endDate);
    guard(!isNaN(endDateObj.getTime()), "End Date is Invalid!");

    // Validate the end date is at least a week from today
    const today = new Date();
    const minDate = new Date(today);
    minDate.setDate(today.getDate() + 7);
    
    guard(endDateObj >= minDate, "End Date must be at least a week from today!");

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

  // Add a method to validate the date on change
  onDateChange(): void {
    // If date is before minimum date, reset it to the minimum date
    if (this.endDate < this.minEndDate) {
      this.endDate = this.minEndDate;
    }
  }

  // Rest of your code remains the same
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