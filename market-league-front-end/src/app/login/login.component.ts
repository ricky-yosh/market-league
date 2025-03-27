import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { LoginService } from './login-service/login.service';
import { LoginResponse } from './login-service/login-response.interface';
import { NavbarComponent } from '../navbar/navbar.component';
import { FooterComponent } from '../footer/footer.component';
import { NgIf } from '@angular/common';
import { devLog } from '../../environments/development/devlog';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [NavbarComponent, FooterComponent, NgIf],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss'
})
export class LoginComponent {
  errorMessage: string | null = null;

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
    devLog('Login successful', response);
    // Session Token Handling is already done in the service
    
    // Navigate to dashboard after user data is loaded
    this.router.navigate(['/dashboard/create-league']);
  }

  handleError(error: any) {
    devLog("Sign Up Error", error);
    this.errorMessage = "Incorrect username or password!";
  }

}
