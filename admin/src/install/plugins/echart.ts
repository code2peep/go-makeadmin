import * as echarts from 'echarts/core'
import { GaugeChart, LineChart, PieChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { LabelLayout } from 'echarts/features'

// Keep chart registration aligned with current core pages: workbench line chart and cache pie/gauge charts.
echarts.use([
    LegendComponent,
    TooltipComponent,
    GridComponent,
    LineChart,
    GaugeChart,
    PieChart,
    CanvasRenderer,
    LabelLayout
])
