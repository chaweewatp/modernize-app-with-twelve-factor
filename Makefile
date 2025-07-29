run:
	docker compose up --build -d

restart:
	docker compose restart

down:
	docker compose down -v