import { Component } from '@angular/core';
import { VerifyUserService } from '../../user-verification/verify-user.service';
import { LeagueService } from '../services/league.service';
import { devLog } from '../../../environments/development/devlog';
import { FormsModule } from '@angular/forms';
import { guard, guardRFC3339 } from '../../utils/guard';

@Component({
  selector: 'app-create-league',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './create-league.component.html',
  styleUrl: './create-league.component.scss'
})
export class CreateLeagueComponent {
  leagueName: string = '';
  endDate: Date | null = null;

  constructor(
    private userService: VerifyUserService,
    private leagueService: LeagueService
  ) {}

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

  private createLeague(leagueName: string, username: number, endDate: string | null): void {    
    guard(leagueName != '', "League name is empty!");
    guardRFC3339(endDate, "End date is required and must be in RFC3339 format");


    devLog('Create League Inputs\n---\nLeague Name: ', leagueName, "\nUsername: ", username, "\nEnd Date: ", endDate)
    this.leagueService.createLeague(leagueName, username, endDate).subscribe({
      next: (response) => {
        // Assuming 'response' has a 'leagues' property that is an array of 'League' objects
        devLog("League Created: ", response)
      },
      error: (error) => {
        devLog('Failed to fetch user leagues:', error);
      }
    });
  }
  
}
