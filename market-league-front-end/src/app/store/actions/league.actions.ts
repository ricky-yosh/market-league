import { createAction, props } from '@ngrx/store';
import { League } from '../../models/league.model';

// Load Leagues action
export const loadLeagues = createAction('[League] Load Leagues');

// Successful load action
export const loadLeaguesSuccess = createAction(
    '[League] Load Leagues Success',
    props<{ leagues: League[] }>()
);

// Failed load action
export const loadLeaguesFailure = createAction(
    '[League] Load Leagues Failure',
    props<{ error: any }>()
);

// Action to select a league
export const selectLeague = createAction(
    '[League] Select League',
    props<{ leagueId: number }>()
);
