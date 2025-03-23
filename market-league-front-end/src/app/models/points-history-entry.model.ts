import { Portfolio } from "./portfolio.model";

// Define an interface for a League
export interface PortfolioPointsHistoryEntry {
    id: number;
    portfolio_id: number;
    portfolio: Portfolio;
    points: number;
    recorded_at: string;
}