export type NavigationItem = {
  label: string
  shortLabel: string
  path: string
}

export const primaryNavigation: NavigationItem[] = [
  { label: '首页概览', shortLabel: '首页', path: '/app/dashboard' },
  { label: '饮食记录', shortLabel: '记录', path: '/app/meals' },
  { label: '食品与菜谱', shortLabel: '食品', path: '/app/foods' },
  { label: '营养报告', shortLabel: '报告', path: '/app/reports' },
  { label: '膳食建议', shortLabel: '建议', path: '/app/recommendations' },
  { label: '个人档案', shortLabel: '我的', path: '/app/profile' },
]

export const pageTitles: Record<string, string> = {
  '/app/dashboard': '首页概览',
  '/app/meals': '饮食记录',
  '/app/meals/new': '新建饮食记录',
  '/app/foods': '食品与菜谱库',
  '/app/reports': '营养报告',
  '/app/recommendations': '膳食建议',
  '/app/profile': '个人档案',
}
