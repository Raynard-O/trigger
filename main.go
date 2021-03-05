
package main

import (
	_ "github.com/go-sql-driver/mysql"
	"trigger/database"
)

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {

	DBReadFrom := "pytest"
	db := database.InterfaceDB(database.Db())

	defer db.Close()

	db.DBQuery(DBReadFrom, "id")

}
