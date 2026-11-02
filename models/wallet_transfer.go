package models

import (
	"time"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/validate"
	"github.com/shopspring/decimal"
)

type WalletTransfer struct {
	Id int `json:"id" gorm:"primaryKey"`
	HasUserId
	FromWalletId   int             `json:"from_wallet_id" gorm:"index;not null"`
	ToWalletId     int             `json:"to_wallet_id" gorm:"index;not null"`
	TransferDate   MyDate          `json:"transfer_date" gorm:"index"`
	TransferAmount decimal.Decimal `json:"transfer_amount"`
	Description    string          `json:"description"`
	FromWallet     Wallet          `gorm:"foreignKey:FromWalletId;->"`
	ToWallet       Wallet          `gorm:"foreignKey:ToWalletId;->"`
	CreatedAt      time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

func (existing *WalletTransfer) UpdateValues(input *WalletTransfer) map[string]any {
	return map[string]any{
		"TransferDate":   input.TransferDate,
		"TransferAmount": input.TransferAmount,
		"Description":    input.Description,
	}
}

func (t *WalletTransfer) ValidationRules(vc *validate.ValidationContext) {
	vc.OnMutate(func(s *app.Session, af validate.AppendFunc) {
		if !t.TransferAmount.GreaterThan(decimal.Zero) {
			af(validate.Failure("transfer amount must be greater than zero"))
			return
		}
		filter := UserFilter(s)
		af(
			validate.ExistsRule("wallets", t.FromWalletId, "from wallet not found", filter),
			validate.ExistsRule("wallets", t.ToWalletId, "to wallet not found", filter),
		)
	})
}
