PROJECT_NAME: book-lending

env-up:
	docker-compose up -d db

env-down:
	docker-compose down

clean:
	rm -rf build/

build:
	CGO_ENABLED=0 GOOS=linux go build -o build/bin cmd/book-lending/main.go

run: env-up
	go run cmd/book-lending/main.go

run-docker: clean build
	docker-compose up --build app

test:
	go test ./...

cover:
	go test -cover ./...