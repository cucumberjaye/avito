migrate:
	migrate -path ./scheme -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up

drun:
	docker run --name=balance -e POSTRGES_PASSWORD=qwerty -p 5436:5432 -d --rm postgres
