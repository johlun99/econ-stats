export interface Transaction {
  id: number
  bookingDate: string
  transactionDate: string
  description: string
  amount: number
  balance: number
  categoryId: number | null
  categoryName: string | null
  categoryColor: string | null
  categoryIcon: string | null
  merchantKey: string
  isTransfer: boolean
}

export interface ImportResult {
  totalRows: number
  newTransactions: number
  duplicatesSkipped: number
  autoCategorized: number
  uncategorized: number
}

export interface MerchantGroup {
  merchantKey: string
  count: number
  totalAmount: number
  incomeTotal: number
  expenseTotal: number
  firstDate: string
  lastDate: string
  transactions: Transaction[]
}

export interface Category {
  id: number
  name: string
  color: string
  icon: string
  isIncome: boolean
  isExpense: boolean
  sortOrder: number
}

export interface CategoryRule {
  id: number
  merchantKey: string
  categoryId: number
  categoryName: string
}

export interface MonthlyStats {
  month: string
  totalExpenses: number
  totalIncome: number
  netSavings: number
  savingsRate: number
  avgDailySpend: number
  monthOverMonth: number
  categoryBreakdown: CategorySpend[]
  topMerchants: MerchantSpend[]
  largestExpenses: Transaction[]
  dailySpending: DailySpend[]
}

export interface CategorySpend {
  categoryId: number
  categoryName: string
  categoryColor: string
  categoryIcon: string
  total: number
  count: number
}

export interface MerchantSpend {
  merchantKey: string
  total: number
  count: number
}

export interface DailySpend {
  date: string
  total: number
}

export interface SpendingTrend {
  month: string
  expenses: number
  income: number
}

export interface AvailableMonth {
  month: string
  label: string
}

export interface YearlyStats {
  year: string
  totalExpenses: number
  totalIncome: number
  netSavings: number
  savingsRate: number
  avgMonthlySpend: number
  yearOverYear: number
  categoryBreakdown: CategorySpend[]
  topMerchants: MerchantSpend[]
  largestExpenses: Transaction[]
  monthlySpending: MonthlySpend[]
}

export interface MonthlySpend {
  month: string
  total: number
  income: number
}

export interface AvailableYear {
  year: string
  label: string
}

export interface Debtor {
  id: number
  name: string
  icon: string
  color: string
}

export interface DebtorDetail extends Debtor {
  merchantKeys: string[]
  balance: number
}

export type Page = 'dashboard' | 'yearly' | 'upload' | 'categorize' | 'categories' | 'transactions' | 'debtors'
