package avito

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
