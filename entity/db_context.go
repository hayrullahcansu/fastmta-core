package entity

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/now"
)

var once sync.Once

func Run() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
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

}

func GetDbContext() (*gorm.DB, error) {

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	once.Do(func() {
		// Migrate the schema
		fmt.Println("migration first")
		db.AutoMigrate(
			&Message{},
			&Transaction{},
			&Domain{},
			&MXRecord{},
		)
		fmt.Println("migration done")
	})
	return db, err
}
