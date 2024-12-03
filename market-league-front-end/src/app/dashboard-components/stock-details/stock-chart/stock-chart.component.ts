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

      const minPrice = Math.min(...prices); // Find the lowest price
      const adjustedMin = minPrice / 2; // Set the Y-axis minimum to half of the lowest price

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
          axisLabel: {
            formatter: (value: number) => `$${value.toFixed(2)}`, // Format to 2 decimal places
          },
          min: adjustedMin, // Apply the adjusted minimum value
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
              color: '#4CAF50', // Green line
            },
            itemStyle: {
              color: '#4CAF50', // Green data points
            },
            areaStyle: {
              color: 'rgba(76, 175, 80, 0.2)', // Greenish transparent area
            },
          },
        ],
        tooltip: {
          trigger: 'axis',
          formatter: (params: any) => {
            const point = params[0];
            return `<div>${point.axisValue}</div>
                    <strong>$${point.data.toFixed(2)}</strong>`;
          }, // Bold styling with green dot and dollar formatting
          textStyle: {
            fontSize: 12,
          },
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
