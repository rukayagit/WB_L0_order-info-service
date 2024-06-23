package cache

import (
	"database/sql"              // Импортируем пакет для работы с SQL базами данных
	"sync"                      // Импортируем пакет для синхронизации goroutine
	"wb_test/internal/database" // Импортируем локальный пакет для работы с базой данных
)

// Cache представляет структуру кэша для хранения заказов
type Cache struct {
	mu     sync.RWMutex               // RWMutex обеспечивает потокобезопасность для кэша
	orders map[string]*database.Order // map для хранения заказов по их UID
}

var cache *Cache // Переменная для хранения кэша заказов

// InitCache инициализирует кэш заказов
func InitCache() {
	cache = &Cache{
		orders: make(map[string]*database.Order), // Инициализируем map для хранения заказов
	}
}

// GetOrderFromCache возвращает заказ из кэша по его UID
func GetOrderFromCache(orderUID string) (*database.Order, bool) {
	cache.mu.RLock()                       // Блокируем кэш для чтения
	defer cache.mu.RUnlock()               // Разблокируем кэш после выполнения функции
	order, found := cache.orders[orderUID] // Ищем заказ в кэше
	return order, found                    // Возвращаем найденный заказ и флаг его наличия
}

// SaveOrderToCache сохраняет заказ в кэш
func SaveOrderToCache(order *database.Order) {
	cache.mu.Lock()                      // Блокируем кэш для записи
	defer cache.mu.Unlock()              // Разблокируем кэш после выполнения функции
	cache.orders[order.OrderUID] = order // Сохраняем заказ в кэш
}

// LoadCacheFromDB загружает кэш из базы данных
func LoadCacheFromDB(db *sql.DB) error {
	cache.mu.Lock()         // Блокируем кэш для записи
	defer cache.mu.Unlock() // Разблокируем кэш после выполнения функции

	orders, err := database.GetAllOrdersFromDB(db) // Получаем все заказы из базы данных
	if err != nil {
		return err // Возвращаем ошибку, если не удалось получить заказы из БД
	}

	for _, order := range orders {
		cache.orders[order.OrderUID] = order // Сохраняем каждый заказ в кэш
	}

	return nil // Возвращаем nil, если загрузка прошла успешно
}
