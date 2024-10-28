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
  private selectedLeagueSource = new BehaviorSubject<string | null>(localStorage.getItem('selectedLeague'));
  selectedLeague$ = this.selectedLeagueSource.asObservable();

  constructor(private http: HttpClient) {
    const storedLeague = localStorage.getItem('selectedLeague');
    console.log('Service Constructor: Stored League:', storedLeague);
    this.selectedLeagueSource = new BehaviorSubject<string | null>(storedLeague ? storedLeague : null);

    this.selectedLeagueSource.subscribe(value => {
      console.log(`Service: selectedLeagueSource emitted value: ${value}`);
    });

  }

  // Method to get the user's leagues using their user ID
  getUserLeagues(userId: number): Observable<any> {
    return this.http.post(this.findUserLeagues, { user_id: userId });
  }

  // Method to set the selected league
  setSelectedLeague(league: string | null): void {
    this.selectedLeagueSource.next(league ? `${league}` : null); // Use a shallow copy or different reference
    if (league) {
      localStorage.setItem('selectedLeague', league);
    } else {
      localStorage.removeItem('selectedLeague');
    }
  }

}
