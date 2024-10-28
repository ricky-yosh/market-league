import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, catchError, tap, throwError } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class VerifyUserService {

  private baseUrl = environment.api_url;
  private verifyUserURL = `${this.baseUrl}/api/auth/user-from-token`;

  constructor(private http: HttpClient) {}

  // Method to get the user information based on JWT token
  getUserFromToken(): Observable<any> {
    const token = localStorage.getItem('token');
    if (!token) {
      return throwError(() => new Error('No token found'));
    }
  
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
  
    console.log('Making HTTP GET request to verify user with token:', token);
    return this.http.get(this.verifyUserURL, { headers }).pipe(
      tap(response => console.log('Received user from token:', response)), // Debugging log
      catchError(error => {
        console.error('Error in getUserFromToken:', error);
        return throwError(() => error);
      })
    );
  }
  
}
