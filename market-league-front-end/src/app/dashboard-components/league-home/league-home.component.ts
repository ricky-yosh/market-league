import { NgFor } from '@angular/common';
import { Component } from '@angular/core';

@Component({
  selector: 'app-league-home',
  standalone: true,
  imports: [NgFor],
  templateUrl: './league-home.component.html',
  styleUrl: './league-home.component.scss'
})
export class LeagueHomeComponent {
  leagueMembers = [
    'John Doe', 'Jane Smith', 'Alice Johnson', 'Bob Lee', 
    'Charlie Brown', 'David Beckham', 'Elon Musk', 
    'Serena Williams', 'Michael Jordan', 'Usain Bolt', 
    'Marie Curie', 'Isaac Newton', 'Albert Einstein'
  ];

  userPortfolio = [
    'John Doe', 'Jane Smith', 'Alice Johnson', 'Bob Lee', 
    'Charlie Brown', 'David Beckham', 'Elon Musk', 
    'Serena Williams', 'Michael Jordan', 'Usain Bolt', 
    'Marie Curie', 'Isaac Newton', 'Albert Einstein'
  ];
  
}
