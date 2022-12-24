package models

import "github.com/shopspring/decimal"

type Icon struct {
	OID   uint64
	Value []byte
}

type User struct {
	ID       string
	Name     string
	Password []byte
}

type UserConfig struct {
	CompareItemsInDifferentShop bool `json:"compare_items_in_different_shop"`
	CompareItemsInSameShop      bool `json:"compare_items_in_same_shop"`
}

type AssetAccount struct {
	OID            uint64
	Name           string
	IconOID        uint64
	InitialBalance decimal.Decimal
	Balance        decimal.Decimal
}

type Category struct {
	OID     uint64
	Name    string
	IconOID uint64
}

type Shop struct {
	OID     uint64
	Name    string
	Address string
}

type Fee struct {
	OID  uint64
	Name string
	// it is one of type:
	//  - FeeValue_Rate
	//  - FeeValue_Fixed
	Value interface{ feeValue() }
}

type FeeValue_Rate decimal.Decimal

func (FeeValue_Rate) feeValue() {}

type FeeValue_Fixed decimal.Decimal

func (FeeValue_Fixed) feeValue() {}
