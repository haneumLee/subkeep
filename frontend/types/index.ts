// ===== Enums =====
export type BillingCycle = 'weekly' | 'monthly' | 'yearly';
export type SubscriptionStatus = 'active' | 'paused' | 'cancelled';
export type SplitType = 'equal' | 'custom_amount' | 'custom_ratio';
export type AuthProvider = 'google' | 'apple' | 'naver' | 'kakao';

// ===== Models =====
export interface User {
  id: string;
  provider: AuthProvider;
  providerUserId: string;
  email: string | null;
  nickname: string | null;
  avatarUrl: string | null;
  createdAt: string;
  updatedAt: string;
  lastLoginAt: string | null;
}

export interface Category {
  id: string;
  userId: string | null;
  name: string;
  color: string | null;
  icon: string | null;
  sortOrder: number;
  isSystem: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Subscription {
  id: string;
  userId: string;
  serviceName: string;
  categoryId: string | null;
  amount: number;
  billingCycle: BillingCycle;
  currency: string;
  nextBillingDate: string;
  autoRenew: boolean;
  status: SubscriptionStatus;
  satisfactionScore: number | null;
  note: string | null;
  serviceUrl: string | null;
  startDate: string;
  createdAt: string;
  updatedAt: string;
  category?: Category;
  monthlyAmount: number;
  annualAmount: number;
}

export interface ShareGroup {
  id: string;
  ownerUserId: string;
  name: string;
  description: string | null;
  createdAt: string;
  updatedAt: string;
  members?: ShareMember[];
}

export interface ShareMember {
  id: string;
  shareGroupId: string;
  nickname: string;
  role: string | null;
  createdAt: string;
}

export interface SubscriptionShare {
  id: string;
  subscriptionId: string;
  shareGroupId: string;
  splitType: SplitType;
  myShareAmount: number | null;
  myShareRatio: number | null;
  totalMembersSnapshot: number;
  createdAt: string;
  updatedAt: string;
  shareGroup?: ShareGroup;
}

// ===== API Requests =====
export interface CreateSubscriptionRequest {
  serviceName: string;
  categoryId?: string;
  amount: number;
  billingCycle: BillingCycle;
  currency?: string;
  nextBillingDate: string;
  autoRenew?: boolean;
  status?: SubscriptionStatus;
  satisfactionScore?: number;
  note?: string;
  serviceUrl?: string;
  startDate: string;
}

export interface UpdateSubscriptionRequest extends Partial<CreateSubscriptionRequest> {}

export interface CancelSimulationRequest {
  subscriptionIds: string[];
}

export interface AddSimulationRequest {
  serviceName: string;
  amount: number;
  billingCycle: BillingCycle;
  categoryId?: string;
}

export interface ApplySimulationRequest {
  action: 'cancel';
  subscriptionIds: string[];
}

export interface LinkShareRequest {
  shareGroupId: string;
  splitType: SplitType;
  myShareAmount?: number;
  myShareRatio?: number;
  totalMembersSnapshot: number;
}

export interface CreateShareGroupRequest {
  name: string;
  description?: string;
  members: { nickname: string; role?: string }[];
}

export interface CreateCategoryRequest {
  name: string;
  color?: string;
  icon?: string;
  sortOrder?: number;
}

// ===== API Responses =====
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

export interface PaginatedResponse<T> {
  success: boolean;
  data: T[];
  meta: {
    page: number;
    perPage: number;
    total: number;
    totalPages: number;
  };
}

export interface DashboardSummary {
  monthlyTotal: number;
  annualTotal: number;
  activeCount: number;
  pausedCount: number;
  categoryBreakdown: CategoryBreakdown[];
}

export interface CategoryBreakdown {
  categoryId: string;
  categoryName: string;
  categoryColor: string;
  amount: number;
  percentage: number;
  count: number;
}

export interface SimulationResult {
  currentMonthlyTotal: number;
  simulatedMonthlyTotal: number;
  monthlyDifference: number;
  annualDifference: number;
  categoryBreakdown: CategoryBreakdown[];
}

export interface Recommendation {
  subscription: Subscription;
  reason: string;
  potentialSaving: number;
}

// ===== Calendar =====
export interface CalendarSubscription {
  subscriptionId: string;
  serviceName: string;
  amount: number;
  monthlyAmount: number;
  personalAmount: number;
  billingCycle: string;
  categoryName: string;
  categoryColor: string;
  autoRenew: boolean;
}

export interface CalendarDay {
  date: string;
  totalAmount: number;
  subscriptions: CalendarSubscription[];
}

export interface MonthlyCalendar {
  year: number;
  month: number;
  totalAmount: number;
  totalCount: number;
  remainingAmount: number;
  remainingCount: number;
  days: CalendarDay[];
}

export interface DayDetail {
  date: string;
  totalAmount: number;
  subscriptions: CalendarSubscription[];
}

export interface UpcomingPayment {
  date: string;
  daysUntil: number;
  subscriptionId: string;
  serviceName: string;
  amount: number;
  personalAmount: number;
  categoryName: string;
  categoryColor: string;
}

// ===== Reports =====
export interface MonthlyTrend {
  year: number;
  month: number;
  amount: number;
  count: number;
}

export interface AverageCost {
  monthly: number;
  annual: number;
  weekly: number;
}

export interface ReportSummary {
  totalSubscriptions: number;
  activeCount: number;
  pausedCount: number;
  mostExpensive: string | null;
  mostExpensiveAmount: number;
  averageSatisfaction: number;
}

export interface ReportCategoryBreakdown {
  categoryId: string;
  categoryName: string;
  color: string;
  monthlyAmount: number;
  percentage: number;
  count: number;
}

export interface ReportOverview {
  categoryBreakdown: ReportCategoryBreakdown[];
  monthlyTrend: MonthlyTrend[];
  averageCost: AverageCost;
  summary: ReportSummary;
}

// ===== Auth =====
export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
}

export interface LoginResponse {
  user: User;
  tokens: AuthTokens;
}
