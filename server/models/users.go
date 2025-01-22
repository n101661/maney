package models

type UserConfigRequestBody struct {
	CompareItemsInDifferentShop bool `json:"compare_items_in_different_shop"`
	CompareItemsInSameShop      bool `json:"compare_items_in_same_shop"`
}
