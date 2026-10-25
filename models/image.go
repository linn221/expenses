package models

type ReferenceType string

const (
	ReferenceTypeUser     ReferenceType = "users"
	ReferenceTypeItem     ReferenceType = "items"
	ReferenceTypeExpense  ReferenceType = "expenses"
	ReferenceTypeIncome   ReferenceType = "incomes"
	ReferenceTypeCategory ReferenceType = "categories"
	ReferenceTypeWallet   ReferenceType = "wallets"
)

type Image struct {
	Id            int           `json:"id" gorm:"primaryKey"`
	Name          string        `json:"name"`
	Path          string        `json:"path"`
	ReferenceType ReferenceType `json:"reference_type" gorm:"index"`
	ReferenceId   int           `json:"reference_id" gorm:"index"`
	HasUserId
}
