import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { User } from '../models/user.model';
import { devLog } from '../../environments/development/devlog';

@Injectable({
  providedIn: 'root'
})
export class VerifyUserService {

  private baseUrl = environment.api_url
  private verifyUserURL = `${this.baseUrl}/api/auth/user-from-token`

  constructor(private http: HttpClient) {}

  // Method to get the user information based on JWT token
  getUserFromToken(): Observable<User> {
    const token = localStorage.getItem('token');
    if (!token) {
      throw new Error('No token found');
    }

    const headers = new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });

    return this.http.get<User>(this.verifyUserURL, { headers });
  }

}
