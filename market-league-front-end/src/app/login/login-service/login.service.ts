import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable, switchMap, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { LoginResponse } from './login-response.interface';
import { VerifyUserService } from '../../dashboard-components/services/verify-user.service'; // Import your user service

@Injectable({
  providedIn: 'root',
})
export class LoginService {
  private baseUrl = environment.api_url
  private loginURL = `${this.baseUrl}/api/auth/login`

  constructor(
    private http: HttpClient,
    private userService: VerifyUserService // Inject user service
  ) {}

  // Login call with user loading
  login(credentials: { username: string; password: string }): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(this.loginURL, credentials).pipe(
      tap(response => {
        // Store token
        sessionStorage.setItem('token', response.token);
      }),
      switchMap(response => {
        // Load user data right after login
        return this.userService.getUserFromToken().pipe(
          // Return original response after user is loaded
          map(() => response)
        );
      })
    );
  }
}