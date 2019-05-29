GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get -u

BINARY_NAME=rds_dbdetails
deps:
	$(GOGET) database/sql
	$(GOGET) github.com/gin-gonic/gin
	$(GOGET) github.com/aws/aws-sdk-go/aws
	$(GOGET) github.com/aws/aws-sdk-go/aws/session
	$(GOGET) github.com/aws/aws-sdk-go/service/rds
	$(GOGET) github.com/lib/pq
    #use verioning using dep (in path)
    # dep update !

build:  deps
	$(GOBUILD) -o $(BINARY_NAME) -v