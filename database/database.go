package database

type DB interface {
	User() UserService
	Account() AccountService
	Category() CategoryService
	Shop() ShopService
	Fee() FeeService
	DailyItem() DailyItemService
	RepeatingItem() RepeatingItemService
}
