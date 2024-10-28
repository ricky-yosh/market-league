import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { VerifyUserService } from '../../services/verify-user/verify-user.service';
import { UserLeaguesService } from '../../services/league-services/user-leagues.service';
import { loadLeagues, loadLeaguesSuccess, loadLeaguesFailure } from '../actions/league.actions';
import { mergeMap, map, catchError, switchMap } from 'rxjs/operators';
import { of } from 'rxjs';

@Injectable()
export class LeagueEffects {
    constructor(
        private actions$: Actions,
        private verifyUserService: VerifyUserService,
        private leagueService: UserLeaguesService
    ) {}

    loadLeagues$ = createEffect(() =>
        this.actions$.pipe(
            ofType(loadLeagues),
            switchMap(() => {
            console.log('loadLeagues action triggered');
            return this.verifyUserService.getUserFromToken().pipe(
                map(user => {
                console.log('User obtained from token:', user);
                return user.id;
                }),
                switchMap(userId => {
                console.log('Fetching leagues for user ID:', userId);
                return this.leagueService.getUserLeagues(userId).pipe(
                    map(leagues => {
                    console.log('Leagues fetched successfully:', leagues);
                    return loadLeaguesSuccess({ leagues });
                    }),
                    catchError(error => {
                    console.error('Error fetching leagues:', error);
                    return of(loadLeaguesFailure({ error }));
                    })
                );
                }),
                catchError(error => {
                console.error('Error in verifyUserService.getUserFromToken:', error);
                return of(loadLeaguesFailure({ error: error.message }));
                })
            );
            })
        )
    );
}
