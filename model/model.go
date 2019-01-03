package model

// premetives
type Account struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Me struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Token        string `json:access_token`
	PasswordHash string
}

type Hisotry struct {
	Text      string `json:"name"`
	Timestamp string `json:"name"`
}

type Borrowing struct {
	ItemName string  `json:"itemName"`
	Uuid     string  `json:"uuid"`
	Account  Account `json:"account"`
	Team     Team    `json:"account"`
}

type Team struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}
