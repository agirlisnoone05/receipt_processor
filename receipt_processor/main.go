package main
import "github.com/google/uuid"
import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Receipt struct {
	ID      string                 `json:"id"`
	Points  int                    `json:"points"`
	Receipt map[string]interface{} `json:"receipt"`
}

var (
	receipts = make(map[string]Receipt)
	mu       sync.Mutex
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/receipts/process", processReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")

	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receiptData map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&receiptData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id := generateID()

	// Calculate points
	points, breakdown := calculatePoints(receiptData)

	// Store receipt
	mu.Lock()
	receipts[id] = Receipt{
		ID:      id,
		Points:  points,
		Receipt: receiptData,
	}
	mu.Unlock()

	fmt.Printf("Breakdown for receipt %s:\n%s\n", id, breakdown)

	// Return the ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	receipt, exists := receipts[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Return the points
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": receipt.Points})
}

func generateID() string {
	return uuid.New().String()
}

func calculatePoints(receipt map[string]interface{}) (int, []string) {
	points := 0
	var breakdown []string

	// Rule 1: One point for every alphanumeric character in the retailer name
	retailer := receipt["retailer"].(string)
	alnumCount := 0
	for _, char := range retailer {
		if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			alnumCount++
		}
	}
	points += alnumCount
	breakdown = append(breakdown, fmt.Sprintf("%d points - retailer name has %d alphanumeric characters", alnumCount, alnumCount))

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	total, _ := strconv.ParseFloat(receipt["total"].(string), 64)
	if total == math.Floor(total) {
		points += 50
		breakdown = append(breakdown, "50 points - total is a round dollar amount")
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
		breakdown = append(breakdown, "25 points - total is a multiple of 0.25")
	}

	// Rule 4: 5 points for every two items on the receipt
	items := receipt["items"].([]interface{})
	itemPairs := len(items) / 2
	points += itemPairs * 5
	breakdown = append(breakdown, fmt.Sprintf("%d points - %d items (%d pairs @ 5 points each)", itemPairs*5, len(items), itemPairs))

	// Rule 5: Points based on item descriptions being multiples of 3
	for _, item := range items {
		itemMap := item.(map[string]interface{})
		description := strings.TrimSpace(itemMap["shortDescription"].(string))
		price, _ := strconv.ParseFloat(itemMap["price"].(string), 64)
		if len(description)%3 == 0 {
			itemPoints := int(math.Ceil(price * 0.2))
			points += itemPoints
			breakdown = append(breakdown, fmt.Sprintf("%d points - \"%s\" description length is a multiple of 3", itemPoints, description))
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	purchaseDate := receipt["purchaseDate"].(string)
	dateParts := strings.Split(purchaseDate, "-")
	day, _ := strconv.Atoi(dateParts[2])
	if day%2 != 0 {
		points += 6
		breakdown = append(breakdown, "6 points - purchase day is odd")
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	purchaseTime := receipt["purchaseTime"].(string)
	t, _ := time.Parse("15:04", purchaseTime)
	if t.Hour() == 14 || (t.Hour() == 15 && t.Minute() == 0) {
		points += 10
		breakdown = append(breakdown, "10 points - purchase time is between 2:00pm and 4:00pm")
	}

	return points, breakdown
}
