import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PortfolioPointsHistoryChartComponent } from './portfolio-points-history-chart.component';

describe('PortfolioPointsHistoryChartComponent', () => {
  let component: PortfolioPointsHistoryChartComponent;
  let fixture: ComponentFixture<PortfolioPointsHistoryChartComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PortfolioPointsHistoryChartComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PortfolioPointsHistoryChartComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
