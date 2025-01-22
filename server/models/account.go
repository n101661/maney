package models

import (
	"github.com/shopspring/decimal"
)

type Account_ struct {
	Name           string          `json:"name" validate:"required"`
	IconOID        uint64          `json:"icon_oid,string"`
	InitialBalance decimal.Decimal `json:"initial_balance"`
}

type GetAccountResponse struct {
	OID uint64 `json:"oid,string"`
	Account_
	Balance decimal.Decimal `json:"balance"`
}
