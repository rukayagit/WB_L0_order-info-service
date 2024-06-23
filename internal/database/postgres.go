package database

import (
	"database/sql"        // импорт стандартного пакета для работы с базой данных
	"fmt"                 // импорт пакета для форматированного вывода
	_ "github.com/lib/pq" // импорт драйвера PostgreSQL
)

// Структура для хранения информации о заказе
type Order struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shardkey"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

// Структура для хранения информации о доставке
type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

// Структура для хранения информации об оплате
type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

// Структура для хранения информации о товарах
type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

// Функция для подключения к базе данных PostgreSQL
func ConnectDB() (*sql.DB, error) {
	connStr := "user=postgres password=12345 dbname=l0db sslmode=disable" // строка подключения к базе данных
	db, err := sql.Open("postgres", connStr)                              // открываем соединение с базой данных
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к базе данных: %v", err) // возвращаем ошибку в случае неудачного подключения
	}

	err = db.Ping() // проверяем соединение с базой данных
	if err != nil {
		return nil, fmt.Errorf("Ошибка проверки соединения с базой данных: %v", err) // возвращаем ошибку в случае неудачной проверки
	}
	fmt.Println("Подключено к базе данных") // выводим сообщение об успешном подключении

	return db, nil // возвращаем объект базы данных и nil в случае успешного подключения
}

// Функция для сохранения заказа в базу данных
func SaveOrder(db *sql.DB, order *Order) error {
	tx, err := db.Begin() // начинаем транзакцию
	if err != nil {
		return fmt.Errorf("Ошибка при запуске транзакции: %v", err) // возвращаем ошибку в случае неудачного начала транзакции
	}

	_, err = tx.Exec(`INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, delivery_service, shardkey, sm_id, date_created, oof_shard)
	                   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	                   ON CONFLICT (order_uid) DO UPDATE SET track_number = EXCLUDED.track_number`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		tx.Rollback()                                        // откатываем транзакцию в случае ошибки
		return fmt.Errorf("Ошибка при вводе order: %v", err) // возвращаем ошибку в случае неудачного ввода данных о заказе
	}

	_, err = tx.Exec(`INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
	                   VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	                   ON CONFLICT (order_uid) DO UPDATE SET name = EXCLUDED.name`,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		tx.Rollback()                                           // откатываем транзакцию в случае ошибки
		return fmt.Errorf("Ошибка при вводе delivery: %v", err) // возвращаем ошибку в случае неудачного ввода данных о доставке
	}

	_, err = tx.Exec(`INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
	                   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	                   ON CONFLICT (order_uid) DO UPDATE SET transaction = EXCLUDED.transaction`,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		tx.Rollback()                                          // откатываем транзакцию в случае ошибки
		return fmt.Errorf("Ошибка при вводе payment: %v", err) // возвращаем ошибку в случае неудачного ввода данных об оплате
	}

	for _, item := range order.Items { // цикл для ввода данных о каждом товаре в заказе
		_, err = tx.Exec(`INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		                   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		                   ON CONFLICT (id) DO UPDATE SET chrt_id = EXCLUDED.chrt_id`,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()                                       // откатываем транзакцию в случае ошибки
			return fmt.Errorf("Ошибка при вводе item: %v", err) // возвращаем ошибку в случае неудачного ввода данных о товаре
		}
	}

	err = tx.Commit() // подтверждаем транзакцию
	if err != nil {
		return fmt.Errorf("Ошибка при совершении транзакции: %v", err) // возвращаем ошибку в случае неудачного подтверждения транзакции
	}

	return nil // возвращаем nil в случае успешного сохранения заказа
}

// Функция для получения заказа из базы данных по его ID
func GetOrderFromDB(db *sql.DB, orderUID string) (*Order, error) {
	var order Order
	order.Items = make([]Item, 0) // инициализируем пустой срез для товаров в заказе

	// Получаем данные о заказе из таблицы orders
	err := db.QueryRow(`SELECT order_uid, track_number, entry, locale, internal_signature, delivery_service, shardkey, sm_id, date_created, oof_shard
	                    FROM orders WHERE order_uid = $1`, orderUID).
		Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // если заказ не найден, возвращаем nil
		}
		return nil, fmt.Errorf("Ошибка получения order: %v", err) // возвращаем ошибку в случае неудачного запроса
	}

	// Получаем данные о доставке из таблицы delivery
	err = db.QueryRow(`SELECT name, phone, zip, city, address, region, email
	                    FROM delivery WHERE order_uid = $1`, orderUID).
		Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения delivery: %v", err) // возвращаем ошибку в случае неудачного запроса
	}

	// Получаем данные об оплате из таблицы payment
	err = db.QueryRow(`SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
	                    FROM payment WHERE order_uid = $1`, orderUID).
		Scan(&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения payment: %v", err) // возвращаем ошибку в случае неудачного запроса
	}

	// Получаем данные о товарах из таблицы items
	rows, err := db.Query(`SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	                       FROM items WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения items: %v", err) // возвращаем ошибку в случае неудачного запроса
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, fmt.Errorf("Ошибка сканирования item: %v", err) // возвращаем ошибку в случае неудачного сканирования строки
		}
		order.Items = append(order.Items, item) // добавляем товар в срез товаров заказа
	}

	return &order, nil // возвращаем указатель на заказ
}

// Функция для получения всех заказов из базы данных
func GetAllOrdersFromDB(db *sql.DB) ([]*Order, error) {
	rows, err := db.Query(`SELECT order_uid, track_number, entry, locale, internal_signature, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders`)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения orders: %v", err) // возвращаем ошибку в случае неудачного запроса
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		var order Order
		order.Items = make([]Item, 0) // инициализируем пустой срез для товаров в заказе

		// Сканируем строку с данными о заказе
		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
		if err != nil {
			return nil, fmt.Errorf("Ошибка сканирования order: %v", err) // возвращаем ошибку в случае неудачного сканирования строки
		}

		// Получаем данные о доставке для текущего заказа
		err = db.QueryRow(`SELECT name, phone, zip, city, address, region, email
		                    FROM delivery WHERE order_uid = $1`, order.OrderUID).
			Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
		if err != nil {
			return nil, fmt.Errorf("Ошибка получения delivery: %v", err) // возвращаем ошибку в случае неудачного запроса
		}

		// Получаем данные об оплате для текущего заказа
		err = db.QueryRow(`SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		                    FROM payment WHERE order_uid = $1`, order.OrderUID).
			Scan(&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
		if err != nil {
			return nil, fmt.Errorf("Ошибка получения payment: %v", err) // возвращаем ошибку в случае неудачного запроса
		}

		// Получаем данные о товарах для текущего заказа
		itemRows, err := db.Query(`SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		                           FROM items WHERE order_uid = $1`, order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("Ошибка получения items: %v", err) // возвращаем ошибку в случае неудачного запроса
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item Item
			err := itemRows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
			if err != nil {
				return nil, fmt.Errorf("Ошибка сканирования item: %v", err) // возвращаем ошибку в случае неудачного сканирования строки
			}
			order.Items = append(order.Items, item) // добавляем товар в срез товаров заказа
		}

		orders = append(orders, &order) // добавляем заказ в срез всех заказов
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка итерации по строкам orders: %v", err) // возвращаем ошибку в случае ошибки итерации по строкам
	}

	return orders, nil // возвращаем срез всех заказов
}
