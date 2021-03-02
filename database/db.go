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
	Date sql.NullTime `json:"date"`
	x sql.NullInt64 `json:"x"`
	y sql.NullInt64 `json:"y"`
	z sql.NullInt64 `json:"z"`
}

//DBQuery
// enter table and close value for query
func (d *DB) DBQuery(table, row string)  {
	///results, err := db.Query("SELECT id, name FROM tags")
	query := fmt.Sprintf("SELECT %v, sample_id, x,y,z FROM %v", row, table)
	results, err := d.Db.Query(query)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	//results.Scan("id", "sensor_id")
	var last_id int
	var changeZ,changeY,changeX int
	var num int

	for results.Next() {
		var tag Tag
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.ID, &tag.SampleID, &tag.x, &tag.y, &tag.z)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		if tag.SampleID== last_id {
			num++
			changeX += int(tag.x.Int64)
			changeY += int(tag.y.Int64)
			changeZ += int(tag.z.Int64)
			//fmt.Println("same values")
		}
		if tag.SampleID!= last_id {
			if num == 0 {
				num = 1
			}
			newX := changeX/num
			newY := changeY/num
			newZ := changeZ/num
			//fmt.Printf("newX %v,newY %v, newZ %v,changeX : %v, changeY : %v, changeZ : %v, num : %v\n", newX, newY, newZ, changeX,changeY,changeZ,  num)
			num = 0
			changeX,changeY,changeZ = 0,0,0
			table_row := fmt.Sprintf("INSERT INTO operation_statistics (id, sample_id, x_mean, y_mean, z_mean) \nVALUES (%v,%v,%v,%v,%v)", last_id, tag.SampleID, newX,newY,newZ)
			d.Insert(table_row)

		}
		//("INSERT INTO test VALUES ( 2, 'TEST' )")
		last_id = tag.SampleID
	}

}

//"INSERT INTO operation_statistics(id, datetime, sample_id, x_mean, y_mean, z_mean) VALUES ( 2, 'TEST' )"
//Insert
//insert data in new database when trigger
//takes in data change struct
func (d *DB) Insert(change string)  {

	//perform a db.Query insert
	insert, err := d.Db.Query(change)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()


}

