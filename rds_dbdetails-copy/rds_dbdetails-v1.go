package main

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"

	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type SaasDbDetail struct {
	DbHost               string   `json:"db_host"`
	CustomerSchemaCount  int      `json:"customer_schema_count"`
	CustomerSchemas      []string `json:"customer_schemas"`
	AvailableSchemaSlots int      `json: "available_schema_slots`
}

//Global variables to entire program

var (
	wg        sync.WaitGroup
	sess, err = session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")})
	sqlStatement   = `SELECT datname FROM pg_catalog.pg_database`
	maxSchemaCount = 5
)

func main() {

	r := gin.Default()
	r.GET("/dbdetails", fetchDetails)
	r.Run(":8080")
}

// 1. one primary function to be called from main returns JSON oject
// 2. separate function for describe db instances to be called from primary function
// 3. separate function for connect sql and return schemas, to be called from primary function
// 4. Apply go rountines and channels to function 2,3  and return the constructed array object in primary function.

func fetchDetails(c *gin.Context) {
	start := time.Now()
	// Create a New rds service
	rdsSvc := rds.New(sess)
	// Call to get detailed information on each DBinstance
	result, err := rdsSvc.DescribeDBInstances(nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	var dbDetails []SaasDbDetail
	// Adding lengths of dbInstances to wait group
	dbInstancesCount := len(result.DBInstances)
	wg.Add(dbInstancesCount)
	// Create a c hannel for final result.
	var SaasDBDetailChannel = make(chan SaasDbDetail, dbInstancesCount)
	for _, db := range result.DBInstances {
		go fetchDBDetails(db, SaasDBDetailChannel)
	}

	go func() {
		wg.Wait()
		close(SaasDBDetailChannel)
	}()

	for dbDetail := range SaasDBDetailChannel {
		log.Println("DB DETAIL IN CHANNEL RECEIVED: ", dbDetail)
		dbDetails = append(dbDetails, dbDetail)
	}

	log.Println("Elapsed time: ", time.Since(start))
	c.JSON(200, dbDetails)
}

func fetchDBDetails(dbInstance *rds.DBInstance, SaasDBDetailChannel chan SaasDbDetail) {
	defer wg.Done()
	dbHost := *dbInstance.Endpoint.Address

	var schemas []string
	if strings.HasPrefix(dbHost, "db-") {
		log.Println(dbHost)
		rows, err := execSqlQuery(sqlStatement, dbHost)
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			var schemaName string
			err := rows.Scan(&schemaName)
			if err != nil {
				panic(err)
			}
			if strings.Contains(schemaName, "dg") {
				schemas = append(schemas, schemaName)
			}
		}
		custSchemaCount := len(schemas)
		availableSlots := maxSchemaCount - custSchemaCount
		// fmt.Println(custSchemaCount)
		// fmt.Println(schemas)
		// fmt.Println(availableSlots)
		dbDetail := SaasDbDetail{DbHost: dbHost, CustomerSchemas: schemas, CustomerSchemaCount: custSchemaCount, AvailableSchemaSlots: availableSlots}
		SaasDBDetailChannel <- dbDetail
	}
}

func execSqlQuery(q string, dbHost string) (*sql.Rows, error) {
	dbPort := 5432
	dbUser := "postgres"
	dbPass := "dataguise"
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, "postgres")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	return rows, err
}
