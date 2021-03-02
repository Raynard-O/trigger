package database

import (
	"database/sql"
	"fmt"
)

type InterfaceDB interface {
	DBQuery(table, row string)
	Close()
}

type DB struct {
	Db	*sql.DB
}

const password = "INdu6990&!"
const user = "rthms"

func Db() *DB{
	data := fmt.Sprintf("%v:%v@tcp(198.71.225.63:3306)/RTHMS", user,password)
	db, err := sql.Open("mysql", data)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Go MySQL Tutorial")
	// defer the close till after the main function has finished
	// executing

	return &DB{Db: db}
}

//Close()
//close database
func (d *DB) Close()  {
	d.Db.Close()
}


type Tag struct {
	ID   int    `json:"id"`
	SensorID sql.NullInt64 `json:"sensor_id"`
	SampleID	int `json:"sample_id"`
}

//DBQuery
// enter table and close value for query
func (d *DB) DBQuery(table, row string)  {
	///results, err := db.Query("SELECT id, name FROM tags")
	query := fmt.Sprintf("SELECT %v, sample_id FROM %v", row, table)
	results, err := d.Db.Query(query)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	//results.Scan("id", "sensor_id")
	var last_id int
	for results.Next() {
		var tag Tag
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.ID, &tag.SampleID)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		if tag.SampleID!= last_id {
			fmt.Println("changes")
			fmt.Printf("new id %v, last id : %v\n", tag.SampleID, last_id)

		}

		last_id = tag.SampleID
	}

}

type Change struct {

}

//Insert
//insert data in new database when trigger
//takes in data change struct
func Insert(change string)  {
	fmt.Println(change)
}
