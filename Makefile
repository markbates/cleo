test:
	go test -count 1 -race -cover ./...

cov:
	go test -short -coverprofile cover.out ./...
	go tool cover -html cover.out
	go mod tidy -v
