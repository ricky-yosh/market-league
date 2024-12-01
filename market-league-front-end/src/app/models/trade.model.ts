import { User } from "./user.model";

export interface Trade {
    id: number;
    league_id: number;
    user1: User;
    user2: User;
    portfolio1_id: number;
    portfolio2_id: number;
    stocks1: any[] | null; // Adjust the type as needed
    stocks2: any[] | null; // Adjust the type as needed
    user1_confirmed: boolean;
    user2_confirmed: boolean;
    status: string;
    created_at: string; // ISO date string
    updated_at: string | null; // ISO date string or null
}