import { Injectable, isDevMode } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../environments/environment';

export interface CountResponse {
  value: number;
}

@Injectable({
  providedIn: 'root',
})
export class ApiService {

  private baseUrl = environment.api_url

  constructor(private http: HttpClient) {}

  // Function to call the /ping endpoint
  increaseCount(): Observable<any> {
    return this.http.get(`${this.baseUrl}/api/increment`);
  }
}