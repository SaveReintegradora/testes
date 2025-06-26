up:
	docker compose -f docker-compose.yml up -d --build --force-recreate --no-deps api

down:
	docker compose -f docker-compose.yml down

clean:
	docker compose -f docker-compose.yml down -v --rmi all

rebuild:
	docker compose -f docker-compose.yml up -d --build --force-recreate --no-deps api