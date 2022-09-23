package balanceAPI

type UserData struct {
	Id      int    `json:"id"`
	Sum     int    `json:"sum"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type Transactions struct {
	Id      int    `db:"id"`
	Sum     int    `db:"sum"`
	Comment string `db:"comment"`
}

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type TwoUsers struct {
	AddMoneyUser      User `json:"add_id"`
	DecreaseMoneyUser User `json:"decrease_id"`
	Sum               int  `json:"sum"`
}
