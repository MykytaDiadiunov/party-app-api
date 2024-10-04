package domain

type User struct {
	Id       uint64
	Name     string
	Email    string
	Password string
	Points   int32
}

type UpdateUserBalanceAmount struct {
	Amount int32
}
