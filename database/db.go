package database

import (
	"database/sql"
	"fmt"
	"math"
	"sync"
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
	Processed sql.NullBool `json:"processed"`
}


type TagOld struct {
	ID   int    `json:"id"`

	SampleID	int `json:"sample_id"`

}

var (
	DBWriteTo = "timitest"
)

//DBQuery
// enter table and close value for query
func (d *DB) DBQuery(table, row string)  {
	///results, err := db.Query("SELECT id, name FROM tags")
	query := fmt.Sprintf("SELECT %v, sample_id, x,y,z FROM %v", row, table)


	/// if id  is lesser than oldID dont do anything
	//else continue the process
	//fmt.Sprintf("SELECT @last_id := MAX(id) FROM table; SELECT * FROM table WHERE id = @last_id;")


	results, err := d.Db.Query(query)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	//results.Scan("id", "sensor_id")
	var last_id int
	var changeZ,changeY,changeX int
	var num int
	 valuesX,valuesY,valuesZ := []int{},[]int{},[]int{}


	f := fmt.Sprintf("SELECT sample_id FROM %v ORDER BY id DESC LIMIT 1", DBWriteTo)
	var tg Tag
	_= d.Db.QueryRow(f).Scan(&tg.SampleID)
	fmt.Printf("the last  processed sample id: %v\n", tg.SampleID)
	//r.NextResultSet()

	for results.Next() {
		//for i := 0; i < 1; i++ {
		//

		var tag Tag
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.ID, &tag.SampleID, &tag.x, &tag.y, &tag.z)
		//if tag.Processed.Bool {
		//	fmt.Println(tag.SampleID)
		if tg.SampleID < tag.SampleID {
			fmt.Println("hello")

			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}
			if tag.SampleID == last_id {
				num++
				//save values for standard cal
				valuesX = append(valuesX, int(tag.x.Int64))
				valuesY = append(valuesY, int(tag.y.Int64))
				valuesZ = append(valuesZ, int(tag.z.Int64))

				changeX += int(tag.x.Int64)
				changeY += int(tag.y.Int64)
				changeZ += int(tag.z.Int64)
				//fmt.Println("same values")
			}
			if tag.SampleID != last_id {
				if num == 0 {
					num = 1
				}
				var difff1 int
				for _, v := range valuesX {
					diff := v - changeX
					//fmt.Printf("diff: %v,,,,,,, Value : %v,,,,,, changeX: %v \n", diff, v, changeX)
					diff2 := diff * diff
					difff1 += diff2
					//neWvaluesX = append(neWvaluesX, diff2)
				}
				var difff2 int
				for _, v := range valuesY {
					diff := v - changeY
					diff2 := diff * diff
					difff2 += diff2
					//neWvaluesX = append(neWvaluesY, diff2)
				}
				var difff3 int
				for _, v := range valuesZ {
					diff := v - changeZ
					diff2 := diff * diff
					difff3 += diff2
					//neWvaluesX = append(neWvaluesZ, diff2)
				}

				varianceX := difff3 / num
				varianceY := difff2 / num
				varianceZ := difff1 / num

				sDX := math.Sqrt(float64(varianceX))
				sDY := math.Sqrt(float64(varianceY))
				sDZ := math.Sqrt(float64(varianceZ))
				//fmt.Printf("diff: %v,,,,,,, Value : %v,,,,,, changeX: %v standare : %v\n", difff3, num, varianceX,sDX)
				//fmt.Println(varianceX, varianceY, varianceZ, sDX, sDY, sDZ)
				newX := changeX / num
				newY := changeY / num
				newZ := changeZ / num

				//fmt.Printf("newX %v,newY %v, newZ %v,changeX : %v, changeY : %v, changeZ : %v, num : %v\n", newX, newY, newZ, changeX,changeY,changeZ,  num)

				tableRow := fmt.Sprintf("INSERT INTO %v (id, sample_id, x_mean,x_stdev, x_ms_ratio, y_mean,y_stdev, y_ms_ratio, z_mean, z_stdev, z_ms_ratio ) \nVALUES (%v,%v,%v,%v,%v,%v,%v,%v,%v, %v,%v)",DBWriteTo, tag.ID, tag.SampleID, newX, sDX, varianceX, newY, sDY, varianceY, newZ, sDZ, varianceZ)

				//tableRow2 := fmt.Sprintf("INSERT INTO pytest (processed) \nVALUES (%v)", true)

				var wg sync.WaitGroup
				wg.Add(1)

				go func() {
					defer wg.Done()
					d.Insert(tableRow)
				}()


				wg.Wait()
				num = 0
				changeX, changeY, changeZ = 0, 0, 0
			}
			//("INSERT INTO test VALUES ( 2, 'TEST' )")
			last_id = tag.SampleID
		}
	}
	//}
}

//"INSERT INTO operation_statistics(id, datetime, sample_id, x_mean, y_mean, z_mean) VALUES ( 2, 'TEST' )"
//Insert
//insert data in new database when trigger
//takes in data change struct
func (d *DB) Insert(change string)  {


	//fmt.Println(change)
	//perform a db.Query insert
	insert, err := d.Db.Query(change)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

}

