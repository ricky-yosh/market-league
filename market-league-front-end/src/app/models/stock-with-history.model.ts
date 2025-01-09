import { PriceHistory } from "./price-history.model";

export interface StockWithHistory {
    id: number;
    ticker_symbol: string;
    company_name: string;
    current_price: number;
    price_histories: PriceHistory[];
}