import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RemoveLeagueComponent } from './remove-league.component';

describe('RemoveLeagueComponent', () => {
  let component: RemoveLeagueComponent;
  let fixture: ComponentFixture<RemoveLeagueComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RemoveLeagueComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(RemoveLeagueComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
