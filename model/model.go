package model

// premetives
type Account struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Me struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Token        string
	PasswordHash string
}

func (me *Me) ToAccount() Account {
	return Account{
		Id:   me.Id,
		Name: me.Name,
	}
}

type Hisotry struct {
	Text      string `json:"name"`
	Timestamp string `json:"timestamp"`
}

type Borrowing struct {
	ItemName string  `json:"itemName"`
	Uuid     string  `json:"idHash"`
	Account  Account `json:"account"`
	Team     Team    `json:"team"`
}

type Team struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}
