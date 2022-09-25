package balanceAPI

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type UserData struct {
	User
	Sum int `json:"sum"`
}

type TwoUsers struct {
	AddMoneyUser      User `json:"add_user"`
	DecreaseMoneyUser User `json:"decrease_user"`
	Sum               int  `json:"sum"`
}

type Transactions struct {
	Id      int    `db:"id"`
	Sum     int    `db:"sum"`
	Comment string `db:"comment"`
}
