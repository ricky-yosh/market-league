import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [
    // Router
    CommonModule, RouterOutlet, RouterLink, RouterLinkActive
  ],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.scss'
})
export class NavbarComponent {
  constructor(private router: Router) {}

  redirectToHome() {
    this.router.navigate(['/home']);
  }

  redirectToAbout() {
    this.router.navigate(['/about']);
  }

  redirectToLogin() {
    this.router.navigate(['/login']);
  }
}
