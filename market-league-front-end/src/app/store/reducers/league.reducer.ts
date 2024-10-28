import { createReducer, on } from '@ngrx/store';
import { loadLeaguesSuccess, selectLeague } from '../actions/league.actions';
import { League } from '../../models/league.model';

export interface LeagueState {
    allLeagues: League[]; // Make sure this is an array type
    selectedLeagueId: number | null;
    }

    export const initialState: LeagueState = {
    allLeagues: [], // Initialize as an empty array
    selectedLeagueId: null
    };

    export const leagueReducer = createReducer(
    initialState,
    on(loadLeaguesSuccess, (state, { leagues }) => ({
        ...state,
        allLeagues: leagues // Ensure this is set to an array of leagues
    }))
);

