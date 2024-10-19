import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LeaguePortfolioComponent } from './league-portfolio.component';

describe('LeaguePortfolioComponent', () => {
  let component: LeaguePortfolioComponent;
  let fixture: ComponentFixture<LeaguePortfolioComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LeaguePortfolioComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LeaguePortfolioComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
