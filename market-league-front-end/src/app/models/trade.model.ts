export interface Trade {
    ID: number;
    league_id: number;
    player1_id: number;
    player2_id: number;
    player1_portfolio_id: number;
    player2_portfolio_id: number;
    player1_stocks: any[] | null; // Adjust the type as needed
    player2_stocks: any[] | null; // Adjust the type as needed
    player1_confirmed: boolean;
    player2_confirmed: boolean;
    created_at: string; // ISO date string
    confirmed_at: string | null; // ISO date string or null
}