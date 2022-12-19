package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Item struct {
	Name          string
	Type          ItemType
	CategoriesOID []uint64
	ShopOID       uint64
	Quantity      ItemQuantity
	FeeOID        uint64
	Price         decimal.Decimal
	Memo          string
}

type ItemType uint8

const (
	ItemType_Expense ItemType = iota
	ItemType_Income
	ItemType_Transfer
)

type ItemQuantity struct {
	Value decimal.Decimal
}

type DailyItem struct {
	OID  uint64
	Date time.Time
	Item Item
}

type RepeatingItem struct {
	OID       uint64
	Item      Item
	Valid     TimeRange
	Frequency interface{ repeatingItemFrequencyValue() }
}

type RepeatingItemFrequency_Duration int

func (RepeatingItemFrequency_Duration) repeatingItemFrequencyValue() {}

type RepeatingItemFrequency_EveryWorkDay bool

func (RepeatingItemFrequency_EveryWorkDay) repeatingItemFrequencyValue() {}
