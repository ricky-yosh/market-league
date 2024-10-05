import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { TestButtonComponent } from "./test-button/test-button.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, TestButtonComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'market-league-front-end';
}
