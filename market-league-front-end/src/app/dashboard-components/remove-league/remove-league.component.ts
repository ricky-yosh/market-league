import { NgFor } from '@angular/common';
import { Component } from '@angular/core';
import { League } from '../../models/league.model';
import { LeagueService } from '../services/league.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-remove-league',
  standalone: true,
  imports: [NgFor],
  templateUrl: './remove-league.component.html',
  styleUrl: './remove-league.component.scss'
})
export class RemoveLeagueComponent {
  leagues: League[] = [];
  selectedLeague: League | null = null;
  user: string = "User"

  private subscription!: Subscription;

  constructor(
    private leagueService: LeagueService
  ) {}

  ngOnInit(): void {

    // * Subscribe to the observables to listen for changes

    this.subscription = this.leagueService.userLeagues$.subscribe((leagues) => {
      this.leagues = leagues;
    });

    this.leagueService.getUserLeagues();
  }

  ngOnDestroy(): void {
    // Unsubscribe to avoid memory leaks
    this.subscription.unsubscribe();
  }

  // * User facing functions

  selectLeague(selectedLeague: League) {
    this.selectedLeague = selectedLeague
  }

  removeLeague(leagueToRemove: League) {
    this.leagueService.removeLeague(leagueToRemove.id)
  }

}
