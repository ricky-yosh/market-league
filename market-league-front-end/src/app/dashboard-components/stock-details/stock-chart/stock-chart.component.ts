import { Component, Input } from '@angular/core';
import { NgxEchartsModule } from 'ngx-echarts';

@Component({
  selector: 'app-stock-chart',
  standalone: true,
  imports: [NgxEchartsModule],
  templateUrl: './stock-chart.component.html',
  styleUrl: './stock-chart.component.scss'
})
export class StockChartComponent {
  @Input() stockData: any; // Input for stock data

  chartOptions: any; // Options for ECharts

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
        },
        yAxis: {
          type: 'value',
        },
        series: [
          {
            data: prices,
            type: 'line',
            smooth: true,
            areaStyle: {},
          },
        ],
        tooltip: {
          trigger: 'axis',
        },
        title: {
          text: `${this.stockData.company_name} Stock Prices`,
          left: 'center',
        },
      };
    }
  }
}