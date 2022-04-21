package models

type BaseReport struct {
	TotalBalance float64
	Debit        float64
	Credit       float64
	NDebit       int
	NCredit      int
}
