package http

import (
	"encoding/json"
	"log"
	"net/http"
	"wb_test/internal/cache"
)

func StartServer() {
	http.HandleFunc("/order", getOrderHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("order_uid")
	if orderUID == "" {
		http.Error(w, "Требуется order_uid", http.StatusBadRequest)
		return
	}

	order, found := cache.GetOrderFromCache(orderUID)
	if !found {
		http.Error(w, "Заказ не найден", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}
