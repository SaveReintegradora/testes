up:
	docker compose -f docker-compose.yml up -d

down:
	docker compose -f docker-compose.yml down

clean:
	teardown

rebuild:
	docker compose -f docker-compose.yml up -d --build --force-recreate --no-deps api