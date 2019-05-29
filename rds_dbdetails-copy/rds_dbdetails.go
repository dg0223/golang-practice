package main

import (
	"database/sql"

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

func main() {

	r := gin.Default()
	r.GET("/dbdetails", func(c *gin.Context) {

		// creating session
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1")},
		)

		// Create a New rds service
		rdsSvc := rds.New(sess)

		// Call to get detailed information on each DBinstance
		result, err := rdsSvc.DescribeDBInstances(nil)
		if err != nil {
			fmt.Println("Error", err)
		}

		dbPort := 5432
		dbUser := "postgres"
		dbPass := "dataguise"

		maxSchemaCount := 5
		var dbDetails []SaasDbDetail

		for _, dbInstances := range result.DBInstances {
			dbHost := *dbInstances.Endpoint.Address
			var schemas []string
			if strings.Contains(dbHost, "db-") {
				fmt.Println("Endpoint: ", dbHost)

				psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
					"password=%s dbname=%s sslmode=disable",
					dbHost, dbPort, dbUser, dbPass, "postgres")
				db, err := sql.Open("postgres", psqlInfo)
				if err != nil {
					panic(err)
				}
				defer db.Close()

				err = db.Ping()
				if err != nil {
					panic(err)
				}

				sqlStatement := `SELECT datname FROM pg_catalog.pg_database`
				rows, err := db.Query(sqlStatement)
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
				fmt.Println(custSchemaCount)
				fmt.Println(schemas)
				fmt.Println(availableSlots)
				dbDetail := SaasDbDetail{DbHost: dbHost, CustomerSchemas: schemas, CustomerSchemaCount: custSchemaCount, AvailableSchemaSlots: availableSlots}
				dbDetails = append(dbDetails, dbDetail)
				c.JSON(200, dbDetails)
			}
		}
	})
	r.Run(":8080")
}
