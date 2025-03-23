import { DraftStatus } from './draft-status.enum';

export interface PlayerDraftStatus {
    player_id: number;
    status: DraftStatus;
}
