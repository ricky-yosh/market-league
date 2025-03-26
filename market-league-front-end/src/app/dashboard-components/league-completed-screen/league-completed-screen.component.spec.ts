import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LeagueCompletedScreenComponent } from './league-completed-screen.component';

describe('LeagueCompletedScreenComponent', () => {
  let component: LeagueCompletedScreenComponent;
  let fixture: ComponentFixture<LeagueCompletedScreenComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LeagueCompletedScreenComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LeagueCompletedScreenComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
