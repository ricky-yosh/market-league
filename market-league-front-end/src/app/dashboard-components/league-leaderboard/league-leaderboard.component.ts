import { Component, OnInit } from '@angular/core';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { CommonModule } from '@angular/common'; // Import CommonModule
import { AppState } from '../../store/app.state';
import { League } from '../../models/league.model';
import { selectAllLeagues } from '../../store/selectors/league.selectors';
import { loadLeagues } from '../../store/actions/league.actions';

@Component({
  standalone: true,
  imports: [CommonModule], // Include CommonModule here
  selector: 'app-league-state-view',
  template: `
    <h2>League State View</h2>
    <button (click)="loadLeagues()">Load Leagues</button>
    <ul>
      <li *ngFor="let league of leagues$ | async">
        {{ league.name }} (ID: {{ league.id }})
      </li>
    </ul>
  `
})
export class LeagueLeaderboardComponent implements OnInit {
  leagues$: Observable<League[]>;

  constructor(private store: Store<AppState>) {
    this.leagues$ = this.store.select(selectAllLeagues);
  }

  ngOnInit(): void {
    // Dispatch an action to load leagues when the component initializes
    this.loadLeagues();
  }

  loadLeagues(): void {
    this.store.dispatch(loadLeagues());
  }
}
