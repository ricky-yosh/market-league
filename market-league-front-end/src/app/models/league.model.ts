import { User } from './user.model';

export interface League {
    id: number;
    name: string;
    startDate: string;        // Start date of the league in ISO format
    endDate: string;          // End date of the league in ISO format
    users: User[];            // List of users participating in the league
}
