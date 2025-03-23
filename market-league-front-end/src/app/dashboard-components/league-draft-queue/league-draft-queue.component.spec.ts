import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LeagueDraftQueueComponent } from './league-draft-queue.component';

describe('LeagueDraftQueueComponent', () => {
  let component: LeagueDraftQueueComponent;
  let fixture: ComponentFixture<LeagueDraftQueueComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LeagueDraftQueueComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LeagueDraftQueueComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
