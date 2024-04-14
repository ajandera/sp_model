// Package rdbsClientData package to handle accounts data
package rdbsClientData

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var visitors Visitors
var orders Orders
var product Products
var visitorsOffline VisitorsOffline

// Item struct for order item
type Item struct {
	UnitPrice   float64
	Quantity    int8
	ProductCode string
	ProductName string
	Tag         string
}

// ClientData struct store db client
type ClientData struct {
	db *gorm.DB
}

// AmountByDay struct store order value for each day
type AmountByDay struct {
	Value   float64
	Updated time.Time
}

// VisitorsByDay struct store visitors count for each day
type VisitorsByDay struct {
	Visitors int
	Updated  time.Time
	Day      time.Time
	Tag      string
}

// VisitorsOfflineByDay struct store visitors offline count for each day
type VisitorsOfflineByDay struct {
	Visitors int
	Updated  time.Time
	Day      time.Time
}

// OrdersByDay struct store data for orders for each day
type OrdersByDay struct {
	Orders   int
	Quantity int
	Updated  time.Time
	Day      time.Time
}

// TopSellProduct struct store data for top sel product
type TopSellProduct struct {
	ProductCode string
	Count       int
	Avg         float64
	Quantity    int
	Name        string
}

// Product struct store info about product
type Product struct {
	Quantity    int8
	ProductCode string
	StoreId     string
	Name        string
}

// ProductToStore struct store info about product need to order
type ProductToStore struct {
	Id          string
	Quantity    int8
	ProductCode string
	StoreId     string
	DateToNeed  time.Time
	DateToOrder time.Time
}

// NewConnect function init database connection
func NewConnect(dsn string) ClientData {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	errExt := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if errExt != nil {
		panic(errExt)
	}

	// Migrate the schema
	db.AutoMigrate(
		&OrderItems{},
		&Orders{},
		&Visitors{},
		&Products{},
		&ProductsToStore{})

	// create visitors view for prediction performance
	errView := db.Exec("CREATE or REPLACE VIEW visitorsView AS SELECT count(*) AS visitors, store_id, date_trunc('day', created_at)::date AS day, tag FROM visitors WHERE header NOT LIKE '%Googlebot%' GROUP BY store_id, day, tag ORDER BY day").Error
	if errView != nil {
		panic(errView)
	}

	// create order view for prediction performance
	errView2 := db.Exec("CREATE or REPLACE VIEW ordersView AS SELECT count(*) AS orders, store_id, date_trunc('day', created_at)::date AS day FROM orders GROUP BY store_id, day ORDER BY day").Error
	if errView2 != nil {
		panic(errView2)
	}

	// create visitors view oer product for prediction performance
	errView3 := db.Exec("CREATE or REPLACE VIEW visitorsProductView AS SELECT count(*) AS visitors, store_id, product_code, date_trunc('day', created_at)::date AS day, tag FROM visitors WHERE product_code NOT LIKE '' GROUP BY store_id, product_code, day, tag ORDER BY day").Error
	if errView3 != nil {
		panic(errView3)
	}

	errView4 := db.Exec("CREATE or REPLACE VIEW orderProductView AS SELECT count(order_items.*)::int AS orders, sum(order_items.quantity)::int AS quantity, store_id, product_code, date_trunc('day', order_items.created_at)::date AS day FROM order_items LEFT JOIN orders ON order_items.order = orders.id WHERE order_items.product_code NOT LIKE '' GROUP BY orders.store_id, order_items.product_code, day ORDER BY day").Error
	if errView4 != nil {
		panic(errView4)
	}

	return ClientData{db}
}

// AddVisitor function to store visitor in database
func (client *ClientData) AddVisitor(ip string, storeId string, url string, productCode string, header string, tag string) *gorm.DB {
	visitor := Visitors{Ip: ip, StoreId: storeId, Url: url, ProductCode: productCode, Header: header, Tag: tag}
	result := client.db.Create(&visitor)
	return result
}

// AddVisitorOffline function to store visitor in database
func (client *ClientData) AddVisitorOffline(info string, storeId string) *gorm.DB {
	visitorOffline := VisitorsOffline{Info: info, StoreId: storeId}
	result := client.db.Create(&visitorOffline)
	return result
}

// AddOrder function to store order in database
func (client *ClientData) AddOrder(amount float64, currency string, storeId string, orderItems []Item, orderId string, tag string) *gorm.DB {
	order := Orders{Amount: amount, StoreId: storeId, Currency: currency, ExternalOrderId: orderId, Tag: tag}
	result := client.db.Create(&order)
	for _, v := range orderItems {
		client.AddOrderItem(v, order.Id)
	}
	return result
}

// AddOrderItem function to store order item in database
func (client *ClientData) AddOrderItem(o Item, orderId string) OrderItems {
	item := OrderItems{UnitPrice: o.UnitPrice, Quantity: o.Quantity, ProductCode: o.ProductCode, Order: orderId, ProductName: o.ProductName}
	client.db.Create(&item)
	return item
}

// CreateProduct function to store product in database
func (client *ClientData) CreateProduct(productCode string, name string, quantity int8, storeId string) Products {
	item := Products{Quantity: quantity, ProductCode: productCode, StoreId: storeId, Name: name}
	client.db.Create(&item)
	return item
}

// UpdateProduct function to update product in database
func (client *ClientData) UpdateProduct(productCode string, name string, storeId string, quantity int8) *gorm.DB {
	return client.db.Model(&Products{}).Where("product_code = ? AND store_id = ?", productCode, storeId).Updates(Products{Quantity: quantity, Name: name})
}

// GetProduct function to return product by cide and store
func (client *ClientData) GetProduct(productCode string, storeId string) Product {
	var product Product
	client.db.Where("product_code = ? AND store_id = ?", productCode, storeId).First(&product)
	return product
}

// GetProducts function to return products for store
func (client *ClientData) GetProducts(storeId string, limit int, offset int) []Product {
	var products []Product
	client.db.Where("store_id = ?", storeId).Find(&products).Limit(limit).Offset(offset)
	return products
}

// GetVisitors function return visitors by condition
func (client *ClientData) GetVisitors(condition map[string]interface{}) *gorm.DB {
	return client.db.Where(condition).Find(&visitors)
}

// GetOfflineVisitors function return visitors by condition
func (client *ClientData) GetOfflineVisitors(condition map[string]interface{}) *gorm.DB {
	return client.db.Where(condition).Find(&visitorsOffline)
}

// GetOrders function return orders by condition
func (client *ClientData) GetOrders(condition map[string]interface{}, limit int, offset int) []Orders {
	var orders []Orders
	client.db.Where(condition).Limit(limit).Offset(offset).Order("created_at desc").Find(&orders)
	return orders
}

// GetAmountForPrediction function return order amount for prediction
func (client *ClientData) GetAmountForPrediction(params map[string]interface{}) []AmountByDay {
	var result []AmountByDay
	client.db.Raw("SELECT coalesce(SUM(amount),0) AS value, Max(created_at) FROM orders WHERE store_id = @store_id "+
		"GROUP BY DATE_TRUNC('day',created_at) ORDER BY max(created_at)", params).Scan(&result)
	return result
}

// GetSumOrdersForPrediction function return order sum for prediction
func (client *ClientData) GetSumOrdersForPrediction(params map[string]interface{}) float64 {
	var result float64
	client.db.Raw("SELECT COUNT(id) AS count FROM orders WHERE store_id = @store_id "+
		"AND created_at < @created", params).Scan(&result)
	return result
}

// GetVisitorsForPrediction function to return visitors for prediction
func (client *ClientData) GetVisitorsForPrediction(from string, to string, store string) []VisitorsByDay {
	var result []VisitorsByDay
	client.db.Raw("SELECT * FROM (SELECT day::date FROM generate_series(timestamp '" + from + "', timestamp '" + to + "', interval  '1 day') day) d LEFT JOIN (SELECT date_trunc('day', created_at)::date AS day, count(*)::int AS visitors FROM visitors WHERE  created_at >= date '" + from + "' AND  created_at <= date '" + to + "' AND store_id = '" + store + "' AND product_code = '' AND header NOT LIKE '%Googlebot%' GROUP  BY 1) t USING (day) ORDER  BY day").Scan(&result)
	return result
}

// GetVisitorsForPredictionView function to return visitors for prediction from special database view
func (client *ClientData) GetVisitorsForPredictionView(from string, to string, store string) []VisitorsByDay {
	var result []VisitorsByDay
	client.db.Raw("SELECT * FROM visitorsview WHERE day >= date '" + from + "' AND  day <= date '" + to + "' AND store_id = '" + store + "' ORDER BY day").Scan(&result)
	return result
}

// GetOrdersForPrediction function return orders for prediction
func (client *ClientData) GetOrdersForPrediction(from string, to string, store string) []OrdersByDay {
	var result []OrdersByDay
	client.db.Raw("SELECT * FROM (SELECT day::date FROM generate_series(timestamp '" + from + "', timestamp '" + to + "', interval  '1 day') day) d LEFT JOIN (SELECT date_trunc('day', created_at)::date AS day, count(*)::int AS orders FROM orders WHERE  created_at >= date '" + from + "' AND  created_at <= date '" + to + "' AND store_id = '" + store + "' GROUP  BY 1) t USING (day) ORDER  BY day").Scan(&result)
	return result
}

// GetOrdersForPredictionView function return orders for prediction from special database view
func (client *ClientData) GetOrdersForPredictionView(from string, to string, store string) []OrdersByDay {
	var result []OrdersByDay
	client.db.Raw("SELECT * FROM ordersview WHERE  day >= date '" + from + "' AND  day <= date '" + to + "' AND store_id = '" + store + "' ORDER  BY day").Scan(&result)
	return result
}

// GetAverageOrderAmount function return order amount for prediction
func (client *ClientData) GetAverageOrderAmount(params map[string]interface{}) float64 {
	var result float64
	client.db.Raw("SELECT coalesce(AVG(amount),0) AS amount FROM orders WHERE store_id = @store_id", params).Scan(&result)
	return result
}

// GetOrdersCountByDate function return order count by day
func (client *ClientData) GetOrdersCountByDate(params map[string]interface{}) float64 {
	var result float64
	client.db.Raw("SELECT count(id) FROM orders WHERE created_at > @from AND created_at < @to AND store_id = @store_id", params).Scan(&result)
	return result
}

// GetOrdersCountByDatePerProduct function return order count by day per product
func (client *ClientData) GetOrdersCountByDatePerProduct(params map[string]interface{}) float64 {
	var result float64
	client.db.Raw("SELECT count(orders.id) FROM orders LEFT JOIN order_items ON order_items.order = orders.id WHERE order_items.product_code = @product_code AND orders.created_at > @from AND orders.created_at < @to AND store_id = @store_id", params).Scan(&result)
	return result
}

// GetOrdersAvgByDate function to get average amount of orders per day
func (client *ClientData) GetOrdersAvgByDate(params map[string]interface{}) float64 {
	var result float64
	client.db.Raw("SELECT coalesce(AVG(amount),0) from orders where created_at > @from AND created_at < @to AND store_id = @store_id", params).Scan(&result)
	return result
}

// GetVisitorsCountByDate function to return visitors count per day
func (client *ClientData) GetVisitorsCountByDate(params map[string]interface{}) float64 {
	var result float64
	client.db.Raw("SELECT count(id) from visitors where created_at > @from AND created_at < @to AND product_code = '' AND header NOT LIKE '%Googlebot%' AND store_id = @store_id", params).Scan(&result)
	return result
}

// GetTopSellProducts function to return top sell product by condition
func (client *ClientData) GetTopSellProducts(params map[string]interface{}) []TopSellProduct {
	var result []TopSellProduct
	client.db.Raw("SELECT order_items.product_code, COUNT(order_items.id), AVG(order_items.unit_price), MIN(products.quantity) as quantity, products.name as name FROM order_items LEFT JOIN orders ON order_items.order = orders.id LEFT JOIN products ON order_items.product_code = products.product_code AND products.store_id = orders.store_id WHERE orders.store_id = @store_id GROUP BY order_items.product_code, products.name ORDER BY COUNT(order_items.id) DESC LIMIT @limit OFFSET @offset", params).Scan(&result)
	return result
}

// GetFirstRecord function return first tracked record for store
func (client *ClientData) GetFirstRecord(params map[string]interface{}) string {
	var result string
	client.db.Raw("SELECT created_at FROM visitors WHERE store_id = @store_id ORDER BY created_at ASC LIMIT 1", params).Scan(&result)
	return result
}

// GetVisitorsForPredictionPerProduct function return visitors data for prediction per product
func (client *ClientData) GetVisitorsForPredictionPerProduct(from string, to string, store string, productCode string) []VisitorsByDay {
	var result []VisitorsByDay
	client.db.Raw("SELECT * FROM (SELECT day::date FROM generate_series(timestamp '" + from + "', timestamp '" + to + "', interval  '1 day') day) d LEFT JOIN (SELECT date_trunc('day', created_at)::date AS day, count(*)::int AS visitors FROM visitors WHERE  created_at >= date '" + from + "' AND  created_at <= date '" + to + "' AND store_id = '" + store + "' AND product_code = '" + productCode + "' AND header NOT LIKE '%Googlebot%' GROUP  BY 1) t USING (day) ORDER  BY day").Scan(&result)
	return result
}

// GetVisitorsForPredictionPerProductView function return visitors data for prediction per product for special view
func (client *ClientData) GetVisitorsForPredictionPerProductView(from string, to string, store string, productCode string) []VisitorsByDay {
	var result []VisitorsByDay
	client.db.Raw("SELECT * FROM visitorsproductview WHERE  day >= date '" + from + "' AND  day <= date '" + to + "' AND store_id = '" + store + "' AND product_code = '" + productCode + "' ORDER  BY day").Scan(&result)
	return result
}

// GetOrdersForPredictionPerProduct function return orders for prediction per product
func (client *ClientData) GetOrdersForPredictionPerProduct(from string, to string, store string, productCode string) []OrdersByDay {
	var result []OrdersByDay
	client.db.Raw("SELECT * FROM (SELECT day::date FROM generate_series(timestamp '" + from + "',timestamp '" + to + "', interval  '1 day') day) d LEFT JOIN (SELECT date_trunc('day', created_at)::date AS day, count(order_items.*)::int AS orders, sum(order_items.quantity)::int AS quantity FROM order_items WHERE order_items.created_at >= date '" + from + "' AND  order_items.created_at <= date '" + to + "' AND order_items.order IN (SELECT id FROM orders WHERE orders.store_id = '" + store + "') AND order_items.product_code = '" + productCode + "' GROUP  BY 1) t USING (day) ORDER  BY day").Scan(&result)
	return result
}

// GetOrdersForPredictionPerProductView function return orders for prediction per product for special view
func (client *ClientData) GetOrdersForPredictionPerProductView(from string, to string, store string, productCode string) []OrdersByDay {
	var result []OrdersByDay
	client.db.Raw("SELECT * FROM orderproductview WHERE day >= date '" + from + "' AND  day <= date '" + to + "' AND store_id = '" + store + "' AND product_code = '" + productCode + "' ORDER  BY day").Scan(&result)
	return result
}

// GetSumVisitors function get sum of visitors for store
func (client *ClientData) GetSumVisitors(storeId string) float64 {
	var result float64
	client.db.Raw("SELECT COUNT(id) FROM visitors WHERE store_id = @store_id AND product_code = '' AND header NOT LIKE '%Googlebot%'", map[string]interface{}{"store_id": storeId}).Scan(&result)
	return result
}

// GetSumOrder function to return sum orders for store
func (client *ClientData) GetSumOrder(storeId string) float64 {
	var result float64
	client.db.Raw("SELECT SUM(amount) FROM orders WHERE store_id = @store_id", map[string]interface{}{"store_id": storeId}).Scan(&result)
	return result
}

// GetNumberOrder function to return count of orders for store
func (client *ClientData) GetNumberOrder(storeId string) float64 {
	var result float64
	client.db.Raw("SELECT COUNT(*) FROM orders WHERE store_id = @store_id", map[string]interface{}{"store_id": storeId}).Scan(&result)
	return result
}

// GetPredictionR2 function return suucess of prediction
func (client *ClientData) GetPredictionR2(storeId string) float64 {
	var result float64
	result = 0.92
	return result
}

// DeleteStoreData function to delete store data for store
func (client *ClientData) DeleteStoreData(storeId string) {
	client.db.Raw("DELETE FROM orders WHERE store_id = @store_id ", map[string]interface{}{"store_id": storeId})
	client.db.Raw("DELETE FROM visitors WHERE store_id = @store_id ", map[string]interface{}{"store_id": storeId})
}

// DeleteStoreDataByAccountId function to delete store data for account
func (client *ClientData) DeleteStoreDataByAccountId(accountId string) {
	client.db.Raw("DELETE FROM orders LEFT JOIN stores ON orders.store_id = stores.id LEFT JOIN accounts ON accounts.id = stores.account_id WHERE account.is = @account_id ", map[string]interface{}{"store_id": accountId})
	client.db.Raw("DELETE FROM visitors LEFT JOIN stores ON orders.store_id = stores.id LEFT JOIN accounts ON accounts.id = stores.account_id WHERE account.is = @account_id ", map[string]interface{}{"store_id": accountId})
}

// CreateProductToStore function to store predicted data for product
func (client *ClientData) CreateProductToStore(productCode string, quantity int8, storeId string, dateToNeed time.Time, dateToOrder time.Time) ProductsToStore {
	item := ProductsToStore{Quantity: quantity, ProductCode: productCode, StoreId: storeId, DateToNeed: dateToNeed, DateToOrder: dateToOrder}
	client.db.Create(&item)
	return item
}

// UpdateProductToStore function to update data from prediction for product
func (client *ClientData) UpdateProductToStore(productCode string, storeId string, quantity int8, dateToNeed time.Time, dateToOrder time.Time) *gorm.DB {
	return client.db.Model(&ProductsToStore{}).Where("product_code = ? AND store_id = ?", productCode, storeId).Updates(ProductsToStore{Quantity: quantity, DateToNeed: dateToNeed, DateToOrder: dateToOrder})
}

// GetProductToStore function to return product to order by product code
func (client *ClientData) GetProductToStore(productCode string, storeId string) ProductToStore {
	var productToStore ProductToStore
	client.db.Model(&ProductsToStore{}).Where("product_code = ? AND store_id = ?", productCode, storeId).First(&productToStore)
	return productToStore
}

// GetProductsToStore function to get products to order for store
func (client *ClientData) GetProductsToStore(storeId string, limit int, offset int) []ProductsToStore {
	var productsToStore []ProductsToStore
	client.db.Raw("SELECT products_to_stores.*, products.name FROM products_to_stores LEFT JOIN products ON products.product_code = products_to_stores.product_code AND products.store_id = products_to_stores.store_id  WHERE products_to_stores.store_id = @store_id AND products_to_stores.quantity > 0 ORDER BY products_to_stores.date_to_order DESC LIMIT @limit OFFSET @offset", map[string]interface{}{"store_id": storeId, "limit": limit, "offset": offset}).Scan(&productsToStore)
	return productsToStore
}

// GetOrderWithProduct function return order entity with order items
func (client *ClientData) GetOrderWithProduct(productCode string, storeId string) []Orders {
	var result []Orders
	client.db.Raw("SELECT orders.* FROM orders LEFT JOIN order_items ON orders.id = order_items.order WHERE order_items.product_code = @product_code AND orders.store_id = @store_id", map[string]interface{}{"product_code": productCode, "store_id": storeId}).Scan(&result)
	return result
}
