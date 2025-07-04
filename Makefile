up:
	docker compose -f docker-compose.yml up -d --build --force-recreate --no-deps api

down:
	docker compose -f docker-compose.yml down

clean:
	docker compose -f docker-compose.yml down -v --rmi all

rebuild:
	docker compose -f docker-compose.yml up -d --build --force-recreate --no-deps api

# ATENÇÃO: 'total' remove TODOS os dados do banco (volumes) e imagens!
total:
	docker-compose down -v
	docker system prune -af
	docker-compose up --build

# Use este comando para reiniciar a API e o banco SEM perder os dados
restart:
	docker compose -f docker-compose.yml down
	docker compose -f docker-compose.yml up -d --build