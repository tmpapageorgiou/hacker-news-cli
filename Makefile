

run:
	go run main.go api.go client.go formatter.go

test:
	go test -race -cover
