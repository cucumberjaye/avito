migrate:
	migrate -path ./scheme -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up
