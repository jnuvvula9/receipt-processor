package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type APiError struct {
	Error string `json:"error"`
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			log.Println(err)
			var httpError *HTTPError
			if errors.As(err, &httpError) {
				WriteJSON(w, httpError.StatusCode, APiError{Error: httpError.Message})
			} else {
				WriteJSON(w, http.StatusInternalServerError, APiError{Error: "Internal server error"})
			}
		}
	}

}

type APIServer struct {
	listenAddr string
	store      *ReceiptStore
}

func NewAPIServer(listenAddr string, store *ReceiptStore) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/receipts/process", makeHTTPHandlerFunc(s.handleProcessReceipts)).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", makeHTTPHandlerFunc(s.handleGetPoints)).Methods("GET")

	log.Println("Starting API server on: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func createReceiptId() string {
	return uuid.New().String()
}

func (s *APIServer) handleProcessReceipts(w http.ResponseWriter, r *http.Request) error {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		return &HTTPError{StatusCode: http.StatusBadRequest, Message: "The receipt is invalid"}
	}
	receiptId := createReceiptId()
	points := calculatePoints(receipt)
	s.store.AddReceipt(receiptId, points)
	return WriteJSON(w, http.StatusOK, map[string]string{"id": receiptId})
}

func (s *APIServer) handleGetPoints(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	receiptId, ok := vars["id"]
	if !ok {
		return fmt.Errorf("receipt ID is required")
	}
	points, found := s.store.GetPoints(receiptId)
	if !found {
		return &HTTPError{StatusCode: http.StatusNotFound, Message: "No receipt found for that id"}
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"points": points})
}

func calculatePoints(receipt Receipt) int {
	var points int

	//Rule 1
	points += countAlphaNumericChars(receipt.RetailerName)

	//Rule 2
	if isWholeNumber(receipt.TotalAmount) {
		points += 50
	}

	//Rule 3
	if isMultipleOf(receipt.TotalAmount, 0.25) {
		points += 25
	}

	//Rule 4
	points += (len(receipt.Items) / 2) * 5

	//Rule 5
	for _, item := range receipt.Items {
		points += pointsFromItem(item)
	}

	//Rule 6
	if isOddDay(receipt.PurchaseDate) {
		points += 6
	}

	//Rule 7
	if isBetween(receipt.PurchaseTime, "14:00", "16:00") {
		points += 10
	}

	return points
}

func countAlphaNumericChars(s string) int {
	return len(regexp.MustCompile(`[a-zA-Z0-9]`).FindAllString(s, -1))
}

func isWholeNumber(f float64) bool {
	return f == float64(int(f))
}

func isMultipleOf(value, factor float64) bool {
	return math.Mod(value, factor) == 0
}

func pointsFromItem(item Item) int {
	if len(strings.TrimSpace(item.Description))%3 == 0 {
		itemPoints := item.Price * 0.2
		return int(math.Ceil(itemPoints))
	}
	return 0
}

func isOddDay(date string) bool {
	purchaseDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		log.Println("Error parsing date: ", err)
		return false
	}
	return purchaseDate.Day()%2 == 1
}

func isBetween(purchaseTime, start, end string) bool {
	t, _ := time.Parse("15:04", purchaseTime)
	startTime, _ := time.Parse("15:04", start)
	endTime, _ := time.Parse("15:04", end)
	return t.After(startTime) && t.Before(endTime)
}
