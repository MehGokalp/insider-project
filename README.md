## API SETUP INSTRUCTIONS

1. Clone the repository
2. Run `docker compose up` to start the application
3. Run `go run main.go populate` inside main container
4. Hit http://localhost:8081 to access the application
5. Hit http://localhost:8080/swagger/index.html to access Swagger API DOC

## HOW TO RUN MESSAGE ENGINE
1. Log into the main container
2. `/root/go/src/github.com/mehgokalp/insider-project/sbin/app engine:message`

## HOW TO RUN TESTS
```go
go test ./...
```