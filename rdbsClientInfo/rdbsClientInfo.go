// Package rdbsClientInfo package to handle communication with account info database
package rdbsClientInfo

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ClientData struct to save gorm instance
type ClientData struct {
	db *gorm.DB
}

// NewConnect function to init database connection
func NewConnect(dsn string) ClientData {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	errExt := db.Raw("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if errExt != nil {
		panic(errExt)
	}
	// Migrate the schema
	db.AutoMigrate(
		&OpenData{},
		&StoreWeights{},
		&Stores{},
		&Plan{},
		&Accounts{},
		&Suppliers{},
		&Invoices{},
		&Orders{})
	return ClientData{db}
}

// GetStoreByUrl function to get store by url
func (client *ClientData) GetStoreByUrl(url string) (string, error) {
	var store Stores
	err := client.db.First(&Stores{}, "url = ?", url).Scan(&store).Error
	return store.Id.String(), err
}

// CheckCode function to check if code belongs to store
func (client *ClientData) CheckCode(code string, url string) (string, error) {
	var store Stores
	var err error
	// sp code
	if code[0:3] == "SP-" {
		err = client.db.First(&Stores{}, "code = ? AND url = ?", code, url).Scan(&store).Error
	} else { // shoptet id
		err = client.db.First(&Stores{}, "shoptet_id = ? AND url = ?", code, url).Scan(&store).Error
	}
	return store.Id.String(), err
}

// CheckCodeOffline function to check if code belongs to store
func (client *ClientData) CheckCodeOffline(code string, url string) (string, error) {
	var store Stores
	err := client.db.First(&Stores{}, "code = ? AND url = ? AND offline = 1", code, url).Scan(&store).Error
	return store.Id.String(), err
}

// Auth function to check username and password
func (client *ClientData) Auth(email string, password string) Accounts {
	var a Accounts
	client.db.Model(&Accounts{}).Where("email = ?", email).Find(&a)
	pw := CheckPasswordHash(password, a.Password)
	if pw == true {
		return a
	} else {
		return Accounts{}
	}
}

// CreateAccount function to create account in db
func (client *ClientData) CreateAccount(email string, password string, newsletter bool) *gorm.DB {

	// hash password
	passwordHash, _ := HashPassword(password)

	// create account
	item := Accounts{
		Email:                  email,
		Password:               passwordHash,
		Newsletter:             newsletter,
		NewsletterConfirmation: time.Now()}
	result := client.db.Create(&item)
	return result
}

// EditAccount function to edit account in db
func (client *ClientData) EditAccount(id string, name string, email string, street string, city string, zip string,
	countryCode string, companyNumber string, vatNumber string, role string, parent string, password string, newsletter bool) Accounts {
	var a Accounts
	client.db.Model(&Accounts{}).Where("id = ?", id).First(&a)

	if len(name) > 0 {
		a.Name = name
	}

	if len(email) > 0 {
		a.Email = email
	}

	if len(street) > 0 {
		a.Street = street
	}

	if len(city) > 0 {
		a.City = city
	}

	if len(zip) > 0 {
		a.Zip = zip
	}

	if len(countryCode) > 0 {
		a.CountryCode = countryCode
	}

	if len(companyNumber) > 0 {
		a.CompanyNumber = companyNumber
	}

	if len(vatNumber) > 0 {
		a.VatNumber = vatNumber
	}

	if len(role) > 0 {
		a.Role = role
	}

	if len(parent) > 0 {
		a.Parent = parent
	}

	if len(password) > 6 {
		hash, _ := HashPassword(password)
		a.Password = hash
	}

	// update newsletter
	if a.Newsletter != newsletter {
		a.Newsletter = newsletter
		a.NewsletterConfirmation = time.Now()
	}
	client.db.Save(&a)
	return a
}

// SetPwToken function to set tojken for pw restore
func (client *ClientData) SetPwToken(id string, token string) Accounts {
	var a Accounts
	client.db.Model(&Accounts{}).Where("id = ?", id).First(&a)
	t := time.Now().Add(time.Hour * 24)
	a.ValidTokenTo = t
	a.RestoreToken = token
	client.db.Save(&a)
	return a
}

// UpdatePw function to update password in databse
func (client *ClientData) UpdatePw(token string, password string) Accounts {
	var a Accounts
	client.db.Model(&Accounts{}).Where("restore_token = ?", token).First(&a)

	t1 := time.Now()
	t2 := a.ValidTokenTo
	hourDiff := t2.Sub(t1).Hours()

	if hourDiff < 24 {
		a.RestoreToken = ""
		if len(password) > 6 {
			hash, _ := HashPassword(password)
			a.Password = hash
		}
	}
	client.db.Save(&a)
	return a
}

// DeleteAccount function to remove account
func (client *ClientData) DeleteAccount(id string) {
	var a Accounts
	var s Stores
	client.db.Model(&Accounts{}).Where("id = ?", id).Delete(&a)
	client.db.Model(&Accounts{}).Where("parent = ?", id).Delete(&a)
	client.db.Model(&Stores{}).Where("account_id = ?", id).Delete(&s)
}

// GetAccountById function to return account by id
func (client *ClientData) GetAccountById(accountId string) Accounts {
	var a Accounts
	client.db.Model(&Accounts{}).Where("id = ?", accountId).Find(&a)
	return a
}

// GetChildAccountById function to get all child accounts for  user
func (client *ClientData) GetChildAccountById(accountId string) []Accounts {
	var a []Accounts
	client.db.Model(&Accounts{}).Where("parent = ?", accountId).Find(&a)
	return a
}

// GetAccountByEmail function return account by email
func (client *ClientData) GetAccountByEmail(email string) Accounts {
	var a Accounts
	client.db.Model(&Accounts{}).Where("email = ?", email).Find(&a)
	return a
}

// GetAccounts function to return all accounts
func (client *ClientData) GetAccounts() []Accounts {
	var a []Accounts
	client.db.Model(&Accounts{}).Find(&a)
	return a
}

// GetAccountsForPrediction function to return accounts ready for prediction
func (client *ClientData) GetAccountsForPrediction() []Accounts {
	var a []Accounts
	client.db.Model(&Accounts{}).Where("parent", "").Order("last_prediction asc").Limit(20).Find(&a)
	return a
}

// CreateStore function to create store in db
func (client *ClientData) CreateStore(countryCode string, url string, code string, accountRefer string, offline bool, shoptetId string, shoptetToken string, feed string, window int8) *gorm.DB {
	var a Accounts
	client.db.Model(&Accounts{}).Where("id = ?", accountRefer).First(&a)

	item := Stores{
		CountryCode:                countryCode,
		Url:                        url,
		Code:                       code,
		AccountRefer:               a.Id.String(),
		MaximalProductPrice:        1000,
		MinimalProductPrice:        100,
		ActualStorePower:           0.9,
		ActualCustomerSatisfaction: 0.9,
		PerceivedValue:             0.85,
		ProductSell:                2000,
		Offline:                    offline,
		ShoptetId:                  shoptetId,
		ShoptetAccessToken:         shoptetToken,
		XmlFeed:                    feed,
		Window:                     window}
	result := client.db.Create(&item)

	// insert default open weights
	sw := StoreWeights{
		StoreRefer:         item.Id.String(),
		Name:               item.Url,
		Beta:               0.3,
		Gama:               0.4,
		Delta:              0.3,
		A:                  0.2,
		B:                  0.2,
		C:                  0.2,
		D:                  0.2,
		E:                  0.2,
		ProbabilityWeights: "[0.3333 0.3333 0.1111;0.3333 0.3333 0.1111;0.3333 0.3333 0.1111]"}
	client.db.Create(&sw)

	return result
}

// EditStore function to edit store in db
func (client *ClientData) EditStore(id string, countryCode string, url string, maximalProductPrice float64, minimalProductPrice float64,
	actualStorePower float64, actualCustomerSatisfaction float64, perceivedValue float64, productSell int, offline bool, feed string, window int8) Stores {
	var s Stores
	client.db.Model(&Stores{}).Where("id = ?", id).First(&s)

	if len(countryCode) > 0 {
		s.CountryCode = countryCode
	}

	if len(url) > 0 {
		s.Url = url
	}

	if maximalProductPrice > 0 {
		s.MaximalProductPrice = maximalProductPrice
	}

	if minimalProductPrice > 0 {
		s.MinimalProductPrice = minimalProductPrice
	}

	if actualStorePower > 0 {
		s.ActualStorePower = actualStorePower
	}

	if actualCustomerSatisfaction > 0 {
		s.ActualCustomerSatisfaction = actualCustomerSatisfaction
	}

	if perceivedValue > 0 {
		s.PerceivedValue = perceivedValue
	}

	if productSell > 0 {
		s.ProductSell = productSell
	}

	if offline != s.Offline {
		s.Offline = offline
	}

	if feed != s.XmlFeed {
		s.XmlFeed = feed
	}

	if window != s.Window {
		s.Window = window
	}

	client.db.Save(&s)
	return s
}

// Update store token and eshop id from shoptet
func (client *ClientData) UpdateShoptetTokenAndId(storeId string, shoptId string, token string) Stores {
	var s Stores
	client.db.Model(&Stores{}).Where("id = ?", storeId).First(&s)
	s.ShoptetId = shoptId
	s.ShoptetAccessToken = token
	client.db.Save(&s)
	return s
}

// DeleteStore function to delte store in db by id
func (client *ClientData) DeleteStore(id string) {
	var s Stores
	client.db.Model(&Stores{}).Where("id = ?", id).Delete(&s)
}

// GetStoresByAccount function to return all stores for account
func (client *ClientData) GetStoresByAccount(accountId string) []Stores {
	var s []Stores
	client.db.Model(&Stores{}).Where("account_refer = ?", accountId).Find(&s)
	return s
}

// GetStores function to return all stores
func (client *ClientData) GetStores() []Stores {
	var s []Stores
	client.db.Model(&Stores{}).Scan(&s)
	return s
}

// GetStoreById function return store by id
func (client *ClientData) GetStoreById(storeId string) Stores {
	var s Stores
	client.db.Model(&Stores{}).Where("Id = ?", storeId).Find(&s)
	return s
}

// CreateStoreWeights function to create weights for store
func (client *ClientData) CreateStoreWeights(storeRefer string, name string, beta float64, gama float64, delta float64,
	a float64, b float64, c float64, d float64, e float64, probabilityWeights string, shift int, longShift int) *gorm.DB {
	var s Stores
	client.db.Model(&Stores{}).Where("id = ?", storeRefer).First(&s)

	item := StoreWeights{
		StoreRefer:         s.Id.String(),
		Name:               name,
		Beta:               beta,
		Gama:               gama,
		Delta:              delta,
		A:                  a,
		B:                  b,
		C:                  c,
		D:                  d,
		E:                  e,
		ProbabilityWeights: probabilityWeights,
		Shift:              shift,
		LongShift:          longShift}
	result := client.db.Create(&item)
	return result
}

// EditStoreWeights function to edit weights for store
func (client *ClientData) EditStoreWeights(storeRefer string, name string, beta float64, gama float64, delta float64,
	a float64, b float64, c float64, d float64, e float64, probabilityWeights string, shift int, longShift int) StoreWeights {
	var storeWeights StoreWeights
	client.db.Model(&StoreWeights{}).Where("store_refer = ?", storeRefer).First(&storeWeights)
	storeWeights.Name = name
	storeWeights.Beta = beta
	storeWeights.Gama = gama
	storeWeights.Delta = delta
	storeWeights.A = a
	storeWeights.B = b
	storeWeights.C = c
	storeWeights.D = d
	storeWeights.E = e
	storeWeights.ProbabilityWeights = probabilityWeights
	storeWeights.Shift = shift
	storeWeights.LongShift = longShift
	client.db.Save(&storeWeights)
	return storeWeights
}

// GetStoreWeights funstion return weights for store
func (client *ClientData) GetStoreWeights(storeId string) StoreWeights {
	var storeWeights StoreWeights
	client.db.Model(&StoreWeights{}).Where("store_refer = ?", storeId).Find(&storeWeights)
	return storeWeights
}

// GetOpenData function return open data for store
func (client *ClientData) GetOpenData(storeRefer string) []OpenData {
	var od []OpenData
	client.db.Model(&OpenData{}).Where("store_refer = ?", storeRefer).Find(&od)
	return od
}

// CreateOpenData function to create open data for store
func (client *ClientData) CreateOpenData(storePower float64, customerSatisfaction float64, maximalProductPrice float64,
	minimalProductPrice float64, perceivedValue float64, storeRefer string) *gorm.DB {
	var s Stores
	client.db.Model(&Stores{}).Where("id = ?", storeRefer).First(&s)

	item := OpenData{
		StorePower:           storePower,
		CustomerSatisfaction: customerSatisfaction,
		MaximalProductPrice:  maximalProductPrice,
		MinimalProductPrice:  minimalProductPrice,
		PerceivedValue:       perceivedValue,
		StoreRefer:           s.Id.String(),
	}
	result := client.db.Create(&item)
	return result
}

// CreatePlan function to create new plan in database
func (client *ClientData) CreatePlan(name string, price float64, period int, products int,
	enabled bool, free bool) *gorm.DB {
	item := Plan{
		Name:     name,
		Price:    price,
		Period:   period,
		Products: products,
		Enabled:  enabled,
		Free:     free}
	result := client.db.Create(&item)
	return result
}

// EditPlan function to edit new plan in database
func (client *ClientData) EditPlan(id string, name string, price float64, period int, products int,
	enabled bool, free bool) Plan {
	var plan Plan
	client.db.Model(&Plan{}).Where("id = ?", id).First(&plan)
	plan.Name = name
	plan.Price = price
	plan.Period = period
	plan.Products = products
	plan.Enabled = enabled
	plan.Free = free
	client.db.Save(&plan)
	return plan
}

// GetPlans function to return all plans
func (client *ClientData) GetPlans() []Plan {
	var p []Plan
	client.db.Model(&Plan{}).Find(&p)
	return p
}

// GetPlanById function to return plan by id
func (client *ClientData) GetPlanById(planId string) Plan {
	var p Plan
	client.db.Model(&Plan{}).Where("id = ?", planId).Find(&p)
	return p
}

// GetPaidPlans function to return all paid plans
func (client *ClientData) GetPaidPlans() []Plan {
	var p []Plan
	client.db.Model(&Plan{}).Where("free = false AND enabled = true").Find(&p)
	return p
}

// DeletePlan function to delete plan from database
func (client *ClientData) DeletePlan(id string) {
	var plan Plan
	client.db.Model(&Plan{}).Where("id = ?", id).Delete(&plan)
}

// IsAvailableToView function to check if account is able to view store
func (client *ClientData) IsAvailableToView(accountId string, storeId string) Stores {
	var s Stores
	var a Accounts
	var id string
	client.db.Model(&Accounts{}).Where("id = ?", accountId).Find(&a)
	if a.Parent != "" {
		id = a.Parent
	} else {
		id = a.Id.String()
	}

	client.db.Model(&Stores{}).Where("id = ?", storeId).Where("account_refer = ?", id).Find(&s)
	return s
}

// GetSuppliers function return all suppliers for store
func (client *ClientData) GetSuppliers(storeId string) []Suppliers {
	var sup []Suppliers
	client.db.Model(&Suppliers{}).Where("store_refer = ?", storeId).Find(&sup)
	return sup
}

// GetSupplier function return supplier by id
func (client *ClientData) GetSupplier(supplierId string) Suppliers {
	var sup Suppliers
	client.db.Model(&Suppliers{}).Where("id = ?", supplierId).Find(&sup)
	return sup
}

// DeleteSupplier function to delete supplier by id
func (client *ClientData) DeleteSupplier(id string) {
	var sup Suppliers
	client.db.Model(&Suppliers{}).Where("id = ?", id).Delete(&sup)
}

// CreateSupplier function to create supplier in db
func (client *ClientData) CreateSupplier(name string, street string, city string, zip string, country string,
	email string, phone string, person string, storeRefer string, template string, subject string) *gorm.DB {
	var s Stores
	client.db.Model(&Stores{}).Where("id = ?", storeRefer).First(&s)

	item := Suppliers{
		Name:       name,
		Street:     street,
		Country:    country,
		City:       city,
		Zip:        zip,
		Email:      email,
		Phone:      phone,
		Person:     person,
		StoreRefer: s.Id.String(),
		Template:   template,
		Subject:    subject}
	result := client.db.Create(&item)
	return result
}

// EditSupplier function to edit supplier in db
func (client *ClientData) EditSupplier(id string, name string, street string, city string, zip string, country string,
	email string, phone string, person string, template string, subject string) Suppliers {
	var supplier Suppliers
	client.db.Model(&Suppliers{}).Where("id = ?", id).First(&supplier)
	supplier.Name = name
	supplier.Street = street
	supplier.Country = country
	supplier.City = city
	supplier.Zip = zip
	supplier.Email = email
	supplier.Phone = phone
	supplier.Person = person
	supplier.Template = template
	supplier.Subject = subject
	client.db.Save(&supplier)
	return supplier
}

// GetInvoices function to return all invoices for store
func (client *ClientData) GetInvoices(storeId string) []Invoices {
	var invoices []Invoices
	client.db.Model(&Invoices{}).Where("store_refer = ?", storeId).Find(&invoices)
	return invoices
}

// GetInvoices function to return all invoices for store
func (client *ClientData) GetInvoicesFilter(storeId string, from int, to int) []Invoices {
	var invoices []Invoices
	startDate := time.Now().AddDate(0, -from, 0)
	startDateString := startDate.Format("2006-01-02") + " 00:00:00"

	endDate := time.Now().AddDate(0, to, 0)
	endDateString := endDate.Format("2006-01-02") + " 00:00:00"

	client.db.Model(&Invoices{}).Where("store_refer = ?", storeId).Where("due_date > ?", startDateString).Where("due_date < ?", endDateString).Find(&invoices)
	return invoices
}

// CreateInvoice function to create invoice in database
func (client *ClientData) CreateInvoice(dueDate time.Time, amount float64, currency string, storeRefer string) *gorm.DB {
	var s Stores
	client.db.Model(&Stores{}).Where("id = ?", storeRefer).First(&s)

	item := Invoices{
		DueDate:    dueDate,
		Amount:     amount,
		Currency:   currency,
		StoreRefer: s.Id.String()}
	result := client.db.Create(&item)
	return result
}

// EditInvoice function to edit invoice in database
func (client *ClientData) EditInvoice(id string, dueDate time.Time, amount float64, currency string) Invoices {
	var invoice Invoices
	client.db.Model(&Invoices{}).Where("id = ?", id).First(&invoice)
	invoice.Amount = amount
	invoice.DueDate = dueDate
	invoice.Currency = currency
	client.db.Save(&invoice)
	return invoice
}

// DeleteInvoice function to delete invoice by id
func (client *ClientData) DeleteInvoice(id string) {
	var inv Invoices
	client.db.Model(&Invoices{}).Where("id = ?", id).Delete(&inv)
}

// GenerateOrder function to generate unique order number
func GenerateOrder() string {
	now := time.Now()
	sec := now.Unix()
	rand.Seed(time.Now().UnixNano())
	suffix := rand.Intn(100)
	return fmt.Sprintf("SP%d%d", sec, suffix)
}

// CreateOrder function to create order in db
func (client *ClientData) CreateOrder(accountRefer string, storeRefer string, planRefer string, amount float64, paid bool) string {
	account := client.GetAccountById(accountRefer)
	var a Accounts
	client.db.Model(&Accounts{}).Where("id = ?", accountRefer).First(&a)

	var s Stores
	client.db.Model(&Stores{}).Where("id = ?", storeRefer).First(&s)

	var p Plan
	client.db.Model(&Plan{}).Where("id = ?", planRefer).First(&p)

	item := Orders{
		AccountRefer:  a.Id.String(),
		StoreRefer:    s.Id.String(),
		PlanRefer:     p.Id.String(),
		Amount:        amount,
		Paid:          paid,
		Name:          account.Name,
		Email:         account.Email,
		Street:        account.Street,
		City:          account.City,
		Zip:           account.Zip,
		CountryCode:   account.CountryCode,
		CompanyNumber: account.CompanyNumber,
		VatNumber:     account.VatNumber,
		Number:        GenerateOrder()}
	client.db.Create(&item)
	return item.Id.String()
}

// GetOrders function to return all orders for account
func (client *ClientData) GetOrders(accountId string) []Orders {
	var ord []Orders
	client.db.Model(&Orders{}).Where("account_refer = ?", accountId).Find(&ord)
	return ord
}

// GetOrderById function to return order by id
func (client *ClientData) GetOrderById(id string) Orders {
	var ord Orders
	client.db.Model(&Orders{}).Where("id = ?", id).Find(&ord)
	return ord
}

// HashPassword function to hash pw string before save in db
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash function to compare hash pw and inserted string
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
