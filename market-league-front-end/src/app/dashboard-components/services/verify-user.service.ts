import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { BehaviorSubject, catchError, Observable, tap, throwError } from 'rxjs';
import { User } from '../../models/user.model';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class VerifyUserService {
  private baseUrl = environment.api_url;
  private verifyUserURL = `${this.baseUrl}/api/auth/user-from-token`;

  // Initialize BehaviorSubject with a value from sessionStorage, if available
  private currentUserSubject = new BehaviorSubject<User | null>(this.getStoredUser());
  public currentUser$ = this.currentUserSubject.asObservable(); // Expose observable

  constructor(
    private http: HttpClient,
    private router: Router,
  ) {}

  // * Getters

  // Fetch user and persist data
  getUserFromToken(): Observable<User> {
    const token = sessionStorage.getItem('token');
    if (!token) {
      this.router.navigate(['/logged-out']);
      throw new Error('No token found');
    }

    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
    });

    return this.http.get<User>(this.verifyUserURL, { headers }).pipe(
      tap((user) => {
        this.setCurrentUser(user); // Update BehaviorSubject and persist to storage
      }),
      catchError((error) => {
        console.error('Error fetching user:', error);
        this.router.navigate(['/logged-out']);
        return throwError(() => new Error('Failed to fetch user.'));
      })
    );
  }

  getStoredUser(): User | null {
    return JSON.parse(sessionStorage.getItem('currentUser') || 'null') // Load from storage
  }

  // Get current user value
  getCurrentUserValue(): User | null {
    return this.currentUserSubject.value;
  }

  // * Setters

  // Set user and persist it
  setCurrentUser(user: User | null): void {
    this.currentUserSubject.next(user); // Update BehaviorSubject
    sessionStorage.setItem('currentUser', JSON.stringify(user)); // Persist to storage
  }

  // * Clear User

  // Clear user and storage (e.g., logout)
  clearUser(): void {
    this.currentUserSubject.next(null); // Clear BehaviorSubject
    // sessionStorage.removeItem('currentUser'); // Clear storage
    // sessionStorage.removeItem('token'); // Optionally clear token
    sessionStorage.clear()// Clear storage
  }

}