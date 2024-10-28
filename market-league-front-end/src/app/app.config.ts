import { ApplicationConfig, provideZoneChangeDetection, isDevMode } from '@angular/core';
import { provideHttpClient } from '@angular/common/http';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideStore } from '@ngrx/store';
import { provideEffects } from '@ngrx/effects';
import { provideStoreDevtools } from '@ngrx/store-devtools';

import { reducers } from './store/reducers'; // Root reducer
import { LeagueEffects } from './store/effects/league.effects';

export const appConfig: ApplicationConfig = {
  providers: [provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes),
    provideHttpClient(),

    // State Manager
    provideStore(reducers), // Register the root reducer
    provideEffects([LeagueEffects]), // Register the NgRx effects

    provideStoreDevtools({ maxAge: 25, logOnly: !isDevMode()})]
};
