# This project is rejected by reviewers

Here is the feedback:

> We think it would be beneficial to further develop certain Go-specific topics, such as pointer usage, goroutine management, channel usage, concurrency structures, process control with context, graceful shutdown procedures, error handling and recovery mechanisms, and memory management.

Not sure how to improve the project according to the feedback. I think the project is already using all the mentioned topics.

Anyway I hope it will be helpful for someone.

## API SETUP INSTRUCTIONS

1. Clone the repository
2. Run `docker compose up` to start the application
3. Run `go run main.go populate` inside main container
4. Hit http://localhost:8081 to access the application
5. Hit http://localhost:8080/swagger/index.html to access Swagger API DOC

## HOW TO RUN MESSAGE ENGINE
1. Log into the main container
2. `/root/go/src/github.com/mehgokalp/insider-project/sbin/app engine:message`

## HOW TO POPULATE DUMMY DATA
1. Log into the main container
2. `/root/go/src/github.com/mehgokalp/insider-project/sbin/app populate`

## HOW TO RUN TESTS
```go
go test ./...
```