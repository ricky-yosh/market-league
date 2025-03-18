import { DraftStatus } from './draft-status.enum';

export interface LeaguePlayer {
    league_id: number;
    player_id: number;
    draft_status: DraftStatus;
}