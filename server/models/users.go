package models

type LoginRequestBody struct {
	ID       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type SignUpRequestBody struct {
	ID       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserConfigRequestBody struct {
	CompareItemsInDifferentShop bool `json:"compare_items_in_different_shop"`
	CompareItemsInSameShop      bool `json:"compare_items_in_same_shop"`
}
