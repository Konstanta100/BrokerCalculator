package dto

type CalculateCommission struct {
	DateFrom        string                        `json:"dateFrom"`
	DateTo          string                        `json:"dateTo"`
	TotalPayments   map[string]float64            `json:"totalPayments"`
	DateCommissions map[string]map[string]float64 `json:"dateCommissions"`
}
