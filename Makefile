default: start
start:
	docker-compose up --build -d
stop:
	docker-compose down
tests:
	go test ./test/... -count=1 -p 1
