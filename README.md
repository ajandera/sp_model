# Database model for storePredictor

- using Gorm 2 and UUID v4 as unique id.
- https://gorm.io/docs/

## Type of database
### rdbsClientInfo
- connector to application database
- store infromation about accounts and stores

### rdbsClientsData
- connector to clients data
- store data about visitor, proucts and orders

### noSqlClientPredictedData
- nosql infux database client
- tore information from prediction module

## contribution
- before create a tag commit all clients data
- after it run go mod tidy to update model.go
- then commit updates and create a new tag

## Multimodule structure
- using new multimodule more at https://go.dev/doc/tutorial/workspaces