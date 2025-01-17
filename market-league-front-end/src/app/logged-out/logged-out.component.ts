import { Component } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-logged-out',
  standalone: true,
  imports: [],
  templateUrl: './logged-out.component.html',
  styleUrl: './logged-out.component.scss'
})
export class LoggedOutComponent {

  constructor(private router: Router) {}

  redirectToLogin() {
    this.router.navigate(['/login']);
  }
}
