import { useNavigate } from 'react-router-dom'
import {
  Badge,
  Button,
  Grid,
  Heading,
  HStack,
  ProgressBar,
  Surface,
  Text,
  VStack,
} from '../components/ui'

// 当前阶段使用独立演示数据塑造完整首页，后续接入 API 时只需替换数据来源。
const nutrientSummary = [
  { label: '热量', value: '1,280', unit: 'kcal', progress: 64, variant: 'green' as const },
  { label: '蛋白质', value: '68', unit: 'g', progress: 76, variant: 'blue' as const },
  { label: '脂肪', value: '42', unit: 'g', progress: 58, variant: 'orange' as const },
  { label: '碳水', value: '156', unit: 'g', progress: 62, variant: 'purple' as const },
]

const todayMeals = [
  { name: '早餐', detail: '燕麦酸奶碗', energy: '420 kcal', state: '已记录' },
  { name: '午餐', detail: '鸡胸肉糙米饭', energy: '610 kcal', state: '已记录' },
  { name: '晚餐', detail: '等待记录', energy: '--', state: '待添加' },
]

export function DashboardPage() {
  const navigate = useNavigate()
  // 日期由浏览器本地时区生成，避免在页面中硬编码演示日期。
  const today = new Intl.DateTimeFormat('zh-CN', {
    month: 'long',
    day: 'numeric',
    weekday: 'long',
  }).format(new Date())

  return (
    <VStack className="dashboard" gap={5}>
      <Surface className="dashboard__hero" variant="green">
        <VStack className="dashboard__hero-copy" gap={2}>
          <Text color="secondary" type="supporting">{today}</Text>
          <Heading level={2}>下午好，演示用户</Heading>
          <Text as="p" color="secondary">
            今天已经完成两餐记录，继续保持稳定、轻松的饮食节奏。
          </Text>
          <Button
            className="dashboard__primary-action"
            label="添加饮食记录"
            onClick={() => navigate('/app/meals/new')}
            size="lg"
            variant="primary"
          />
        </VStack>

        <Surface className="dashboard__profile-card" variant="default">
          <VStack gap={3}>
            <HStack align="center" justify="between">
              <VStack gap={0.5}>
                <Text type="label" weight="bold">健康档案</Text>
                <Text color="secondary" type="supporting">再补充活动水平即可完成</Text>
              </VStack>
              <Badge label="75%" variant="green" />
            </HStack>
            <ProgressBar
              isLabelHidden
              label="健康档案完善进度"
              max={100}
              value={75}
              variant="success"
            />
            <Button
              label="继续完善"
              onClick={() => navigate('/app/profile')}
              size="sm"
              variant="secondary"
            />
          </VStack>
        </Surface>
      </Surface>

      <section aria-labelledby="nutrition-title">
        <HStack align="end" className="dashboard__section-heading" justify="between">
          <VStack gap={0.5}>
            <Text color="secondary" type="supporting">今日摄入</Text>
            <Heading id="nutrition-title" level={3}>营养概览</Heading>
          </VStack>
          <Text color="secondary" type="supporting">目标完成度</Text>
        </HStack>

        <Grid columns={{ minWidth: 190, max: 4, repeat: 'fit' }} gap={3}>
          {nutrientSummary.map((item) => (
            <Surface className="dashboard__metric" key={item.label} variant={item.variant}>
              <VStack gap={3}>
                <HStack align="center" justify="between">
                  <Text type="label" weight="bold">{item.label}</Text>
                  <Badge label={`${item.progress}%`} variant={item.variant} />
                </HStack>
                <HStack align="end" gap={1}>
                  <Text className="dashboard__metric-value" weight="bold">{item.value}</Text>
                  <Text color="secondary" type="supporting">{item.unit}</Text>
                </HStack>
                <ProgressBar
                  isLabelHidden
                  label={`${item.label}目标完成度`}
                  max={100}
                  value={item.progress}
                  variant={item.label === '脂肪' ? 'warning' : 'success'}
                />
              </VStack>
            </Surface>
          ))}
        </Grid>
      </section>

      <Grid columns={{ minWidth: 320, max: 2, repeat: 'fit' }} gap={4}>
        <Surface className="dashboard__content-card">
          <VStack gap={3}>
            <HStack align="center" justify="between">
              <Heading level={3}>今日餐次</Heading>
              <Badge label="2 / 3 已记录" variant="green" />
            </HStack>
            <VStack className="dashboard__meal-list" gap={0}>
              {todayMeals.map((meal) => (
                <HStack className="dashboard__meal-row" gap={3} justify="between" key={meal.name}>
                  <HStack align="center" gap={3}>
                    <span aria-hidden="true" className="dashboard__meal-dot" />
                    <VStack gap={0.5}>
                      <Text type="label" weight="bold">{meal.name}</Text>
                      <Text color="secondary" type="supporting">{meal.detail}</Text>
                    </VStack>
                  </HStack>
                  <VStack align="end" gap={0.5}>
                    <Text type="label">{meal.energy}</Text>
                    <Text color="secondary" type="supporting">{meal.state}</Text>
                  </VStack>
                </HStack>
              ))}
            </VStack>
          </VStack>
        </Surface>

        <Surface className="dashboard__content-card">
          <VStack gap={4}>
            <VStack gap={1}>
              <Text color="secondary" type="supporting">最近记录</Text>
              <Heading level={3}>鸡胸肉糙米饭</Heading>
              <Text as="p" color="secondary">今天 12:36 · 午餐 · 约 420 克</Text>
            </VStack>
            <Grid columns={3} gap={2}>
              <VStack className="dashboard__mini-stat" gap={0.5}>
                <Text weight="bold">610</Text><Text color="secondary" type="supporting">千卡</Text>
              </VStack>
              <VStack className="dashboard__mini-stat" gap={0.5}>
                <Text weight="bold">38g</Text><Text color="secondary" type="supporting">蛋白质</Text>
              </VStack>
              <VStack className="dashboard__mini-stat" gap={0.5}>
                <Text weight="bold">72g</Text><Text color="secondary" type="supporting">碳水</Text>
              </VStack>
            </Grid>
            <Button label="查看饮食记录" onClick={() => navigate('/app/meals')} variant="secondary" />
          </VStack>
        </Surface>
      </Grid>

      <Surface className="dashboard__tip" variant="muted">
        <HStack align="center" className="dashboard__tip-layout" gap={4} justify="between">
          <HStack align="center" gap={3}>
            <Badge label="TIP" variant="green" />
            <VStack gap={0.5}>
              <Text type="label" weight="bold">今日膳食提示</Text>
              <Text color="secondary">晚餐可增加深色蔬菜，并适当减少精制主食。</Text>
            </VStack>
          </HStack>
          <Button label="查看建议" onClick={() => navigate('/app/recommendations')} size="sm" variant="ghost" />
        </HStack>
      </Surface>
    </VStack>
  )
}
