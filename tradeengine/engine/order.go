package engine

import "encoding/json"

type Order struct {
	Quantity uint32 `json:"amount"`
	Price    uint64 `json:"price"`
	Id       string `json:"id"`
	Type     int8   `json:"side"`
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
