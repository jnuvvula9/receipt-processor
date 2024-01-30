build:
	@go build -o bin/receipt-processor
run: build
	@./bin/receipt-processor
test:
	@go test -v ./...