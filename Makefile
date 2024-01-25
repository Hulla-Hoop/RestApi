run:
	@go run ./cmd/app

#Команда для запуска сервиса локально
run-local:
	@docker-compose up -d postgres
	@time sleep 5
	@go run ./cmd/app

#Команда для запуска сервиса в докер образе
run-docker:
	@docker-compose up -d postgres
	@time sleep 5
	@docker-compose up sobes-api
stop-docker:
	@docker-compose down