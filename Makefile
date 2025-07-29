run:
	docker compose up --build -d

restart:
	docker compose restart

down:
	docker compose down -v

create_env:
	@export APP_PORT=8080 DATABASE_DSN="host=db user=postgres password=postgres dbname=mydb port=5432 sslmode=disable" REDIS_ADDR="cache:6380"; \
	echo "Env variables set for this shell"

list_env:
	env | grep APP_PORT
	env | grep DATABASE_DSN
	env | grep REDIS_ADDR

