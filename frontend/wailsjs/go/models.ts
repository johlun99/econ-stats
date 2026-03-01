export namespace models {
	
	export class AvailableMonth {
	    month: string;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new AvailableMonth(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.label = source["label"];
	    }
	}
	export class AvailableYear {
	    year: string;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new AvailableYear(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.year = source["year"];
	        this.label = source["label"];
	    }
	}
	export class Category {
	    id: number;
	    name: string;
	    color: string;
	    icon: string;
	    isIncome: boolean;
	    isExpense: boolean;
	    sortOrder: number;
	
	    static createFrom(source: any = {}) {
	        return new Category(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.color = source["color"];
	        this.icon = source["icon"];
	        this.isIncome = source["isIncome"];
	        this.isExpense = source["isExpense"];
	        this.sortOrder = source["sortOrder"];
	    }
	}
	export class CategoryRule {
	    id: number;
	    merchantKey: string;
	    categoryId: number;
	    categoryName: string;
	
	    static createFrom(source: any = {}) {
	        return new CategoryRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.merchantKey = source["merchantKey"];
	        this.categoryId = source["categoryId"];
	        this.categoryName = source["categoryName"];
	    }
	}
	export class CategorySpend {
	    categoryId: number;
	    categoryName: string;
	    categoryColor: string;
	    categoryIcon: string;
	    total: number;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new CategorySpend(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.categoryId = source["categoryId"];
	        this.categoryName = source["categoryName"];
	        this.categoryColor = source["categoryColor"];
	        this.categoryIcon = source["categoryIcon"];
	        this.total = source["total"];
	        this.count = source["count"];
	    }
	}
	export class DailySpend {
	    date: string;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new DailySpend(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.total = source["total"];
	    }
	}
	export class Debtor {
	    id: number;
	    name: string;
	    icon: string;
	    color: string;
	
	    static createFrom(source: any = {}) {
	        return new Debtor(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.icon = source["icon"];
	        this.color = source["color"];
	    }
	}
	export class DebtorDetail {
	    id: number;
	    name: string;
	    icon: string;
	    color: string;
	    merchantKeys: string[];
	    balance: number;
	    pinnedToDashboard: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DebtorDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.icon = source["icon"];
	        this.color = source["color"];
	        this.merchantKeys = source["merchantKeys"];
	        this.balance = source["balance"];
	        this.pinnedToDashboard = source["pinnedToDashboard"];
	    }
	}
	export class ImportResult {
	    totalRows: number;
	    newTransactions: number;
	    duplicatesSkipped: number;
	    updated: number;
	    autoCategorized: number;
	    uncategorized: number;
	
	    static createFrom(source: any = {}) {
	        return new ImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalRows = source["totalRows"];
	        this.newTransactions = source["newTransactions"];
	        this.duplicatesSkipped = source["duplicatesSkipped"];
	        this.updated = source["updated"];
	        this.autoCategorized = source["autoCategorized"];
	        this.uncategorized = source["uncategorized"];
	    }
	}
	export class Transaction {
	    id: number;
	    bookingDate: string;
	    transactionDate: string;
	    description: string;
	    amount: number;
	    balance: number;
	    categoryId?: number;
	    categoryName?: string;
	    categoryColor?: string;
	    categoryIcon?: string;
	    merchantKey: string;
	    isTransfer: boolean;
	    isManual: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Transaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.bookingDate = source["bookingDate"];
	        this.transactionDate = source["transactionDate"];
	        this.description = source["description"];
	        this.amount = source["amount"];
	        this.balance = source["balance"];
	        this.categoryId = source["categoryId"];
	        this.categoryName = source["categoryName"];
	        this.categoryColor = source["categoryColor"];
	        this.categoryIcon = source["categoryIcon"];
	        this.merchantKey = source["merchantKey"];
	        this.isTransfer = source["isTransfer"];
	        this.isManual = source["isManual"];
	    }
	}
	export class MerchantGroup {
	    merchantKey: string;
	    count: number;
	    totalAmount: number;
	    incomeTotal: number;
	    expenseTotal: number;
	    firstDate: string;
	    lastDate: string;
	    transactions: Transaction[];
	
	    static createFrom(source: any = {}) {
	        return new MerchantGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.merchantKey = source["merchantKey"];
	        this.count = source["count"];
	        this.totalAmount = source["totalAmount"];
	        this.incomeTotal = source["incomeTotal"];
	        this.expenseTotal = source["expenseTotal"];
	        this.firstDate = source["firstDate"];
	        this.lastDate = source["lastDate"];
	        this.transactions = this.convertValues(source["transactions"], Transaction);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MerchantSpend {
	    merchantKey: string;
	    total: number;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new MerchantSpend(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.merchantKey = source["merchantKey"];
	        this.total = source["total"];
	        this.count = source["count"];
	    }
	}
	export class MonthlySpend {
	    month: string;
	    total: number;
	    income: number;
	
	    static createFrom(source: any = {}) {
	        return new MonthlySpend(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.total = source["total"];
	        this.income = source["income"];
	    }
	}
	export class MonthlyStats {
	    month: string;
	    totalExpenses: number;
	    totalIncome: number;
	    netSavings: number;
	    savingsRate: number;
	    avgDailySpend: number;
	    monthOverMonth: number;
	    categoryBreakdown: CategorySpend[];
	    topMerchants: MerchantSpend[];
	    largestExpenses: Transaction[];
	    dailySpending: DailySpend[];
	
	    static createFrom(source: any = {}) {
	        return new MonthlyStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.totalExpenses = source["totalExpenses"];
	        this.totalIncome = source["totalIncome"];
	        this.netSavings = source["netSavings"];
	        this.savingsRate = source["savingsRate"];
	        this.avgDailySpend = source["avgDailySpend"];
	        this.monthOverMonth = source["monthOverMonth"];
	        this.categoryBreakdown = this.convertValues(source["categoryBreakdown"], CategorySpend);
	        this.topMerchants = this.convertValues(source["topMerchants"], MerchantSpend);
	        this.largestExpenses = this.convertValues(source["largestExpenses"], Transaction);
	        this.dailySpending = this.convertValues(source["dailySpending"], DailySpend);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SpendingTrend {
	    month: string;
	    expenses: number;
	    income: number;
	
	    static createFrom(source: any = {}) {
	        return new SpendingTrend(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.expenses = source["expenses"];
	        this.income = source["income"];
	    }
	}
	
	export class YearlyStats {
	    year: string;
	    totalExpenses: number;
	    totalIncome: number;
	    netSavings: number;
	    savingsRate: number;
	    avgMonthlySpend: number;
	    yearOverYear: number;
	    categoryBreakdown: CategorySpend[];
	    topMerchants: MerchantSpend[];
	    largestExpenses: Transaction[];
	    monthlySpending: MonthlySpend[];
	
	    static createFrom(source: any = {}) {
	        return new YearlyStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.year = source["year"];
	        this.totalExpenses = source["totalExpenses"];
	        this.totalIncome = source["totalIncome"];
	        this.netSavings = source["netSavings"];
	        this.savingsRate = source["savingsRate"];
	        this.avgMonthlySpend = source["avgMonthlySpend"];
	        this.yearOverYear = source["yearOverYear"];
	        this.categoryBreakdown = this.convertValues(source["categoryBreakdown"], CategorySpend);
	        this.topMerchants = this.convertValues(source["topMerchants"], MerchantSpend);
	        this.largestExpenses = this.convertValues(source["largestExpenses"], Transaction);
	        this.monthlySpending = this.convertValues(source["monthlySpending"], MonthlySpend);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

