import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { LoginService } from './services/login.service';
import { LoginResponse } from './services/login-response.interface';
import { NavbarComponent } from '../navbar/navbar.component';
import { FooterComponent } from '../footer/footer.component';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [NavbarComponent, FooterComponent],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss'
})
export class LoginComponent {

  constructor(private loginService: LoginService, private router: Router) {}

  // Change page to sign up
  redirectToSignUp() {
    this.router.navigate(['/sign-up']);
  }

  onSubmit(event: Event, username: string, password: string) {
    event.preventDefault();
    this.login(username, password);
  }

  // Backend call to get jwt token
  login(username: string, password: string) {
    const credentials = {
      username: username,
      password: password,
    };

    this.loginService.login(credentials).subscribe({
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
    
  }

  handleError(error: any) {
    console.error('Login failed', error);
  }

}
