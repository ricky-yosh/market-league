import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LeagueHomeComponent } from './league-home.component';

describe('LeagueHomeComponent', () => {
  let component: LeagueHomeComponent;
  let fixture: ComponentFixture<LeagueHomeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LeagueHomeComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LeagueHomeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
