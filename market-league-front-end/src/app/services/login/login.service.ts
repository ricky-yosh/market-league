import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { LoginResponse } from './login-response.interface';

@Injectable({
  providedIn: 'root',
})

export class LoginService {

  private baseUrl = environment.api_url
  private loginURL = `${this.baseUrl}/api/auth/login`

  constructor(private http: HttpClient) {}

  // Login call
  login(credentials: { username: string; password: string }): Observable<any> {
    return this.http.post<LoginResponse>(this.loginURL, credentials);
  }
}