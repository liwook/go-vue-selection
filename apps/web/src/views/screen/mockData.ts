/** 实时游客统计 */
export const realTimeStats = {
  total: 216908,
  trend: 0.83, // 水球图填充比例 0 ~ 1
}

/** 男女比例 */
export const genderRatio = {
  male: 60,
  female: 40,
}

/** 年龄比例 */
export const ageRatio = [
  { name: '0-18岁', value: 16 },
  { name: '19-30岁', value: 8 },
  { name: '31-40岁', value: 12 },
  { name: '41-60岁', value: 24 },
  { name: '60岁以上', value: 20 },
]

/** 热门景区排行 TOP5 */
export const hotScenicSpots = [
  { name: '峨眉山', value: 8.0 },
  { name: '稻城亚丁', value: 6.5 },
  { name: '九寨沟', value: 5.0 },
  { name: '万里长城', value: 4.0 },
  { name: '北京故宫', value: 3.0 },
]

/** 年度游客量对比（按月） */
export const yearlyComparison = {
  years: ['2021年', '2022年', '2023年'],
  months: Array.from({ length: 12 }, (_, i) => `${i + 1}月`),
  data: [
    [120, 132, 101, 134, 90, 230, 210, 120, 132, 101, 134, 90],
    [220, 182, 191, 234, 290, 330, 310, 220, 182, 191, 234, 290],
    [150, 232, 201, 154, 190, 330, 410, 150, 232, 201, 154, 190],
  ],
}

/** 未来 30 天游客量趋势 */
export const monthlyTrend = {
  dates: Array.from({ length: 30 }, (_, i) => `04/${String(i + 10).padStart(2, '0')}`),
  values: Array.from({ length: 30 }, () => Math.floor(Math.random() * 2000) + 1000),
}

/** 预约渠道数据统计 */
export const channelStats = [
  { name: '微信公众号', value: 40 },
  { name: '携程', value: 10 },
  { name: '飞猪', value: 20 },
  { name: '其他渠道', value: 30 },
]

/** 平台高峰预警信息（地图高亮省份） */
export const peakWarning = [
  { name: '北京市', value: 100 },
  { name: '上海市', value: 80 },
  { name: '广东省', value: 70 },
  { name: '四川省', value: 60 },
  { name: '陕西省', value: 50 },
]
