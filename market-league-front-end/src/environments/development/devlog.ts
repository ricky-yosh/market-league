import { environment } from "../environment";


export function devLog(...args: any[]): void {
    if (!environment.production && environment.enableDebug) {
        console.log(...args);
    }
}