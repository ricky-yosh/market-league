import { Component } from '@angular/core';
import { ApiService } from '../api.service'; // Adjust the path as necessary

@Component({
  selector: 'app-test-button',
  standalone: true,
  templateUrl: './test-button.component.html',
  styleUrls: ['./test-button.component.scss']
})
export class TestButtonComponent {
  responseMessage: string = '';

  constructor(private apiService: ApiService) {}

  // Method to call the Go backend endpoint
  sendRequest() {
    this.apiService.getPing().subscribe(
      (response) => {
        this.responseMessage = response.message || 'Success!';
        console.log('Response from Go server:', response);
      },
      (error) => {
        this.responseMessage = 'Error occurred while fetching data';
        console.error('Error:', error);
      }
    );
  }
}