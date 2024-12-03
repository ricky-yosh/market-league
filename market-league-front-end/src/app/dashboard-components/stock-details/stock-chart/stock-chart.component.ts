import { Component, Input } from '@angular/core';
import { NgxEchartsModule } from 'ngx-echarts';

@Component({
  selector: 'app-stock-chart',
  standalone: true,
  imports: [NgxEchartsModule],
  templateUrl: './stock-chart.component.html',
  styleUrls: ['./stock-chart.component.scss'],
})
export class StockChartComponent {
  @Input() stockData: any;

  chartOptions: any;

  ngOnChanges(): void {
    if (this.stockData?.price_histories) {
      const timestamps = this.stockData.price_histories.map((history: any) =>
        new Date(history.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
      );

      const prices = this.stockData.price_histories.map((history: any) => history.price);

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
          axisLine: { lineStyle: { color: '#ccc' } },
          axisLabel: { formatter: '${value}' },
        },
        grid: {
          left: '5%',
          right: '5%',
          bottom: '15%',
          containLabel: true,
        },
        series: [
          {
            data: prices,
            type: 'line',
            smooth: false, // Disable smooth lines to make them jagged
            lineStyle: {
              width: 2,
              color: '#4CAF50',
            },
            itemStyle: {
              color: '#4CAF50',
            },
            areaStyle: {
              color: 'rgba(76, 175, 80, 0.2)',
            },
          },
        ],
        tooltip: {
          trigger: 'axis',
        },
        dataZoom: [
          {
            type: 'inside',
            start: 80,
            end: 100,
          },
          {
            type: 'slider',
            start: 80,
            end: 100,
          },
        ],
        title: {
          text: `${this.stockData.company_name} Stock Price`,
          left: 'center',
          textStyle: {
            fontSize: 18,
            fontWeight: 'bold',
          },
        },
        legend: {
          show: true,
          data: [this.stockData.company_name],
          top: 'bottom',
        },
      };
      
    }
  }
}
