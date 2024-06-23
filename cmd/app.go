package main

import (
	"database/sql"              // импорт стандартного пакета для работы с базой данных
	"encoding/json"             // импорт пакета для работы с json
	"fmt"                       // импорт пакета для форматированного вывода
	"log"                       // импорт пакета для логирования
	"net/http"                  // импорт пакета для работы с http протоколом
	"wb_test/internal/cache"    // импорт пакета для работы с кэшем
	"wb_test/internal/database" // импорт локального пакета для работы с базой данных

	"github.com/gorilla/mux" // импорт библиотеки gorilla/mux для маршрутизации http запросов
)

var db *sql.DB // объявляем переменную для работы с базой данных

func main() {
	cache.InitCache() // инициализируем кэш

	var err error
	db, err = database.ConnectDB() // подключаемся к базе данных
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err) // выбрасываем ошибку, если не получилось подключиться к базе данных
	}

	err = loadCacheFromDB(db) // восстанавливаем кэш из базы данных
	if err != nil {
		log.Fatalf("Не удалось загрузить кэш из базы данных: %v", err) // выбрасываем ошибку, если не получилось восстановить кэш
	}

	r := mux.NewRouter()                                         // создаем новый роутер с использованием библиотеки gorilla/mux
	r.HandleFunc("/orders/{id}", getOrderHandler).Methods("GET") // добавляем обработчик GET запроса по пути /orders/{id}
	r.HandleFunc("/orders", createOrderHandler).Methods("POST")  // добавляем обработчик POST запроса по пути  /orders

	fmt.Println("Сервер работает на порту 8000")
	log.Fatal(http.ListenAndServe(":8000", r)) // запускаем сервер на порту  8000
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                               // получаем переменные из URL
	orderUID := vars["id"]                            // получаем ID заказа из переменных
	log.Printf("Получение заказа с ID: %s", orderUID) // логируем получение заказа по ID

	order, found := cache.GetOrderFromCache(orderUID) // получаем заказ из кэша
	if found {
		log.Printf("Заказ найден в кэше: %+v", order) // логируем нахождение заказа в кэше
		json.NewEncoder(w).Encode(order)              // отпраляем json ответа с найденным заказом
		return
	}

	order, err := database.GetOrderFromDB(db, orderUID) // получаем заказ из базы данных
	if err != nil {
		log.Printf("Ошибка получения заказа из БД.: %v", err)      // логируем ошибку получения заказа из БД
		http.Error(w, err.Error(), http.StatusInternalServerError) // возвращаем http ошибки в случае ошибки БД
		return
	}
	if order == nil {
		log.Printf("Заказ не найден в БД") // логируем ошибку получения заказа из БД
		http.NotFound(w, r)                // возвращаем ошибку 404
		return
	}

	cache.SaveOrderToCache(order)                                   // сохраняем заказ в кэш
	log.Printf("Заказ получен из БД и сохранен в кэше: %+v", order) // логируем успешное получение и сохранение заказа
	json.NewEncoder(w).Encode(order)                                // отпраляем json ответа с найденным заказом
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order database.Order                      // объявляем переменную для нового заказа типа database.Order
	err := json.NewDecoder(r.Body).Decode(&order) // декодируем JSON тела запроса в структуру заказа
	if err != nil {
		log.Printf("Ошибка получения заказа из JSON: %v", err) // логируем ошибку декодирования JSON
		http.Error(w, err.Error(), http.StatusBadRequest)      // возвращаем http ошибки 400 в случае некорректного запроса
		return
	}

	err = database.SaveOrder(db, &order) // сохраняем заказ в БД
	if err != nil {
		log.Printf("Ошибка сохранения заказа в БД: %v", err)       // логируем ошибку сохранения заказа в БД
		http.Error(w, err.Error(), http.StatusInternalServerError) // возвращаем http ошибки в случае ошибки сохранения заказа в БД
		return
	}

	cache.SaveOrderToCache(&order) // сохраняем заказ в кэш

	w.WriteHeader(http.StatusCreated) // устанавливаем HTTP код 201 - созданный заказ
	json.NewEncoder(w).Encode(order)  // отпраляем json ответа с созданным заказом
}

func loadCacheFromDB(db *sql.DB) error {
	orders, err := database.GetAllOrdersFromDB(db) // получаем все заказы из базы данных
	if err != nil {
		return fmt.Errorf("Ошибка загрузки заказов из базы данных: %v", err)
	}

	for _, order := range orders {
		cache.SaveOrderToCache(order) // сохраняем каждый заказ в кэш
	}
	return nil
}
