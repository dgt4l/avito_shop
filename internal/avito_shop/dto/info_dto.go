package dto

type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []Inventory `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type Inventory struct {
	Type     string `json:"type" db:"name"`
	Quantity int    `json:"quantity" db:"quantity"`
}

type CoinHistory struct {
	Received []Received `json:"received"`
	Sent     []Sent     `json:"sent"`
}

type Received struct {
	FromUser string `json:"fromUser" db:"from_user"`
	Amount   int    `json:"amount" db:"amount"`
}

type Sent struct {
	ToUser string `json:"toUser" db:"to_user"`
	Amount int    `json:"amount" db:"amount"`
}
