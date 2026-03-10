package model

type ReportSummary struct {
	TotalIncome   float64 `json:"total_income"`
	TotalExpenses float64 `json:"total_expenses"`
	NetBalance    float64 `json:"net_balance"`
	PeriodFrom    string  `json:"period_from"`
	PeriodTo      string  `json:"period_to"`
}

type CategorySpending struct {
	CategoryID   *string `json:"category_id,omitempty"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
	Count        int     `json:"count"`
}

type CashFlowPoint struct {
	Month    string  `json:"month"`
	Income   float64 `json:"income"`
	Expenses float64 `json:"expenses"`
	Net      float64 `json:"net"`
}

type CashFlow struct {
	Points     []CashFlowPoint `json:"points"`
	PeriodFrom string          `json:"period_from"`
	PeriodTo   string          `json:"period_to"`
}

type PeriodSummary struct {
	TotalIncome   float64
	TotalExpenses float64
}
