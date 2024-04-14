// Package model for database model
package model

import (
	"regexp"
	"time"

	"model/noSqlClientPredictedData"
	"model/rdbsClientData"
	"model/rdbsClientInfo"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"gorm.io/gorm"
)

// Repository struct to store psql database clients
type Repository struct {
	cld rdbsClientData.ClientData
	cli rdbsClientInfo.ClientData
}

// Influx struct to store influx client
type Influx struct {
	db noSqlClientPredictedData.ClientData
}

// ClientsInit function to connect to psql databases
func ClientsInit(dataDsn string, clientDsn string) Repository {
	return Repository{
		rdbsClientData.NewConnect(dataDsn),
		rdbsClientInfo.NewConnect(clientDsn),
	}
}

// ClientPredictedDataInit function to connect to Influx
func ClientPredictedDataInit(url string, token string) Influx {
	return Influx{noSqlClientPredictedData.NewConnect(url, token)}
}

// SaveVisitor function to save Visitors
func (r Repository) SaveVisitor(ip string, storeId string, url string, header string, productCode string, tag string) *gorm.DB {
	return r.cld.AddVisitor(ip, storeId, url, productCode, header, tag)
}

// SaveVisitorOffline function to save offline Visitors
func (r Repository) SaveVisitorOffline(info string, storeId string) *gorm.DB {
	return r.cld.AddVisitorOffline(info, storeId)
}

// SaveOrder function to save order
func (r Repository) SaveOrder(amount float64, currency string, storeId string, orderItems []rdbsClientData.Item, orderId string, tag string) *gorm.DB {
	return r.cld.AddOrder(amount, currency, storeId, orderItems, orderId, tag)
}

// GetVisitors function to return visitors by condition
func (r Repository) GetVisitors(condition map[string]interface{}) *gorm.DB {
	return r.cld.GetVisitors(condition)
}

// GetVisitorsOffline function to return visitors by condition
func (r Repository) GetVisitorsOffline(condition map[string]interface{}) *gorm.DB {
	return r.cld.GetOfflineVisitors(condition)
}

// GetOrders function to return orders by condition
func (r Repository) GetOrders(condition map[string]interface{}, limit int, offset int) []rdbsClientData.Orders {
	return r.cld.GetOrders(condition, limit, offset)
}

// GetAmountForPrediction function to return day orders amount for prediction
func (r Repository) GetAmountForPrediction(params map[string]interface{}) []rdbsClientData.AmountByDay {
	return r.cld.GetAmountForPrediction(params)
}

// GetVisitorsForPredictionView function to return viditors day count for prediction by special view
func (r Repository) GetVisitorsForPredictionView(from string, to string, store string) []rdbsClientData.VisitorsByDay {
	return r.cld.GetVisitorsForPredictionView(from, to, store)
}

// GetOrdersForPredictionView get orders count per day for prediction by special view
func (r Repository) GetOrdersForPredictionView(from string, to string, store string) []rdbsClientData.OrdersByDay {
	return r.cld.GetOrdersForPrediction(from, to, store)
}

// GetVisitorsForPrediction function to return viditors day count for prediction
func (r Repository) GetVisitorsForPrediction(from string, to string, store string) []rdbsClientData.VisitorsByDay {
	return r.cld.GetVisitorsForPrediction(from, to, store)
}

// GetOrdersForPrediction get orders count per day for prediction
func (r Repository) GetOrdersForPrediction(from string, to string, store string) []rdbsClientData.OrdersByDay {
	return r.cld.GetOrdersForPrediction(from, to, store)
}

// GetVisitorsForPredictionPerProduct function to count day visitors per product
func (r Repository) GetVisitorsForPredictionPerProduct(from string, to string, store string, productCode string) []rdbsClientData.VisitorsByDay {
	return r.cld.GetVisitorsForPredictionPerProduct(from, to, store, productCode)
}

// GetVisitorsForPredictionPerProductView GetVisitorsForPredictionPerProduct function to count day visitors per product
func (r Repository) GetVisitorsForPredictionPerProductView(from string, to string, store string, productCode string) []rdbsClientData.VisitorsByDay {
	return r.cld.GetVisitorsForPredictionPerProductView(from, to, store, productCode)
}

// GetOrdersForPredictionPerProduct function to count orders per product per day
func (r Repository) GetOrdersForPredictionPerProduct(from string, to string, store string, productCode string) []rdbsClientData.OrdersByDay {
	return r.cld.GetOrdersForPredictionPerProduct(from, to, store, productCode)
}

// GetOrdersForPredictionPerProductView function to count orders per product per day
func (r Repository) GetOrdersForPredictionPerProductView(from string, to string, store string, productCode string) []rdbsClientData.OrdersByDay {
	return r.cld.GetOrdersForPredictionPerProductView(from, to, store, productCode)
}

// GetAvgAmountForPrediction average order amount for prediction
func (r Repository) GetAvgAmountForPrediction(params map[string]interface{}) float64 {
	return r.cld.GetAverageOrderAmount(params)
}

// GetSumOrdersForPrediction get sum of orders for prediction by params
func (r Repository) GetSumOrdersForPrediction(params map[string]interface{}) float64 {
	return r.cld.GetSumOrdersForPrediction(params)
}

// CheckStoreCode function to check if code belongs to store request
func (r Repository) CheckStoreCode(code string, url string) string {
	store, err := r.cli.CheckCode(code, url)
	if err != nil {
		return ""
	} else {
		return store
	}
}

// CheckStoreCodeOffline function to check if code belongs to store request
func (r Repository) CheckStoreCodeOffline(code string, url string) string {
	store, err := r.cli.CheckCodeOffline(code, url)
	if err != nil {
		return ""
	} else {
		return store
	}
}

// CreateAccount function to create account
func (r Repository) CreateAccount(email string, password string, newsletter bool) *gorm.DB {
	account := r.cli.CreateAccount(email, password, newsletter)
	return account
}

// EditAccount function to edit account
func (r Repository) EditAccount(id string, name string, email string, street string, city string, zip string,
	countryCode string, companyNumber string, vatNumber string, paidTo string, planRefer string, role string, parent string, password string, newsletter bool) rdbsClientInfo.Accounts {
	return r.cli.EditAccount(id, name, email, street, city, zip, countryCode, companyNumber, vatNumber, role, parent, password, newsletter)
}

// SetRestorePw function to send restore password tokens
func (r Repository) SetRestorePw(id string, token string) rdbsClientInfo.Accounts {
	return r.cli.SetPwToken(id, token)
}

// UpdatePw function to update password in databse
func (r Repository) UpdatePw(token string, password string) rdbsClientInfo.Accounts {
	return r.cli.UpdatePw(token, password)
}

// DeleteAccount function to delete account form database
func (r Repository) DeleteAccount(id string) {
	r.cli.DeleteAccount(id)
	r.cld.DeleteStoreDataByAccountId(id)
}

// GetAccountById function to get account by id
func (r Repository) GetAccountById(accountId string) rdbsClientInfo.Accounts {
	return r.cli.GetAccountById(accountId)
}

// GetChildAccountById function to get child accounts for main account
func (r Repository) GetChildAccountById(accountId string) []rdbsClientInfo.Accounts {
	return r.cli.GetChildAccountById(accountId)
}

// GetAccountByEmail function to get account by email
func (r Repository) GetAccountByEmail(email string) rdbsClientInfo.Accounts {
	return r.cli.GetAccountByEmail(email)
}

// GetAccounts functionto get all accounts
func (r Repository) GetAccounts() []rdbsClientInfo.Accounts {
	return r.cli.GetAccounts()
}

// GetAccountsForPrediction function to get accounts ready for prediction
func (r Repository) GetAccountsForPrediction() []rdbsClientInfo.Accounts {
	return r.cli.GetAccountsForPrediction()
}

// CreateStore function to create store
func (r Repository) CreateStore(countryCode string, url string, code string, accountRefer string, offline bool, shoptetId string, shoptetToken string, feed string, window int8) *gorm.DB {
	return r.cli.CreateStore(countryCode, url, code, accountRefer, offline, shoptetId, shoptetToken, feed, window)
}

// EditStore function to edit store
func (r Repository) EditStore(id string, countryCode string, url string, maximalProductPrice float64, minimalProductPrice float64,
	actualStorePower float64, actualCustomerSatisfaction float64, perceivedValue float64, productSell int, offline bool, feed string, window int8) rdbsClientInfo.Stores {
	return r.cli.EditStore(id, countryCode, url, maximalProductPrice, minimalProductPrice, actualStorePower,
		actualCustomerSatisfaction, perceivedValue, productSell, offline, feed, window)
}

// Update shoptet info
func (r Repository) UpdateShoptetTokenAndId(storeId string, shoptId string, token string) rdbsClientInfo.Stores {
	return r.cli.UpdateShoptetTokenAndId(storeId, shoptId, token)
}

// DeleteStore function to remove store
func (r Repository) DeleteStore(id string) {
	r.cli.DeleteStore(id)
	r.cld.DeleteStoreData(id)
}

// GetStoresByAccount function to get stores for account
func (r Repository) GetStoresByAccount(accountId string) []rdbsClientInfo.Stores {
	return r.cli.GetStoresByAccount(accountId)
}

// GetStoreById function to get store by id
func (r Repository) GetStoreById(storeId string) rdbsClientInfo.Stores {
	return r.cli.GetStoreById(storeId)
}

// GetStoresByAccount function to get stores for account
func (r Repository) GetStores() []rdbsClientInfo.Stores {
	return r.cli.GetStores()
}

// CreateStoreWeights function to create store weights for prediction
func (r Repository) CreateStoreWeights(storeRefer string, name string, beta float64, gama float64, delta float64,
	a float64, b float64, c float64, d float64, e float64, probabilityWeights string, shift int, longShift int) *gorm.DB {
	return r.cli.CreateStoreWeights(storeRefer, name, beta, gama, delta, a, b, c, d, e, probabilityWeights, shift, longShift)
}

// EditStoreWeights function to edit store weights for prediction
func (r Repository) EditStoreWeights(storeRefer string, name string, beta float64, gama float64, delta float64,
	a float64, b float64, c float64, d float64, e float64, probabilityWeights string, shift int, longShift int) rdbsClientInfo.StoreWeights {
	return r.cli.EditStoreWeights(storeRefer, name, beta, gama, delta, a, b, c, d, e, probabilityWeights, shift, longShift)
}

// GetStoreWeights function to return store weights by store id
func (r Repository) GetStoreWeights(storeId string) rdbsClientInfo.StoreWeights {
	return r.cli.GetStoreWeights(storeId)
}

// GetOpenData function to return open data for store id
func (r Repository) GetOpenData(storeRefer string) []rdbsClientInfo.OpenData {
	return r.cli.GetOpenData(storeRefer)
}

// CreateOpenData function to store parsed open data in database
func (r Repository) CreateOpenData(storePower float64, customerSatisfaction float64, maximalProductPrice float64,
	minimalProductPrice float64, perceivedValue float64, storeRefer string) *gorm.DB {
	return r.cli.CreateOpenData(storePower, customerSatisfaction, maximalProductPrice, minimalProductPrice, perceivedValue, storeRefer)
}

// Auth function to authenticate user
func (r Repository) Auth(email string, password string) rdbsClientInfo.Accounts {
	return r.cli.Auth(email, password)
}

// CreatePlan function to create new plan
func (r Repository) CreatePlan(name string, price float64, period int, products int,
	enabled bool, free bool) *gorm.DB {
	return r.cli.CreatePlan(name, price, period, products, enabled, free)
}

// EditPlan function to edit plan
func (r Repository) EditPlan(id string, name string, price float64, period int, products int,
	enabled bool, free bool) rdbsClientInfo.Plan {
	return r.cli.EditPlan(id, name, price, period, products, enabled, free)
}

// GetPlans function to return all plans
func (r Repository) GetPlans() []rdbsClientInfo.Plan {
	return r.cli.GetPlans()
}

// GetPaidPlans function to return only paid plans
func (r Repository) GetPaidPlans() []rdbsClientInfo.Plan {
	return r.cli.GetPaidPlans()
}

// GetPlanById function return plans by id
func (r Repository) GetPlanById(planId string) rdbsClientInfo.Plan {
	return r.cli.GetPlanById(planId)
}

// DeletePlan function to remove plan from database
func (r Repository) DeletePlan(id string) {
	r.cli.DeletePlan(id)
}

// StoreData function to store predicted data in influx
func (i Influx) StoreData(measurement string, dayIndex string, value int,
	setAverageOrderAmount float64, time time.Time, bucket string, org string) (bool, error) {
	return i.db.StoreData(measurement, dayIndex, value, setAverageOrderAmount, time, bucket, org)
}

// Flush function to flush influx data prepared to store in bucket
func (i Influx) Flush(bucket string, org string) (bool, error) {
	return i.db.Flush(bucket, org)
}

// GetInfluxData function to returned predicted data as string
func (i Influx) GetInfluxData(query string, org string) (string, error) {
	return i.db.GetData(query, org)
}

// GetInfluxQuery function to returned predicted data as query result table
func (i Influx) GetInfluxQuery(query string, org string) (*api.QueryTableResult, error) {
	return i.db.GetQuery(query, org)
}

// GetProducts function to return products by condition
func (r Repository) GetProducts(condition map[string]interface{}) []rdbsClientData.TopSellProduct {
	return r.cld.GetTopSellProducts(condition)
}

// GetOrdersCountByDate function return count orders per specified day
func (r Repository) GetOrdersCountByDate(condition map[string]interface{}) float64 {
	return r.cld.GetOrdersCountByDate(condition)
}

// GetOrdersCountByDatePerProduct function to return orders for specific day and product
func (r Repository) GetOrdersCountByDatePerProduct(condition map[string]interface{}) float64 {
	return r.cld.GetOrdersCountByDatePerProduct(condition)
}

// GetOrdersAvgByDate function return average order by specified day
func (r Repository) GetOrdersAvgByDate(condition map[string]interface{}) float64 {
	return r.cld.GetOrdersAvgByDate(condition)
}

// GetVisitorsCountByDate function return visitors for specified day
func (r Repository) GetVisitorsCountByDate(condition map[string]interface{}) float64 {
	return r.cld.GetVisitorsCountByDate(condition)
}

// GetFirstRecord function return first tracked record for store
func (r Repository) GetFirstRecord(condition map[string]interface{}) string {
	return r.cld.GetFirstRecord(condition)
}

// IsPermitted function check if store is belongs to account
func (r Repository) IsPermitted(accountId string, storeId string) bool {
	id := r.cli.IsAvailableToView(accountId, storeId).Id
	return IsValidUUID(id.String())
}

// GetSumVisitors function return number of visitors for store
func (r Repository) GetSumVisitors(storeId string) float64 {
	return r.cld.GetSumVisitors(storeId)
}

// GetSumOrder function return sum of order for specific store
func (r Repository) GetSumOrder(storeId string) float64 {
	return r.cld.GetSumOrder(storeId)
}

// GetNumberOrders function return count number of orders for specified store
func (r Repository) GetNumberOrders(storeId string) float64 {
	return r.cld.GetNumberOrder(storeId)
}

// GetPredictionR2 function return prediction success for store
func (r Repository) GetPredictionR2(storeId string) float64 {
	return r.cld.GetPredictionR2(storeId)
}

// CreateProduct function to create product in database
func (r Repository) CreateProduct(productCode string, name string, quantity int8, storeId string) rdbsClientData.Products {
	return r.cld.CreateProduct(productCode, name, quantity, storeId)
}

// UpdateProduct function to update product in database
func (r Repository) UpdateProduct(productCode string, name string, storeId string, quantity int8) *gorm.DB {
	return r.cld.UpdateProduct(productCode, name, storeId, quantity)
}

// GetProduct function to return product by product code in specified store
func (r Repository) GetProduct(productCode string, storeId string) rdbsClientData.Product {
	return r.cld.GetProduct(productCode, storeId)
}

// GetProductsWarehouse function to return products in warehouse for each store
func (r Repository) GetProductsWarehouse(storeId string, limit int, offset int) []rdbsClientData.Product {
	return r.cld.GetProducts(storeId, limit, offset)
}

// CreateProductToStore function to save prediction results about products needed to order
func (r Repository) CreateProductToStore(productCode string, quantity int8, storeId string, dateToNeed time.Time, dateToOrder time.Time) rdbsClientData.ProductsToStore {
	return r.cld.CreateProductToStore(productCode, quantity, storeId, dateToNeed, dateToOrder)
}

// UpdateProductToStore function to update prediction results about products needed to order
func (r Repository) UpdateProductToStore(productCode string, storeId string, quantity int8, dateToNeed time.Time, dateToOrder time.Time) *gorm.DB {
	return r.cld.UpdateProductToStore(productCode, storeId, quantity, dateToNeed, dateToOrder)
}

// GetProductToStore function to return product by code need to be ordered
func (r Repository) GetProductToStore(productCode string, storeId string) rdbsClientData.ProductToStore {
	return r.cld.GetProductToStore(productCode, storeId)
}

// GetProductsToStore function to return products need to be ordered
func (r Repository) GetProductsToStore(storeId string, limit int, offset int) []rdbsClientData.ProductsToStore {
	return r.cld.GetProductsToStore(storeId, limit, offset)
}

// GetOrdersWithProduct funcition return order entity with order items
func (r Repository) GetOrdersWithProduct(productCode string, storeId string) []rdbsClientData.Orders {
	return r.cld.GetOrderWithProduct(productCode, storeId)
}

// CreateSupplier function to create supplier in databse
func (r Repository) CreateSupplier(name string, street string, city string, zip string, country string,
	email string, phone string, person string, storeRefer string, template string, subject string) *gorm.DB {
	return r.cli.CreateSupplier(name, street, city, zip, country, email, phone, person, storeRefer, template, subject)
}

// EditSupplier function to edit supplier in databse
func (r Repository) UpdateSupplier(id string, name string, street string, city string, zip string, country string,
	email string, phone string, person string, template string, subject string) rdbsClientInfo.Suppliers {
	return r.cli.EditSupplier(id, name, street, city, zip, country, email, phone, person, template, subject)
}

// GetSupplier function to return suppliers by id
func (r Repository) GetSupplier(supplierId string) rdbsClientInfo.Suppliers {
	return r.cli.GetSupplier(supplierId)
}

// GetSuppliers function to return all suppliers for store
func (r Repository) GetSuppliers(storeId string) []rdbsClientInfo.Suppliers {
	return r.cli.GetSuppliers(storeId)
}

// DeleteSupplier function to delete supplier
func (r Repository) DeleteSupplier(id string) {
	r.cli.DeleteSupplier(id)
}

// CreateInvoice function to create invoice in database
func (r Repository) CreateInvoice(dueDate time.Time, amount float64, currency string, storeRefer string) *gorm.DB {
	return r.cli.CreateInvoice(dueDate, amount, currency, storeRefer)
}

// UpdateInvoice function to create invoice in database
func (r Repository) UpdateInvoice(id string, dueDate time.Time, amount float64, currency string) rdbsClientInfo.Invoices {
	return r.cli.EditInvoice(id, dueDate, amount, currency)
}

// GetInvoices function return all invoices for store
func (r Repository) GetInvoices(storeId string) []rdbsClientInfo.Invoices {
	return r.cli.GetInvoices(storeId)
}

// GetInvoicesFilter function return all invoices for store
func (r Repository) GetInvoicesFilter(storeId string, from int, to int) []rdbsClientInfo.Invoices {
	return r.cli.GetInvoicesFilter(storeId, from, to)
}

// DeleteInvoice function to delete invoice
func (r Repository) DeleteInvoice(id string) {
	r.cli.DeleteInvoice(id)
}

// CreateOrder function to create new plan order
func (r Repository) CreateOrder(accountRefer string, storeRefer string, planRefer string, amount float64, paid bool) string {
	return r.cli.CreateOrder(accountRefer, storeRefer, planRefer, amount, paid)
}

// GetAccountOrders function to return all orders for account
func (r Repository) GetAccountOrders(accountId string) []rdbsClientInfo.Orders {
	return r.cli.GetOrders(accountId)
}

// GetOrderById function return order by id
func (r Repository) GetOrderById(id string) rdbsClientInfo.Orders {
	return r.cli.GetOrderById(id)
}

// GetOrderById function return order by id
func (r Repository) GetStoreByUrl(url string) (string, error) {
	return r.cli.GetStoreByUrl(url)
}

// IsValidUUID function to validate uuid v4
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
