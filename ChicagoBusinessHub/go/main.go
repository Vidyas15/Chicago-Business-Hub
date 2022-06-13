package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

//type HandlerFunc func(ResponseWriter, *Request) // ServeHTTP calls f(w, r).

type CovidCCVI []struct {
	GeographyType      string   `json:"geoar_type"`
	CommunityAreaOrZip string   `json:"community_area_or_zip"`
	CcviScore          string   `json:"ccvi_score"`
	CcviCategory       string   `json:"ccvi_category"`
	Location           Location `json:"location"`
}

type Location struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type TAXI_TRIPS_DATA []struct {
	TRIP_ID         string `json:"trip_id"`
	TRIP_START_TIME string `json:"trip_start_timestamp"`
	TRIP_END_TIME   string `json:"trip_end_timestamp"`
	PICKUP_LAT      string `json:"pickup_centroid_latitude"`
	PICKUP_LONG     string `json:"pickup_centroid_longitude"`
	DROPOFF_LAT     string `json:"dropoff_centroid_latitude"`
	DROPOFF_LONG    string `json:"dropoff_centroid_longitude"`
}

type CovidDaily []struct {
	Date        string `json:"lab_report_date"`
	TotalCases  string `json:"cases_total"`
	TotalDeaths string `json:"deaths_total"`
}

type Covidweekly []struct {
	zipcode1           string `json:"zip_code"`
	week1              string `json:"week_end"`
	cases_weekly1      string `json:"cases_weekly"`
	cases_cumulative1  string `json:"cases_cumulative"`
	deaths_weekly1     string `json:"deaths_weekly"`
	deaths_cumulative1 string `json:"deaths_cumulative"`
	population1        string `json:"population"`
	tests_weekly1      string `json:"tests_weekly"`
}

type building_permits []struct {
	permit_type            string `json:"permit_type"`
	APPLICATION_START_DATE string `json:"application_start_date"`
	ISSUE_DATE             string `json:"issue_date"`
	BUILDING_FEE_PAID      string `json:"building_fee_paid"`
	SUBTOTAL_PAID          string `json:"subtotal_paid"`
	SUBTOTAL_UNPAID        string `json:"subtotal_unpaid"`
	COMMUNITY_AREA         string `json:"community_area"`
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there %s!", r.URL.Path[1:])
}

const (
	host     = "host.docker.internal"
	port     = 54320
	user     = "postgres"
	password = "HaveFun"
	dbname   = "chicago_bi"
)

func main() {
	//***********connetion to postgres****change pwd

	//connStr := "user=postgres dbname=chicago_bi password=1404 host=host.docker.internal port=54320 sslmode=disable"
	//db, err := sql.Open("postgres".connStr)
	fmt.Printf(" Start updating of DB Table")
	connStr_0 := "user=postgres password=HaveFun host=host.docker.internal port=54320 sslmode=disable"
	db, err := sql.Open("postgres", connStr_0)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connection to postgres established")

	createDB := `CREATE DATABASE chicago_bi`
	_, err = db.Exec(createDB)
	// if err != nil {
	// 	panic(err)
	// }

	db.Close()

	connStr := "user=postgres dbname=chicago_bi password=HaveFun host=host.docker.internal port=54320 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connection to postgres established")

	//***********Getting covid ccvi data from chicago web portal

	dropSql2 := `drop table if exists covid_ccvi`
	_, err12 := db.Exec(dropSql2)
	if err != nil {
		panic(err12)
	}

	createSql := `CREATE TABLE IF NOT EXISTS "covid_ccvi" (
	"id" SERIAL,	
	"geographytype" VARCHAR(255),
	"zipcode" VARCHAR(255),
	"ccviscore" DOUBLE PRECISION,
	"ccvicategory" VARCHAR(255),
	"latitude" DOUBLE PRECISION,
	"longitude" DOUBLE PRECISION,
	PRIMARY KEY ("id"));`
	_, createSqlErr := db.Exec(createSql)
	if createSqlErr != nil {
		panic(createSqlErr)
	}
	var url1 = "https://data.cityofchicago.org/resource/xhc6-88s9.json"
	res, err := http.Get(url1)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(res.Body)

	var covidDataArray CovidCCVI
	json.Unmarshal(body, &covidDataArray)

	for i := 0; i < len(covidDataArray); i++ {
		CommunityAreaOrZipcode := covidDataArray[i].CommunityAreaOrZip
		ccviScore, _ := strconv.ParseFloat(covidDataArray[i].CcviScore, 64)
		//createdAt := time.Now()
		lat := covidDataArray[i].Location.Coordinates[1]
		lng := covidDataArray[i].Location.Coordinates[0]
		sql := `INSERT into covid_ccvi("geographytype","zipcode","ccviscore","ccvicategory","latitude","longitude") VALUES($1, $2, $3, $4, $5, $6)`
		_, err := db.Exec(sql, covidDataArray[i].GeographyType, CommunityAreaOrZipcode, ccviScore, covidDataArray[i].CcviCategory, lat, lng)
		if err != nil {
			panic(err)
		}

	}

	fmt.Printf("done with covid_ccvi")

	//*************data engineered ccvi final with community area read into db using GO

	createSql244d := `CREATE TABLE IF NOT EXISTS "covid_final" ("id" SERIAL,
	"zipcode" VARCHAR(255),
	"neighborhood" VARCHAR(255),
	"ccviscore" DOUBLE PRECISION,
	"ccvicategory" VARCHAR(255),
	"latitude" VARCHAR(255),
	"longitude" VARCHAR(255),
	PRIMARY KEY ("id"));`
	_, createSqlErr244c := db.Exec(createSql244d)
	if createSqlErr244c != nil {
		panic(createSqlErr244c)
	}
	csvfile45ff, err99 := os.Open("ccvi_final2.csv")

	if err99 != nil {
		log.Fatalln("Couldn't open the csv file", err99)
	}

	// Parse the file
	r4466 := csv.NewReader(csvfile45ff)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r4466.Read()
		if err == io.EOF {
			break
		}
		sql := `INSERT into covid_final("zipcode","neighborhood","ccviscore","ccvicategory","latitude","longitude") VALUES($1, $2, $3, $4, $5, $6)`
		_, err7 := db.Exec(sql, record[0], record[1], record[2], record[3], record[4], record[5])
		if err7 != nil {
			panic(err7)
		}
	}

	//after processing it in python for community area reading and writing it postgres

	// dropSql2 := `drop table if exists taxi_trips`
	// _, err12 := db.Exec(dropSql2)
	// if err != nil {
	// 	panic(err12)
	// }

	//***********getting taxi_ trips from chicago data portal

	// createSql2 := `CREATE TABLE IF NOT EXISTS "taxi_trips" ("tripid" VARCHAR(255),
	// "tripstarttime" VARCHAR(255),
	// "tripendtime" VARCHAR(255),
	// "pickuplat" VARCHAR(255),
	// "pickuplong" VARCHAR(255),
	// "dropofflat" VARCHAR(255),
	// "dropofflong" VARCHAR(255),
	// PRIMARY KEY ("tripid"));`
	// _, createSqlErr2 := db.Exec(createSql2)
	// if createSqlErr2 != nil {
	// 	panic(createSqlErr2)
	// }
	// var url2 = "https://data.cityofchicago.org/resource/wrvz-psew.json"
	// res, err22 := http.Get(url2)
	// if err22 != nil {
	// 	log.Fatal(err22)
	// }

	// body2, _ := ioutil.ReadAll(res.Body)

	// var taxiarray TAXI_TRIPS_DATA
	// json.Unmarshal(body2, &taxiarray)

	// for i := 0; i < len(taxiarray); i++ {

	// 	tripid := taxiarray[i].TRIP_ID
	// 	tripstarttime := taxiarray[i].TRIP_START_TIME
	// 	tripendtime := taxiarray[i].TRIP_END_TIME

	// 	//createdAt := time.Now()
	// 	pickuplat := taxiarray[i].PICKUP_LAT
	// 	pickuplong := taxiarray[i].PICKUP_LONG
	// 	dropofflat := taxiarray[i].DROPOFF_LAT
	// 	dropofflong := taxiarray[i].DROPOFF_LONG
	// 	sql := `INSERT into taxi_trips("tripid","tripstarttime","tripendtime","pickuplat","pickuplong","dropofflat","dropofflong") VALUES($1, $2, $3, $4, $5, $6, $7)`
	// 	_, err6 := db.Exec(sql, tripid, tripstarttime, tripendtime, pickuplat, pickuplong, dropofflat, dropofflong)
	// 	if err6 != nil {
	// 		panic(err6)
	// 	}

	// }

	// fmt.Printf("done with Taxi trips")

	// csvfile, err9 := os.Open("CommAreas.csv")
	// csvfile2, err10 := os.Open("zip.csv")

	// if err9 != nil {
	// 	log.Fatalln("Couldn't open the csv file", err9)
	// }

	// if err10 != nil {
	// 	log.Fatalln("Couldn't open the csv file", err10)
	// }

	// // Parse the file
	// r := csv.NewReader(csvfile)
	// r2 := csv.NewReader(csvfile2)

	// //r := csv.NewReader(bufio.NewReader(csvfile))
	// createSql3 := `CREATE TABLE IF NOT EXISTS "community_area" ("areano" VARCHAR(255),
	// "communityareas" VARCHAR(255));`
	// _, createSqlErr3 := db.Exec(createSql3)
	// if createSqlErr3 != nil {
	// 	panic(createSqlErr3)
	// }

	// // Iterate through the records
	// for {
	// 	// Read each record from csv
	// 	record, err := r.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	sql := `INSERT into community_area("areano","communityareas") VALUES($1, $2)`
	// 	_, err7 := db.Exec(sql, record[5], record[6])
	// 	if err7 != nil {
	// 		panic(err7)
	// 	}

	// }
	// createSql4 := `CREATE TABLE IF NOT EXISTS "zipcodes" ("zipcode" VARCHAR(255),
	// "communityarea" VARCHAR(255));`
	// _, createSqlErr4 := db.Exec(createSql4)
	// if createSqlErr4 != nil {
	// 	panic(createSqlErr4)
	// }

	// // Iterate through the records
	// for {
	// 	// Read each record from csv
	// 	record2, err := r2.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	sql := `INSERT into zipcodes("zipcode","communityarea") VALUES($1, $2)`
	// 	_, err8 := db.Exec(sql, record2[0], record2[1])
	// 	if err8 != nil {
	// 		panic(err8)
	// 	}

	// }

	// // Iterate through the records

	// sql12 := `Select z.zipcode,c.communityareas into zipcodes_and_CA from zipcodes z join community_area c on z.communityarea=c.areano;`
	// _, err15 := db.Exec(sql12)
	// if err15 != nil {
	// 	panic(err15)
	// }

	// fmt.Printf("done with zipcodes and community area")

	// dropSql24 := `drop table if exists covid_daily`
	// _, err123 := db.Exec(dropSql24)
	// if err123 != nil {
	// 	panic(err123)
	// }

	//***********getting covid daily data from chicago web portal

	createSql5 := `CREATE TABLE IF NOT EXISTS "covid_daily" ("reportdate" VARCHAR(255),
	"cases_total" integer,
	"deaths_total" integer);`
	_, createSqlErr5 := db.Exec(createSql5)
	if createSqlErr5 != nil {
		panic(createSqlErr5)
	}

	var url3 = "https://data.cityofchicago.org/resource/naz8-j4nc.json"
	res, err23 := http.Get(url3)
	if err23 != nil {
		log.Fatal(err23)
	}
	body3, _ := ioutil.ReadAll(res.Body)

	var covidarry1 CovidDaily
	json.Unmarshal(body3, &covidarry1)

	for i := 0; i < len(covidarry1); i++ {
		//date1, _ := covidarry1[i].Date
		sql := `INSERT into covid_daily("reportdate","cases_total","deaths_total") VALUES($1, $2, $3)`
		_, err18 := db.Exec(sql, covidarry1[i].Date, covidarry1[i].TotalCases, covidarry1[i].TotalDeaths)
		if err18 != nil {
			panic(err18)
		}

	}

	//***********performing data processing on covid daily

	dropSql24444 := `drop table if exists covid_daily_report`
	_, err1244ddd := db.Exec(dropSql24444)
	if err1244ddd != nil {
		panic(err1244ddd)
	}

	sql1234 := `SELECT substr(reportdate,1,10) as report_date,cases_total,deaths_total into covid_daily_report FROM covid_daily;`
	_, err155 := db.Exec(sql1234)
	if err155 != nil {
		panic(err155)
	}
	fmt.Printf("done with covid_daily")
	dropSql244 := `drop table if exists covid_zip`
	_, err1244 := db.Exec(dropSql244)
	if err != nil {
		panic(err1244)
	}

	createSql244 := `CREATE TABLE IF NOT EXISTS "covid_zip" ("id" SERIAL,
	"communityareas" VARCHAR(255),
	"zipcode" VARCHAR(255),
	"week" VARCHAR(255),
	"cases_weekly" VARCHAR(255),
	"cases_cumulative" VARCHAR(255),
	"deaths_weekly" VARCHAR(255),
	"deaths_cumulative" VARCHAR(255),
	"population" VARCHAR(255),
	"tests_weekly" VARCHAR(255),
	PRIMARY KEY ("id"));`
	_, createSqlErr244 := db.Exec(createSql244)
	if createSqlErr244 != nil {
		panic(createSqlErr244)
	}
	csvfile45, err99 := os.Open("covid_week_fin.csv")

	if err99 != nil {
		log.Fatalln("Couldn't open the csv file", err99)
	}

	// Parse the file
	r44 := csv.NewReader(csvfile45)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r44.Read()
		if err == io.EOF {
			break
		}
		sql := `INSERT into covid_zip("communityareas","zipcode","week","cases_weekly","cases_cumulative","deaths_weekly","deaths_cumulative","population","tests_weekly") VALUES($1, $2, $3, $4, $5, $6, $7, $8,$9)`
		_, err7 := db.Exec(sql, record[0], record[2], record[3], record[4], record[5], record[6], record[7], record[8], record[9])
		if err7 != nil {
			panic(err7)
		}
	}

	// /*// for i := 0; i < len(covidweekarr); i++ {
	// // 	//fmt.Printf("\n", i)

	// // 	zip := covidweekarr[i].zipcode1
	// // 	wk := covidweekarr[i].week1
	// // 	cw := covidweekarr[i].cases_weekly1
	// // 	fmt.Printf(covidweekarr[i].zipcode1)

	// // 	//createdAt := time.Now()
	// // 	cc := covidweekarr[i].cases_cumulative1
	// // 	dw := covidweekarr[i].deaths_weekly1
	// // 	dc := covidweekarr[i].deaths_cumulative1
	// // 	p := covidweekarr[i].population1
	// // 	tw := covidweekarr[i].tests_weekly1

	// // 	sql101 := `INSERT into covid_zip("zipcode","week","cases_weekly","cases_cumulative","deaths_weekly","deaths_cumulative","population","tests_weekly") VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	// // 	_, err611 := db.Exec(sql101, zip, wk, cw, cc, dw, dc, p, tw)
	// // 	if err611 != nil {
	// // 		panic(err611)
	// // 	}

	// // }
	// */

	fmt.Printf("done with covid_zip")
	// dropSql255 := `drop table if exists poverty`
	// _, err12445 := db.Exec(dropSql255)
	// if err != nil {
	// 	panic(err12445)
	// }

	// createSql2445 := `CREATE TABLE IF NOT EXISTS "poverty" ("id" SERIAL,
	// "community_area" VARCHAR(255),
	// "below_poverty_percent" VARCHAR(255),
	// "per_capita_income" VARCHAR(255),
	// "hardship_index" VARCHAR(255),
	// PRIMARY KEY ("id"));`
	// _, createSqlErr2445 := db.Exec(createSql2445)
	// if createSqlErr2445 != nil {
	// 	panic(createSqlErr2445)
	// }
	// csvfile451, err999 := os.Open("poverty.csv")

	// if err999 != nil {
	// 	log.Fatalln("Couldn't open the csv file", err999)
	// }

	// // Parse the file
	// r445 := csv.NewReader(csvfile451)

	// // Iterate through the records
	// for {
	// 	// Read each record from csv
	// 	record, err := r445.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	sql := `INSERT into poverty("community_area","below_poverty_percent","per_capita_income","hardship_index") VALUES($1, $2, $3, $4)`
	// 	_, err71 := db.Exec(sql, record[1], record[3], record[7], record[8])
	// 	if err71 != nil {
	// 		panic(err71)
	// 	}
	// }

	// fmt.Printf("done with poverty")
	// dropSql25532 := `drop table if exists taxi_zip_ca`
	// _, err1244533 := db.Exec(dropSql25532)
	// if err1244533 != nil {
	// 	panic(err1244533)
	// }

	//************running file processed from python after data engineering for taxi trips***********

	createSql24453 := `CREATE TABLE IF NOT EXISTS "taxi_zip_ca" ( "id" SERIAL,
	"tripid" VARCHAR(255),
	"tripstarttime" VARCHAR(255),
	"tripendtime" VARCHAR(255),
	"pickuplat" VARCHAR(255),
	"pickuplong" VARCHAR(255),
	"dropofflat" VARCHAR(255),
	"dropofflong" VARCHAR(255),
	"pick_zipcode" VARCHAR(255),
	"drop_zipcode" VARCHAR(255),
	"pick_ca" VARCHAR(255),
	"drop_ca" VARCHAR(255),
	PRIMARY KEY ("id"));`
	_, createSqlErr24453 := db.Exec(createSql24453)
	if createSqlErr24453 != nil {
		panic(createSqlErr24453)
	}
	csvfile4512, err9991 := os.Open("taxi_zip_ca.csv")

	if err9991 != nil {
		log.Fatalln("Couldn't open the csv file", err9991)
	}

	// Parse the file
	r4451 := csv.NewReader(csvfile4512)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r4451.Read()
		if err == io.EOF {
			break
		}
		sql := `INSERT into taxi_zip_ca("tripid","tripstarttime","tripendtime","pickuplat","pickuplong","dropofflat","dropofflong","pick_zipcode","drop_zipcode","pick_ca","drop_ca") VALUES($1, $2, $3, $4,$5, $6, $7, $8,$9, $10, $11)`
		_, err712 := db.Exec(sql, record[2], record[3], record[4], record[5], record[6], record[7], record[8], record[9], record[10], record[11], record[12])
		if err712 != nil {
			panic(err712)
		}
	}

	//**********getting the data for the taxi_trips 2021 after it was data engineered in python and sending it to postgres

	dropSql25532s := `drop table if exists taxi_2021`
	_, err1244533d := db.Exec(dropSql25532s)
	if err1244533d != nil {
		panic(err1244533d)
	}
	createSql244522 := `CREATE TABLE taxi_2021 (
		id1 serial,
		trip_id VARCHAR(255),
		taxi_id VARCHAR(255),
		trip_start_timestamp VARCHAR(255),
		trip_end_timestamp VARCHAR(255),
		trip_miles VARCHAR(255),
		pickup_community_area VARCHAR(255),
		dropoff_community_area VARCHAR(255),
		trip_total	VARCHAR(255),
		company VARCHAR(255),
		pick_ca VARCHAR(255),
		drop_ca VARCHAR(255),
		PRIMARY KEY (id1)
	  );`
	_, createSqlErr24452 := db.Exec(createSql244522)
	if createSqlErr24452 != nil {
		panic(createSqlErr24452)
	}
	csvfile45122, err9992 := os.Open("taxi_pred_m.csv")

	if err9992 != nil {
		log.Fatalln("Couldn't open the csv file", err9992)
	}

	// Parse the file
	r44523 := csv.NewReader(csvfile45122)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r44523.Read()
		if err == io.EOF {
			break
		}
		sql := `INSERT into taxi_2021("trip_id","taxi_id","trip_start_timestamp","trip_end_timestamp","trip_miles","pickup_community_area","dropoff_community_area","trip_total","company","pick_ca","drop_ca") VALUES($1, $2, $3, $4, $5,$6, $7, $8, $9, $10, $11)`
		_, err71 := db.Exec(sql, record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7], record[8], record[9], record[10])
		if err71 != nil {
			panic(err71)
		}
	}
	// dropSql25544 := `drop table if exists poverty`
	// _, err1244544 := db.Exec(dropSql25544)
	// if err != nil {
	// 	panic(err1244544)
	// }

	// createSql2445 := `CREATE TABLE IF NOT EXISTS "poverty" ("id" SERIAL,
	// "community_area" VARCHAR(255),
	// "below_poverty_percent" VARCHAR(255),
	// "percent_unemployed" VARCHAR(255),
	// "per_capita_income" VARCHAR(255),
	// "hardship_index" VARCHAR(255),
	// PRIMARY KEY ("id"));`
	// _, createSqlErr2445 := db.Exec(createSql2445)
	// if createSqlErr2445 != nil {
	// 	panic(createSqlErr2445)
	// }
	// csvfile451, err999 := os.Open("poverty.csv")

	// if err999 != nil {
	// 	log.Fatalln("Couldn't open the csv file", err999)
	// }

	// // Parse the file
	// r445 := csv.NewReader(csvfile451)

	// // Iterate through the records
	// for {
	// 	// Read each record from csv
	// 	record, err := r445.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	sql := `INSERT into poverty("community_area","below_poverty_percent","percent_unemployed","per_capita_income","hardship_index") VALUES($1, $2, $3, $4, $5)`
	// 	_, err71 := db.Exec(sql, record[1], record[3], record[4], record[7], record[8])
	// 	if err71 != nil {
	// 		panic(err71)
	// 	}
	// }
	////https://data.cityofchicago.org/resource/building_permits.json
	// dropSql24566 := `drop table if exists building_permits`
	// _, err12345 := db.Exec(dropSql24566)
	// if err12345 != nil {
	// 	panic(err12345)
	// }

	//*************fetching building permits from chicago data portal exporting it to csv and python for community area processing

	createSql512 := `CREATE TABLE IF NOT EXISTS "building_permits2" ("ID" SERIAL,
	"permit_type" VARCHAR(255),
	"APPLICATION_START_DATE" VARCHAR(255),
	"ISSUE_DATE" VARCHAR(255),
	"BUILDING_FEE_PAID" VARCHAR(255),
	"SUBTOTAL_PAID" VARCHAR(255),
	"SUBTOTAL_UNPAID"  VARCHAR(255),
	"CA_code" VARCHAR(255));`
	_, createSqlErr533 := db.Exec(createSql512)
	if createSqlErr533 != nil {
		panic(createSqlErr533)
	}

	var url334 = "https://data.cityofchicago.org/resource/building-permits.json"
	res, err2324 := http.Get(url334)
	if err2324 != nil {
		log.Fatal(err2324)
	}
	body333, _ := ioutil.ReadAll(res.Body)

	var bparr building_permits
	json.Unmarshal(body333, &bparr)

	for i := 0; i < len(bparr); i++ {
		//date1, _ := bparr[i].Date
		sql := `INSERT into building_permits2("permit_type","APPLICATION_START_DATE","ISSUE_DATE","BUILDING_FEE_PAID","SUBTOTAL_PAID","SUBTOTAL_UNPAID","CA_code") VALUES($1, $2, $3,$4,$5,$6,$7)`
		_, err18 := db.Exec(sql, bparr[i].permit_type, bparr[i].APPLICATION_START_DATE, bparr[i].ISSUE_DATE, bparr[i].BUILDING_FEE_PAID, bparr[i].SUBTOTAL_PAID, bparr[i].SUBTOTAL_UNPAID, bparr[i].COMMUNITY_AREA)
		if err18 != nil {
			panic(err18)
		}

	}

	dropSql25534 := `drop table if exists building_permits`
	_, err124453 := db.Exec(dropSql25534)
	if err != nil {
		panic(err124453)
	}

	//*******************after community area is processed and written into the building permits using python it is read back into database

	createSql5123 := `CREATE TABLE IF NOT EXISTS "building_permits" ("ID" SERIAL,
	"permit_no" VARCHAR(255),
	"permit_type" VARCHAR(255),
	"APPLICATION_START_DATE" VARCHAR(255),
	"ISSUE_DATE" VARCHAR(255),
	"LATITUDE" VARCHAR(255),
	"LONGITUDE" VARCHAR(255),
	"zipcode"  VARCHAR(255),
	"community_area" VARCHAR(255));`
	_, createSqlErr5333 := db.Exec(createSql5123)
	if createSqlErr5333 != nil {
		panic(createSqlErr5333)
	}

	csvfile4521, err9991 := os.Open("building_main.csv")

	if err9991 != nil {
		log.Fatalln("Couldn't open the csv file", err9991)
	}

	// Parse the file
	r445 := csv.NewReader(csvfile4521)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r445.Read()
		if err == io.EOF {
			break
		}
		sql := `INSERT into building_permits("permit_no","permit_type","APPLICATION_START_DATE","ISSUE_DATE","LATITUDE","LONGITUDE","zipcode","community_area") VALUES($1, $2, $3,$4,$5,$6,$7,$8)`
		_, err71 := db.Exec(sql, record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7])
		if err71 != nil {
			panic(err71)
		}
	}
	dropSql255dd := `drop table if exists poverty`
	_, err12445d := db.Exec(dropSql255dd)
	if err12445d != nil {
		panic(err12445d)
	}

	//************poverty with community area file is read and exported to  db

	createSql2445 := `CREATE TABLE IF NOT EXISTS "poverty" ("id" SERIAL,
	"community_area" VARCHAR(255),
	"below_poverty_percent" VARCHAR(255),
	"percent_unemployed" VARCHAR(255),
	"per_capita_income" VARCHAR(255),
	"hardship_index" VARCHAR(255),
	"zipcode" VARCHAR(255),
	PRIMARY KEY ("id"));`
	_, createSqlErr2445 := db.Exec(createSql2445)
	if createSqlErr2445 != nil {
		panic(createSqlErr2445)
	}
	csvfile451, err999 := os.Open("poverty_main.csv")

	if err999 != nil {
		log.Fatalln("Couldn't open the csv file", err999)
	}

	// Parse the file
	r445a := csv.NewReader(csvfile451)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r445a.Read()
		if err == io.EOF {
			break
		}
		sql := `INSERT into poverty("community_area","below_poverty_percent","percent_unemployed","per_capita_income","hardship_index","zipcode") VALUES($1, $2, $3, $4, $5,$6)`
		_, err71 := db.Exec(sql, record[1], record[3], record[4], record[7], record[8], record[9])
		if err71 != nil {
			panic(err71)
		}
	}

	//***********facebook prophet predicted data exported from python and read into db using go lang

	dropSql25534d := `drop table if exists mo_oh`
	_, err124453d := db.Exec(dropSql25534d)
	if err124453d != nil {
		panic(err124453d)
	}
	createSql2445a := `CREATE TABLE IF NOT EXISTS "mo_oh" ("id" SERIAL,
	"date1" DATE,
	"traffic" DOUBLE PRECISION,
	PRIMARY KEY ("id"));`
	_, createSqlErr2445c := db.Exec(createSql2445a)
	if createSqlErr2445c != nil {
		panic(createSqlErr2445c)
	}
	csvfile451s, err999d := os.Open("mon_oh2.csv")

	if err999d != nil {
		log.Fatalln("Couldn't open the csv file", err999d)
	}

	// Parse the file
	r445r := csv.NewReader(csvfile451s)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r445r.Read()
		if err == io.EOF {
			break
		}
		sql := `INSERT into mo_oh("date1","traffic") VALUES($1, $2)`
		_, err71 := db.Exec(sql, record[1], record[16])
		if err71 != nil {
			panic(err71)
		}
	}

	//***********facebook prophet predicted data exported from python and read into db using go lang

	// createSql2445433 := `CREATE TABLE IF NOT EXISTS "mo_lk" ("id" SERIAL,
	// "Date" DATE,
	// "Traffic" VARCHAR(255),
	// PRIMARY KEY ("id"));`
	// _, createSqlErr244533 := db.Exec(createSql2445433)
	// if createSqlErr244533 != nil {
	// 	panic(createSqlErr244533)
	// }
	// csvfile45133, err99933 := os.Open("mon_lk.csv")

	// if err99933 != nil {
	// 	log.Fatalln("Couldn't open the csv file", err99933)
	// }

	// // Parse the file
	// r44533 := csv.NewReader(csvfile45133)

	// // Iterate through the records
	// for {
	// 	// Read each record from csv
	// 	record, err := r44533.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	sql := `INSERT into mo_lk("Date","Traffic") VALUES($1, $2)`
	// 	_, err71 := db.Exec(sql, record[1], record[16])
	// 	if err71 != nil {
	// 		panic(err71)
	// 	}
	// }

	// //***********facebook prophet predicted data exported from python and read into db using go lang

	// createSql244543333 := `CREATE TABLE IF NOT EXISTS "mo_lp" ("id" SERIAL,
	// "Date" DATE,
	// "Traffic" VARCHAR(255),
	// PRIMARY KEY ("id"));`
	// _, createSqlErr24453333 := db.Exec(createSql244543333)
	// if createSqlErr24453333 != nil {
	// 	panic(createSqlErr24453333)
	// }
	// csvfile451334, err99933 := os.Open("mon_lp.csv")

	// if err99933 != nil {
	// 	log.Fatalln("Couldn't open the csv file", err99933)
	// }

	// // Parse the file
	// r445331 := csv.NewReader(csvfile451334)

	// // Iterate through the records
	// for {
	// 	// Read each record from csv
	// 	record, err := r445331.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	sql := `INSERT into mo_lp("Date","Traffic") VALUES($1, $2)`
	// 	_, err71 := db.Exec(sql, record[1], record[16])
	// 	if err71 != nil {
	// 		panic(err71)
	// 	}
	// }

	// //***********facebook prophet predicted data exported from python and read into db using go lang

	// createSql2445433334 := `CREATE TABLE IF NOT EXISTS "mo_sc" ("id" SERIAL,
	// "Date" DATE,
	// "Traffic" VARCHAR(255),
	// PRIMARY KEY ("id"));`
	// _, createSqlErr244533333 := db.Exec(createSql2445433334)
	// if createSqlErr244533333 != nil {
	// 	panic(createSqlErr244533333)
	// }
	// csvfile4512334, err99933 := os.Open("mon_sc.csv")

	// if err99933 != nil {
	// 	log.Fatalln("Couldn't open the csv file", err99933)
	// }

	// // Parse the file
	// r4453331 := csv.NewReader(csvfile4512334)

	// // Iterate through the records
	// for {
	// 	// Read each record from csv
	// 	record, err := r4453331.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	sql := `INSERT into mo_sc("Date","Traffic") VALUES($1, $2)`
	// 	_, err71 := db.Exec(sql, record[1], record[16])
	// 	if err71 != nil {
	// 		panic(err71)
	// 	}
	// }

	//***********SQL table for monthly traffic pattern

	sql555 := `select sum(traffic) as total_traffic,EXTRACT(MONTH FROM date1) as month  into mh_oh_f2 from mo_oh
	group by EXTRACT(MONTH FROM date1)
	order by month`
	_, err71e := db.Exec(sql555)
	if err71e != nil {
		panic(err71e)
	}

	//***********SQL table for weekly traffic pattern

	sql555w := `select sum(traffic) as total_traffic,EXTRACT(week FROM date1) as week  into we_oh2 from mo_oh
	group by EXTRACT(week FROM date1)
	order by week`
	_, err71ed := db.Exec(sql555w)
	if err71ed != nil {
		panic(err71ed)
	}

}
