import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SignUpService {

  private baseUrl = environment.api_url
  private loginURL = `${this.baseUrl}/api/auth/signup`

  constructor(private http: HttpClient) {}

  // Login call
  signUp(credentials: { username: string; password: string }): Observable<any> {
    return this.http.post<any>(this.loginURL, credentials);
  }

}
