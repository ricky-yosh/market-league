import { Component } from '@angular/core';
import { VerifyUserService } from '../../user-verification/verify-user.service';
import { UserLeaguesService } from '../league-services/user-leagues/user-leagues.service';
import { devLog } from '../../../environments/development/devlog';
import { FormsModule } from '@angular/forms';

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
    private leagueService: UserLeaguesService
  ) {}

  onSubmit() {
    devLog('League Created:', {
      leagueName: this.leagueName,
      endDate: this.endDate,
    });
  }

}
