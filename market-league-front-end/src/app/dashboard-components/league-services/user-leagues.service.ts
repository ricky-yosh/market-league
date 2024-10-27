import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { BehaviorSubject, Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class UserLeaguesService {

  private baseUrl = environment.api_url
  private findUserLeagues = `${this.baseUrl}/api/users/user-leagues`
  private selectedLeagueSubject = new BehaviorSubject<string | null>(null);
  selectedLeague$ = this.selectedLeagueSubject.asObservable();

  constructor(private http: HttpClient) {
    const storedLeague = localStorage.getItem('selectedLeague');
    console.log('Service Constructor: Stored League:', storedLeague);
    this.selectedLeagueSubject = new BehaviorSubject<string | null>(storedLeague ? storedLeague : null);
  }

  // Method to get the user's leagues using their user ID
  getUserLeagues(userId: number): Observable<any> {
    return this.http.post(this.findUserLeagues, { user_id: userId });
  }

  // Method to set the selected league
  setSelectedLeague(league: string) {
    this.selectedLeagueSubject.next(league);
    localStorage.setItem('selectedLeague', league); // Optionally persist to local storage
  }

  // Method to get the current selected league (synchronously)
  getSelectedLeague(): string | null {
    return this.selectedLeagueSubject.value;
  }

}
