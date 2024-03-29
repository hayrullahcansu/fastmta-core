package entity

import (
	"fmt"
	"sync"

	"github.com/hayrullahcansu/fastmta-core/global"

	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/now"
)

var once sync.Once

func GetDbContext() (*gorm.DB, error) {
	return openDb(
		global.StaticConfig.Database.Driver,
		global.StaticConfig.Database.Connection,
	)
}

// Create
//db.Create(&Product{Code: "L1212", Price: 1000})

// Read
//var product Product
//db.First(&product, 1)                   // find product with id 1
//db.First(&product, "code = ?", "L1212") // find product with code l1212

// Update - update product's price to 2000
//db.Model(&product).Update("Price", 2000)

// Delete - delete product
//db.Delete(&product)
func openDb(driver string, connection string) (*gorm.DB, error) {
	db, err := gorm.Open(driver, connection)
	if err != nil {
		panic("failed to connect database")
	}
	once.Do(func() {
		// Migrate the schema
		fmt.Println("migration first")
		db.AutoMigrate(
			&Message{},
			&Transaction{},
			&TransactionLog{},
			&Domain{},
			&MXRecord{},
			&Dkimmer{},
		)
		fmt.Println("migration done")
	})
	return db, err
}

func PanicOnError(err error) {
	if err != nil {
		logger.Panicf("Database exception: %s", err.Error())
	}
}
