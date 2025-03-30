package dto

type CalculateCommission struct {
	DateFrom        string                        `json:"date_from"`
	DateTo          string                        `json:"date_to"`
	TotalPayments   map[string]float64            `json:"total_payments"`
	DateCommissions map[string]map[string]float64 `json:"date_commissions"`
}
