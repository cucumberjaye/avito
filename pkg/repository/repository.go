package repository

type Balance interface {
	Add()
	Decrease()
	Transfer()
	GetBalance()
}

type Repository struct {
	Balance
}
