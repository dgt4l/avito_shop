package dto

type BuyItemRequest struct {
	Id   int    `json:"id"`
	Item string `query:"item"`
}
