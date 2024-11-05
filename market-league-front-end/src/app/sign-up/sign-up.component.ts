import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { NavbarComponent } from '../navbar/navbar.component';
import { FooterComponent } from '../footer/footer.component';
import { SignUpService } from './sign-up-service/sign-up.service';
import { devLog } from '../../environments/development/devlog';
import { NgIf } from '@angular/common';

@Component({
  selector: 'app-sign-up',
  standalone: true,
  imports: [NavbarComponent, FooterComponent, NgIf],
  templateUrl: './sign-up.component.html',
  styleUrl: './sign-up.component.scss'
})
export class SignUpComponent {
  errorMessage: string | null = null;

  constructor(private signUpService: SignUpService, private router: Router) {}

  redirectToLogin() {
    this.router.navigate(['/login']);
  }

  onSubmit(event: Event, username: string, email: string, password: string, confirm_password: string) {
    event.preventDefault();
    this.signUp(username, email, password, confirm_password);
  }

  // Backend call to get jwt token
  signUp(username: string, email: string, password: string, confirm_password: string) {

    // If passwords don't match stop
    if (password !== confirm_password) {
      devLog("Passwords do not match!");
      this.errorMessage = "Passwords do not match!";
      return;
    }

    // If passwords match create an account
    const credentials = {
      username: username,
      email: email,
      password: password,
    };

    this.signUpService.signUp(credentials).subscribe({
      next: (response) => {
        this.handleNext(response);
      },
      error: (error) => {
        this.handleError(error);
      }
    });
  }

  handleNext(response: any) {
    devLog("Sign Up Response", response.message);
    this.redirectToLogin();
  }

  handleError(error: any) {
    devLog("Sign Up Error", error);
    if (error.status === 409) {
      this.errorMessage = "Username already in use!";
    } else {
      this.errorMessage = "Failed to register. Please try again.";
    }
  }

}
