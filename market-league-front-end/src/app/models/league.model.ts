// Define an interface for a League
export interface League {
    id: number;
    league_name: string;
    start_date: string;
    end_date: string;
    users: any[] | null;
}