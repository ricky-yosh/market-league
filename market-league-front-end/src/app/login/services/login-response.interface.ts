export interface LoginResponse {
    token: string;
    username: string;
    message?: string; // Optional message, could be for errors or information
}