package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type Product struct {
	ProductID string `json:"id"`
	Quantity  int    `json:"quantity"`
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

var redisClient *redis.Client

func main() {
	redisClient = redis.NewClient(&redis.Options{Addr: "redis:6379"})
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Unable to connect to Redis", err)
	}

	defer redisClient.Close()

	router := http.NewServeMux()

	router.HandleFunc("GET /cart/{id}/", getCart)
	router.HandleFunc("POST /cart/{id}/", addProductCart)
	router.HandleFunc("DELETE /cart/{id}/", removeProductFromCart)

	err := http.ListenAndServe(":8082", router)
	if err != nil {
		log.Fatal("Unable to connext server", err)
	}

}

func addProductCart(w http.ResponseWriter, r *http.Request) {

	var newProduct Product

	err1 := ParseJSON(r, &newProduct)

	if err1 != nil {
		response := map[string]any{
			"message": "Error",
			"product": map[string]any{"message": "Invalid request format"},
		}
		WriteJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	id := newProduct.ProductID
	quantity := newProduct.Quantity

	customerID := r.PathValue("id")
	ctx := r.Context()

	_, err2 := redisClient.HSet(ctx, "cart:"+customerID, id, quantity).Result()

	if err2 != nil {
		response := map[string]any{
			"message": "Error",
			"product": map[string]any{"message": "Failed to save product to cart"},
		}
		WriteJSONResponse(w, http.StatusInternalServerError, response)
		return
	}
	err1 = redisClient.Expire(ctx, "cart:"+customerID, 30*time.Minute).Err()

	if err1 != nil {
		log.Println("Warning: Failed to set TTL for cart:", err1)
	}

	response := map[string]any{
		"message": "Done",
		"product": newProduct,
	}
	WriteJSONResponse(w, http.StatusCreated, response)
}

func getCart(w http.ResponseWriter, r *http.Request) {

	customerID := r.PathValue("id")
	ctx := r.Context()
	cartData, err1 := redisClient.HGetAll(ctx, "cart:"+customerID).Result()

	if err1 != nil {
		response := map[string]any{
			"message": "Error",
			"product": map[string]any{"message": "Failed to show cart"},
		}
		WriteJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]any{
		"message": "Done",
		"product": cartData,
	}
	WriteJSONResponse(w, http.StatusOK, response)

}

func removeProductFromCart(w http.ResponseWriter, r *http.Request) {

	var newProduct Product

	err1 := ParseJSON(r, &newProduct)

	if err1 != nil {
		response := map[string]any{
			"message": "Error",
			"product": map[string]any{"message": "Invalid request format"},
		}
		WriteJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	customerID := r.PathValue("id")
	ctx := r.Context()

	product_id := newProduct.ProductID

	_, err2 := redisClient.HDel(ctx, "cart:"+customerID, product_id).Result()
	if err2 != nil {
		response := map[string]any{
			"message": "Error",
			"product": map[string]any{"message": "Failed to delete product from cart"},
		}
		WriteJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]any{
		"message": "Done",
		"product": map[string]any{"message": fmt.Sprintf("Product successfully deleted %s", product_id)},
	}
	WriteJSONResponse(w, http.StatusOK, response)

}
