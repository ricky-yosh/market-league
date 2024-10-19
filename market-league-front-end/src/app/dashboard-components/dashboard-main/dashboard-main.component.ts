import { CommonModule, NgFor } from '@angular/common';
import { Component, EventEmitter, Output } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-dashboard-main',
  standalone: true,
  imports: [RouterOutlet, CommonModule, NgFor],
  templateUrl: './dashboard-main.component.html',
  styleUrl: './dashboard-main.component.scss'
})
export class DashboardMainComponent {
  constructor(private router: Router) {}

  leagues = ['Tech Titans', 'Green Energy League', 'Wall Street Wizards'];

  // Routing
  redirectToHome() {
    this.router.navigate(['/home']);
  }

  redirectToDashboard() {
    this.router.navigate(['/dashboard']);
  }
  
  redirectToLeaderboard() {
    this.router.navigate(['dashboard/leaderboard']);
  }

  redirectToPortfolio() {
    this.router.navigate(['dashboard/portfolio']);
  }  

  redirectToTrades() {
    this.router.navigate(['dashboard/trades']);
  }
}
