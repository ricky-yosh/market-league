import { Portfolio } from "./portfolio.model";
import { Stock } from "./stock.model";

// Define an interface for a League
export interface StockHistoryEntry {
    id: number;
    portfolio_id: number;
    portfolio: Portfolio;
    stock_id: number;
    stock: Stock;
    starting_value: number;
    current_value: number;
    start_date: string;
    end_date: string | null;
}