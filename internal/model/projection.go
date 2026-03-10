package model

type ProjectionPoint struct {
	Month              string  `json:"month"` // "2026-04"
	ProjectedBalance   float64 `json:"projected_balance"`
	Income             float64 `json:"income"`
	Expenses           float64 `json:"expenses"`
}

type BalanceProjection struct {
	CurrentBalance      float64           `json:"current_balance"`
	MonthlyIncomeAvg    float64           `json:"monthly_income_avg"`
	MonthlyExpenseAvg   float64           `json:"monthly_expense_avg"`
	ProjectedPoints     []ProjectionPoint `json:"projected_points"`
}
