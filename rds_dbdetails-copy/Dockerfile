FROM golang:latest
COPY ./rds_dbdetails.go Makefile /app/
WORKDIR /app
RUN make build
EXPOSE 8080
CMD ["/bin/bash", "-c", "./rds_dbdetails"]