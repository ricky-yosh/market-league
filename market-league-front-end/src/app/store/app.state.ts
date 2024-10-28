import { LeagueState } from './reducers/league.reducer';  // Import the LeagueState type

// Define the shape of your app's root state
export interface AppState {
  leagues: LeagueState;  // Add other feature states here as your app grows
}
