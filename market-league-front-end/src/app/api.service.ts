import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface CountResponse {
  value: number;
}

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private baseUrl = 'http://localhost:9000'; // Go server URL

  constructor(private http: HttpClient) {}

  // Function to call the /ping endpoint
  increaseCount(): Observable<any> {
    return this.http.get(`${this.baseUrl}/api/increment`);
  }
}