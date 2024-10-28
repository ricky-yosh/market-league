import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { BehaviorSubject, Observable, catchError, map, throwError } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { League } from '../../models/league.model';

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
  getUserLeagues(userId: number): Observable<League[]> {
    return this.http.post<{ leagues: League[] | undefined }>(this.findUserLeagues, { user_id: userId }).pipe(
      map(response => {
        // Add a check to ensure that `response` and `response.leagues` are not undefined
        if (response && response.leagues) {
          return response.leagues; // Return the leagues array if it exists
        } else {
          console.warn('Leagues property is missing or undefined in the response', response);
          return []; // Return an empty array as a fallback
        }
      }),
      catchError(error => {
        console.error('Error fetching user leagues:', error);
        return throwError(() => error);
      })
    );
  }  

  // Method to set the selected league
  setSelectedLeague(league: string | null): void {
    console.log(`Setting selected league to: ${league}`);
    this.selectedLeagueSource.next(league ? `${league}` : null); // Use a shallow copy or different reference
    console.log(`Emitted new league: ${this.selectedLeagueSource.value}`);
    if (league) {
      localStorage.setItem('selectedLeague', league);
    } else {
      localStorage.removeItem('selectedLeague');
    }
  }

}
