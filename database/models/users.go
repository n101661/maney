package models

import "github.com/shopspring/decimal"

type User struct {
	OID      uint64
	ID       string
	Password []byte
}

type UserAccount struct {
	OID            uint64
	Name           string
	IconOID        uint64
	InitialBalance decimal.Decimal
	Balance        decimal.Decimal
}

type UserCategory struct {
	OID     uint64
	Name    string
	IconOID uint64
}

type UserShop struct {
	OID     uint64
	Name    string
	Address string
}

type UserFee struct {
	OID  uint64
	Name string
	// it is one of type:
	//  - UserFeeValueRate
	//  - UserFeeValueFixed
	Value interface{ userFeeValue() }
}

type UserFeeValueRate decimal.Decimal

func (UserFeeValueRate) userFeeValue() {}

type UserFeeValueFixed decimal.Decimal

func (UserFeeValueFixed) userFeeValue() {}
