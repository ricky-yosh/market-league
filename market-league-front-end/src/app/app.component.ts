import { Component, OnInit } from '@angular/core';
import { NavigationEnd, Router, RouterOutlet } from '@angular/router';
import { filter } from 'rxjs';
import { WebSocketService } from './dashboard-components/services/websocket.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent implements OnInit {

  constructor(
    private router: Router,
    private websocketService: WebSocketService,
  ) {}

  ngOnInit() {
    // Subscribe to route changes and filter for NavigationEnd events
    this.router.events
      .pipe(
        filter((event): event is NavigationEnd => event instanceof NavigationEnd)
      )
      .subscribe(event => {
        // Now "event" is correctly typed as "NavigationEnd"
        if (event.url.startsWith('/dashboard')) {
          // Disable scrolling for /dashboard
          document.body.style.overflow = 'hidden';
        } else {
          // Enable scrolling for other routes
          document.body.style.overflow = 'auto';
        }
      });
    // Load Websocket
    this.websocketService.connect();
  }

  ngOnDestroy(): void {
    // Optionally close the connection when the app is destroyed
    this.websocketService.closeSocketConnection();
  }

}
