package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Merhaba! Go ile yazdığın ilk HTTP servisi çalışıyor.")
}

func WriteJSONResponse(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing body request")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func main() {

	router := http.NewServeMux()

	router.HandleFunc("/", helloHandler)

	router.HandleFunc("GET /products/", getProducts)
	router.HandleFunc("GET /products/{id}/", getProductByID)
	router.HandleFunc("POST /products/", createProduct)

	// 3. Sunucuyu Başlatma: 8080 portunu dinlemeye alıyoruz.
	fmt.Println("Sunucu çalışmaya başladı http://localhost:8080 adresinden erişebilirsin.")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Sunucu başlatılamadı", err)
	}
}

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

var products = []Product{
	{
		ID:          "1",
		Name:        "Wireless Headphones",
		Description: "Bluetooth headphones with active noise cancellation",
		Price:       79.90,
		Stock:       42,
		Category:    "electronics",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          "2",
		Name:        "Ceramic Coffee Mug",
		Description: "350 ml, dishwasher safe",
		Price:       12.50,
		Stock:       120,
		Category:    "home",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          "3",
		Name:        "Running Shoes",
		Description: "Lightweight running shoes with cushioned sole",
		Price:       59.00,
		Stock:       8,
		Category:    "sports",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          "4",
		Name:        "Spiral Notebook",
		Description: "A5 size, 120 sheets, ruled",
		Price:       4.90,
		Stock:       0,
		Category:    "stationery",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}
var mu sync.Mutex

func getProducts(w http.ResponseWriter, r *http.Request) {
	response := map[string]any{
		"message":  "Done",
		"products": products,
	}
	WriteJSONResponse(w, http.StatusOK, response)
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	for i, product := range products {
		if product.ID == id {
			response := map[string]any{
				"message": "Done",
				"product": products[i],
			}
			WriteJSONResponse(w, http.StatusOK, response)
			return

		}

	}
	response := map[string]any{
		"message": "Error",
		"product": map[string]any{"message": "The product couldn't found"},
	}
	WriteJSONResponse(w, http.StatusNotFound, response)

}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product

	err := ParseJSON(r, &newProduct)
	if err != nil {
		response := map[string]any{
			"message": "Error",
			"product": map[string]any{"message": "Invalid request format"},
		}
		WriteJSONResponse(w, http.StatusBadRequest, response)
		return
	}
	mu.Lock()

	newProduct.ID = strconv.Itoa(len(products) + 1)
	newProduct.CreatedAt = time.Now()
	newProduct.UpdatedAt = time.Now()
	products = append(products, newProduct)

	mu.Unlock()

	response := map[string]any{
		"message": "Done",
		"product": newProduct,
	}
	WriteJSONResponse(w, http.StatusCreated, response)

}
