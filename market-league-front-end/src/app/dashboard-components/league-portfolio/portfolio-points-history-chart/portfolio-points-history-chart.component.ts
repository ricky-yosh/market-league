import { Component, Input } from '@angular/core';
import { NgxEchartsModule } from 'ngx-echarts';
import { PriceHistory } from '../../../models/price-history.model';
import { PortfolioPointsHistoryEntry } from '../../../models/points-history-entry.model';

@Component({
  selector: 'app-portfolio-points-history-chart',
  standalone: true,
  imports: [NgxEchartsModule],
  templateUrl: './portfolio-points-history-chart.component.html',
  styleUrl: './portfolio-points-history-chart.component.scss'
})
export class PortfolioPointsHistoryChartComponent {
  @Input() portfolioPointsHistory: PortfolioPointsHistoryEntry[] | null = null;
  chartOptions: any;
  echartsInstance: any;

  ngOnChanges(): void {
    if (this.portfolioPointsHistory) {
      const sortedHistory = [...this.portfolioPointsHistory].sort((a, b) => new Date(a.recorded_at).getTime() - new Date(b.recorded_at).getTime());
      
      const timestamps = sortedHistory.map((entry: PortfolioPointsHistoryEntry) =>
        new Date(entry.recorded_at).toLocaleDateString([], {
          day: 'numeric', month: 'short', year: 'numeric'
        })
      );

      const points = sortedHistory.map((entry: PortfolioPointsHistoryEntry) => entry.points);

      this.chartOptions = {
        xAxis: {
          type: 'category',
          data: timestamps,
          boundaryGap: false,
          axisLine: { lineStyle: { color: '#ccc' } },
          axisLabel: { rotate: 45 },
        },
        yAxis: {
          type: 'value',
          scale: true,
          axisLine: { lineStyle: { color: '#ccc' } },
          axisLabel: {
            formatter: (value: number) => `${Math.round(value)} pts`,
          },
        },
        grid: {
          left: '5%',
          right: '5%',
          bottom: '15%',
          containLabel: true,
        },
        series: [
          {
            name: 'Portfolio Points',
            data: points,
            type: 'line',
            smooth: false,
            lineStyle: {
              width: 2,
              color: '#2196F3',
            },
            itemStyle: {
              color: '#2196F3',
            },
            areaStyle: {
              color: 'rgba(33, 150, 243, 0.2)',
            },
          },
        ],
        tooltip: {
          trigger: 'axis',
          formatter: (params: any) => {
            const point = params[0];
            return `<div>${point.axisValue}</div><strong>${Math.round(point.data)} pts</strong>`;
          },
          textStyle: {
            fontSize: 12,
          },
        },
        dataZoom: [
          {
            type: 'inside',
            start: 0,
            end: 100,
          },
        ],
      };
    }
  }

  onChartInit(ec: any): void {
    this.echartsInstance = ec;
  }
}

