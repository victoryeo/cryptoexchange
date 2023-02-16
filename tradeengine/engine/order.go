package engine

import "encoding/json"

type Order struct {
	Id        string `json:"id"`
	TokenType string `json:"tokenType"`
	TokenName string `json:"tokenName"`
	Type      int8   `json:"side"`
	Price     uint64 `json:"price"`
	Quantity  uint32 `json:"quantity"`
}

func (order *Order) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, order)
}

func (order *Order) FromMap(msg map[string]string) error {
	return nil
}

func (order *Order) ToJSON() []byte {
	str, _ := json.Marshal(order)
	return str
}
