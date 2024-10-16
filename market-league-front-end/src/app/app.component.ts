import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { TestButtonComponent } from "./test-button/test-button.component";
import { HomeComponent } from "./home/home.component";
import { FooterComponent } from "./footer/footer.component";
import { NavbarComponent } from "./navbar/navbar.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, TestButtonComponent, HomeComponent, FooterComponent, NavbarComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'market-league-front-end';
}
