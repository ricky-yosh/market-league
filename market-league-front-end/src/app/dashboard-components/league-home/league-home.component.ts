import { NgFor } from '@angular/common';
import { ChangeDetectorRef, Component } from '@angular/core';
import { UserLeaguesService } from '../league-services/user-leagues/user-leagues.service';

@Component({
  selector: 'app-league-home',
  standalone: true,
  imports: [NgFor],
  templateUrl: './league-home.component.html',
  styleUrl: './league-home.component.scss'
})
export class LeagueHomeComponent {
  leagueMembers = [
    'John Doe', 'Jane Smith', 'Alice Johnson', 'Bob Lee', 
    'Charlie Brown', 'David Beckham', 'Elon Musk', 
    'Serena Williams', 'Michael Jordan', 'Usain Bolt', 
    'Marie Curie', 'Isaac Newton', 'Albert Einstein'
  ];

  userPortfolio = [
    'John Doe', 'Jane Smith', 'Alice Johnson', 'Bob Lee', 
    'Charlie Brown', 'David Beckham', 'Elon Musk', 
    'Serena Williams', 'Michael Jordan', 'Usain Bolt', 
    'Marie Curie', 'Isaac Newton', 'Albert Einstein'
  ];

  selectedLeague: string | null = null;
  constructor(
    private leagueService: UserLeaguesService,
    private cd: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    // Subscribe to the league changes
    this.leagueService['selectedLeagueSource'].subscribe(league => {
      this.selectedLeague = league;
      this.cd.detectChanges();
    });
  }
  
}
