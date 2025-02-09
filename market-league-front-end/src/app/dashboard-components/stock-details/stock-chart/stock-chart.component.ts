import { Component, Input } from '@angular/core';
import { NgxEchartsModule } from 'ngx-echarts';
import { PriceHistory } from '../../../models/price-history.model';

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
  echartsInstance: any;

  ngOnChanges(): void {
    if (this.stockData?.price_histories) {
      const timestamps = this.stockData.price_histories.map((history: PriceHistory) =>
        new Date(history.timestamp).toLocaleDateString([], { day: 'numeric', month: 'numeric', year: 'numeric'})
      );

      const prices = this.stockData.price_histories.map((history: PriceHistory) => history.price);

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
          scale: true, // Ensures the axis doesn't include zero automatically
          axisLine: { lineStyle: { color: '#ccc' } },
          axisLabel: {
            formatter: (value: number) => `$${value.toFixed(2)}`,
          },
          // Remove min: 'dataMin' to prevent automatic adjustment
        },
        grid: {
          left: '5%',
          right: '5%',
          bottom: '15%',
          containLabel: true,
        },
        series: [
          {
            name: this.stockData.company_name, // Added name property
            data: prices,
            type: 'line',
            smooth: false,
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
          formatter: (params: any) => {
            const point = params[0];
            return `<div>${point.axisValue}</div><strong>$${point.data.toFixed(2)}</strong>`;
          },
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
        }
      };
    }
  }

  onChartInit(ec: any): void {
    this.echartsInstance = ec;
    this.echartsInstance.on('dataZoom', this.handleDataZoom.bind(this));
  }

  handleDataZoom(params: any): void {
    if (params.batch && params.batch.length > 0) {
      const { start, end } = params.batch[0];
      this.updateYAxisMin(start, end);
    }
  }

  updateYAxisMin(startPercent: number, endPercent: number): void {
    if (this.stockData?.price_histories) {
      const prices = this.stockData.price_histories.map((history: PriceHistory) => history.price);
      const totalPoints = prices.length;

      const startIndex = Math.floor((startPercent / 100) * totalPoints);
      const endIndex = Math.ceil((endPercent / 100) * totalPoints);

      const visiblePrices = prices.slice(startIndex, endIndex);

      if (visiblePrices.length > 0) {
        const minVisiblePrice = Math.min(...visiblePrices);
        const adjustedMin = minVisiblePrice / 2;

        // Ensure adjustedMin is not negative
        const finalMin = adjustedMin < 0 ? 0 : adjustedMin;

        // Update the entire yAxis configuration
        this.echartsInstance.setOption({
          yAxis: {
            type: 'value',
            scale: true,
            axisLine: { lineStyle: { color: '#ccc' } },
            axisLabel: {
              formatter: (value: number) => `$${value.toFixed(2)}`,
            },
            min: finalMin,
          },
        });
      }
    }
  }
}
