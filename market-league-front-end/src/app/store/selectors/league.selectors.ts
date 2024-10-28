import { createFeatureSelector, createSelector } from '@ngrx/store';
import { LeagueState } from '../reducers/league.reducer';

export const selectLeaguesState = createFeatureSelector<LeagueState>('leagues');

export const selectAllLeagues = createSelector(
    selectLeaguesState,
    (state: LeagueState) => state.allLeagues
);

export const selectSelectedLeagueId = createSelector(
    selectLeaguesState,
    (state: LeagueState) => state.selectedLeagueId
);

export const selectSelectedLeague = createSelector(
    selectAllLeagues,
    selectSelectedLeagueId,
    (leagues, selectedLeagueId) => leagues.find(league => league.id === selectedLeagueId)
);
