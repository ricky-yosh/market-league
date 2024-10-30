import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { NavbarComponent } from '../navbar/navbar.component';
import { FooterComponent } from '../footer/footer.component';
import { SignUpService } from './sign-up-service/sign-up.service';

@Component({
  selector: 'app-sign-up',
  standalone: true,
  imports: [NavbarComponent, FooterComponent],
  templateUrl: './sign-up.component.html',
  styleUrl: './sign-up.component.scss'
})
export class SignUpComponent {
  constructor(private signUpService: SignUpService, private router: Router) {}

  redirectToLogin() {
    this.router.navigate(['/login']);
  }

  onSubmit(event: Event, username: string, password: string, confirm_password: string) {
    event.preventDefault();
    this.signUp(username, password, confirm_password);
  }

  // Backend call to get jwt token
  signUp(username: string, password: string, confirm_password: string) {

    // Check that passwords match

    const credentials = {
      username: username,
      password: password,
    };

    this.signUpService.signUp(credentials).subscribe({
      next: (response) => {
        this.handleUpdateResponse(response);
      },
      error: (error) => {
        this.handleError(error);
      }
    });
  }

  handleUpdateResponse(response: LoginResponse) {
    // Success
    console.log('Login successful', response);
    // Session Token Handling
    localStorage.setItem('token', response.token);
    this.router.navigate(['/dashboard']);
  }

  handleError(error: any) {
    console.error('Login failed', error);
  }

}
