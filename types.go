package main

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message

}

type Receipt struct {
	RetailerName string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	TotalAmount  float64 `json:"total,string"`
	Items        []Item  `json:"items"`
}

type Item struct {
	Description string  `json:"shortDescription"`
	Price       float64 `json:"price,string"`
}

type ReceiptStore struct {
	receipts map[string]int
}

func NewReceiptStore() *ReceiptStore {
	return &ReceiptStore{
		receipts: make(map[string]int),
	}
}

func (s *ReceiptStore) AddReceipt(id string, points int) {
	s.receipts[id] = points
}

func (s *ReceiptStore) GetPoints(id string) (int, bool) {
	points, found := s.receipts[id]
	return points, found
}
