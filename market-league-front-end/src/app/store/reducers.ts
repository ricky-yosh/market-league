import { ActionReducerMap } from '@ngrx/store';
import { AppState } from './app.state';
import { leagueReducer } from './reducers/league.reducer';

// Combine all feature reducers here
export const reducers: ActionReducerMap<AppState> = {
    leagues: leagueReducer
};
