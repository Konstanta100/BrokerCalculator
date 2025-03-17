package operation

type Commissions []Commission

type Commission struct {
	Currency string  `json:"currency"`
	Payment  float64 `json:"payment"`
	Date     string  `json:"date"`
}

type SumPayment struct {
	Currency string  `json:"currency"`
	Payment  float64 `json:"payment"`
}

type CalculateCommission struct {
	DateFrom    string       `json:"date_from"`
	DateTo      string       `json:"date_to"`
	SumPayment  []SumPayment `json:"sum_payment"`
	Commissions Commissions  `json:"commissions"`
}
