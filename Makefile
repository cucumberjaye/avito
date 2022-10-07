migrate:
	migrate -path ./scheme -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up

run:
	docker-compose up

test:
	go test -v ./...


