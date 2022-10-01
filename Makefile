migrate:
	migrate -path ./scheme -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up

drun:
	docker run --name balance -e POSTGRES_PASSWORD=qwerty -p 5432:5432 -d --rm postgres

run:
	go run cmd/main.go
