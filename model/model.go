package model

// premetives
type Account struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Hisotry struct {
	text      string `json:"name"`
	timestamp string `json:"name"`
}

type Borrowing struct {
	ItemName string  `json:"itemName"`
	Uuid     string  `json:"uuid"`
	Account  Account `json:"account"`
}

type Team struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

//api models
type ErrorResponse struct {
	Message string `json:"message"`
}

type AccountCreateRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AccountCreateResponse struct {
	AccessToken string `json:"accessToken"`
}

type AccessTokenHolder interface{}
