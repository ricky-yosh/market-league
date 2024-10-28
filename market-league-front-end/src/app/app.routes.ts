import { Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { AboutComponent } from './about/about.component';
import { HomeComponent } from './home/home.component';
import { SignUpComponent } from './sign-up/sign-up.component';
import { DashboardMainComponent } from './dashboard-components/dashboard-main/dashboard-main.component';
import { LeagueHomeComponent } from './dashboard-components/league-home/league-home.component';
import { LeagueLeaderboardComponent } from './dashboard-components/league-leaderboard/league-leaderboard.component';
import { LeagueTradesComponent } from './dashboard-components/league-trades/league-trades.component';
import { LeaguePortfolioComponent } from './dashboard-components/league-portfolio/league-portfolio.component';

export const routes: Routes = [
    { path: '', component: HomeComponent },
    { path: 'home', component: HomeComponent },
    { path: 'login', component: LoginComponent },
    { path: 'about', component: AboutComponent},
    { path: 'sign-up', component: SignUpComponent},
    { path: 'dashboard', component: DashboardMainComponent,
        children: [
            { path: '', component: LeagueHomeComponent },
            { path: 'home', component: LeagueHomeComponent },
            { path: 'leaderboard', component: LeagueLeaderboardComponent },
            { path: 'portfolio', component: LeaguePortfolioComponent },
            { path: 'trades', component: LeagueTradesComponent },
        ]
    },
];