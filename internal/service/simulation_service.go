package service

import (
	"math"

	"finapp/internal/model"
)

type simulationService struct{}

func NewSimulationService() SimulationService {
	return &simulationService{}
}

func (s *simulationService) CompoundInterest(req model.CompoundInterestRequest) model.CompoundInterestResponse {
	frequency := req.Frequency
	if frequency <= 0 {
		frequency = 12 // default: monthly compounding
	}
	r := (req.AnnualRate / 100) / float64(frequency)
	n := float64(frequency) * req.Years

	fv := req.Principal * math.Pow(1+r, n)
	interest := fv - req.Principal

	return model.CompoundInterestResponse{
		FutureValue:   round2(fv),
		TotalInterest: round2(interest),
		Principal:     req.Principal,
		AnnualRate:    req.AnnualRate,
		Years:         req.Years,
	}
}

func (s *simulationService) Loan(req model.LoanRequest) model.LoanResponse {
	monthlyRate := (req.AnnualRate / 100) / 12
	n := req.Terms

	switch req.Table {
	case model.LoanTableSAC:
		return s.calcSAC(req.Principal, monthlyRate, n, req.AnnualRate)
	default:
		return s.calcPrice(req.Principal, monthlyRate, n, req.AnnualRate)
	}
}

// calcPrice calculates the PRICE (constant installment) amortization table.
func (s *simulationService) calcPrice(principal, monthlyRate float64, n int, annualRate float64) model.LoanResponse {
	var pmt float64
	if monthlyRate == 0 {
		pmt = principal / float64(n)
	} else {
		pmt = principal * monthlyRate / (1 - math.Pow(1+monthlyRate, float64(-n)))
	}

	installments := make([]model.LoanInstallment, n)
	balance := principal
	totalPaid := 0.0

	for i := 0; i < n; i++ {
		interest := balance * monthlyRate
		amort := pmt - interest
		balance -= amort
		if i == n-1 {
			// last installment: adjust for rounding
			amort += balance
			balance = 0
		}
		installments[i] = model.LoanInstallment{
			Number:      i + 1,
			Installment: round2(pmt),
			Principal:   round2(amort),
			Interest:    round2(interest),
			Balance:     round2(math.Max(balance, 0)),
		}
		totalPaid += pmt
	}

	return model.LoanResponse{
		Table:            model.LoanTablePrice,
		Principal:        principal,
		AnnualRate:       annualRate,
		Terms:            n,
		FirstInstallment: round2(pmt),
		LastInstallment:  round2(pmt),
		TotalPaid:        round2(totalPaid),
		TotalInterest:    round2(totalPaid - principal),
		Installments:     installments,
	}
}

// calcSAC calculates the SAC (constant amortization) table.
func (s *simulationService) calcSAC(principal, monthlyRate float64, n int, annualRate float64) model.LoanResponse {
	amort := principal / float64(n)
	installments := make([]model.LoanInstallment, n)
	balance := principal
	totalPaid := 0.0

	for i := 0; i < n; i++ {
		interest := balance * monthlyRate
		installment := amort + interest
		balance -= amort
		installments[i] = model.LoanInstallment{
			Number:      i + 1,
			Installment: round2(installment),
			Principal:   round2(amort),
			Interest:    round2(interest),
			Balance:     round2(math.Max(balance, 0)),
		}
		totalPaid += installment
	}

	return model.LoanResponse{
		Table:            model.LoanTableSAC,
		Principal:        principal,
		AnnualRate:       annualRate,
		Terms:            n,
		FirstInstallment: installments[0].Installment,
		LastInstallment:  installments[n-1].Installment,
		TotalPaid:        round2(totalPaid),
		TotalInterest:    round2(totalPaid - principal),
		Installments:     installments,
	}
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
