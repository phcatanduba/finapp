package model

type CompoundInterestRequest struct {
	Principal  float64 `json:"principal"`
	AnnualRate float64 `json:"annual_rate"` // percentage, e.g. 12.5 for 12.5%
	Years      float64 `json:"years"`
	Frequency  int     `json:"frequency"` // compounding per year: 1=annual, 12=monthly, etc.
}

type CompoundInterestResponse struct {
	FutureValue    float64 `json:"future_value"`
	TotalInterest  float64 `json:"total_interest"`
	Principal      float64 `json:"principal"`
	AnnualRate     float64 `json:"annual_rate"`
	Years          float64 `json:"years"`
}

type LoanTable string

const (
	LoanTablePrice LoanTable = "PRICE"
	LoanTableSAC   LoanTable = "SAC"
)

type LoanRequest struct {
	Principal  float64   `json:"principal"`
	AnnualRate float64   `json:"annual_rate"` // percentage
	Terms      int       `json:"terms"`        // number of monthly installments
	Table      LoanTable `json:"table"`        // "PRICE" or "SAC"
}

type LoanInstallment struct {
	Number       int     `json:"number"`
	Installment  float64 `json:"installment"`
	Principal    float64 `json:"principal"`
	Interest     float64 `json:"interest"`
	Balance      float64 `json:"balance"`
}

type LoanResponse struct {
	Table              LoanTable         `json:"table"`
	Principal          float64           `json:"principal"`
	AnnualRate         float64           `json:"annual_rate"`
	Terms              int               `json:"terms"`
	FirstInstallment   float64           `json:"first_installment"`
	LastInstallment    float64           `json:"last_installment"`
	TotalPaid          float64           `json:"total_paid"`
	TotalInterest      float64           `json:"total_interest"`
	Installments       []LoanInstallment `json:"installments"`
}
