package balanceAPI

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type UserData struct {
	User
	Sum     int    `json:"sum"`
	Comment string `json:"comment"`
}

type TwoUsers struct {
	AddMoneyUser      User   `json:"add_user"`
	DecreaseMoneyUser User   `json:"decrease_user"`
	Sum               int    `json:"sum"`
	Comment           string `json:"comment"`
}

type Transactions struct {
	No      int    `db:"no"`
	Sum     int    `db:"sum"`
	Comment string `db:"comment"`
	Date    string `db:"transaction_date"`
}
