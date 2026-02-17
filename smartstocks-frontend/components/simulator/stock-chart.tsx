'use client'

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Area,
  AreaChart,
} from 'recharts'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import type { ChartData } from '@/types'

interface StockChartProps {
  data: ChartData
  showFullData?: boolean
  title?: string
}

export function StockChart({ data, showFullData = false, title }: StockChartProps) {
  const prices = showFullData && data.full_prices ? data.full_prices : data.prices

  const chartData = data.labels.slice(0, prices.length).map((label, index) => ({
    date: label,
    price: prices[index],
  }))

  const minPrice = Math.min(...prices) * 0.95
  const maxPrice = Math.max(...prices) * 1.05

  // Determine trend color
  const isUptrend = prices[prices.length - 1] > prices[0]
  const trendColor = isUptrend ? '#10B981' : '#EF4444'

  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">
            {title || `${data.asset_name} (${data.ticker})`}
          </CardTitle>
          <div className="text-right">
            <p className="text-2xl font-bold" style={{ color: trendColor }}>
              ${prices[prices.length - 1].toFixed(2)}
            </p>
            <p className="text-sm text-muted-foreground">
              {isUptrend ? '+' : ''}
              {((prices[prices.length - 1] - prices[0]) / prices[0] * 100).toFixed(2)}%
            </p>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="h-[300px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={chartData}>
              <defs>
                <linearGradient id="colorPrice" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor={trendColor} stopOpacity={0.3} />
                  <stop offset="95%" stopColor={trendColor} stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
              <XAxis
                dataKey="date"
                tick={{ fontSize: 12 }}
                tickLine={false}
                axisLine={{ stroke: '#e5e7eb' }}
              />
              <YAxis
                domain={[minPrice, maxPrice]}
                tick={{ fontSize: 12 }}
                tickLine={false}
                axisLine={{ stroke: '#e5e7eb' }}
                tickFormatter={(value) => `$${value.toFixed(0)}`}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: 'white',
                  border: '1px solid #e5e7eb',
                  borderRadius: '8px',
                  boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
                }}
                formatter={(value) => [`$${Number(value).toFixed(2)}`, 'Precio']}
                labelFormatter={(label) => `Fecha: ${label}`}
              />
              <Area
                type="monotone"
                dataKey="price"
                stroke={trendColor}
                strokeWidth={2}
                fill="url(#colorPrice)"
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>

        {!showFullData && (
          <p className="text-center text-sm text-muted-foreground mt-4">
            Los proximos 6 meses se revelaran despues de tu decision
          </p>
        )}
      </CardContent>
    </Card>
  )
}

export default StockChart
