package dtos

import (
	"net/http"
	"strconv"
	"time"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/utils"
	"github.com/shopspring/decimal"
)

type NewWalletTransfer struct {
	FromWalletId   int            `json:"from_wallet_id" validate:"required,gte=1"`
	ToWalletId     int            `json:"to_wallet_id" validate:"required,gte=1"`
	TransferDate   models.MyDate  `json:"transfer_date" validate:"required"`
	TransferAmount NaturalDecimal `json:"transfer_amount" validate:"required"`
	Description    string         `json:"description" validate:"omitempty,min=3,max=255"`
}

func (input NewWalletTransfer) ToModel(s *app.Session) models.WalletTransfer {
	t := models.WalletTransfer{
		FromWalletId:   input.FromWalletId,
		ToWalletId:     input.ToWalletId,
		TransferDate:   input.TransferDate,
		TransferAmount: input.TransferAmount.Decimal,
		Description:    input.Description,
	}
	t.UserId = s.UserId
	return t
}

type WalletTransferFilter struct {
	Search       string
	TransferDate models.MyDate
	FromWalletId int
	ToWalletId   int
}

func (f WalletTransferFilter) Parse(r *http.Request) (*WalletTransferFilter, error) {
	q := r.URL.Query()
	if search := q.Get("search"); search != "" {
		f.Search = search
	}
	if dateq := q.Get("transfer_date"); dateq != "" {
		date, err := time.Parse("2006-01-02", dateq)
		if err != nil {
			return nil, utils.NewBadRequest("transfer_date is not a valid date")
		}
		f.TransferDate = models.MyDate{Time: date}
	}
	if raw := q.Get("from_wallet_id"); raw != "" {
		id, err := strconv.Atoi(raw)
		if err != nil {
			return nil, utils.NewBadRequest("from_wallet_id is not a valid int")
		}
		f.FromWalletId = id
	}
	if raw := q.Get("to_wallet_id"); raw != "" {
		id, err := strconv.Atoi(raw)
		if err != nil {
			return nil, utils.NewBadRequest("to_wallet_id is not a valid int")
		}
		f.ToWalletId = id
	}
	return &f, nil
}

type walletTransferResponse struct {
	Id             int             `json:"id"`
	FromWallet     WalletResource  `json:"from_wallet"`
	ToWallet       WalletResource  `json:"to_wallet"`
	TransferDate   models.MyDate   `json:"transfer_date"`
	TransferAmount decimal.Decimal `json:"transfer_amount"`
	Description    string          `json:"description"`
}

func toWalletTransferResponse(transfer *models.WalletTransfer) walletTransferResponse {
	return walletTransferResponse{
		Id: transfer.Id,
		FromWallet: WalletResource{
			Id:   transfer.FromWalletId,
			Name: transfer.FromWallet.Name,
		},
		ToWallet: WalletResource{
			Id:   transfer.ToWalletId,
			Name: transfer.ToWallet.Name,
		},
		TransferDate:   transfer.TransferDate,
		TransferAmount: transfer.TransferAmount,
		Description:    transfer.Description,
	}
}

func RespondWalletTransfer(transfer *models.WalletTransfer) any {
	response := toWalletTransferResponse(transfer)
	return &response
}

type walletTransferCollection struct {
	WalletTransfers []walletTransferResponse `json:"wallet_transfers"`
}

func RespondWalletTransferCollection(results []models.WalletTransfer) any {
	transfers := addToCollection(results, func(transfer models.WalletTransfer) walletTransferResponse {
		return toWalletTransferResponse(&transfer)
	})
	return &walletTransferCollection{WalletTransfers: transfers}
}
