# Golang REST API to get rds-db-details

**FileName:** rds_dbdetails.go

A simple rest api that gets rds-db-details.

*Use below command to run the api:*

    #go run rds_dbdetails.go	
    #curl localhost:8080/dbdetails


**FileName:** rds_dbdetails-v1.go

Rest API(concurrency developed) that gets rds-db-details.

*Use below command to run the api:*

    #go run rds_dbdetails-v1.go
	#curl localhost:8080/dbdetails


**FileName:** main.go

An API that gets rds-db-details developed to deploy in AWS Lambda function.

**Lambda function name:** saas-lambda-dbdetails

**API gateway name:** saas-dbdetails-api

*Use below command to run the api:*

    #curl https://b6pjzp3nmc.execute-api.us-east-1.amazonaws.com/dbdetails/dbdetails

*Use below url for postman(GET) or browser:*

    https://b6pjzp3nmc.execute-api.us-east-1.amazonaws.com/dbdetails/dbdetails
