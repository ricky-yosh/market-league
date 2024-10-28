import { createReducer, on } from '@ngrx/store';
import { loadLeaguesSuccess, selectLeague } from '../actions/league.actions';
import { League } from '../../models/league.model';

export interface LeagueState {
    allLeagues: League[];
    selectedLeagueId: number | null;
}

export const initialState: LeagueState = {
    allLeagues: [],
    selectedLeagueId: null
};

export const leagueReducer = createReducer(
    initialState,
    on(loadLeaguesSuccess, (state, { leagues }) => ({
        ...state,
        allLeagues: leagues
    })),
        on(selectLeague, (state, { leagueId }) => ({
            ...state,
            selectedLeagueId: leagueId
    }))
);
