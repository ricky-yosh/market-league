import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LeagueLeaderboardComponent } from './league-leaderboard.component';

describe('LeagueLeaderboardComponent', () => {
  let component: LeagueLeaderboardComponent;
  let fixture: ComponentFixture<LeagueLeaderboardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LeagueLeaderboardComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LeagueLeaderboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
